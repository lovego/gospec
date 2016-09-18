package main

import (
	"os"
	"path"
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
}

func processArgs() (traverseDirs, dirs, files []string) {
	for i, p := range os.Args {
		if i == 0 {
			continue
		}
		finfo, err := os.Stat(p)
		if err != nil {
			panic(err)
		}
		mode := finfo.Mode()
		switch {
		case mode.IsDir():
			if p[len(p)-1] == '/' {
				traverseDirs = append(traverseDirs, p)
			} else {
				dirs = append(dirs, p)
			}
		case mode.IsRegular():
			if path.Ext(p) == `.go` {
				files = append(files, p)
			}
		}
	}
	return
}
