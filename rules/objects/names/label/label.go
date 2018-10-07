package labelpkg

import (
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/rules/name"
)

var Label = Rule{
	thing: "label",
	key:   "label",
	Rule: name.Rule{
		MaxLen: 30,
		Style:  "lowerCamelCase",
	},
}

func Check(node ast.Node, fileSet *token.FileSet) {
	switch label := node.(type) {
	case *ast.LabeledStmt:
		Label.check(label.Label, fileSet)
	}
}
