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

var Rule = RuleT{
	key:  "func",
	Name: namepkg.Rule{MaxLen: 30, Style: "camelCase"},
	Size: sizeRule{MaxParams: 5, MaxResults: 3, MaxStatements: 30},
}

var RuleInTest = RuleT{
	key:  "funcInTest",
	Name: namepkg.Rule{MaxLen: 50, Style: "camelCase"},
	Size: sizeRule{MaxParams: 5, MaxResults: 3, MaxStatements: 30},
}

type RuleT struct {
	key  string
	Name namepkg.Rule
	Size sizeRule
}

type sizeRule struct {
	MaxParams     uint `yaml:"maxParams"`
	MaxResults    uint `yaml:"maxResults"`
	MaxStatements uint `yaml:"maxStatements"`
}

func Check(isTest bool, node ast.Node, fileSet *token.FileSet) {
	if isTest {
		RuleInTest.Check(node, fileSet)
	} else {
		Rule.Check(node, fileSet)
	}
}

func (r *RuleT) Check(node ast.Node, fileSet *token.FileSet) {
	switch n := node.(type) {
	case *ast.FuncDecl:
		r.checkName(n, fileSet)
		thing := "func " + n.Name.Name
		r.checkTypeSize(thing, n.Type, fileSet)
		r.checkBodySize(thing, n.Body, fileSet)
	case *ast.FuncLit:
		r.checkTypeSize("literal func", n.Type, fileSet)
		r.checkBodySize("literal func", n.Body, fileSet)
	case *ast.InterfaceType:
		r.checkInterface(n, fileSet)
	}
}

func (r *RuleT) checkName(fun *ast.FuncDecl, fileSet *token.FileSet) {
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

func (r *RuleT) checkInterface(ifc *ast.InterfaceType, fileSet *token.FileSet) {
	if ifc.Methods == nil {
		return
	}
	for _, f := range ifc.Methods.List {
		for _, ident := range f.Names {
			r.Name.Exec(ident.Name, "interface method", r.key+".name", fileSet.Position(ident.Pos()))
		}
		r.checkTypeSize("interface method", f.Type.(*ast.FuncType), fileSet)
	}
}
