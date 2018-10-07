package constpkg

import (
	"go/ast"
	"go/token"

	namepkg "github.com/lovego/gospec/rules/name"
)

var Const = Rule{
	thing: "const",
	key:   "const",
	Rule: namepkg.Rule{
		MaxLen: 30,
		Style:  "camelCase",
	},
}

var LocalConst = Rule{
	thing: "local const",
	key:   "localConst",
	Rule: namepkg.Rule{
		MaxLen: 20,
		Style:  "lowerCamelCase",
	},
}

func Check(local bool, node ast.Node, fileSet *token.FileSet) {
	switch value := node.(type) {
	case *ast.ValueSpec:
		if value.Names[0].Obj.Kind == ast.Con {
			if local {
				LocalConst.check(value, fileSet)
			} else {
				Const.check(value, fileSet)
			}
		}
	}
}
