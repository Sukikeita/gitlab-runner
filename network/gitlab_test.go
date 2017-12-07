package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "gitlab.com/gitlab-org/gitlab-runner/common"
)

var brokenCredentials = RunnerCredentials{
	URL: "broken",
}

var brokenConfig = RunnerConfig{
	RunnerCredentials: brokenCredentials,
}

func TestClients(t *testing.T) {
	c := NewGitLabClient()
	c1, _ := c.getClient(&RunnerCredentials{
		URL: "http://test/",
	})
	c2, _ := c.getClient(&RunnerCredentials{
		URL: "http://test2/",
	})
	c4, _ := c.getClient(&RunnerCredentials{
		URL:       "http://test/",
		TLSCAFile: "ca_file",
	})
	c5, _ := c.getClient(&RunnerCredentials{
		URL:       "http://test/",
		TLSCAFile: "ca_file",
	})
	c6, _ := c.getClient(&RunnerCredentials{
		URL:         "http://test/",
		TLSCAFile:   "ca_file",
		TLSCertFile: "cert_file",
		TLSKeyFile:  "key_file",
	})
	c7, _ := c.getClient(&RunnerCredentials{
		URL:         "http://test/",
		TLSCAFile:   "ca_file",
		TLSCertFile: "cert_file",
		TLSKeyFile:  "key_file2",
	})
	c8, c8err := c.getClient(&brokenCredentials)
	assert.NotEqual(t, c1, c2)
	assert.NotEqual(t, c1, c4)
	assert.Equal(t, c4, c5)
	assert.NotEqual(t, c5, c6)
	assert.Equal(t, c6, c7)
	assert.Nil(t, c8)
	assert.Error(t, c8err)
}

func testRegisterRunnerHandler(w http.ResponseWriter, r *http.Request, t *testing.T) {
	if r.URL.Path != "/api/v4/runners" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	var req map[string]interface{}
	err = json.Unmarshal(body, &req)
	assert.NoError(t, err)

	res := make(map[string]interface{})

	switch req["token"].(string) {
	case "valid":
		if req["description"].(string) != "test" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res["token"] = req["token"].(string)
	case "invalid":
		w.WriteHeader(http.StatusForbidden)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Accept") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(output)
}

func TestRegisterRunner(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testRegisterRunnerHandler(w, r, t)
	}))
	defer s.Close()

	validToken := RunnerCredentials{
		URL:   s.URL,
		Token: "valid",
	}

	invalidToken := RunnerCredentials{
		URL:   s.URL,
		Token: "invalid",
	}

	otherToken := RunnerCredentials{
		URL:   s.URL,
		Token: "other",
	}

	c := NewGitLabClient()

	res := c.RegisterRunner(validToken, "test", "tags", true, true)
	if assert.NotNil(t, res) {
		assert.Equal(t, validToken.Token, res.Token)
	}

	res = c.RegisterRunner(validToken, "invalid description", "tags", true, true)
	assert.Nil(t, res)

	res = c.RegisterRunner(invalidToken, "test", "tags", true, true)
	assert.Nil(t, res)

	res = c.RegisterRunner(otherToken, "test", "tags", true, true)
	assert.Nil(t, res)

	res = c.RegisterRunner(brokenCredentials, "test", "tags", true, true)
	assert.Nil(t, res)
}

