package dirpkg

import (
	namepkg "github.com/lovego/gospec/rules/name"
)

var Dir = Rule{
	key:  "dir",
	Name: namepkg.Rule{MaxLen: 30, Style: "lower_case"},
	Size: sizeRule{MaxEntries: 20},
}

func Check(path string) {
	Dir.check(path)
}
