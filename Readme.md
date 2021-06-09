# gospec
a configurable golang coding specification checker.

[![Build Status](https://github.com/lovego/gospec/actions/workflows/go.yml/badge.svg)](https://github.com/lovego/gospec/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/lovego/gospec/badge.svg?branch=master)](https://coveralls.io/github/lovego/gospec)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/gospec)](https://goreportcard.com/report/github.com/lovego/gospec)
[![Documentation](https://pkg.go.dev/badge/github.com/lovego/gospec)](https://pkg.go.dev/github.com/lovego/gospec@v1.0.0)

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

