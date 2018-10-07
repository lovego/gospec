package walker

import (
	"go/ast"
)

type visitor struct {
	fun     func(isLocal bool, node ast.Node)
	isLocal bool
}

func (v visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	v.fun(v.isLocal, node)

	if !v.isLocal {
		if _, ok := node.(*ast.BlockStmt); ok {
			return visitor{fun: v.fun, isLocal: true}
		}
	}
	return v
}
