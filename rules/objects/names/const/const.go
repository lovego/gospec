package constpkg

import (
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/rules/name"
)

var Rule = name.Rule{
	MaxLen: 30,
	Style:  "camelCase",
}
var LocalRule = name.Rule{
	MaxLen: 20,
	Style:  "lowerCamelCase",
}

func Check(node ast.Node, local bool, fileSet *token.FileSet) {
	switch con := node.(type) {
	case *ast.ValueSpec:
		checkConst(con, local, fileSet)
	}
}

func checkConst(con *ast.ValueSpec, local bool, fileSet *token.FileSet) {
	if con.Names[0].Obj.Kind != ast.Con {
		return
	}
	for _, ident := range con.Names {
		if local {
			LocalRule.Exec(ident.Name, `local const`, `localConst`, fileSet.Position(ident.Pos()))
		} else {
			Rule.Exec(ident.Name, `const`, `const`, fileSet.Position(ident.Pos()))
		}
	}
}
