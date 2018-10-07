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
	key:  "file",
	Name: name.Rule{MaxLen: 20, Style: "lower_case"},
	Size: sizeRule{MaxLineWidth: 100, MaxLines: 300},
}

var TestFileRule = RuleT{
	key:  "testFile",
	Name: name.Rule{MaxLen: 50, Style: "lower_case"},
	Size: sizeRule{MaxLineWidth: 100, MaxLines: 600},
}

type RuleT struct {
	key  string
	Name name.Rule
	Size sizeRule
}

type sizeRule struct {
	MaxLines     uint `yaml:"maxLines"`
	MaxLineWidth uint `yaml:"maxLineWidth"`
}

func Check(isTest bool, path, src string, astFile *ast.File, fileSet *token.FileSet) {
	if isTest {
		TestFileRule.Check(path, src, astFile, fileSet)
	} else {
		Rule.Check(path, src, astFile, fileSet)
	}
}

func (r *RuleT) Check(path, src string, astFile *ast.File, fileSet *token.FileSet) {
	filename := filepath.Base(path)
	r.checkName(filename, path)
	lines := scanLines(src)
	r.checkLines(uint(len(lines)), filename, path)
	r.checkLineWidth(lines, path, astFile, fileSet)
}

func (r *RuleT) checkName(filename, path string) {
	r.Name.Exec(
		strings.TrimSuffix(strings.TrimSuffix(filename, `.go`), `_test`),
		"file", r.key+".name", token.Position{Filename: path},
	)
}

func (r RuleT) checkLines(lineCount uint, filename, path string) {
	if lineCount <= r.Size.MaxLines {
		return
	}
	problems.Add(
		token.Position{Filename: path}, fmt.Sprintf(
			`file %s size: %d lines, limit: %d`, filename, lineCount, r.Size.MaxLines,
		), r.key+".size.maxLines",
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
