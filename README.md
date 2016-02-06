# wbs [![Build Status](https://travis-ci.org/achiku/wbs.svg?branch=master)](https://travis-ci.org/achiku/wbs)

Watch, build, and (re)start Go net/http server, customizeable by toml configuration file


## Why?

Because fresh is not fresh anymore.


## Installation

```
go get -u github.com/achiku/wbs
```

## Start

```
wbs
```

## Configuration

```
root_path = "."

watch_target_dirs = ["."]
watch_exclude_dirs = [".git", "vendor", "tmp"]
watch_file_ext = [".go", ".tmpl", ".html"]

build_target_dir = "tmp"
build_target_name = "myserver"
build_command = "go"
build_options = ["build", "-v"]

start_options = ["-v"]
pid_file = "./tmp/pid"
```
