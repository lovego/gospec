package structpkg

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/problems"
)

type sizeRule struct {
	MaxFields uint `yaml:"maxFields"`
}

func (r sizeRule) check(strut *ast.StructType, key string, fileSet *token.FileSet) {
	if num := strut.Fields.NumFields(); uint(num) > r.MaxFields {
		problems.Add(fileSet.Position(strut.Pos()), fmt.Sprintf(
			"struct size: %d fields, limit: %d", num, r.MaxFields,
		), key+".maxFields")
	}
}