func testUnregisterRunnerHandler(w http.ResponseWriter, r *http.Request, t *testing.T) {
	if r.URL.Path != "/api/v4/runners" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	var req map[string]interface{}
	err = json.Unmarshal(body, &req)
	assert.NoError(t, err)

	switch req["token"].(string) {
	case "valid":
		w.WriteHeader(http.StatusNoContent)
	case "invalid":
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func TestUnregisterRunner(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testUnregisterRunnerHandler(w, r, t)
	}

	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	validToken := RunnerCredentials{
		URL:   s.URL,
		Token: "valid",
	}

	invalidToken := RunnerCredentials{
		URL:   s.URL,
		Token: "invalid",
	}

	otherToken := RunnerCredentials{
		URL:   s.URL,
		Token: "other",
	}

	c := NewGitLabClient()

	state := c.UnregisterRunner(validToken)
	assert.True(t, state)

	state = c.UnregisterRunner(invalidToken)
	assert.False(t, state)

	state = c.UnregisterRunner(otherToken)
	assert.False(t, state)

	state = c.UnregisterRunner(brokenCredentials)
	assert.False(t, state)
}

func testVerifyRunnerHandler(w http.ResponseWriter, r *http.Request, t *testing.T) {
	if r.URL.Path != "/api/v4/runners/verify" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	var req map[string]interface{}
	err = json.Unmarshal(body, &req)
	assert.NoError(t, err)

	switch req["token"].(string) {
	case "valid":
		w.WriteHeader(http.StatusOK) // since the job id is broken, we should not find this job
	case "invalid":
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func TestVerifyRunner(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testVerifyRunnerHandler(w, r, t)
	}

	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	validToken := RunnerCredentials{
		URL:   s.URL,
		Token: "valid",
	}

	invalidToken := RunnerCredentials{
		URL:   s.URL,
		Token: "invalid",
	}

	otherToken := RunnerCredentials{
		URL:   s.URL,
		Token: "other",
	}

	c := NewGitLabClient()

	state := c.VerifyRunner(validToken)
	assert.True(t, state)

	state = c.VerifyRunner(invalidToken)
	assert.False(t, state)

	state = c.VerifyRunner(otherToken)
	assert.True(t, state, "in other cases where we can't explicitly say that runner is valid we say that it's")

	state = c.VerifyRunner(brokenCredentials)
	assert.True(t, state, "in other cases where we can't explicitly say that runner is valid we say that it's")
}

func getRequestJobResponse() (res map[string]interface{}) {
	jobToken := "job-token"

	res = make(map[string]interface{})
	res["id"] = 10
	res["token"] = jobToken
	res["allow_git_fetch"] = false

	jobInfo := make(map[string]interface{})
	jobInfo["name"] = "test-job"
	jobInfo["stage"] = "test"
	jobInfo["project_id"] = 123
	jobInfo["project_name"] = "test-project"
	res["job_info"] = jobInfo

	gitInfo := make(map[string]interface{})
	gitInfo["repo_url"] = "https://gitlab-ci-token:testTokenHere1234@gitlab.example.com/test/test-project.git"
	gitInfo["ref"] = "master"
	gitInfo["sha"] = "abcdef123456"
	gitInfo["before_sha"] = "654321fedcba"
	gitInfo["ref_type"] = "branch"
	res["git_info"] = gitInfo

	runnerInfo := make(map[string]interface{})
	runnerInfo["timeout"] = 3600
	res["runner_info"] = runnerInfo

	variables := make([]map[string]interface{}, 1)
	variables[0] = make(map[string]interface{})
	variables[0]["key"] = "CI_REF_NAME"
	variables[0]["value"] = "master"
	variables[0]["public"] = true
	variables[0]["file"] = true
	res["variables"] = variables

	steps := make([]map[string]interface{}, 2)
	steps[0] = make(map[string]interface{})
	steps[0]["name"] = "script"
	steps[0]["script"] = []string{"date", "ls -ls"}
	steps[0]["timeout"] = 3600
	steps[0]["when"] = "on_success"
	steps[0]["allow_failure"] = false
	steps[1] = make(map[string]interface{})
	steps[1]["name"] = "after_script"
	steps[1]["script"] = []string{"ls -ls"}
	steps[1]["timeout"] = 3600
	steps[1]["when"] = "always"
	steps[1]["allow_failure"] = true
	res["steps"] = steps

	image := make(map[string]interface{})
	image["name"] = "ruby:2.0"
	image["entrypoint"] = []string{"/bin/sh"}
	res["image"] = image

	services := make([]map[string]interface{}, 2)
	services[0] = make(map[string]interface{})
	services[0]["name"] = "postgresql:9.5"
	services[0]["entrypoint"] = []string{"/bin/sh"}
	services[0]["command"] = []string{"sleep", "30"}
	services[0]["alias"] = "db-pg"
	services[1] = make(map[string]interface{})
	services[1]["name"] = "mysql:5.6"
	services[1]["alias"] = "db-mysql"
	res["services"] = services

	artifacts := make([]map[string]interface{}, 1)
	artifacts[0] = make(map[string]interface{})
	artifacts[0]["name"] = "artifact.zip"
	artifacts[0]["untracked"] = false
	artifacts[0]["paths"] = []string{"out/*"}
	artifacts[0]["when"] = "always"
	artifacts[0]["expire_in"] = "7d"
	res["artifacts"] = artifacts

	cache := make([]map[string]interface{}, 1)
	cache[0] = make(map[string]interface{})
	cache[0]["key"] = "$CI_COMMIT_REF"
	cache[0]["untracked"] = false
	cache[0]["paths"] = []string{"vendor/*"}
	cache[0]["policy"] = "push"
	res["cache"] = cache

	credentials := make([]map[string]interface{}, 1)
	credentials[0] = make(map[string]interface{})
	credentials[0]["type"] = "Registry"
	credentials[0]["url"] = "http://registry.gitlab.example.com/"
	credentials[0]["username"] = "gitlab-ci-token"
	credentials[0]["password"] = jobToken
	res["credentials"] = credentials

	dependencies := make([]map[string]interface{}, 1)
	dependencies[0] = make(map[string]interface{})
	dependencies[0]["id"] = 9
	dependencies[0]["name"] = "other-job"
	dependencies[0]["token"] = "other-job-token"
	artifactsFile0 := make(map[string]interface{})
	artifactsFile0["filename"] = "binaries.zip"
	artifactsFile0["size"] = 13631488
	dependencies[0]["artifacts_file"] = artifactsFile0
	res["dependencies"] = dependencies

	return
}

func testRequestJobHandler(w http.ResponseWriter, r *http.Request, t *testing.T) {
	if r.URL.Path != "/api/v4/jobs/request" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	var req map[string]interface{}
	err = json.Unmarshal(body, &req)
	assert.NoError(t, err)

	switch req["token"].(string) {
	case "valid":
	case "no-jobs":
		w.Header().Add("X-GitLab-Last-Update", "a nice timestamp")
		w.WriteHeader(http.StatusNoContent)
		return
	case "invalid":
		w.WriteHeader(http.StatusForbidden)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Accept") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := json.Marshal(getRequestJobResponse())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(output)
	t.Logf("JobRequest response: %s\n", output)
}

func TestRequestJob(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testRequestJobHandler(w, r, t)
	}))
	defer s.Close()

	validToken := RunnerConfig{
		RunnerCredentials: RunnerCredentials{
			URL:   s.URL,
			Token: "valid",
		},
	}

	noJobsToken := RunnerConfig{
		RunnerCredentials: RunnerCredentials{
			URL:   s.URL,
			Token: "no-jobs",
		},
	}

	invalidToken := RunnerConfig{
		RunnerCredentials: RunnerCredentials{
			URL:   s.URL,
			Token: "invalid",
		},
	}

	c := NewGitLabClient()

	res, ok := c.RequestJob(validToken)
	if assert.NotNil(t, res) {
		assert.NotEmpty(t, res.ID)
	}
	assert.True(t, ok)

	assert.Equal(t, "ruby:2.0", res.Image.Name)
	assert.Equal(t, []string{"/bin/sh"}, res.Image.Entrypoint)
	require.Len(t, res.Services, 2)
	assert.Equal(t, "postgresql:9.5", res.Services[0].Name)
	assert.Equal(t, []string{"/bin/sh"}, res.Services[0].Entrypoint)
	assert.Equal(t, []string{"sleep", "30"}, res.Services[0].Command)
	assert.Equal(t, "db-pg", res.Services[0].Alias)
	assert.Equal(t, "mysql:5.6", res.Services[1].Name)
	assert.Equal(t, "db-mysql", res.Services[1].Alias)

	assert.Empty(t, c.getLastUpdate(&noJobsToken.RunnerCredentials), "Last-Update should not be set")
	res, ok = c.RequestJob(noJobsToken)
	assert.Nil(t, res)
	assert.True(t, ok, "If no jobs, runner is healthy")
	assert.Equal(t, c.getLastUpdate(&noJobsToken.RunnerCredentials), "a nice timestamp", "Last-Update should be set")

	res, ok = c.RequestJob(invalidToken)
	assert.Nil(t, res)
	assert.False(t, ok, "If token is invalid, the runner is unhealthy")

	res, ok = c.RequestJob(brokenConfig)
	assert.Nil(t, res)
	assert.False(t, ok)
}

