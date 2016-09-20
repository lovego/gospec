package main

import (
	"os"
	"path"

	"github.com/bughou-go/spec/problems"
)

func main() {
	dirsR, dirs, files := processArgs()
	for _, d := range dirsR {
		traverseDir(d)
	}
	for _, d := range dirs {
		doDir(d)
	}
	if len(files) > 0 {
		doFiles(files)
	}

	if problems.Count() > 0 {
		problems.Render()
		os.Exit(1)
	}
}

func processArgs() (traverseDirs, dirs, files []string) {
	for i, p := range os.Args {
		if i == 0 {
			continue
		}
		switch mode := fileMode(p); {
		case mode.IsDir():
			if p[len(p)-1] == '/' {
				traverseDirs = append(traverseDirs, path.Clean(p))
			} else {
				dirs = append(dirs, path.Clean(p))
			}
		case mode.IsRegular():
			if willBuild(p) {
				files = append(files, path.Clean(p))
			}
		}
	}
	return
}

func willBuild(p string) bool {
	return path.Ext(p) == `.go` && p[0] != '.' && p[0] != '_'
}

func fileMode(p string) os.FileMode {
	if fi, err := os.Stat(p); err == nil {
		return fi.Mode()
	} else {
		panic(err)
	}
}
