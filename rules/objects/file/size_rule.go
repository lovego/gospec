package filepkg

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/lovego/gospec/problems"
	"github.com/mattn/go-runewidth"
)

type sizeRule struct {
	MaxLines     uint `yaml:"maxLines"`
	MaxLineWidth uint `yaml:"maxLineWidth"`
}

func (r sizeRule) check(src, path, key string, astFile *ast.File, fileSet *token.FileSet) {
	lines := scanLines(src)
	r.checkLines(uint(len(lines)), path, key)
	r.checkLineWidth(lines, path, key, astFile, fileSet)
}

func (r sizeRule) checkLines(lineCount uint, path, key string) {
	if lineCount <= r.MaxLines {
		return
	}
	problems.Add(
		token.Position{Filename: path}, fmt.Sprintf(
			`file %s size: %d lines, limit: %d`, filepath.Base(path), lineCount, r.MaxLines,
		), key+".maxLines",
	)
}

func (r sizeRule) checkLineWidth(
	lines []string, path, key string, astFile *ast.File, fileSet *token.FileSet,
) {
	for i, line := range lines {
		lineNum := i + 1
		width := uint(runewidth.StringWidth(line))
		if width <= r.MaxLineWidth || isComment(lines, lineNum, astFile, fileSet) {
			continue
		}
		problems.Add(
			token.Position{Filename: path, Line: lineNum},
			fmt.Sprintf(`line %d width: %d, limit: %d`, lineNum, width, r.MaxLineWidth),
			key+`.maxLineWidth`,
		)
	}
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
