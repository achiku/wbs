package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var runnerLog = NewLogFunc("runner")

type runnerLogWriter struct{}

func (a runnerLogWriter) Write(p []byte) (n int, err error) {
	runnerLog(string(p))
	return len(p), nil
}

// Runner runner struct
type Runner struct {
	Pid          int
	StartCommand string
	StartOptions []string
}

// Serve execute binary with configured options
func (r *Runner) Serve() error {
	evaledCommand, err := shellParser.Parse(r.StartCommand)
	if err != nil {
		return errors.Wrapf(err, "failed to parse command: %s", r.StartCommand)
	}
	evaledOptions, err := shellParser.Parse(strings.Join(r.StartOptions, " "))
	if err != nil {
		return errors.Wrapf(err, "failed to parse options: %s", r.StartOptions)
	}
	runnerLog(fmt.Sprintf("starting server: %s %s", evaledCommand, evaledOptions))
	cmd := exec.Command(evaledCommand[0], evaledOptions...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		runnerLog(err.Error())
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		runnerLog(err.Error())
		return err
	}
	rl := runnerLogWriter{}
	go io.Copy(rl, stderr)
	go io.Copy(rl, stdout)

	err = cmd.Start()
	if err != nil {
		runnerLog(err.Error())
		return err
	}
	r.Pid = cmd.Process.Pid
	runnerLog(fmt.Sprintf("server started: PID %d", r.Pid))
	return nil
}

// Stop stops running process
func (r *Runner) Stop() error {
	runnerLog(fmt.Sprintf("stopping server: PID %d", r.Pid))
	p, err := os.FindProcess(r.Pid)
	if err != nil {
		runnerLog(err.Error())
		return err
	}
	if err = p.Kill(); err != nil {
		runnerLog(err.Error())
		return err
	}
	runnerLog(fmt.Sprintf("server stopped: PID %d", r.Pid))
	return nil
}

// NewRunner create runner
func NewRunner(config *Config) (*Runner, error) {
	targetBinary := filepath.Join(config.BuildTargetDir, config.BuildTargetName)
	r := &Runner{
		Pid:          -1,
		StartCommand: targetBinary,
		StartOptions: config.StartOptions,
	}
	return r, nil
}
