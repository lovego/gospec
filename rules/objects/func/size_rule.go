package funcpkg

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/problems"
)

type sizeRule struct {
	MaxParams     uint `yaml:"maxParams"`
	MaxResults    uint `yaml:"maxResults"`
	MaxStatements uint `yaml:"maxStatements"`
}

func (r sizeRule) checkType(typ *ast.FuncType, thing, key string, fileSet *token.FileSet) {
	if typ.Params != nil && uint(typ.Params.NumFields()) > r.MaxParams {
		problems.Add(
			fileSet.Position(typ.Params.Pos()), fmt.Sprintf(
				`%s params size: %d, limit: %d`, thing, typ.Params.NumFields(), r.MaxParams,
			), key+`.maxParams`,
		)
	}

	if typ.Results != nil && uint(typ.Results.NumFields()) > r.MaxResults {
		problems.Add(
			fileSet.Position(typ.Results.Pos()), fmt.Sprintf(
				`%s results size: %d, limit: %d`, thing, typ.Results.NumFields(), r.MaxResults,
			), key+`.maxResults`,
		)
	}
}

func (r sizeRule) checkBody(body *ast.BlockStmt, thing, key string, fileSet *token.FileSet) {
	if size := stmtsCount(body); size > r.MaxStatements {
		problems.Add(
			fileSet.Position(body.Pos()), fmt.Sprintf(
				`%s body size: %d statements, limit: %d`, thing, size, r.MaxStatements,
			), key+`.maxStatements`,
		)
	}
}

func stmtsCount(node ast.Node) uint {
	w := &stmtsWalker{}
	ast.Walk(w, node)
	return w.count
}

type stmtsWalker struct {
	count uint
}

func (w *stmtsWalker) Visit(node ast.Node) ast.Visitor {
	if _, ok := node.(ast.Stmt); ok {
		if _, ok := node.(*ast.BlockStmt); !ok {
			w.count++
		}
	}
	return w
}
