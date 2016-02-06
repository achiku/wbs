package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

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
			log.Fatal(err)
			return err
		}
	}
	return nil
}

// Build execute build command with configured options
func (r *WbsRunner) Build() error {
	log.Printf("starting build: %s %s", r.BuildCommand, r.BuildOptions)
	createBuildTargetDir(r.BuildTargetDir)
	cmd := exec.Command(r.BuildCommand, r.BuildOptions...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("\n" + string(out))
	return nil
}

// Server execute binary with configured options
func (r *WbsRunner) Serve() error {
	log.Printf("starting server: %s %s", r.StartCommand, r.StartOptions)
	cmd := exec.Command(r.StartCommand, r.StartOptions...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
		return err
	}
	r.Pid = cmd.Process.Pid
	log.Printf("starting server: PID %d", r.Pid)
	return nil
}

func (r *WbsRunner) Stop() error {
	log.Printf("stopping server: PID %d", r.Pid)
	p, err := os.FindProcess(r.Pid)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if err = p.Kill(); err != nil {
		log.Fatal(err)
		return err
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
