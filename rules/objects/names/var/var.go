package varpkg

import (
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/rules/name"
)

var Var = Rule{
	thing: "var",
	key:   "var",
	Rule: name.Rule{
		MaxLen: 40,
		Style:  "camelCase",
	},
}

var LocalVar = Rule{
	thing: "local var",
	key:   "localVar",
	Rule: name.Rule{
		MaxLen: 30,
		Style:  "lowerCamelCase",
	},
}

func Check(local bool, node ast.Node, fileSet *token.FileSet) {
	switch n := node.(type) {
	case *ast.ValueSpec:
		if local {
			LocalVar.checkDecl(n, fileSet)
		} else {
			Var.checkDecl(n, fileSet)
		}
		return
	case *ast.FuncType:
		LocalVar.checkFuncType(n, fileSet)
		return
	}
	if local {
		switch n := node.(type) {
		case *ast.AssignStmt:
			LocalVar.checkShortDecl(n, fileSet)
		case *ast.RangeStmt:
			LocalVar.checkRangeStmt(n, fileSet)
		}
	}
}
