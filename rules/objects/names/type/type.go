package typepkg

import (
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/rules/name"
)

var Type = Rule{
	thing: "type",
	key:   "type",
	Rule: name.Rule{
		MaxLen: 40,
		Style:  "camelCase",
	},
}

var LocalType = Rule{
	thing: "local type",
	key:   "localType",
	Rule: name.Rule{
		MaxLen: 30,
		Style:  "lowerCamelCase",
	},
}

func Check(local bool, node ast.Node, fileSet *token.FileSet) {
	switch typ := node.(type) {
	case *ast.TypeSpec:
		if local {
			LocalType.check(typ, fileSet)
		} else {
			Type.check(typ, fileSet)
		}
	}
}
