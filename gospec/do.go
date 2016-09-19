package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"

	"github.com/bughou-go/spec/c"
	"github.com/bughou-go/spec/check"
)

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
	doDir(p)
}

func doDir(dir string) {
	var fset = token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	check.Check(&c.Dir{Path: dir, Fset: fset, Pkgs: pkgs})
}

func doFiles(paths []string) {
	dirs := make(map[string][]string)
	for _, p := range paths {
		dir := path.Dir(p)
		dirs[dir] = append(dirs[dir], p)
	}
	for dir, files := range dirs {
		doDirFiles(dir, files)
	}
}

func doDirFiles(dir string, files []string) {
	var fset = token.NewFileSet()
	pkgs := make(map[string]*ast.Package)
	for _, p := range files {
		f, err := parser.ParseFile(fset, p, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		setupPkgs(p, f, pkgs)
	}
	check.Check(&c.Dir{Path: dir, Fset: fset, Pkgs: pkgs})
}

func setupPkgs(p string, f *ast.File, pkgs map[string]*ast.Package) {
	name := f.Name.Name
	if pkg := pkgs[name]; pkg != nil {
		pkg.Files[p] = f
		return
	}
	pkgs[name] = &ast.Package{
		Name:  name,
		Files: map[string]*ast.File{p: f},
	}
}
