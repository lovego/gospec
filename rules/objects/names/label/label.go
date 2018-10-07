package labelpkg

import (
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/rules/name"
)

var Rule = name.Rule{
	MaxLen: 30,
	Style:  "lowerCamelCase",
}

func Check(node ast.Node, fileSet *token.FileSet) {
	switch label := node.(type) {
	case *ast.LabeledStmt:
		checkLabel(label.Label, fileSet)
	}
}

func checkLabel(ident *ast.Ident, fileSet *token.FileSet) {
	Rule.Exec(ident.Name, `label`, `label`, fileSet.Position(ident.Pos()))
}
