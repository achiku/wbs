package main

import (
	"os"
	"testing"
)

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

func TestBuilderBuild(t *testing.T) {
	targetDir := "test_tmp_dir"
	err := createBuildTargetDir(targetDir)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanBuildTargetDir(targetDir)

	os.Setenv("WBS_TEST_SERVER_ROOT", "test_tmp_dir")
	config, err := NewConfig("./wbs.example.toml")
	if err != nil {
		t.Fatal(err)
	}
	builder, err := NewBuilder(config)
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build(); err != nil {
		t.Errorf("failed to execute: %s", err)
	}
}
