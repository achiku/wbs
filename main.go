package main

import (
	"flag"
	"log"

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
		log.Fatal(err)
	}
	runner, err := NewWbsRunner(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := runner.Build(); err != nil {
		log.Fatal(err)
	}
	err = runner.Serve()
	if err != nil {
		log.Fatal(err)
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
					runner.Build()
					runner.Serve()
				}
			case err := <-watcher.w.Errors:
				mainLogger("error:", err)
			}
		}
	}()
	<-done
}
