package spec

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	//"os"

	"github.com/bughou-go/spec/names"
	"github.com/bughou-go/spec/sizes"
)

type Dir struct {
	Path string
	Fset *token.FileSet
	Pkgs map[string]*ast.Package
}

func doDirRecursively(p string) {
}

func doDir(p string) {
	var fset = token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, p, nil, parser.ParseComments)
	checkDir(&Dir{Path: p, Fset: fset, Pkgs: pkgs})
}

func doFiles(ps []string) {
	var fset = token.NewFileSet()
	pkgs := make(map[string]*ast.Package)
	for _, p := range ps {
		f := openFile(p, fset)
		name = f.Name.Name
		if pkg := pkgs[name]; pkg != nil {
			pkg.Files[p] = f
			continue
		}
		pkgs[name] = &ast.Package{
			Name:  name,
			Files: map[string]*ast.File{p: f},
		}
	}
	checkDir(&Dir{Fset: fset, Pkgs: pkgs})
}

func openFile(p string, fset *token.FileSet) *ast.File {
	if path.Ext(p) != `.go` {
		return
	}
	f, err := parser.ParseFile(fset, p, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return f
}
