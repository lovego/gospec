package funcpkg

import (
	"go/ast"
	"go/token"

	namepkg "github.com/lovego/gospec/rules/name"
)

var Func = Rule{
	key: "func",
	Name: namepkg.Rule{
		MaxLen: 30, Style: "camelCase",
	},
	Size: sizeRule{
		MaxParams: 5, MaxResults: 3, MaxStatements: 30,
	},
}

var FuncInTest = Rule{
	key: "funcInTest",
	Name: namepkg.Rule{
		MaxLen: 50, Style: "camelCase",
	},
	Size: sizeRule{
		MaxParams: 5, MaxResults: 3, MaxStatements: 30,
	},
}

func Check(isTest bool, node ast.Node, fileSet *token.FileSet) {
	if isTest {
		FuncInTest.Check(node, fileSet)
	} else {
		Func.Check(node, fileSet)
	}
}
