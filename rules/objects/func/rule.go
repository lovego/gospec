package funcpkg

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"

	"github.com/lovego/gospec/problems"
	namepkg "github.com/lovego/gospec/rules/name"
)

var exampleTestCase = regexp.MustCompile(`^Example[_A-Z]`)

type Rule struct {
	key  string
	Name namepkg.Rule
	Size sizeRule
}

func (r *Rule) Check(node ast.Node, fileSet *token.FileSet) {
	switch n := node.(type) {
	case *ast.FuncDecl:
		r.checkName(n, fileSet)
		thing := "func " + n.Name.Name
		r.Size.checkType(n.Type, thing, r.key+".size", fileSet)
		r.Size.checkBody(n.Body, thing, r.key+".size", fileSet)
	case *ast.FuncLit:
		r.Size.checkType(n.Type, "literal func", r.key+".size", fileSet)
		r.Size.checkBody(n.Body, "literal func", r.key+".size", fileSet)
	case *ast.InterfaceType:
		r.checkInterface(n, fileSet)
	}
}

func (r *Rule) checkName(fun *ast.FuncDecl, fileSet *token.FileSet) {
	name := fun.Name.Name
	if fun.Recv == nil && r.key == "funcInTest" && exampleTestCase.MatchString(name) {
		if uint(len(name)) > r.Name.MaxLen {
			problems.Add(fileSet.Position(fun.Name.Pos()), fmt.Sprintf(
				"func name %s %d chars long, limit: %d", name, len(name), r.Name.MaxLen,
			), "funcInTest.name.maxLen")
		}
		return
	}
	r.Name.Exec(name, "func", r.key+".name", fileSet.Position(fun.Name.Pos()))
}

func (r *Rule) checkInterface(ifc *ast.InterfaceType, fileSet *token.FileSet) {
	for _, f := range ifc.Methods.List {
		thing := "interface method"
		for _, ident := range f.Names {
			r.Name.Exec(ident.Name, thing, r.key+".name", fileSet.Position(ident.Pos()))
		}
		if len(f.Names) > 0 {
			thing += " " + f.Names[0].Name
		}
		if typ, ok := f.Type.(*ast.FuncType); ok {
			r.Size.checkType(typ, thing, r.key+".size", fileSet)
		}
	}
}