func testUpdateJobHandler(w http.ResponseWriter, r *http.Request, t *testing.T) {
	if r.URL.Path != "/api/v4/jobs/10" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "PUT" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	var req map[string]interface{}
	err = json.Unmarshal(body, &req)
	assert.NoError(t, err)

	assert.Equal(t, "token", req["token"])
	assert.Equal(t, "trace", req["trace"])

	switch req["state"].(string) {
	case "running":
		w.WriteHeader(http.StatusOK)
	case "forbidden":
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func TestUpdateJob(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testUpdateJobHandler(w, r, t)
	}

	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	config := RunnerConfig{
		RunnerCredentials: RunnerCredentials{
			URL: s.URL,
		},
	}

	jobCredentials := &JobCredentials{
		Token: "token",
	}

	trace := "trace"
	c := NewGitLabClient()

	state := c.UpdateJob(config, jobCredentials, 10, "running", &trace)
	assert.Equal(t, UpdateSucceeded, state, "Update should continue when running")

	state = c.UpdateJob(config, jobCredentials, 10, "forbidden", &trace)
	assert.Equal(t, UpdateAbort, state, "Update should if the state is forbidden")

	state = c.UpdateJob(config, jobCredentials, 10, "other", &trace)
	assert.Equal(t, UpdateFailed, state, "Update should fail for badly formatted request")

	state = c.UpdateJob(config, jobCredentials, 4, "state", &trace)
	assert.Equal(t, UpdateAbort, state, "Update should abort for unknown job")

	state = c.UpdateJob(brokenConfig, jobCredentials, 4, "state", &trace)
	assert.Equal(t, UpdateAbort, state)
}

var patchToken = "token"
var patchTraceString = "trace trace trace"

func getPatchServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request, body string, offset, limit int)) (*httptest.Server, *GitLabClient, RunnerConfig) {
	patchHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v4/jobs/1/trace" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Method != "PATCH" {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		assert.Equal(t, patchToken, r.Header.Get("JOB-TOKEN"))

		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		contentRange := r.Header.Get("Content-Range")
		ranges := strings.Split(contentRange, "-")

		offset, err := strconv.Atoi(ranges[0])
		assert.NoError(t, err)

		limit, err := strconv.Atoi(ranges[1])
		assert.NoError(t, err)

		handler(w, r, string(body), offset, limit)
	}

	server := httptest.NewServer(http.HandlerFunc(patchHandler))

	config := RunnerConfig{
		RunnerCredentials: RunnerCredentials{
			URL: server.URL,
		},
	}

	return server, NewGitLabClient(), config
}

