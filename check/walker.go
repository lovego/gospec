package check

import (
	"go/ast"
	"go/token"

	"github.com/lovego/spec/check/names"
	"github.com/lovego/spec/check/sizes"
)

type walker struct {
	f    *ast.File
	file *token.File
	src  []string
}

func (w walker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		names.CheckFuncDecl(n, w.file)
		sizes.CheckFunc(n, w.f, w.file, w.src)
	case *ast.BlockStmt:
		return walkerLocal{w}
	default:
		w.process(n, false)
	}
	return w
}

type walkerLocal struct {
	walker
}

func (w walkerLocal) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	// statement
	case *ast.AssignStmt:
		names.CheckShortVarDecl(n, w.file)
	case *ast.RangeStmt:
		names.CheckRangeStmt(n, w.file)
	default:
		w.process(n, true)
	}
	return w
}

func (w walker) process(node ast.Node, local bool) {
	switch n := node.(type) {
	// declare
	case *ast.GenDecl:
		names.CheckGenDecl(n, local, w.file)
	// type define
	case *ast.StructType:
		names.CheckStruct(n, local, w.file)
	case *ast.InterfaceType:
		names.CheckInterface(n, local, w.file)
	// func literal
	case *ast.FuncLit:
		names.CheckFuncLit(n, w.file)
		sizes.CheckFunc(n, w.f, w.file, w.src)
	}
}
