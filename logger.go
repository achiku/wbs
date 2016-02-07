package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type logFunc func(string, ...interface{})

var logger = log.New(os.Stderr, "", 0)

func NewLogFunc(prefix string) func(string, ...interface{}) {
	prefix = fmt.Sprintf("%-11s", prefix)
	return func(format string, v ...interface{}) {
		now := time.Now()
		timeString := now.Format("15:04:05")
		format = fmt.Sprintf("%s %s | %s", timeString, prefix, format)
		logger.Printf(format, v...)
	}
}
