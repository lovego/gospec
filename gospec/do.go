package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	//"os"

	"github.com/bughou-go/spec/c"
	"github.com/bughou-go/spec/check"
)

func traverseDir(p string) {
}

func doDir(p string) {
	var fset = token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, p, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	check.Check(&c.Dir{Path: p, Fset: fset, Pkgs: pkgs})
}

func doFiles(ps []string) {
	var fset = token.NewFileSet()
	pkgs := make(map[string]*ast.Package)
	for _, p := range ps {
		if path.Ext(p) != `.go` {
			continue
		}
		f, err := parser.ParseFile(fset, p, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		setupPkgs(p, f, pkgs)
	}
	check.Check(&c.Dir{Fset: fset, Pkgs: pkgs})
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
