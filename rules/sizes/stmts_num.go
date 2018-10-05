package sizes

import (
	"go/ast"
)

type stmtsWalker struct {
	count int
}

func (w *stmtsWalker) Visit(node ast.Node) ast.Visitor {
	if _, ok := node.(ast.Stmt); ok {
		if _, ok := node.(*ast.BlockStmt); !ok {
			w.count++
		}
	}
	return w
}

func stmtsNum(node ast.Node) int {
	w := &stmtsWalker{}
	ast.Walk(w, node)
	return w.count
}
