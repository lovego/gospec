package check

import (
	"go/ast"
	"go/token"

	//"github.com/bughou-go/spec/check/names"
	"github.com/bughou-go/spec/check/sizes"
)

type Dir struct {
	Path string
	Fset *token.FileSet
	Pkgs map[string]*ast.Package
}

func Check(dir *Dir) {
	sizes.CheckDir(dir.Path)
	for _, pkg := range dir.Pkgs {
		for _, f := range pkg.Files {
			file := dir.Fset.File(f.Package)
			sizes.CheckFile(file)
			sizes.CheckLines(file)
			ast.Walk(walker{file}, f)
		}
	}
}

type walker struct {
	*token.File
}

func (w walker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl, *ast.FuncLit:
		sizes.CheckFunc(n, w.File)
	}
	return w
}
