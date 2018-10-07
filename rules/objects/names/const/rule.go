package constpkg

import (
	"go/ast"
	"go/token"

	namepkg "github.com/lovego/gospec/rules/name"
)

type Rule struct {
	thing, key   string
	namepkg.Rule `yaml:",inline"`
}

func (r Rule) check(value *ast.ValueSpec, fileSet *token.FileSet) {
	for _, ident := range value.Names {
		r.Exec(ident.Name, r.thing, r.key, fileSet.Position(ident.Pos()))
	}
}
