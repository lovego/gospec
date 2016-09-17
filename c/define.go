package c

import (
	"go/ast"
	"go/token"
	"os"
)

type Dir struct {
	Path string
	Fset *token.FileSet
	Pkgs map[string]*ast.Package
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

var problemsLimit, problems uint = 10, 0

func Problem() {
	problems++
	if problemsLimit > 0 && problems > problemsLimit {
		os.Exit(1)
	}
}

func Problems() uint {
	return problems
}
