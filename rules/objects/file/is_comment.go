package filepkg

import (
	"go/ast"
	"go/token"
	"strings"
)

func isComment(lines []string, lineNum int, astFile *ast.File, fileSet *token.FileSet) bool {
	for _, commentGroup := range astFile.Comments {
		start, end := fileSet.Position(commentGroup.Pos()), fileSet.Position(commentGroup.End())
		if lineNum < start.Line || lineNum > end.Line {
			continue
		}
		if lineNum == start.Line && !blankPrefix(lines[start.Line-1], start) ||
			lineNum == end.Line && !blankSuffix(lines[end.Line-1], end) {
			return false
		}
		return true
	}
	return false
}

func blankPrefix(line string, start token.Position) bool {
	if start.Column <= 1 {
		return true
	}
	suffix := line[:start.Column-1]
	return len(strings.TrimSpace(suffix)) == 0
}

func blankSuffix(line string, end token.Position) bool {
	// end is the next pos after ending token
	if end.Column > len(line) {
		return true
	}
	suffix := line[end.Column-1:]
	return len(strings.TrimSpace(suffix)) == 0
}
