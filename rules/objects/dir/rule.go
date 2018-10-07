package dirpkg

import (
	"go/token"
	"path/filepath"

	namepkg "github.com/lovego/gospec/rules/name"
)

type Rule struct {
	key  string
	Name namepkg.Rule
	Size sizeRule
}

func (r Rule) check(path string) {
	if path == "" {
		return
	}
	r.checkName(path)
	r.Size.check(path, r.key+".size")
}

func (r Rule) checkName(path string) {
	if path == `.` || path == `..` || path == `/` {
		return
	}
	r.Name.Exec(filepath.Base(path), r.key, r.key+".name", token.Position{Filename: path})
}
