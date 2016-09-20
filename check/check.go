package check

import (
	"go/ast"
	"go/token"

	"github.com/bughou-go/spec/check/names"
	"github.com/bughou-go/spec/check/sizes"
)

type Dir struct {
	Path string
	Fset *token.FileSet
	Pkgs map[string]*ast.Package
}

func Check(dir *Dir) {
	names.CheckDir(dir.Path)
	sizes.CheckDir(dir.Path)
	for _, pkg := range dir.Pkgs {
		first := true
		for p, f := range pkg.Files {
			file := dir.Fset.File(f.Package)
			if first {
				names.CheckIdent(f.Name, file, ``)
				first = false
			}
			names.CheckFile(p)
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
	// declare
	case *ast.GenDecl:
		names.CheckGenDecl(n, w.File)
	case *ast.FuncDecl, *ast.FuncLit:
		names.CheckFunc(n, w.File)
		sizes.CheckFunc(n, w.File)
	// statement
	case *ast.AssignStmt:
		names.CheckShortVarDecl(n, w.File)
	case *ast.RangeStmt:
		names.CheckRangeStmt(n, w.File)
	// type define
	case *ast.StructType:
		names.CheckStruct(n, w.File)
	case *ast.InterfaceType:
		names.CheckInterface(n, w.File)
	}
	return w
}
