package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

var builderLog = NewLogFunc("builder")

// createBuildTargetDir create dir for build binary
func createBuildTargetDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			builderLog(err.Error())
		}
	}
	return nil
}

type WbsBuilder struct {
	BuildTargetDir string
	BuildCommand   string
	BuildOptions   []string
}

// Build execute build command with configured options
func (b *WbsBuilder) Build() error {
	builderLog("starting build: %s %s", b.BuildCommand, b.BuildOptions)
	createBuildTargetDir(b.BuildTargetDir)
	cmd := exec.Command(b.BuildCommand, b.BuildOptions...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		builderLog(err.Error())
		builderLog("\n" + string(out))
		return err
	}
	builderLog("\n" + string(out))
	builderLog("build completed")
	return nil
}

// NewWbsRunner create runner
func NewWbsBuilder(config *WbsConfig) (*WbsBuilder, error) {
	targetBinary := filepath.Join(config.BuildTargetDir, config.BuildTargetName)
	buildOptions := append(config.BuildOptions, "-o", targetBinary)
	b := &WbsBuilder{
		BuildTargetDir: config.BuildTargetDir,
		BuildCommand:   config.BuildCommand,
		BuildOptions:   buildOptions,
	}
	return b, nil
}
