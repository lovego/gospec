package funcpkg

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/problems"
)

func (r *RuleT) checkTypeSize(thing string, typ *ast.FuncType, fileSet *token.FileSet) {
	if typ.Params != nil && uint(typ.Params.NumFields()) > r.Size.MaxParams {
		problems.Add(
			fileSet.Position(typ.Params.Pos()), fmt.Sprintf(
				`%s params size: %d, limit: %d`, thing, typ.Params.NumFields(), r.Size.MaxParams,
			), r.key+`.size.maxParams`,
		)
	}

	if typ.Results != nil && uint(typ.Results.NumFields()) > r.Size.MaxResults {
		problems.Add(
			fileSet.Position(typ.Results.Pos()), fmt.Sprintf(
				`%s results size: %d, limit: %d`, thing, typ.Results.NumFields(), r.Size.MaxResults,
			), r.key+`.size.maxResults`,
		)
	}
}

func (r *RuleT) checkBodySize(thing string, body *ast.BlockStmt, fileSet *token.FileSet) {
	if size := stmtsCount(body); size > r.Size.MaxStatements {
		problems.Add(
			fileSet.Position(body.Pos()), fmt.Sprintf(
				`%s body size: %d statements, limit: %d`, thing, size, r.Size.MaxStatements,
			), r.key+`.size.maxStatements`,
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
