package spec

import (
	"go/parser"
	"go/token"
	"os"
	"path"
)

func checkDirRecursively(dirname string) {
}

func checkDir(dirname string) {
}

func checkFiles(filenames []string) {
}

func checkFile(p string) {
	if path.Ext(p) != `.go` {
		return
	}
	var fs = token.NewFileSet()
	f, err := parser.ParseFile(fs, p, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	checkFileNode(f, fs)
}

func checkFileNode(f *ast.File, fs *token.FileSet) {
}
