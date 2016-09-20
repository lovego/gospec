# gospec
a configurable golang coding spec checker.

## Installation
    go get github.com/bughou-go/spec/gospec

## Usage
    gospec [ <dir>/ | <dir> | <file> ] ...
- dir with trailing "/" means check all the packages in dir and it's subdir.
- dir without trailing "/" means check only th package in dir, not include it's subdir.
- file means check only the file.

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
