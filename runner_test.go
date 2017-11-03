package main

import (
	"os"
	"testing"
)

func TestNewRunner(t *testing.T) {
	config, err := NewConfig("./wbs.example.toml")
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewRunner(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(runner.StartCommand)
	t.Log(runner.StartOptions)
	if runner.Pid != -1 {
		t.Errorf("expected '-1' but got %d", runner.Pid)
	}
}

func TestRunnerServe(t *testing.T) {
	config, err := NewConfig("./wbs.example.toml")
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewRunner(config)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("TEST_RUNNER_CMD", "nc")
	os.Setenv("TEST_RUNNER_PORT", "8508")
	runner.StartCommand = "$TEST_RUNNER_CMD"
	runner.StartOptions = []string{"-l", "$TEST_RUNNER_PORT"}
	err = runner.Serve()
	if err != nil {
		t.Fatal(err)
	}
	defer runner.Stop()
}
