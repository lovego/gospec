package main

import (
	"flag"
	"os"
	"path"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules"
)

func main() {
	traDirs, dirs, files := processArgs()
	for _, d := range traDirs {
		traverseDir(d)
	}
	for _, d := range dirs {
		checkDir(d)
	}
	if len(files) > 0 {
		checkFiles(files)
	}

	if problems.Count() > 0 {
		problems.Render()
		os.Exit(1)
	}
}

func traverseDir(p string) {
	f, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	list, err := f.Readdir(-1)
	if err != nil {
		panic(err)
	}
	for _, d := range list {
		if d.IsDir() && d.Name()[0] != '.' {
			traverseDir(path.Join(p, d.Name()))
		}
	}
	checkDir(p)
}

func checkDir(dir string) {
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	names, err := f.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	files := make([]string, 0, len(names))
	for _, name := range names {
		if willBuild(name) {
			files = append(files, path.Join(dir, name))
		}
	}
	if len(files) > 0 {
		rules.Check(dir, files)
	}
}

func checkFiles(paths []string) {
	dirs := make(map[string][]string)
	for _, p := range paths {
		dir := path.Dir(p)
		dirs[dir] = append(dirs[dir], p)
	}
	for dir, files := range dirs {
		rules.Check(dir, files)
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
