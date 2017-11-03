package main

import "testing"

func TestNewWatcher(t *testing.T) {
	config, err := NewConfig("./wbs.example.toml")
	if err != nil {
		t.Fatal(err)
	}
	watcher, err := NewWatcher(config)
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
