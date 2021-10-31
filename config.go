package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Config wbs configuration struct
type Config struct {
	RootPath                 string   `toml:"root_path"`
	RestartProcess           bool     `toml:"restart_process"`
	BuildTargetDir           string   `toml:"build_target_dir"`
	BuildTargetName          string   `toml:"build_target_name"`
	BuildCommand             string   `toml:"build_command"`
	BuildOptions             []string `toml:"build_options"`
	StartOptions             []string `toml:"start_options"`
	WatchTargetDirs          []string `toml:"watch_target_dirs"`
	WatchExcludeDirs         []string `toml:"watch_exclude_dirs"`
	WatchFileExt             []string `toml:"watch_file_ext"`
	WatchFileExcludePatterns []string `toml:"watch_file_exclude_pattern"`
	WatchDirOnly             bool     `toml:"watch_dir_only"`
}

// NewConfig create wbs config struct
func NewConfig(configFilePath string) (*Config, error) {
	var config Config
	// set default value
	config.RestartProcess = true
	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		log.Fatalf("failed to create Config from file: %s", err)
		return nil, err
	}
	return &config, nil
}

// NewDefaultConfig create wbs default config
func NewDefaultConfig() *Config {
	config := Config{
		RootPath:                 ".",
		RestartProcess:           true,
		BuildTargetDir:           "tmp",
		BuildTargetName:          "server",
		BuildCommand:             "go",
		BuildOptions:             []string{"build", "-v"},
		StartOptions:             []string{},
		WatchFileExt:             []string{".go", ".tmpl", ".html"},
		WatchFileExcludePatterns: []string{"*_gen.go"},
		WatchTargetDirs:          []string{"."},
		WatchExcludeDirs:         []string{".git", "tmp", "bin"},
	}
	return &config
}
