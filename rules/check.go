package rules

import (
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
	astFiles, sources := parseFiles(fileSet, files)

	names.Check(dir, astFiles, fileSet)
	sizes.Check(dir, astFiles, sources, fileSet)
}

func parseFiles(fset *token.FileSet, filesPath []string) ([]*ast.File, []string) {
	asts := make([]*ast.File, 0, len(filesPath))
	sources := make([]string, 0, len(filesPath))
	for _, path := range filesPath {
		src, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		astFile, err := parser.ParseFile(fset, path, src, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		asts = append(asts, astFile)
		sources = append(sources, string(src))
	}
	return asts, sources
}
