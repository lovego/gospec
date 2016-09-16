package names

import (
	"go/ast"
	"go/token"

	"github.com/bughou-go/spec/d"
)

type TConfig struct {
	Dir, File, Pkg, Func, Const, Var, LocalConst, LocalVar tConfig
}
type tConfig struct {
	Style  string
	MaxLen uint8
}

var Config = TConfig{
	Dir:        tConfig{Style: `lower`, MaxLen: 10},
	File:       tConfig{Style: `lower`, MaxLen: 10},
	Pkg:        tConfig{Style: `lower`, MaxLen: 10},
	Func:       tConfig{Style: `camel`, MaxLen: 20},
	Const:      tConfig{Style: `camel`, MaxLen: 20},
	Var:        tConfig{Style: `camel`, MaxLen: 20},
	LocalConst: tConfig{Style: `camel`, MaxLen: 10},
	LocalVar:   tConfig{Style: `camel`, MaxLen: 10},
}

func Check(dir *d.Dir) {
}

func doFileNode(f *ast.File, fs *token.FileSet) {
	// ast.Walk(walker)
}
