package sizes

import (
	"fmt"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/lovego/gospec/problems"
)

func checkFile(filename string, lines []string) {
	limit := Rules.File

	isTest := strings.HasSuffix(filename, "_test.go")
	if isTest {
		limit = Rules.TestFile
	}

	if len(lines) <= limit {
		return
	}

	var rule = "sizes.file"
	if isTest {
		rule = "sizes.testFile"
	}

	problems.Add(
		token.Position{Filename: filename},
		fmt.Sprintf(
			`file %s size: %d statements, limit: %d`, filepath.Base(filename), len(lines), limit,
		),
		rule,
	)
}
