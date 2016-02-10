package main

import (
	"log"
	"os"
	"regexp"
	"sort"

	"gopkg.in/fsnotify.v1"

	"path/filepath"
)

var watcherLog = NewLogFunc("watcher")

// WbsWatcher file wather struct
type WbsWatcher struct {
	w                   *fsnotify.Watcher
	TargetDirs          []string
	ExcludeDirs         []string
	TargetFileExt       []string
	ExcludeFilePatterns []string
}

// Close close watcher
func (w *WbsWatcher) Close() {
	w.w.Close()
}

func contains(v string, l []string) bool {
	sort.Strings(l)
	i := sort.SearchStrings(l, v)
	if i < len(l) && l[i] == v {
		return true
	}
	return false
}

func matchContains(v string, excludeStrs []string) bool {
	for _, es := range excludeStrs {
		match, err := filepath.Match(es, v)
		if err != nil {
			log.Fatal(err)
		}
		if match {
			return true
		}
	}
	return false
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
			fileExt := filepath.Ext(path)
			if contains(fileExt, w.TargetFileExt) && !matchContains(path, w.ExcludeFilePatterns) {
				watcherLog("start watching %s", path)
				err := w.w.Add(path)
				if err != nil {
					watcherLog("failed to watch file: %s: %s", path, err)
					return err
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
		w:                   w,
		TargetDirs:          config.WatchTargetDirs,
		ExcludeDirs:         config.WatchExcludeDirs,
		TargetFileExt:       config.WatchFileExt,
		ExcludeFilePatterns: config.WatchFileExcludePatterns,
	}
	watcher.initWatcher()
	return watcher, nil
}
