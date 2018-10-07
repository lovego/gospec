package typepkg

import (
	"go/ast"
	"go/token"

	namepkg "github.com/lovego/gospec/rules/name"
)

type Rule struct {
	thing, key   string
	namepkg.Rule `yaml:",inline"`
}

func (r Rule) check(typ *ast.TypeSpec, fileSet *token.FileSet) {
	ident := typ.Name
	r.Exec(ident.Name, r.thing, r.key, fileSet.Position(ident.Pos()))
}
