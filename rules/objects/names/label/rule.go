package labelpkg

import (
	"go/ast"
	"go/token"

	namepkg "github.com/lovego/gospec/rules/name"
)

type Rule struct {
	thing, key   string
	namepkg.Rule `yaml:",inline"`
}

func (r Rule) check(name *ast.Ident, fileSet *token.FileSet) {
	r.Exec(name.Name, r.thing, r.key, fileSet.Position(name.Pos()))
}
