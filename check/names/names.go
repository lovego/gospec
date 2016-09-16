package names

import (
	"go/ast"
	"go/token"

	"github.com/bughou-go/spec/d"
)

type Config struct {
	Dir, File, Pkg, Func, Const, Var, LocalConst, LocalVar nameConfig
}

var DefaultConfig = Config{
	Dir:        nameConfig{Style: `lower`, MaxLen: 10},
	File:       nameConfig{Style: `lower`, MaxLen: 10},
	Pkg:        nameConfig{Style: `lower`, MaxLen: 10},
	Func:       nameConfig{Style: `camel`, MaxLen: 20},
	Const:      nameConfig{Style: `camel`, MaxLen: 20},
	Var:        nameConfig{Style: `camel`, MaxLen: 20},
	LocalConst: nameConfig{Style: `camel`, MaxLen: 10},
	LocalVar:   nameConfig{Style: `camel`, MaxLen: 10},
}

type nameConfig struct {
	Style  string
	MaxLen uint8
}

func Check(dir *d.Dir) {
}

func doFileNode(f *ast.File, fs *token.FileSet) {
	// ast.Walk(walker)
}