func getTracePatch(traceString string, offset int) *tracePatch {
	trace := bytes.Buffer{}
	trace.WriteString(traceString)
	tracePatch, _ := newTracePatch(trace, offset)

	return tracePatch
}

func TestUnknownPatchTrace(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request, body string, offset, limit int) {
		w.WriteHeader(http.StatusNotFound)
	}

	server, client, config := getPatchServer(t, handler)
	defer server.Close()

	tracePatch := getTracePatch(patchTraceString, 0)
	state := client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateNotFound, state)
}

func TestForbiddenPatchTrace(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request, body string, offset, limit int) {
		w.WriteHeader(http.StatusForbidden)
	}

	server, client, config := getPatchServer(t, handler)
	defer server.Close()

	tracePatch := getTracePatch(patchTraceString, 0)
	state := client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateAbort, state)
}

func TestPatchTrace(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request, body string, offset, limit int) {
		assert.Equal(t, patchTraceString[offset:limit], body)

		w.WriteHeader(http.StatusAccepted)
	}

	server, client, config := getPatchServer(t, handler)
	defer server.Close()

	tracePatch := getTracePatch(patchTraceString, 0)
	state := client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateSucceeded, state)

	tracePatch = getTracePatch(patchTraceString, 3)
	state = client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateSucceeded, state)

	tracePatch = getTracePatch(patchTraceString[:10], 3)
	state = client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateSucceeded, state)
}

func TestRangeMismatchPatchTrace(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request, body string, offset, limit int) {
		if offset > 10 {
			w.Header().Set("Range", "0-10")
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		}

		w.WriteHeader(http.StatusAccepted)
	}

	server, client, config := getPatchServer(t, handler)
	defer server.Close()

	tracePatch := getTracePatch(patchTraceString, 11)
	state := client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateRangeMismatch, state)

	tracePatch = getTracePatch(patchTraceString, 15)
	state = client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateRangeMismatch, state)

	tracePatch = getTracePatch(patchTraceString, 5)
	state = client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateSucceeded, state)
}

func TestResendPatchTrace(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request, body string, offset, limit int) {
		if offset > 10 {
			w.Header().Set("Range", "0-10")
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		}

		w.WriteHeader(http.StatusAccepted)
	}

	server, client, config := getPatchServer(t, handler)
	defer server.Close()

	tracePatch := getTracePatch(patchTraceString, 11)
	state := client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateRangeMismatch, state)

	state = client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateSucceeded, state)
}

