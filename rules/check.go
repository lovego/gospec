package rules

import (
	"go/ast"
	"strings"

	"github.com/lovego/gospec/rules/walker"

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

// check rules
func Check(dir string, files []string) {
	dirPkg.Check(dir)

	pkgChecker := pkgPkg.NewChecker()
	for _, path := range files {
		isTest := strings.HasSuffix(path, "_test.go")
		w := walker.New(path)

		pkgChecker.Check(w.AstFile.Name, w.FileSet)
		filePkg.Check(isTest, path, w.SrcFile, w.AstFile, w.FileSet)

		w.Walk(func(isLocal bool, node ast.Node) {
			funcPkg.Check(isTest, node, w.FileSet)
			structPkg.Check(node, w.FileSet)

			constPkg.Check(isLocal, node, w.FileSet)
			varPkg.Check(isLocal, node, w.FileSet)
			typePkg.Check(isLocal, node, w.FileSet)
			labelPkg.Check(node, w.FileSet)
		})
	}
}
