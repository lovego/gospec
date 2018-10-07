package filepkg

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	namepkg "github.com/lovego/gospec/rules/name"
)

type Rule struct {
	key  string
	Name namepkg.Rule
	Size sizeRule
}

func (r *Rule) Check(path, src string, astFile *ast.File, fileSet *token.FileSet) {
	r.checkName(path)
	r.Size.check(src, path, r.key+".size", astFile, fileSet)
}

func (r *Rule) checkName(path string) {
	r.Name.Exec(
		strings.TrimSuffix(strings.TrimSuffix(filepath.Base(path), `.go`), `_test`),
		"file", r.key+".name", token.Position{Filename: path},
	)
}
