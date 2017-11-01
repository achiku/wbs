package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type logFunc func(string, ...interface{})

var logger = log.New(os.Stderr, "", 0)

// NewLogFunc create log func
func NewLogFunc(prefix string) func(string) {
	prefix = fmt.Sprintf("%-11s", prefix)
	return func(message string) {
		now := time.Now()
		timeString := now.Format("15:04:05")
		message = fmt.Sprintf("%s %s | ", timeString, prefix) + message
		logger.Print(message)
	}
}
