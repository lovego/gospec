package sizes

import (
	"go/ast"
	"go/token"
)

var Rules = RulesT{
	Dir: 20, File: 300, TestFile: 600, Line: 100,
	Func: 30, Params: 5, Results: 3,
}

type RulesT struct {
	Dir, File, TestFile, Line, Func, Params, Results int
}

func Check(dir string, astFiles []*ast.File, sources []string, fileSet *token.FileSet) {
	checkDir(dir)
	w := walker{fileSet: fileSet}
	for i, astFile := range astFiles {
		checkFile(astFile, fileSet)
		checkLines(sources[i], astFile, fileSet)
		ast.Walk(w, astFile)
	}
}
