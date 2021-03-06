package processes

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/rbeuque74/jagozzi/plugins"
	"github.com/stretchr/testify/assert"
)

func startProcesses(t *testing.T) func() {
	var pids []*os.Process
	cmd := exec.Command("/bin/sleep", "5", "2", "1")

	if err := cmd.Start(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	pids = append(pids, cmd.Process)

	cmd = exec.Command("/bin/sleep", "10")

	if err := cmd.Start(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	pids = append(pids, cmd.Process)

	cmd = exec.Command("/bin/sleep", "13")

	if err := cmd.Start(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	pids = append(pids, cmd.Process)

	cmd = exec.Command("/bin/sleep", "13")

	if err := cmd.Start(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	pids = append(pids, cmd.Process)

	return func() {
		for _, process := range pids {
			if process == nil {
				t.Error("process pid is nil")
				t.FailNow()
			}
			if err := process.Kill(); err != nil {
				t.Error(err)
				t.FailNow()
			}
		}
	}
}

func TestProcesses(t *testing.T) {
	teardown := startProcesses(t)
	defer teardown()

	// creating processes checker
	cfg := map[string]interface{}{
		"exec": "/bin/sleep",
		"args": "10",
		"name": "test-1",
	}
	checker, err := NewProcessesChecker(cfg, nil)
	assert.Nilf(t, err, "processes checker instantiation failed: %q", err)

	assert.Equal(t, "Processes", checker.Name())
	assert.Equal(t, "test-1", checker.ServiceName())

	ctxRun, cancelFunc1 := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc1()
	result := checker.Run(ctxRun)
	assert.Equal(t, plugins.STATE_OK, result.Status)
	assert.Equal(t, "Process /bin/sleep 10 is running", result.Message, "processes bad message: %q", result.Message)

	// multiple args
	cfg["args"] = "5 2 1"

	checker, err = NewProcessesChecker(cfg, nil)
	assert.Nilf(t, err, "processes checker instantiation failed: %q", err)

	ctxRun, cancelFunc1 = context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc1()
	result = checker.Run(ctxRun)
	assert.Equal(t, plugins.STATE_OK, result.Status)
	assert.Equal(t, "Process /bin/sleep 5 2 1 is running", result.Message, "processes bad message: %q", result.Message)

	// partial args
	cfg["args"] = "5"

	checker, err = NewProcessesChecker(cfg, nil)
	assert.Nilf(t, err, "processes checker instantiation failed: %q", err)

	ctxRun, cancelFunc1 = context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc1()
	result = checker.Run(ctxRun)
	assert.Equal(t, plugins.STATE_CRITICAL, result.Status)
	assert.Equal(t, "Process /bin/sleep 5 is not running", result.Message, "processes bad message: %q", result.Message)

	// different cmd
	cfg["exec"] = "/bin/false"
	cfg["args"] = "5 2 1"

	checker, err = NewProcessesChecker(cfg, nil)
	assert.Nilf(t, err, "processes checker instantiation failed: %q", err)

	ctxRun, cancelFunc1 = context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc1()
	result = checker.Run(ctxRun)
	assert.Equal(t, plugins.STATE_CRITICAL, result.Status)
	assert.Equal(t, "Process /bin/false 5 2 1 is not running", result.Message, "processes bad message: %q", result.Message)

	// missing all args
	cfg["args"] = ""
	cfg["exec"] = "/bin/sleep"

	checker, err = NewProcessesChecker(cfg, nil)
	assert.Nilf(t, err, "processes checker instantiation failed: %q", err)

	ctxRun, cancelFunc1 = context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancelFunc1()
	result = checker.Run(ctxRun)
	assert.Equal(t, plugins.STATE_CRITICAL, result.Status)
	assert.Equal(t, "Process /bin/sleep  is not running", result.Message, "processes bad message: %q", result.Message)

	// partial args
	cfg["args"] = "13"

	checker, err = NewProcessesChecker(cfg, nil)
	assert.Nilf(t, err, "processes checker instantiation failed: %q", err)

	ctxRun, cancelFunc1 = context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc1()
	result = checker.Run(ctxRun)
	assert.Equal(t, plugins.STATE_WARNING, result.Status)
	assert.Equal(t, "Process /bin/sleep 13 have too many instances running", result.Message)
}
