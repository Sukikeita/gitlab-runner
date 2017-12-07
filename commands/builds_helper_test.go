package commands

import (
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/gitlab-runner/common"
)

var fakeRunner = &common.RunnerConfig{
	RunnerCredentials: common.RunnerCredentials{
		Token: "a1b2c3d4e5f6",
	},
}

func TestBuildsHelperCollect(t *testing.T) {
	ch := make(chan prometheus.Metric, 50)
	b := &buildsHelper{}
	b.builds = append(b.builds, &common.Build{
		CurrentState: common.BuildRunStatePending,
		CurrentStage: common.BuildStagePrepare,
		Runner:       fakeRunner,
	})
	b.Collect(ch)
	assert.Len(t, ch, 1)
}

func TestBuildsHelperAcquireRequestWithLimit(t *testing.T) {
	runner := common.RunnerConfig{
		RequestConcurrency: 2,
	}

	b := &buildsHelper{}
	result := b.acquireRequest(&runner)
	require.True(t, result)

	result = b.acquireRequest(&runner)
	require.True(t, result)

	result = b.acquireRequest(&runner)
	require.False(t, result, "allow only two requests")

	result = b.releaseRequest(&runner)
	require.True(t, result)

	result = b.releaseRequest(&runner)
	require.True(t, result)

	result = b.releaseRequest(&runner)
	require.False(t, result, "release only two requests")
}

func TestBuildsHelperAcquireRequestWithDefault(t *testing.T) {
	runner := common.RunnerConfig{
		RequestConcurrency: 0,
	}

	b := &buildsHelper{}
	result := b.acquireRequest(&runner)
	require.True(t, result)

	result = b.acquireRequest(&runner)
	require.False(t, result, "allow only one request")

	result = b.releaseRequest(&runner)
	require.True(t, result)

	result = b.releaseRequest(&runner)
	require.False(t, result, "release only one request")

	result = b.acquireRequest(&runner)
	require.True(t, result)

	result = b.releaseRequest(&runner)
	require.True(t, result)

	result = b.releaseRequest(&runner)
	require.False(t, result, "nothing to release")
}

func TestBuildsHelperAcquireBuildWithLimit(t *testing.T) {
	runner := common.RunnerConfig{
		Limit: 1,
	}

	b := &buildsHelper{}
	result := b.acquireBuild(&runner)
	require.True(t, result)

	result = b.acquireBuild(&runner)
	require.False(t, result, "allow only one build")

	result = b.releaseBuild(&runner)
	require.True(t, result)

	result = b.releaseBuild(&runner)
	require.False(t, result, "release only one build")
}

func TestBuildsHelperAcquireBuildUnlimited(t *testing.T) {
	runner := common.RunnerConfig{
		Limit: 0,
	}

	b := &buildsHelper{}
	result := b.acquireBuild(&runner)
	require.True(t, result)

	result = b.acquireBuild(&runner)
	require.True(t, result)

	result = b.releaseBuild(&runner)
	require.True(t, result)

	result = b.releaseBuild(&runner)
	require.True(t, result)
}

var testBuildCurrentID int

func getTestBuild() *common.Build {
	testBuildCurrentID++

	runner := common.RunnerConfig{}
	runner.Token = "a1b2c3d4"
	jobInfo := common.JobInfo{
		ProjectID: 1,
	}

	build := &common.Build{}
	build.ID = testBuildCurrentID
	build.Runner = &runner
	build.JobInfo = jobInfo

	return build
}

func concurrentlyCallStatesAndStages(b *buildsHelper, finished chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-finished:
			wg.Done()
			return
		case <-time.After(1 * time.Millisecond):
			b.statesAndStages()
		}
	}
}

func handleTestBuild(wg *sync.WaitGroup, b *buildsHelper, build *common.Build) {
	stages := []common.BuildStage{
		common.BuildStagePrepare,
		common.BuildStageGetSources,
		common.BuildStageRestoreCache,
		common.BuildStageDownloadArtifacts,
		common.BuildStageUserScript,
		common.BuildStageAfterScript,
		common.BuildStageArchiveCache,
		common.BuildStageUploadArtifacts,
	}
	states := []common.BuildRuntimeState{
		common.BuildRunStatePending,
		common.BuildRunRuntimeRunning,
		common.BuildRunRuntimeFinished,
	}

	b.addBuild(build)
	time.Sleep(10 * time.Millisecond)
	for _, stage := range stages {
		for _, state := range states {
			build.CurrentStage = stage
			build.CurrentState = state
			time.Sleep(5 * time.Millisecond)
		}
	}
	time.Sleep(5 * time.Millisecond)
	b.removeBuild(build)

	time.Sleep(5 * time.Millisecond)
	wg.Done()
}

func TestBuildHelperCollectWhenRemovingBuild(t *testing.T) {
	t.Log("This test tries to simulate concurrency. It can be false-positive! But failure is always a sign that something is wrong.")
	b := &buildsHelper{}
	statesAndStagesCallConcurrency := 10
	finished := make(chan struct{})

	wg1 := &sync.WaitGroup{}
	wg1.Add(statesAndStagesCallConcurrency)
	for i := 0; i < statesAndStagesCallConcurrency; i++ {
		go concurrentlyCallStatesAndStages(b, finished, wg1)
	}

	var builds []*common.Build
	for i := 1; i < 30; i++ {
		builds = append(builds, getTestBuild())
	}

	wg2 := &sync.WaitGroup{}
	wg2.Add(len(builds))
	for _, build := range builds {
		go handleTestBuild(wg2, b, build)
	}
	wg2.Wait()

	close(finished)
	wg1.Wait()
}
