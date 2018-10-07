package pkgpkg

import (
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/rules/name"
)

var Pkg = name.Rule{
	MaxLen: 20,
	Style:  "lower_case",
}

type Checker struct {
	m map[string]bool
}

func NewChecker() Checker {
	return Checker{m: make(map[string]bool)}
}

func (c Checker) Check(name *ast.Ident, fileSet *token.FileSet) {
	if !c.m[name.Name] {
		Pkg.Exec(name.Name, "pkg", "pkg", fileSet.Position(name.Pos()))
		c.m[name.Name] = true
	}
}
