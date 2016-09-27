package main

import (
	"flag"
	"os"
	"path"

	"github.com/bughou-go/spec/problems"
)

func main() {
	traDirs, dirs, files := processArgs()
	for _, d := range traDirs {
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

func processArgs() (traDirs, dirs, files []string) {
	for _, p := range flag.Args() {
		switch mode := fileMode(p); {
		case mode.IsDir():
			if p[len(p)-1] == '/' {
				traDirs = append(traDirs, path.Clean(p))
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
	return path.Ext(p) == `.go` && path.Base(p)[0] != '.' && p[0] != '_'
}

func fileMode(p string) os.FileMode {
	if fi, err := os.Stat(p); err == nil {
		return fi.Mode()
	} else {
		panic(err)
	}
}
