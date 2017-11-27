# gospec
a configurable golang coding spec checker.

## Installation
    go get github.com/lovego/spec/gospec

## Usage
    gospec [ <dir>/ | <dir> | <file> ] ...
- dir with trailing "/" means check all the packages in dir and it&apos;s subdir.
- dir without trailing "/" means check only the package in dir, not include it&apos;s subdir.
- file means check only the file.

exmaple:
```
ubuntu@ubuntu:~/go/src/github.com/lovego/my_project$ gospec services/
+---------------------------------------+--------------------------------------------------+------------+
|               position                |                     problem                      |    rule    |
+---------------------------------------+--------------------------------------------------+------------+
| services/event/on.go:76               | line 76 shouldn not be more than 100 chars       | sizes.line |
| services/event/on.go:49:1             | func listen should not be more than 20 lines     | sizes.func |
| services/permission/perm_org.go:15:2  | var perm_org should be camelCase style           | names.var  |
| services/permission/perm_org.go:26:1  | func LcaPermTree should not be more than 20 lines| sizes.func |
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
2. line max width(chars) check.
3. file max lines count check (ignore comment lines).
4. function max statements count check.


##### name check
1. dir, file, package name check.
2. func, type, const, variable name check.
3. local type, local const, local variable, label name check.

## Config File
gospec find the config file named "gospec.json" from current working directory upwards.
It use the first one it find. If none is found, it uses the <a href="gospec.json">default config</a>.

