package rules

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"

	dirPkg "github.com/lovego/gospec/rules/objects/dir"
	filePkg "github.com/lovego/gospec/rules/objects/file"
	funPkg "github.com/lovego/gospec/rules/objects/fun"
	pkgPkg "github.com/lovego/gospec/rules/objects/pkg"
)

func Check(dir string, files []string) {
	dirPkg.Check(dir)

	fileSet := token.NewFileSet()
	packages := make(map[string]bool)
	w := walker{fileSet: fileSet}
	for _, path := range files {
		src, astFile := loadFile(path, fileSet)

		if !packages[astFile.Name.Name] {
			pkgPkg.Check(astFile.Name, fileSet)
			packages[astFile.Name.Name] = true
		}

		filePkg.Check(path, src, astFile, fileSet)
		ast.Walk(w, astFile)
	}
}

type walker struct {
	local   bool
	fileSet *token.FileSet
}

func (w walker) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return w
	}
	funPkg.Check(node, w.fileSet)

	/*
		switch n := node.(type) {
		case *ast.StructType:
			w.checkStruct(n)
		case *ast.GenDecl:
			w.checkGenDecl(n)
		case *ast.InterfaceType:
			w.checkInterface(n)
		case *ast.AssignStmt:
			w.checkShortVarDecl(n)
		case *ast.RangeStmt:
			w.checkRangeStmt(n)
		}
	*/

	if _, ok := node.(*ast.BlockStmt); ok && !w.local {
		return walker{local: true, fileSet: w.fileSet}
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
