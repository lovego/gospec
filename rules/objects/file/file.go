package filepkg

import (
	"go/ast"
	"go/token"

	namepkg "github.com/lovego/gospec/rules/name"
)

var File = Rule{
	key:  "file",
	Name: namepkg.Rule{MaxLen: 20, Style: "lower_case"},
	Size: sizeRule{MaxLineWidth: 100, MaxCommentLineWidth: 120, MaxLines: 300},
}

var TestFile = Rule{
	key:  "testFile",
	Name: namepkg.Rule{MaxLen: 50, Style: "lower_case"},
	Size: sizeRule{MaxLineWidth: 100, MaxCommentLineWidth: 120, MaxLines: 600},
}

func Check(isTest bool, path, src string, astFile *ast.File, fileSet *token.FileSet) {
	if isTest {
		TestFile.Check(path, src, astFile, fileSet)
	} else {
		File.Check(path, src, astFile, fileSet)
	}
}
