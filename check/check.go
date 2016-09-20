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
				names.CheckIdent(File.Name, file)
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
	case *ast.FuncDecl, *ast.FuncLit:
		names.CheckFunc(n, w.File)
		sizes.CheckFunc(n, w.File)
	case *ast.AssignStmt:
		names.CheckShortVarDecl(n, w.File)
	case *ast.GenDecl:
		names.CheckGenDecl(n, w.File)
	case *ast.InterfaceType:
		// Do not check interface method names.
		// They are often constrainted by the method names of concrete types.
		for _, x := range v.Methods.List {
			ft, ok := x.Type.(*ast.FuncType)
			if !ok { // might be an embedded interface name
				continue
			}
			checkList(ft.Params, "interface method parameter")
			checkList(ft.Results, "interface method result")
		}
	case *ast.RangeStmt:
		if v.Tok == token.ASSIGN {
			return true
		}
		if id, ok := v.Key.(*ast.Ident); ok {
			check(id, "range var")
		}
		if id, ok := v.Value.(*ast.Ident); ok {
			check(id, "range var")
		}
	case *ast.StructType:
		for _, f := range v.Fields.List {
			for _, id := range f.Names {
				check(id, "struct field")
			}
		}
	}
	return w
}
