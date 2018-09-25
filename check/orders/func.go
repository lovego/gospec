package orders

import (
	"go/ast"
)

type funcInfo struct {
	packageName string
	funcName    string
	end         int
	begin       int
	called      map[string]map[string]bool //same callee appears only onece
	calledFuncs []string
	funcDecl    *ast.FuncDecl
}
