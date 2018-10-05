package sizes

import (
	"bufio"
	"go/ast"
	"go/token"
	"strings"
)

var Rules = RulesT{
	Dir: 20, File: 300, TestFile: 600, Line: 100,
	Func: 30, Params: 5, Results: 3,
}

type RulesT struct {
	Dir, File, TestFile, Line int
	Func, Params, Results     int
}

func Check(dir string, astFiles []*ast.File, sources []string, fileSet *token.FileSet) {
	checkDir(dir)
	w := walker{fileSet: fileSet}
	for i, astFile := range astFiles {
		filename := fileSet.Position(astFile.Pos()).Filename
		lines := scanLines(sources[i])
		checkFile(filename, lines)
		checkLines(filename, lines, astFile, fileSet)
		ast.Walk(w, astFile)
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
