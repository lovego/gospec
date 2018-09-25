package names

import (
	"go/ast"
	"go/token"
)

type walker struct {
	local   bool
	fileSet *token.FileSet
}

func (w walker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Ident:
		checkIdent(n, w.local, w.fileSet)
	case *ast.StructType:
		w.visitStruct(n)
		return nil
	case *ast.BlockStmt:
		if !w.local {
			return walker{local: true, fileSet: w.fileSet}
		}
	}
	return w
}

func (w walker) visitStruct(st *ast.StructType) {
	for _, f := range st.Fields.List {
		for _, ident := range f.Names {
			Rules.StructField.exec(ident.Name, w.fileSet.Position(ident.Pos()))
		}
		ast.Walk(w, f.Type)
	}
}
