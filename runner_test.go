package main

import "testing"

func TestNewWbsRunner(t *testing.T) {
	config, err := NewWbsConfig("./wbs.example.toml")
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewWbsRunner(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(runner.StartCommand)
	t.Log(runner.StartOptions)
	if runner.Pid != -1 {
		t.Errorf("expected '-1' but got %d", runner.Pid)
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
