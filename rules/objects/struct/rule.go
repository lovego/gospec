package structpkg

import (
	"go/ast"
	"go/token"

	namepkg "github.com/lovego/gospec/rules/name"
)

type Rule struct {
	key       string
	FieldName namepkg.Rule `yaml:"fieldName"`
	Size      sizeRule
}

func (r Rule) check(strut *ast.StructType, fileSet *token.FileSet) {
	r.checkFieldsName(strut, fileSet)
	r.Size.check(strut, r.key+".size", fileSet)
}

func (r Rule) checkFieldsName(strut *ast.StructType, fileSet *token.FileSet) {
	for _, f := range strut.Fields.List {
		for _, ident := range f.Names {
			r.FieldName.Exec(ident.Name, `struct field`, r.key+`.fieldName`, fileSet.Position(ident.Pos()))
		}
	}
}
