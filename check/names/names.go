package names

import (
	"go/ast"
	"go/token"

	"github.com/bughou-go/spec/specs"
)

func Check(dir string, pkgs []*ast.Pakcage) {
}

func doFileNode(f *ast.File, fs *token.FileSet) {
	// ast.Walk(walker)
}
