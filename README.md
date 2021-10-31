# wbs

[![test](https://github.com/achiku/wbs/actions/workflows/test.yml/badge.svg)](https://github.com/achiku/wbs/actions/workflows/test.yml)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/achiku/wbs/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/achiku/wbs)](https://goreportcard.com/report/github.com/achiku/wbs)

Watch, build, and (re)start Go net/http server, customizeable by toml configuration file


## Why created

This tool is yet another auto-rebuild-reloader for go net/http. This tool uses flexible toml configuration file, makes it possible to use `gb`, `gom`, `make`, or any other command with environmental variables to build binary. It also have bunch of tests, and keeping code structure simple to make it easy to contribute.


## Installation

```
go get -u github.com/achiku/wbs
```


## Configuration Example

```toml
root_path = "."

watch_target_dirs = ["."]
watch_exclude_dirs = [".git", "vendor", "bin"]
watch_file_ext = [".go", ".tmpl", ".html"]
# usually used when we need to ignore files generated by `go generate`
# need to include dir path
watch_file_exclude_pattern = [
    "example/lib/*_gen.go",
    "example/lib/bindata.go"
]
# this option forces wbs only watch directory, not files.
# if you are using vim/neovim set this option `true`
# vim/neovim do not fire events that fsnotify can capture for file.
watch_dir_only = false

# Env vars can be used in build_target_dir, build_target_name,  build_command,
# build_options, start_options
build_target_dir = "$GOPATH/bin"
build_target_name = "myserver"
build_command = "go"
build_options = ["build", "-v"]

# start command will be `build_target_dir/build_target_name start_options`
# in this case $GOPATH/bin/myserver -v
start_options = ["-v", "-p", "$APP_PORT"]

# default true, but it's possible to make this fale, when
# running a program that doesn't persist as a process
restart_process = true
```


## Quick Start

Save the following net/http application to `main.go`.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!\n")
}

func main() {
	helloWorldHandler := http.HandlerFunc(helloWorld)
	http.Handle("/hello", loggingMiddleware(helloWorldHandler))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
```

```
≫≫ ls -l
total 8
-rw-r--r--  1 achiku  staff  636  2  7 14:24 main.go
```

Then, just type `wbs`. This will build binary, and start watching files specified.

```
≫≫ wbs
14:50:58 watcher     |start watching main.go
14:50:58 builder     |starting build: go [build -v -o tmp/server]
14:50:59 builder     |
github.com/achiku/wbs/example
14:50:59 runner      |starting server: tmp/server [-v]
14:50:59 runner      |starting server: PID 41797
```

You can specify configuration file using `-c` option.

```
≫≫ wbs -c wbs.example.toml
14:51:58 watcher     |start watching main.go
14:51:58 builder     |starting build: go [build -v -o tmp/server]
14:51:59 builder     |
github.com/achiku/wbs/example
14:51:59 runner      |starting server: tmp/server [-v]
14:51:59 runner      |starting server: PID 41797
```

Application stdout and stderr goes to wbs stderr.

```
≫≫ for i in $(seq 1 3); do curl http://localhost:8080/hello; done
Hello, world!
Hello, world!
Hello, world!
```

```
14:55:23 runner      |2016/02/07 14:55:23 [GET] "/hello" 5.74µs
14:55:25 runner      |2016/02/07 14:55:25 [GET] "/hello" 5.692µs
14:55:26 runner      |2016/02/07 14:55:26 [GET] "/hello" 7.323µs
```

When `main.go` is modified, wbs will rebuild binary, and restart server.

```
14:59:54 main        |file modified: "main.go": WRITE
14:59:54 runner      |stopping server: PID 44000
14:59:54 builder     |starting build: go [build -v -o tmp/server]
14:59:56 builder     |
github.com/achiku/wbs/example
14:59:56 runner      |starting server: tmp/server [-v]
14:59:56 runner      |starting server: PID 44036
```

### Test

```
$ go test -v
```

### Inspired by

- https://github.com/pilu/fresh
