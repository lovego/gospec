# gospec
a configurable golang coding spec checker.

## Installation
    go get github.com/bughou-go/spec/gospec

## Usage
    gospec [ <dir>/ | <dir> | <file> ] ...
- dir with trailing "/" means check all the packages in dir and it's subdir.
- dir without trailing "/" means check only th package in dir, not include it's subdir.
- file means check only the file.

exmaple:
```
ubuntu@ubuntu:~/go/src/github.com/bughou-go/my_project$ gospec services/
+---------------------------------------+--------------------------------------------------+------------+
|               position                |                     problem                      |    rule    |
+---------------------------------------+--------------------------------------------------+------------+
| services/event/on.go:76               | line 76 shouldn't be more than 100 chars         | sizes.line |
| services/event/on.go:49:1             | func listen shouldn't be more than 20 lines      | sizes.func |
| services/permission/perm_org.go:15:2  | var perm_org should be camelCase style           | names.var  |
| services/permission/perm_org.go:26:1  | func LcaPermTree shouldn't be more than 20 lines | sizes.func |
| services/permission/perm_org.go:34:2  | var perm_orgs should be camelCase style          | names.var  |
| services/permission/perm_org.go:37:9  | range var perm_org should be camelCase style     | names.var  |
| services/permission/perm_org.go:43:9  | range var perm_org should be camelCase style     | names.var  |
| services/permission/perm_org.go:53:3  | var current_type should be camelCase style       | names.var  |
| services/permission/perm_org.go:54:10 | range var perm_org should be camelCase style     | names.var  |
| services/permission/perm_org.go:64:20 | func param perm_orgs should be camelCase style   | names.var  |
| services/permission/perm_org.go:64:51 | func param perm_org should be camelCase style    | names.var  |
+---------------------------------------+--------------------------------------------------+------------+
```

## Checking Rules

##### size check
1. dir max files count check.
2. file max lines check.
3. line max length check.
4. function max lines check.


##### name check
1. dir name check.
2. file name check.
3. package name check.
4. type, func name check.
5. const, variable name check.

## Config File
gospec find the config file named "gospec.json" from current working directory upwards. It use the first one it find. If none is found, it uses the following default config:
```
{
  "names": {
    "dir":        { "style": "lowercase", "maxLen": 20 },
    "pkg":        { "style": "lowercase", "maxLen": 20 },
    "file":       { "style": "lowercase", "maxLen": 20 },

    "type":       { "style": "camelCase", "maxLen": 20 },
    "func":       { "style": "camelCase", "maxLen": 20 },

    "const":      { "style": "camelCase", "maxLen": 20 },
    "var":        { "style": "camelCase", "maxLen": 20 },
    "label":      { "style": "camelCase", "maxLen": 20 }
  },
  "sizes": {
    "dir": 20, "file": 200, "line": 100, "func": 20
  }
}
```
