package lines

import (
	"bufio"
	"go/ast"
	"go/token"
	"strings"
)

type Lines struct {
	srcLines []string // source lines with line termination stripped
	astFile  *ast.File
	fileSet  *token.FileSet
	comments map[int]bool
}

func New(src string, astFile *ast.File, fileSet *token.FileSet) *Lines {
	return &Lines{srcLines: scanLines(src), astFile: astFile, fileSet: fileSet}
}

func (l *Lines) Num() int {
	return len(l.srcLines)
}

func (l *Lines) Get(lineNum int) string {
	return l.srcLines[lineNum-1]
}

func (l *Lines) IsComment(lineNum int) bool {
	if l.comments == nil {
		l.initComments()
	}
	return l.comments[lineNum]
}

func (l *Lines) initComments() {
	l.comments = make(map[int]bool)
	for _, commentGroup := range l.astFile.Comments {
		start, end := l.fileSet.Position(commentGroup.Pos()), l.fileSet.Position(commentGroup.End())

		if start.Line == end.Line {
			if l.blankPrefix(start) && l.blankSuffix(end) {
				l.comments[start.Line] = true
			}
		} else {
			if l.blankPrefix(start) {
				l.comments[start.Line] = true
			}
			for lineNum := start.Line + 1; lineNum < end.Line; lineNum++ {
				l.comments[lineNum] = true
			}

			// end is the next char after ending token
			// so even it's line termination char, it's always on the same line
			if l.blankSuffix(end) {
				l.comments[end.Line] = true
			}
		}
	}
}

func (l *Lines) blankPrefix(position token.Position) bool {
	if position.Column <= 1 { // position is the first char of the line
		return true
	}
	prefix := l.srcLines[position.Line-1][:position.Column-1]
	return len(strings.TrimSpace(prefix)) == 0
}

func (l *Lines) blankSuffix(position token.Position) bool {
	line := l.srcLines[position.Line-1]
	if position.Column > len(line) { // position is the  termination char of the line
		return true
	}
	suffix := line[position.Column-1:]
	return len(strings.TrimSpace(suffix)) == 0
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
