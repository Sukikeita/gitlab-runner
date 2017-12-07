package network

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"gitlab.com/gitlab-org/gitlab-runner/common"
	"gitlab.com/gitlab-org/gitlab-runner/helpers"
)

var traceUpdateInterval = common.UpdateInterval
var traceForceSendInterval = common.ForceTraceSentInterval
var traceFinishRetryInterval = common.UpdateRetryInterval

type tracePatch struct {
	trace  bytes.Buffer
	offset int
	limit  int
}

func (tp *tracePatch) Patch() []byte {
	return tp.trace.Bytes()[tp.offset:tp.limit]
}

func (tp *tracePatch) Offset() int {
	return tp.offset
}

func (tp *tracePatch) Limit() int {
	return tp.limit
}

func (tp *tracePatch) SetNewOffset(newOffset int) {
	tp.offset = newOffset
}

func (tp *tracePatch) ValidateRange() bool {
	if tp.limit >= tp.offset {
		return true
	}

	return false
}

func newTracePatch(trace bytes.Buffer, offset int) (*tracePatch, error) {
	patch := &tracePatch{
		trace:  trace,
		offset: offset,
		limit:  trace.Len(),
	}

	if !patch.ValidateRange() {
		return nil, errors.New("Range is invalid, limit can't be less than offset")
	}

	return patch, nil
}

type clientJobTrace struct {
	*io.PipeWriter

	client         common.Network
	config         common.RunnerConfig
	jobCredentials *common.JobCredentials
	id             int
	limit          int64
	cancelFunc     context.CancelFunc

	log      bytes.Buffer
	lock     sync.RWMutex
	state    common.JobState
	finished chan bool

	sentTrace int
	sentTime  time.Time
	sentState common.JobState
}

func (c *clientJobTrace) Success() {
	c.Fail(nil)
}

func (c *clientJobTrace) Fail(err error) {
	c.lock.Lock()
	if c.state != common.Running {
		c.lock.Unlock()
		return
	}
	if err == nil {
		c.state = common.Success
	} else {
		c.state = common.Failed
	}
	c.lock.Unlock()

	c.finish()
}

func (c *clientJobTrace) SetCancelFunc(cancelFunc context.CancelFunc) {
	c.cancelFunc = cancelFunc
}

func (c *clientJobTrace) IsStdout() bool {
	return false
}

func (c *clientJobTrace) start() {
	reader, writer := io.Pipe()
	c.PipeWriter = writer
	c.finished = make(chan bool)
	c.state = common.Running
	go c.process(reader)
	go c.watch()
}

func (c *clientJobTrace) finish() {
	c.Close()
	c.finished <- true

	// Do final upload of job trace
	for {
		if c.fullUpdate() != common.UpdateFailed {
			return
		}
		time.Sleep(traceFinishRetryInterval)
	}
}

func (c *clientJobTrace) writeRune(r rune, limit int) (n int, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	n, err = c.log.WriteRune(r)
	if c.log.Len() < limit {
		return
	}

	output := fmt.Sprintf("\n%sJob's log exceeded limit of %v bytes.%s\n",
		helpers.ANSI_BOLD_RED,
		limit,
		helpers.ANSI_RESET,
	)
	c.log.WriteString(output)
	err = io.EOF
	return
}

func (c *clientJobTrace) process(pipe *io.PipeReader) {
	defer pipe.Close()

	stopped := false
	limit := c.config.OutputLimit
	if limit == 0 {
		limit = common.DefaultOutputLimit
	}
	limit *= 1024

	reader := bufio.NewReader(pipe)
	for {
		r, s, err := reader.ReadRune()
		if s <= 0 {
			break
		} else if stopped {
			// ignore symbols if job log exceeded limit
			continue
		} else if err == nil {
			_, err = c.writeRune(r, limit)
			if err == io.EOF {
				stopped = true
			}
		} else {
			// ignore invalid characters
			continue
		}
	}
}

func (c *clientJobTrace) incrementalUpdate() common.UpdateState {
	c.lock.RLock()
	state := c.state
	trace := c.log
	c.lock.RUnlock()

	if c.sentState == state &&
		c.sentTrace == trace.Len() &&
		time.Since(c.sentTime) < traceForceSendInterval {
		return common.UpdateSucceeded
	}

	if c.sentState != state {
		c.client.UpdateJob(c.config, c.jobCredentials, c.id, state, nil)
		c.sentState = state
	}

	tracePatch, err := newTracePatch(trace, c.sentTrace)
	if err != nil {
		c.config.Log().Errorln("Error while creating a tracePatch", err.Error())
	}

	update := c.client.PatchTrace(c.config, c.jobCredentials, tracePatch)
	if update == common.UpdateNotFound {
		return update
	}

	if update == common.UpdateRangeMismatch {
		update = c.resendPatch(c.jobCredentials.ID, c.config, c.jobCredentials, tracePatch)
	}

	if update == common.UpdateSucceeded {
		c.sentTrace = tracePatch.Limit()
		c.sentTime = time.Now()
	}

	return update
}

func (c *clientJobTrace) resendPatch(id int, config common.RunnerConfig, jobCredentials *common.JobCredentials, tracePatch common.JobTracePatch) (update common.UpdateState) {
	if !tracePatch.ValidateRange() {
		config.Log().Warningln(id, "Full job update is needed")
		fullTrace := c.log.String()

		return c.client.UpdateJob(c.config, jobCredentials, c.id, c.state, &fullTrace)
	}

	config.Log().Warningln(id, "Resending trace patch due to range mismatch")

	update = c.client.PatchTrace(config, jobCredentials, tracePatch)
	if update == common.UpdateRangeMismatch {
		config.Log().Errorln(id, "Appending trace to coordinator...", "failed due to range mismatch")

		return common.UpdateFailed
	}

	return
}

func (c *clientJobTrace) fullUpdate() common.UpdateState {
	c.lock.RLock()
	state := c.state
	trace := c.log.String()
	c.lock.RUnlock()

	if c.sentState == state &&
		c.sentTrace == len(trace) &&
		time.Since(c.sentTime) < traceForceSendInterval {
		return common.UpdateSucceeded
	}

	upload := c.client.UpdateJob(c.config, c.jobCredentials, c.id, state, &trace)
	if upload == common.UpdateSucceeded {
		c.sentTrace = len(trace)
		c.sentState = state
		c.sentTime = time.Now()
	}

	return upload
}

func (c *clientJobTrace) abort() bool {
	if c.cancelFunc != nil {
		c.cancelFunc()
		c.cancelFunc = nil
		return true
	}
	return false
}

func (c *clientJobTrace) watch() {
	for {
		select {
		case <-time.After(traceUpdateInterval):
			state := c.incrementalUpdate()
			if state == common.UpdateAbort && c.abort() {
				<-c.finished
				return
			}
			break

		case <-c.finished:
			return
		}
	}
}

func newJobTrace(client common.Network, config common.RunnerConfig, jobCredentials *common.JobCredentials) *clientJobTrace {
	return &clientJobTrace{
		client:         client,
		config:         config,
		jobCredentials: jobCredentials,
		id:             jobCredentials.ID,
	}
}
