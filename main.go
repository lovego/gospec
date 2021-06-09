package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules"
)

func main() {
	processArgs()
	traDirs, dirs, files := getTargets()
	rules.LoadConfig()
	for _, dir := range traDirs {
		traverseDir(dir)
	}
	for _, dir := range dirs {
		checkDir(dir)
	}
	if len(files) > 0 {
		checkFiles(files)
	}

	if problems.Count() > 0 {
		problems.Render()
		os.Exit(1)
	}
}

func traverseDir(dir string) {
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	list, err := f.Readdir(-1)
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	for _, d := range list {
		if d.IsDir() && (d.Name()[0] != '.' && d.Name()[0] != '_') {
			traverseDir(filepath.Join(dir, d.Name()))
		}
	}
	checkDir(dir)
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
	if err := f.Close(); err != nil {
		panic(err)
	}
	files := make([]string, 0, len(names))
	for _, name := range names {
		if willBuild(name) {
			files = append(files, filepath.Join(dir, name))
		}
	}
	rules.Check(dir, files)
}

func checkFiles(paths []string) {
	dirs := make(map[string][]string)
	for _, p := range paths {
		dir := filepath.Dir(p)
		dirs[dir] = append(dirs[dir], p)
	}
	for dir, files := range dirs {
		rules.Check(dir, files)
	}
}

func getTargets() (traDirs, dirs, files []string) {
	if len(flag.Args()) == 0 {
		dirs = []string{"."}
		return
	}
	for _, path := range flag.Args() {
		traverse := strings.HasSuffix(path, "/...")
		if traverse {
			path = strings.TrimSuffix(path, "/...")
		}

		switch mode := fileMode(path); {
		case mode.IsDir():
			if traverse {
				traDirs = append(traDirs, filepath.Clean(path))
			} else {
				dirs = append(dirs, filepath.Clean(path))
			}
		case mode.IsRegular():
			if willBuild(filepath.Base(path)) {
				files = append(files, filepath.Clean(path))
			}
		}
	}
	return
}

func willBuild(name string) bool {
	return filepath.Ext(name) == `.go` && name[0] != '.' && name[0] != '_'
}

func fileMode(path string) os.FileMode {
	if fi, err := os.Stat(path); err == nil {
		return fi.Mode()
	} else {
		panic(err)
	}
}

func processArgs() {
	var version bool
	flag.BoolVar(&version, `version`, false, `display gopsec version.`)
	flag.Parse()
	if version {
		fmt.Println("gospec version v1.0.0 20210219")
		os.Exit(0)
	}
}
