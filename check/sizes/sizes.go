package sizes

import (
	"go/ast"
	"go/token"

	"github.com/bughou-go/spec/c"
)

type TConfig struct {
	Dir, File, Row, Func uint
}

var Config = TConfig{Dir: 20, File: 200, Row: 100, Func: 20}

func Check(dir *c.Dir) {
	for _, pkg := range dir.Pkgs {
		for _, f := range pkg.Files {
		}
	}
}

func checkFileSize(dir *c.Dir) {
}
