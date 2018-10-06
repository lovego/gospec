package funcpkg

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/name"
)

var Rule = RuleT{
	Name: name.Rule{MaxLen: 30, Style: "camelCase"},
	Size: sizeRule{MaxParams: 5, MaxResults: 3, MaxStatements: 30},
}

var RuleInTest = RuleT{
	Name: name.Rule{MaxLen: 50, Style: "camelCase"},
	Size: sizeRule{MaxParams: 5, MaxResults: 3, MaxStatements: 30},
}

var exampleTestCase = regexp.MustCompile(`^Example[_A-Z]`)

type RuleT struct {
	Name name.Rule
	Size sizeRule
}

type sizeRule struct {
	MaxParams     uint `yaml:"maxParams"`
	MaxResults    uint `yaml:"maxResults"`
	MaxStatements uint `yaml:"maxStatements"`
}

func Check(node ast.Node, isTest bool, fileSet *token.FileSet) {
	switch n := node.(type) {
	case *ast.FuncDecl:
		checkName(n, isTest, fileSet)
		thing := "func " + n.Name.Name
		checkTypeSize(thing, n.Type, isTest, fileSet)
		checkBodySize(thing, n.Body, isTest, fileSet)
	case *ast.FuncLit:
		checkTypeSize("literal func", n.Type, isTest, fileSet)
		checkBodySize("literal func", n.Body, isTest, fileSet)
	case *ast.InterfaceType:
		checkInterface(n, isTest, fileSet)
	}
}

func checkName(fun *ast.FuncDecl, isTest bool, fileSet *token.FileSet) {
	position := fileSet.Position(fun.Name.Pos())

	name := fun.Name.Name
	if isTest {
		if fun.Recv == nil && exampleTestCase.MatchString(name) {
			if uint(len(name)) > RuleInTest.Name.MaxLen {
				problems.Add(position, fmt.Sprintf(
					"func name %s %d chars long, limit: %d", name, len(name), RuleInTest.Name.MaxLen,
				), "funcInTest.name.maxLen")
			}
		} else {
			RuleInTest.Name.Exec(name, "func", "funcInTest.name", position)
		}
	} else {
		Rule.Name.Exec(name, "func", "func.name", position)
	}
}

func checkInterface(ifc *ast.InterfaceType, isTest bool, fileSet *token.FileSet) {
	if ifc.Methods == nil {
		return
	}
	rule, ruleName := getNameRule(isTest)
	for _, f := range ifc.Methods.List {
		for _, ident := range f.Names {
			rule.Exec(ident.Name, "interface method", ruleName, fileSet.Position(ident.Pos()))
		}
		checkTypeSize("interface method", f.Type.(*ast.FuncType), isTest, fileSet)
	}
}

func getNameRule(isTest bool) (name.Rule, string) {
	if isTest {
		return RuleInTest.Name, "funcInTest.name"
	} else {
		return Rule.Name, "func.name"
	}
}
