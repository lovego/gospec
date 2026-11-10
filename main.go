package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/gospec/gitignore"
	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules"
)

func main() {
	processArgs()
	rules.LoadConfig()
	traDirs, dirs, files := getTargets()
	for _, dir := range traDirs {
		traverseDir(dir)
	}
	for _, dir := range dirs {
		checkDir(dir, gitignore.New(dir))
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
	ig := gitignore.New(dir)
	if ig.Ignored(dir, true) {
		return
	}
	walkDir(dir, ig)
}

func walkDir(dir string, ig *gitignore.Set) {
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
		name := d.Name()
		if d.IsDir() && name[0] != '.' && name[0] != '_' {
			child := filepath.Join(dir, name)
			if ig.Ignored(child, true) {
				continue
			}
			walkDir(child, ig.Enter(child))
		}
	}
	checkDir(dir, ig)
}

func checkDir(dir string, ig *gitignore.Set) {
	if ig.Ignored(dir, true) {
		return
	}
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
		path := filepath.Join(dir, name)
		if willBuild(name) && !ig.Ignored(path, false) {
			files = append(files, path)
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
	args := flag.Args()
	if rules.Targets() != nil {
		args = rules.Targets()
	}
	if len(args) == 0 {
		dirs = []string{"."}
		return
	}
	for _, path := range args {
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
			path = filepath.Clean(path)
			if willBuild(filepath.Base(path)) && !ignoredFile(path) {
				files = append(files, path)
			}
		}
	}
	return
}

func willBuild(name string) bool {
	return filepath.Ext(name) == `.go` && name[0] != '.' && name[0] != '_'
}

func ignoredFile(path string) bool {
	return gitignore.New(filepath.Dir(path)).Ignored(path, false)
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
		fmt.Println("gospec version v1.0.1 20261110")
		os.Exit(0)
	}
}
