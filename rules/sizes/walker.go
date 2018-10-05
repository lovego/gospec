package sizes

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/problems"
)

type walker struct {
	fileSet *token.FileSet
}

func (w walker) Visit(n ast.Node) ast.Visitor {
	switch node := n.(type) {
	case *ast.FuncDecl:
		w.checkFunc("func "+node.Name.Name, node.Type, node.Body)
	case *ast.FuncLit:
		w.checkFunc("literal func", node.Type, node.Body)
	}
	return w
}

func (w walker) checkFunc(name string, typ *ast.FuncType, body *ast.BlockStmt) {
	if typ.Params != nil && typ.Params.NumFields() > Rules.Params {
		problems.Add(
			w.fileSet.Position(typ.Params.Pos()),
			fmt.Sprintf(`%s params size: %d, limit: %d`, name, typ.Params.NumFields(), Rules.Params),
			`sizes.params`,
		)
	}

	if typ.Results != nil && typ.Results.NumFields() > Rules.Results {
		problems.Add(
			w.fileSet.Position(typ.Results.Pos()),
			fmt.Sprintf(`%s results size: %d, limit: %d`, name, typ.Results.NumFields(), Rules.Results),
			`sizes.results`,
		)
	}

	if size := stmtsNum(body); size > Rules.Func {
		problems.Add(
			w.fileSet.Position(typ.Params.Pos()),
			fmt.Sprintf(`%s body size: %d statements, limit: %d`, name, size, Rules.Func),
			`sizes.func`,
		)
	}
}

func stmtsNum(node ast.Node) int {
	w := &stmtsWalker{}
	ast.Walk(w, node)
	return w.count
}

type stmtsWalker struct {
	count int
}

func (w *stmtsWalker) Visit(node ast.Node) ast.Visitor {
	if _, ok := node.(ast.Stmt); ok {
		if _, ok := node.(*ast.BlockStmt); !ok {
			w.count++
		}
	}
	return w
}
