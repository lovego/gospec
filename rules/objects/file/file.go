package filepkg

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/name"
)

var Rule = RuleT{
	Name: name.Rule{MaxLen: 20, Style: "lower_case"},
	Size: sizeRule{MaxLine: 120, MaxLines: 300, TestMaxLines: 600},
}

type RuleT struct {
	Name name.Rule
	Size sizeRule
}

type sizeRule struct {
	MaxLine      uint `yaml:"maxLine"`
	MaxLines     uint `yaml:"maxLines"`
	TestMaxLines uint `yaml:"testMaxLines"`
}

func Check(path string, src string, isTest bool, astFile *ast.File, fileSet *token.FileSet) {
	name := filepath.Base(path)
	checkName(name, path)

	lines := scanLines(src)
	checkSize(name, path, isTest, uint(len(lines)))
	checkLineSize(path, lines, astFile, fileSet)
}

func checkName(name, path string) {
	Rule.Name.Exec(
		strings.TrimSuffix(strings.TrimSuffix(name, `.go`), `_test`),
		"file", "file.name", token.Position{Filename: path},
	)
}

func checkSize(name, path string, isTest bool, lineCount uint) {
	limit := Rule.Size.MaxLines
	if isTest {
		limit = Rule.Size.TestMaxLines
	}

	if lineCount <= limit {
		return
	}

	var rule = "file.size.maxLines"
	if isTest {
		rule = "file.size.testMaxLines"
	}

	problems.Add(
		token.Position{Filename: path},
		fmt.Sprintf(
			`file %s size: %d lines, limit: %d`, name, lineCount, limit,
		),
		rule,
	)
}

func scanLines(src string) (lines []string) {
	scanner := bufio.NewScanner(strings.NewReader(src))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return
}
