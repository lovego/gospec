package walker

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

type Walker struct {
	SrcFile string
	AstFile *ast.File
	FileSet *token.FileSet
}

func New(path string) *Walker {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return Parse(path, string(src))
}

func Parse(path, src string) *Walker {
	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, path, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return &Walker{SrcFile: src, AstFile: astFile, FileSet: fileSet}
}

func (w *Walker) Walk(fun func(isLocal bool, node ast.Node)) {
	ast.Walk(visitor{fun: fun}, w.AstFile)
}
