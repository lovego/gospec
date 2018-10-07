package structpkg

import (
	"go/ast"
	"go/token"

	namepkg "github.com/lovego/gospec/rules/name"
)

var Struct = Rule{
	key: "struct",
	FieldName: namepkg.Rule{
		MaxLen: 30,
		Style:  "camelCase",
	},
	Size: sizeRule{
		MaxFields: 100,
	},
}

func Check(node ast.Node, fileSet *token.FileSet) {
	switch strut := node.(type) {
	case *ast.StructType:
		Struct.check(strut, fileSet)
	}
}
