package sizes

import (
	"go/ast"
	"go/token"

	"github.com/bughou-go/spec/d"
)

type TConfig struct {
	Dir, File, Row, Func uint
}

var Config = TConfig{Dir: 20, File: 200, Row: 100, Func: 20}

func Check(dir *d.Dir) {
}

func checkFileSize(f *ast.File, fs *token.FileSet) {
}
