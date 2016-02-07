package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/fsnotify.v1"
)

func main() {
	configFile := flag.String("c", "", "configuration file path")
	flag.Parse()

	mainLogger := NewLogFunc("main")
	var (
		config *WbsConfig
		err    error
	)
	if *configFile != "" {
		config, err = NewWbsConfig(*configFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		config = NewWbsDefaultConfig()
	}

	watcher, err := NewWbsWatcher(config)
	if err != nil {
		mainLogger("failed to initialize watcher: %s", err)
		os.Exit(1)
	}
	runner, err := NewWbsRunner(config)
	if err != nil {
		mainLogger("failed to initialize runner: %s", err)
		os.Exit(1)
	}
	builder, err := NewWbsBuilder(config)
	if err != nil {
		mainLogger("failed to initialize builder: %s", err)
		os.Exit(1)
	}

	if err := builder.Build(); err != nil {
		mainLogger("failed to build: %s", err)
		os.Exit(1)
	}
	err = runner.Serve()
	if err != nil {
		mainLogger("failed to start server: %s", err)
		os.Exit(1)
	}
	defer runner.Stop()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.w.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					e := event.String()
					mainLogger("file modified: %s", e)
					runner.Stop()
					builder.Build()
					runner.Serve()
				}
			case err := <-watcher.w.Errors:
				mainLogger("error: %s", err)
			}
		}
	}()
	<-done
}
