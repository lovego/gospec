# gospec
a configurable golang coding specification checker.

## Installation
    go get github.com/lovego/gospec

## Usage
    gospec [ <dir>/... | <dir> | <file> ] ...
- dir/... means check the dir and the ".go" files in dir recursively.
- dir     means check the dir and the ".go" files in dir.
- file    means check only the file.

exmaple:
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

## Config File
gospec find the config file named ".gospec.yml" from current working directory upwards.
It use the first one it find. If none is found, it uses the <a href="gospec.yml">default config</a>.

