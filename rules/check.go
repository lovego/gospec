package rules

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"

	dirPkg "github.com/lovego/gospec/rules/objects/dir"
	filePkg "github.com/lovego/gospec/rules/objects/file"
	funcPkg "github.com/lovego/gospec/rules/objects/func"
	structPkg "github.com/lovego/gospec/rules/objects/struct"

	constPkg "github.com/lovego/gospec/rules/objects/names/const"
	labelPkg "github.com/lovego/gospec/rules/objects/names/label"
	pkgPkg "github.com/lovego/gospec/rules/objects/names/pkg"
	typePkg "github.com/lovego/gospec/rules/objects/names/type"
	varPkg "github.com/lovego/gospec/rules/objects/names/var"
)

func Check(dir string, files []string) {
	dirPkg.Check(dir)

	fileSet := token.NewFileSet()
	packages := make(map[string]bool)
	for _, path := range files {
		src, astFile := loadFile(path, fileSet)

		if !packages[astFile.Name.Name] {
			pkgPkg.Check(astFile.Name, fileSet)
			packages[astFile.Name.Name] = true
		}

		isTest := strings.HasSuffix(path, "_test.go")

		filePkg.Check(path, src, isTest, astFile, fileSet)

		ast.Walk(walker{isTest: isTest, fileSet: fileSet}, astFile)
	}
}

type walker struct {
	isTest, isLocal bool
	fileSet         *token.FileSet
}

func (w walker) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return w
	}
	funcPkg.Check(node, w.isTest, w.fileSet)
	structPkg.Check(node, w.fileSet)

	constPkg.Check(node, w.isLocal, w.fileSet)
	varPkg.Check(node, w.isLocal, w.fileSet)
	typePkg.Check(node, w.isLocal, w.fileSet)
	labelPkg.Check(node, w.fileSet)

	if _, ok := node.(*ast.BlockStmt); ok && !w.isLocal {
		return walker{isTest: w.isTest, fileSet: w.fileSet, isLocal: true}
	}
	return w
}

func loadFile(path string, fileSet *token.FileSet) (string, *ast.File) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	astFile, err := parser.ParseFile(fileSet, path, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return string(src), astFile
}
