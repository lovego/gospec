package structpkg

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/name"
)

var Rule = RuleT{
	Name: name.Rule{MaxLen: 30, Style: "camelCase"},
	Size: sizeRule{MaxFields: 100},
}

type RuleT struct {
	Name name.Rule
	Size sizeRule
}

type sizeRule struct {
	MaxFields uint `yaml:"maxFields"`
}

func Check(node ast.Node, fileSet *token.FileSet) {
	switch strut := node.(type) {
	case *ast.StructType:
		checkStruct(strut, fileSet)
	}
}

func checkStruct(strut *ast.StructType, fileSet *token.FileSet) {
	for _, f := range strut.Fields.List {
		for _, ident := range f.Names {
			Rule.Name.Exec(ident.Name, `struct field`, `struct.fieldName`, fileSet.Position(ident.Pos()))
		}
	}
	if num := strut.Fields.NumFields(); uint(num) > Rule.Size.MaxFields {
		problems.Add(fileSet.Position(strut.Pos()), fmt.Sprintf(
			"struct size: %d fields, limit: %d", num, Rule.Size.MaxFields,
		), "struct.size.maxFields")
	}
}
