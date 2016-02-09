package main

import (
	"os"
	"regexp"

	"gopkg.in/fsnotify.v1"

	"path/filepath"
)

var watcherLog = NewLogFunc("watcher")

// WbsWatcher file wather struct
type WbsWatcher struct {
	w             *fsnotify.Watcher
	TargetDirs    []string
	ExcludeDirs   []string
	TargetFileExt []string
}

// Close close watcher
func (w *WbsWatcher) Close() {
	w.w.Close()
}

// initWatcher add watch target files to watcher
func (w *WbsWatcher) initWatcher() {
	var excludeDirRegexps []*regexp.Regexp
	for _, excludeDirStr := range w.ExcludeDirs {
		r := regexp.MustCompile(excludeDirStr)
		excludeDirRegexps = append(excludeDirRegexps, r)
	}

	for _, targetDir := range w.TargetDirs {
		filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
			for _, e := range excludeDirRegexps {
				if e.MatchString(path) {
					return nil
				}
			}
			for _, s := range w.TargetFileExt {
				if filepath.Ext(path) == s {
					watcherLog("start watching %s", path)
					err := w.w.Add(path)
					if err != nil {
						watcherLog("failed to watch file: %s: %s", path, err)
						return err
					}
				}
			}
			return nil
		})
	}
}

// NewWbsWatcher create target file watcher
func NewWbsWatcher(config *WbsConfig) (*WbsWatcher, error) {
	var watcher *WbsWatcher
	w, err := fsnotify.NewWatcher()
	if err != nil {
		watcherLog("failed to create watcher: %s", err)
		return watcher, err
	}
	watcher = &WbsWatcher{
		w:             w,
		TargetDirs:    config.WatchTargetDirs,
		ExcludeDirs:   config.WatchExcludeDirs,
		TargetFileExt: config.WatchFileExt,
	}
	watcher.initWatcher()
	return watcher, nil
}
