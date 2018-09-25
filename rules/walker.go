package rules

import (
	"go/ast"
	"go/token"
)

type NodeWalker interface {
	Visit(ast.Node) ast.Node
}

type walker struct {
	f       *ast.File
	file    *token.File
	src     []string
	walkers []NodeWalker
}

func (w walker) Visit(node ast.Node) ast.Visitor {
	for _, nw := range w.walkers {
		nw.Visit(node)
	}
	return w
}
