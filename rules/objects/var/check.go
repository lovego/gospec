package names

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

func Check(dir string, astFiles []*ast.File, fileSet *token.FileSet) {
	w := walker{fileSet: fileSet}
	for i, astFile := range astFiles {
		if i == 0 {
		}
		ast.Walk(w, astFile)
	}
	if fun.Recv != nil {
		// checkFieldList(`func receiver`, fun.Recv, true, w.fileSet)
	}
}

func checkIdent(thing string, ident *ast.Ident, local bool, fileSet *token.FileSet) {
	if ident == nil || ident.Obj == nil {
		return
	}
	if rule := getRuleForIdent(ident, local, fileSet); rule.valid() {
		rule.Exec(thing, ident.Name, fileSet.Position(ident.Pos()))
	}
}
