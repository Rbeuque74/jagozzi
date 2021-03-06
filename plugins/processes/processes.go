package processes

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	processlib "github.com/mitchellh/go-ps"
	"github.com/rbeuque74/jagozzi/plugins"
	log "github.com/sirupsen/logrus"
)

const pluginName = "Processes"

var backslashZero = []byte("\x00")

func init() {
	plugins.Register(pluginName, NewProcessesChecker)
}

// ProcessesChecker is a plugin to check status code of command
type ProcessesChecker struct {
	cfg            processesConfig
	executableName string
}

// Name returns the name of the checker
func (c ProcessesChecker) Name() string {
	return pluginName
}

// ServiceName returns the name of the NSCA service associated to the checker
func (c ProcessesChecker) ServiceName() string {
	return c.cfg.Name
}

// Periodicity returns the delay between two checks
func (c ProcessesChecker) Periodicity() *time.Duration {
	return c.cfg.Periodicity()
}

// Run is performing the checker protocol
func (c *ProcessesChecker) Run(ctx context.Context) plugins.Result {
	processes, err := processlib.Processes()
	if err != nil {
		return plugins.ResultFromError(c, err, "unable to retrieve processes")
	}

	var candidatesProcesses []processlib.Process
	for _, proc := range processes {
		if proc.Executable() != c.executableName {
			continue
		}
		candidatesProcesses = append(candidatesProcesses, proc)
	}

	var selectedProcesses []processlib.Process
	for _, proc := range candidatesProcesses {
		path, err := filepath.EvalSymlinks(fmt.Sprintf("/proc/%d/exe", proc.Pid()))
		if err != nil {
			return plugins.ResultFromError(c, err, "can't open executable symlink from pid")
		}

		if path != c.cfg.Command {
			log.Debugf("processes: pid [%d] %q doesn't match command line", proc.Pid(), path)
			continue
		}

		b, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", proc.Pid()))
		if err != nil {
			return plugins.ResultFromError(c, err, "can't open cmdline")
		}

		// Removing last NUL characters
		b = bytes.TrimSuffix(b, backslashZero)
		// Splitting all parts
		args := bytes.Split(b, backslashZero)
		// Removing first segment as it's the name of launched executable
		if len(args) != 0 {
			args = args[1:]
		}
		// Putting back in one piece
		b = bytes.Join(args, []byte(" "))

		cmdargs := string(b[:])

		if cmdargs != c.cfg.Args {
			log.Debugf("processes: pid [%d] %q %q doesn't match cmdargs", proc.Pid(), path, cmdargs)
			continue
		}

		selectedProcesses = append(selectedProcesses, proc)
	}

	if len(selectedProcesses) == 0 {
		return plugins.Result{
			Status:  plugins.STATE_CRITICAL,
			Message: fmt.Sprintf("Process %s %s is not running", c.cfg.Command, c.cfg.Args),
			Checker: c,
		}
	} else if len(selectedProcesses) > 1 {
		return plugins.Result{
			Status:  plugins.STATE_WARNING,
			Message: fmt.Sprintf("Process %s %s have too many instances running", c.cfg.Command, c.cfg.Args),
			Checker: c,
		}
	} else {
		return plugins.Result{
			Status:  plugins.STATE_OK,
			Message: fmt.Sprintf("Process %s %s is running", c.cfg.Command, c.cfg.Args),
			Checker: c,
		}
	}
}

// NewProcessesChecker create a Processes checker
func NewProcessesChecker(conf interface{}, pluginConf interface{}) (plugins.Checker, error) {
	out, err := yaml.Marshal(conf)
	if err != nil {
		return nil, err
	}

	cfg := processesConfig{}
	err = yaml.Unmarshal(out, &cfg)
	if err != nil {
		return nil, err
	}

	checker := &ProcessesChecker{
		cfg: cfg,
	}

	checker.executableName = path.Base(checker.cfg.Command)

	log.Infof("processes: Checker activated for watching %q", checker.cfg.Command)
	return checker, nil
}
