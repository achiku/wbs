package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var runnerLog = NewLogFunc("runner")
var builderLog = NewLogFunc("builder")

type runnerLogWriter struct{}

func (a runnerLogWriter) Write(p []byte) (n int, err error) {
	runnerLog(string(p))
	return len(p), nil
}

// WbsRunner runner struct
type WbsRunner struct {
	Pid            int
	PidFile        string
	BuildTargetDir string
	BuildCommand   string
	BuildOptions   []string
	StartCommand   string
	StartOptions   []string
}

// createBuildTargetDir create dir for build binary
func createBuildTargetDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			builderLog(err.Error())
		}
	}
	return nil
}

// Build execute build command with configured options
func (r *WbsRunner) Build() error {
	builderLog("starting build: %s %s", r.BuildCommand, r.BuildOptions)
	createBuildTargetDir(r.BuildTargetDir)
	cmd := exec.Command(r.BuildCommand, r.BuildOptions...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		builderLog(err.Error())
	}
	builderLog("\n" + string(out))
	return nil
}

// Server execute binary with configured options
func (r *WbsRunner) Serve() error {
	runnerLog("starting server: %s %s", r.StartCommand, r.StartOptions)
	cmd := exec.Command(r.StartCommand, r.StartOptions...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		runnerLog(err.Error())
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		runnerLog(err.Error())
	}
	rl := runnerLogWriter{}
	go io.Copy(rl, stderr)
	go io.Copy(rl, stdout)

	err = cmd.Start()
	if err != nil {
		runnerLog(err.Error())
	}
	r.Pid = cmd.Process.Pid
	runnerLog("starting server: PID %d", r.Pid)
	return nil
}

func (r *WbsRunner) Stop() error {
	runnerLog("stopping server: PID %d", r.Pid)
	p, err := os.FindProcess(r.Pid)
	if err != nil {
		runnerLog(err.Error())
	}
	if err = p.Kill(); err != nil {
		runnerLog(err.Error())
	}
	return nil
}

// NewWbsRunner create runner
func NewWbsRunner(config *WbsConfig) (*WbsRunner, error) {
	targetBinary := filepath.Join(config.BuildTargetDir, config.BuildTargetName)
	buildOptions := append(config.BuildOptions, "-o", targetBinary)
	r := &WbsRunner{
		Pid:            -1,
		PidFile:        config.PidFile,
		BuildTargetDir: config.BuildTargetDir,
		BuildCommand:   config.BuildCommand,
		BuildOptions:   buildOptions,
		StartCommand:   targetBinary,
		StartOptions:   config.StartOptions,
	}
	return r, nil
}
