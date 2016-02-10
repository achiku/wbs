package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/naoina/toml"
)

// WbsConfig wbs configuration struct
type WbsConfig struct {
	RootPath                 string   `toml:"root_path"`
	RestartProcess           bool     `toml:"restart_process"`
	BuildTargetDir           string   `toml:"build_target_dir"`
	BuildTargetName          string   `toml:"build_target_name"`
	BuildCommand             string   `toml:"build_command"`
	BuildOptions             []string `toml:"build_options"`
	StartOptions             []string `toml:"start_options"`
	WatchTargetDirs          []string `toml:"watch_target_dir"`
	WatchExcludeDirs         []string `toml:"watch_exclude_dir"`
	WatchFileExt             []string `toml:"watch_file_ext"`
	WatchFileExcludePatterns []string `toml:"watch_file_exclude_pattern"`
}

// NewWbsConfig create wbs config struct
func NewWbsConfig(configFilePath string) (*WbsConfig, error) {
	f, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("failed to open config file: %s: %s", configFilePath, err)
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("failed to read config file: %s", err)
		return nil, err
	}
	var config WbsConfig
	// set default value
	config.RestartProcess = true
	if err := toml.Unmarshal(buf, &config); err != nil {
		log.Fatalf("failed to create Config from file: %s", err)
		return nil, err
	}
	return &config, nil
}

// NewWbsDefaultConfig create wbs default config
func NewWbsDefaultConfig() *WbsConfig {
	config := WbsConfig{
		RootPath:                 ".",
		RestartProcess:           true,
		BuildTargetDir:           "tmp",
		BuildTargetName:          "server",
		BuildCommand:             "go",
		BuildOptions:             []string{"build", "-v"},
		StartOptions:             []string{"-v"},
		WatchFileExt:             []string{".go", ".tmpl", ".html"},
		WatchFileExcludePatterns: []string{"*_gen.go"},
		WatchTargetDirs:          []string{"."},
		WatchExcludeDirs:         []string{".git", "tmp", "bin"},
	}
	return &config
}
