package check

import (
	"bufio"
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"

	"github.com/bughou-go/spec/check/names"
	"github.com/bughou-go/spec/check/sizes"
)

func Check(dir string, files []string) {
	names.CheckDir(dir)
	sizes.CheckDir(dir)
	var fset = token.NewFileSet()
	for i, p := range files {
		src, err1 := ioutil.ReadFile(p)
		if err1 != nil {
			panic(err1)
		}
		f, err2 := parser.ParseFile(fset, p, src, parser.ParseComments)
		if err2 != nil {
			panic(err2)
		}
		checkFile(f, fset.File(f.Package), i, scanLines(src))
	}
}

func checkFile(f *ast.File, file *token.File, i int, src []string) {
	if i == 0 {
		names.CheckIdent(f.Name, false, file, ``)
	}
	names.CheckFile(file.Name())
	sizes.CheckFile(f, file, src)

	sizes.CheckLines(file.Name(), f, file, src)
	ast.Walk(walker{f, file, src}, f)
}

func scanLines(src []byte) (lines []string) {
	scanner := bufio.NewScanner(bytes.NewReader(src))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return
}
