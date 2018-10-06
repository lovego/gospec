package pkg

import (
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/rules/name"
)

var Rule = name.Rule{
	MaxLen: 20,
	Style:  "lower_case",
}

func Check(name *ast.Ident, fileSet *token.FileSet) {
	Rule.Exec(name.Name, "pkg", "pkg", fileSet.Position(name.Pos()))
}
