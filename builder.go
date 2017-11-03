package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// Builder builder struct
type Builder struct {
	BuildTargetDir string
	BuildCommand   string
	BuildOptions   []string
}

// Build execute build command with configured options
func (b *Builder) Build() error {
	evaledCommand, err := shellParser.Parse(b.BuildCommand)
	evaledTargetDir, err := shellParser.Parse(b.BuildTargetDir)
	evaledOptions, err := shellParser.Parse(strings.Join(b.BuildOptions, " "))
	if err != nil {
		return err
	}

	builderLog(fmt.Sprintf("starting build: %s %s", evaledCommand, evaledOptions))
	createBuildTargetDir(evaledTargetDir[0])
	cmd := exec.Command(evaledCommand[0], evaledOptions...)
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

// NewBuilder create runner
func NewBuilder(config *Config) (*Builder, error) {
	targetBinary := filepath.Join(config.BuildTargetDir, config.BuildTargetName)
	buildOptions := append(config.BuildOptions, "-o", targetBinary)
	b := &Builder{
		BuildTargetDir: config.BuildTargetDir,
		BuildCommand:   config.BuildCommand,
		BuildOptions:   buildOptions,
	}
	return b, nil
}
