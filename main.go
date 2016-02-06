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

	var (
		config *WbsConfig
		err    error
	)
	if configFile != nil {
		config, err = NewWbsConfig(*configFile)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
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
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("file modified:", event.Name)
					log.Println("restarting serevr")
					runner.Stop()
					log.Println("serevr stopped")
					runner.Build()
					log.Println("build completed")
					runner.Serve()
					log.Println("serevr restarted")
				}
			case err := <-watcher.w.Errors:
				log.Println("error:", err)
			}
		}
	}()
	<-done
}
