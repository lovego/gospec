package c

import (
	"fmt"
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

var problemsLimit, problemsCount uint = 10, 0

func Problem(pos token.Position, desc, rule string) {
	fmt.Printf("%s: %s (%s)\n", pos, desc, rule)

	problemsCount++
	if problemsLimit > 0 && problemsCount > problemsLimit {
		os.Exit(1)
	}
}

func ProblemsCount() uint {
	return problemsCount
}
