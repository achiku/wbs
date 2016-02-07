# wbs [![Build Status](https://travis-ci.org/achiku/wbs.svg?branch=master)](https://travis-ci.org/achiku/wbs)

Watch, build, and (re)start Go net/http server, customizeable by toml configuration file


## Why created

This is certainly yet another auto-build-reload tool for go net/http. I had been using [fresh](https://github.com/pilu/fresh), but development of this tool got really inactive from middle of the last year. So, I decided to rewrite from the ground up, adding more flexibility by using toml configuration file, making it possible to use `gb` or `gom` to biuld binary, adding bunch of tests, and keeping code structure simple so that it will be easy to contribute.

This tool was built by just spending 5 hours of my weekend, and not really sophisticated at this stage, so pull-requests and issue reports are all very welcomed :)


## Installation

```
go get -u github.com/achiku/wbs
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


## Configuration

```toml
root_path = "."

watch_target_dirs = ["."]
watch_exclude_dirs = [".git", "vendor", "tmp"]
watch_file_ext = [".go", ".tmpl", ".html"]

build_target_dir = "tmp"
build_target_name = "myserver"
build_command = "go"
build_options = ["build", "-v"]

start_options = ["-v"]
```
