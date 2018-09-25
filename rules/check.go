package rules

import (
	"bufio"
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"

	"github.com/lovego/gospec/rules/names"
	// "github.com/lovego/gospec/rules/order"
	"github.com/lovego/gospec/rules/sizes"
)

func Check(dir string, files []string) {
	fileSet := token.NewFileSet()
	astFiles := parseFiles(fileSet, filePath)

	names.Check(dir, astFiles, fileSet)
	sizes.Check(dir, astFiles, fileSet)
}

func parseFiles(fset token.FileSet, filesPath []string) []*ast.File {
	files := make([]*ast.File, 0, len(filesPath))
	for _, path := range files {
		src, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		fileAst, err := parser.ParseFile(fset, path, src, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		files = append(files, fileAst)
	}
	return files
}
