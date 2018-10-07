# gospec
a configurable golang coding specification checker.

[![Build Status](https://travis-ci.org/lovego/gospec.svg?branch=master)](https://travis-ci.org/lovego/gospec)
[![Coverage Status](https://coveralls.io/repos/github/lovego/gospec/badge.svg?branch=master)](https://coveralls.io/github/lovego/gospec?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/gospec)](https://goreportcard.com/report/github.com/lovego/gospec)
[![GoDoc](https://godoc.org/github.com/lovego/gospec?status.svg)](https://godoc.org/github.com/lovego/gospec)

## Installation
    go get github.com/lovego/gospec

## Usage
    gospec [ <dir>/... | <dir> | <file> ] ...
- dir/... means check the dir and the ".go" files in dir recursively.
- dir     means check the dir and the ".go" files in dir.
- file    means check only the file.

### exmaple:
```
MacBook:~/go/src/my_project$ gospec models/...
+---------------------------------------+---------------------------------------------------+----------------------+
|               position                |                         problem                   |         rule         |
+---------------------------------------+---------------------------------------------------+----------------------+
| models/users/list.go:36:10            | func List params size: 7, limit: 5                | func.size.maxParams  |
| models/users/create_or_update.go:28:3 | func CreateOrUpdate results size: 4, limit: 3     | func.size.maxResults |
| models/users/delete.go:111:2          | local var name SQL should be lowerCamelCase style | localVar.style       |
+---------------------------------------+---------------------------------------------------+----------------------+
```

## Configuration
gospec find the config file named ".gospec.yml" from current working directory upwards. It use the first one it find.

If no one is found, it uses the <a href=".gospec.yml">default configuration</a>.

