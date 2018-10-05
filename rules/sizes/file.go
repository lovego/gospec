package sizes

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/lovego/gospec/problems"
)

func checkFile(astFile *ast.File, fileSet *token.FileSet) {
	limit := Rules.File
	filename := fileSet.Position(astFile.Pos()).Filename

	isTest := strings.HasSuffix(filename, "_test.go")
	if isTest {
		limit = Rules.TestFile
	}

	num := stmtsNum(astFile)
	if num <= limit {
		return
	}

	var rule = "sizes.file"
	if isTest {
		rule = "sizes.testFile"
	}

	problems.Add(
		token.Position{Filename: filename},
		fmt.Sprintf(`file %s size: %d statements, limit: %d`, filepath.Base(filename), num, limit),
		rule,
	)
}