// We've had a situation where the same job was triggered second time by GItLab. In GitLab the job trace
// was 17041 bytes long while the repeated job trace was only 66 bytes long. We've had a `RangeMismatch`
// response, so the offset was updated (to 17041) and `client.PatchTrace` was repeated, at it was planned.
// Unfortunately the `tracePatch` struct was  not resistant to a situation when the offset is set to a
// value bigger than trace's length. This test simulates such situation.
func TestResendDoubledJobPatchTrace(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request, body string, offset, limit int) {
		if offset > 10 {
			w.Header().Set("Range", "0-100")
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		}

		w.WriteHeader(http.StatusAccepted)
	}

	server, client, config := getPatchServer(t, handler)
	defer server.Close()

	tracePatch := getTracePatch(patchTraceString, 11)
	state := client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateRangeMismatch, state)
	assert.False(t, tracePatch.ValidateRange())
}

func TestJobFailedStatePatchTrace(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request, body string, offset, limit int) {
		w.Header().Set("Job-Status", "failed")
		w.WriteHeader(http.StatusAccepted)
	}

	server, client, config := getPatchServer(t, handler)
	defer server.Close()

	tracePatch := getTracePatch(patchTraceString, 0)
	state := client.PatchTrace(config, &JobCredentials{ID: 1, Token: patchToken}, tracePatch)
	assert.Equal(t, UpdateAbort, state)
}

func testArtifactsUploadHandler(w http.ResponseWriter, r *http.Request, t *testing.T) {
	if r.URL.Path != "/api/v4/jobs/10/artifacts" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if r.Header.Get("JOB-TOKEN") != "token" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(file)
	assert.NoError(t, err)

	if string(body) != "content" {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func TestArtifactsUpload(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testArtifactsUploadHandler(w, r, t)
	}

	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	config := JobCredentials{
		ID:    10,
		URL:   s.URL,
		Token: "token",
	}
	invalidToken := JobCredentials{
		ID:    10,
		URL:   s.URL,
		Token: "invalid-token",
	}

	tempFile, err := ioutil.TempFile("", "artifacts")
	assert.NoError(t, err)
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	c := NewGitLabClient()

	fmt.Fprint(tempFile, "content")
	state := c.UploadArtifacts(config, tempFile.Name())
	assert.Equal(t, UploadSucceeded, state, "Artifacts should be uploaded")

	fmt.Fprint(tempFile, "too large")
	state = c.UploadArtifacts(config, tempFile.Name())
	assert.Equal(t, UploadTooLarge, state, "Artifacts should be not uploaded, because of too large archive")

	state = c.UploadArtifacts(config, "not/existing/file")
	assert.Equal(t, UploadFailed, state, "Artifacts should fail to be uploaded")

	state = c.UploadArtifacts(invalidToken, tempFile.Name())
	assert.Equal(t, UploadForbidden, state, "Artifacts should be rejected if invalid token")
}

func testArtifactsDownloadHandler(w http.ResponseWriter, r *http.Request, t *testing.T) {
	if r.URL.Path != "/api/v4/jobs/10/artifacts" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if r.Header.Get("JOB-TOKEN") != "token" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes.NewBufferString("Test artifact file content").Bytes())
}

func TestArtifactsDownload(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		testArtifactsDownloadHandler(w, r, t)
	}

	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	credentials := JobCredentials{
		ID:    10,
		URL:   s.URL,
		Token: "token",
	}
	invalidTokenCredentials := JobCredentials{
		ID:    10,
		URL:   s.URL,
		Token: "invalid-token",
	}
	fileNotFoundTokenCredentials := JobCredentials{
		ID:    11,
		URL:   s.URL,
		Token: "token",
	}

	c := NewGitLabClient()

	tempDir, err := ioutil.TempDir("", "artifacts")
	assert.NoError(t, err)

	artifactsFileName := filepath.Join(tempDir, "downloaded-artifact")
	defer os.Remove(artifactsFileName)

	state := c.DownloadArtifacts(credentials, artifactsFileName)
	assert.Equal(t, DownloadSucceeded, state, "Artifacts should be downloaded")

	state = c.DownloadArtifacts(invalidTokenCredentials, artifactsFileName)
	assert.Equal(t, DownloadForbidden, state, "Artifacts should be not downloaded if invalid token is used")

	state = c.DownloadArtifacts(fileNotFoundTokenCredentials, artifactsFileName)
	assert.Equal(t, DownloadNotFound, state, "Artifacts should be bit downloaded if it's not found")
}
