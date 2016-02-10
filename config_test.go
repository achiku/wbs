package main

import "testing"

func TestNewConfig(t *testing.T) {
	config, err := NewWbsConfig("./wbs.example.toml")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("restart_process: %t", config.RestartProcess)
	t.Logf("watch_file_exclude_pattern: %s", config.WatchFileExcludePatterns)

	if config.BuildTargetName != "myserver" {
		t.Errorf("expected myserver but got %s", config.BuildTargetName)
	}
}

func TestNewDefaultConfig(t *testing.T) {
	config := NewWbsDefaultConfig()
	if config.BuildTargetName != "server" {
		t.Errorf("expected server but got %s", config.BuildTargetName)
	}
}
