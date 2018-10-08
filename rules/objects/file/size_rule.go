package filepkg

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"github.com/lovego/gospec/problems"
	linesPkg "github.com/lovego/gospec/rules/objects/file/lines"
	"github.com/mattn/go-runewidth"
)

type sizeRule struct {
	MaxLines            uint `yaml:"maxLines"`
	MaxLineWidth        uint `yaml:"maxLineWidth"`
	MaxCommentLineWidth uint `yaml:"maxCommentLineWidth"`
}

func (r sizeRule) check(src string, astFile *ast.File, fileSet *token.FileSet, path, key string) {
	lines := linesPkg.New(src, astFile, fileSet)
	r.checkLines(uint(lines.Num()), path, key)
	r.checkLineWidth(lines, path, key)
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

func (r sizeRule) checkLineWidth(lines *linesPkg.Lines, path, key string) {
	for lineNum := 1; lineNum <= lines.Num(); lineNum++ {
		width := uint(runewidth.StringWidth(lines.Get(lineNum)))
		if width <= r.MaxLineWidth && width <= r.MaxCommentLineWidth {
			continue
		}
		if lines.IsComment(lineNum) {
			if width > r.MaxCommentLineWidth {
				problems.Add(
					token.Position{Filename: path, Line: lineNum},
					fmt.Sprintf(`line %d width: %d, limit: %d`, lineNum, width, r.MaxCommentLineWidth),
					key+`.maxCommentLineWidth`,
				)
			}
		} else {
			if width > r.MaxLineWidth {
				problems.Add(
					token.Position{Filename: path, Line: lineNum},
					fmt.Sprintf(`line %d width: %d, limit: %d`, lineNum, width, r.MaxLineWidth),
					key+`.maxLineWidth`,
				)
			}
		}
	}
}
