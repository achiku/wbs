package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var runnerLog = NewLogFunc("runner")

type runnerLogWriter struct{}

func (a runnerLogWriter) Write(p []byte) (n int, err error) {
	runnerLog(string(p))
	return len(p), nil
}

// WbsRunner runner struct
type WbsRunner struct {
	Pid          int
	StartCommand string
	StartOptions []string
}

// Serve execute binary with configured options
func (r *WbsRunner) Serve() error {
	evaledCommand, err := shellParser.Parse(r.StartCommand)
	evaledOptions, err := shellParser.Parse(strings.Join(r.StartOptions, " "))
	if err != nil {
		return err
	}
	runnerLog("starting server: %s %s", evaledCommand, evaledOptions)
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
	runnerLog("server started: PID %d", r.Pid)
	return nil
}

// Stop stops running process
func (r *WbsRunner) Stop() error {
	runnerLog("stopping server: PID %d", r.Pid)
	p, err := os.FindProcess(r.Pid)
	if err != nil {
		runnerLog(err.Error())
		return err
	}
	if err = p.Kill(); err != nil {
		runnerLog(err.Error())
		return err
	}
	runnerLog("server stopped: PID %d", r.Pid)
	return nil
}

// NewWbsRunner create runner
func NewWbsRunner(config *WbsConfig) (*WbsRunner, error) {
	targetBinary := filepath.Join(config.BuildTargetDir, config.BuildTargetName)
	r := &WbsRunner{
		Pid:          -1,
		StartCommand: targetBinary,
		StartOptions: config.StartOptions,
	}
	return r, nil
}
