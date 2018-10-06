package fun

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

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
	MaxParams, MaxResults, MaxStatements uint
}

func Check(node ast.Node, fileSet *token.FileSet) {
	switch fun := node.(type) {
	case *ast.FuncDecl:
		checkFuncName(fun, fileSet)
		thing := "func " + fun.Name.Name
		checkFuncType(thing, fun.Type, fileSet)
		checkFuncBody(thing, fun.Body, fileSet)
	case *ast.FuncLit:
		checkFuncType("literal func", fun.Type, fileSet)
		checkFuncBody("literal func", fun.Body, fileSet)
	}
}

func checkFuncName(fun *ast.FuncDecl, fileSet *token.FileSet) {
	position := fileSet.Position(fun.Name.Pos())

	name := fun.Name.Name
	if strings.HasSuffix(position.Filename, "_test.go") {
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
