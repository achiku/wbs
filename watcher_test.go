package main

import "testing"

func TestNewWbsWatcher(t *testing.T) {
	config, err := NewWbsConfig("./wbs.example.toml")
	if err != nil {
		t.Fatal(err)
	}
	watcher, err := NewWbsWatcher(config)
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()

	t.Log(watcher.ExcludeDirs)
	t.Log(watcher.TargetDirs)
	t.Log(watcher.TargetFileExt)
	t.Log(watcher.ExcludeFilePatterns)

	if len(watcher.TargetFileExt) != 3 {
		t.Errorf("expected 3 but got %d", len(watcher.TargetFileExt))
	}
}
