package sizes

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path"

	"github.com/bughou-go/spec/problems"
	"github.com/mattn/go-runewidth"
)

type ConfigT struct {
	Dir, File, Line, Func int
}

var Config = ConfigT{Dir: 20, File: 200, Line: 100, Func: 20}

func CheckDir(dir string) {
	if dir == `` || Config.Dir <= 0 {
		return
	}
	if count := entriesCount(dir); count > Config.Dir {
		problems.Add(
			token.Position{Filename: dir}, fmt.Sprintf(
				`dir %s has %d entries, limit %d`, path.Base(dir), count, Config.Dir,
			), `sizes.dir`,
		)
	}
}

func CheckFile(file *token.File) {
	if Config.File <= 0 {
		return
	}
	if lines := file.LineCount(); lines > Config.File {
		problems.Add(
			token.Position{Filename: file.Name()}, fmt.Sprintf(
				`file %s has %d lines, limit %d`, path.Base(file.Name()), lines, Config.File,
			), `sizes.file`,
		)
	}
}

func CheckLines(p string) {
	if Config.Line <= 0 {
		return
	}
	lines := readLines(p)
	for i, line := range lines {
		if width := runewidth.StringWidth(line); width > Config.Line {
			problems.Add(token.Position{Filename: p, Line: i}, fmt.Sprintf(
				`line %d width %d, limit %d`, i, width, Config.Line), `sizes.line`,
			)
		}
	}
}

func CheckFunc(funct ast.Node, file *token.File) {
	if Config.Func <= 0 {
		return
	}
	position := file.Position(funct.Pos())
	lines := file.Position(funct.End()).Line - position.Line
	if lines <= Config.Func {
		return
	}
	var name string
	if fun, ok := funct.(*ast.FuncDecl); ok {
		name = fun.Name.Name
	}
	problems.Add(position,
		fmt.Sprintf(`func %s has %d lines(max: %d)`, name, lines, Config.Func), `sizes.func`,
	)
}

func entriesCount(dir string) int {
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if names, err := f.Readdirnames(-1); err != nil {
		panic(err)
	} else {
		return len(names)
	}
}

func readLines(path string) (lines []string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return
}
