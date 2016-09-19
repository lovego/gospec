package c

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strconv"
)

type Dir struct {
	Path string
	Fset *token.FileSet
	Pkgs map[string]*ast.Package
}

type Walker func(ast.Node) bool

func (w Walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

var problemsLimit, problemsCount uint = 10, 0

func Problem(position token.Position, ident string, rule [2]string) {
	pos := position.Filename
	if position.Line > 0 {
		pos += `:` + strconv.Itoa(position.Line)
		if position.Column > 0 {
			pos += `:` + strconv.Itoa(position.Column)
		}
	}
	fmt.Printf("%s:%s \t%s (%s)\n", pos, ident, rule[1], rule[0])

	problemsCount++
	if problemsLimit > 0 && problemsCount > problemsLimit {
		os.Exit(1)
	}
}

func ProblemsCount() uint {
	return problemsCount
}
