package main

import (
	"os"
	"testing"
)

func TestNewWbsRunner(t *testing.T) {
	config, err := NewWbsConfig("./wbs.example.toml")
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewWbsRunner(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(runner.BuildCommand)
	t.Log(runner.BuildOptions)
	t.Log(runner.StartCommand)
	if runner.Pid != -1 {
		t.Errorf("expected '-1' but got %d", runner.Pid)
	}
}

func cleanBuildTargetDir(path string) {
	os.RemoveAll(path)
}

func TestCreateBuildTargetDir(t *testing.T) {
	targetDir := "test_tmp_dir"
	err := createBuildTargetDir(targetDir)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanBuildTargetDir(targetDir)

	s, err := os.Stat(targetDir)
	if err != nil {
		t.Fatal(err)
	}
	if s.Name() != targetDir {
		t.Errorf("expected %s but got %s", targetDir, s.Name())
	}
}

func TestWbsRunnerBuild(t *testing.T) {
	config, err := NewWbsConfig("./wbs.example.toml")
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewWbsRunner(config)
	if err != nil {
		t.Fatal(err)
	}

	if err := runner.Build(); err != nil {
		t.Errorf("failed to execute: %s", err)
	}
}

func TestWbsRunnerServe(t *testing.T) {
	config, err := NewWbsConfig("./wbs.example.toml")
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewWbsRunner(config)
	if err != nil {
		t.Fatal(err)
	}
	runner.StartCommand = "nc"
	runner.StartOptions = []string{"-l", "8508"}
	err = runner.Serve()
	if err != nil {
		t.Fatal(err)
	}
	defer runner.Stop()
}
