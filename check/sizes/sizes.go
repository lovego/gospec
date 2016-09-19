package sizes

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path"

	"github.com/bughou-go/spec/problems"
)

type ConfigT struct {
	Dir, File, Line, Func int
}

var Config = ConfigT{Dir: 20, File: 200, Line: 100, Func: 20}

func CheckDir(dir string) {
	if dir == `` || Config.Dir <= 0 {
		return
	}
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	names, err := f.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	if len(names) > Config.Dir {
		problems.Add(
			token.Position{Filename: dir},
			fmt.Sprintf(`dir %s shouldn't be more than %d items`, path.Base(dir), Config.Dir),
			`sizes.dir`,
		)
	}
}

func CheckFile(file *token.File) {
	if Config.File <= 0 {
		return
	}
	if file.LineCount() > Config.File {
		problems.Add(
			token.Position{Filename: file.Name()},
			fmt.Sprintf(`file %s shouldn't be more than %d lines`, path.Base(file.Name()), Config.File),
			`sizes.file`,
		)
	}
}

func CheckLines(file *token.File) {
	if Config.Line <= 0 {
		return
	}
	pos := file.Base()
	end := file.Base() + file.Size()
	curLine := 1
	for pos <= end {
		// move forward maxLine + 1, if it stay on the same line, then it's too long
		pos += Config.Line + 1
		if pos > end {
			break
		}
		position := file.Position(token.Pos(pos))
		if position.Line == curLine {
			problems.Add(
				token.Position{Filename: file.Name(), Line: curLine},
				fmt.Sprintf(`line %d shouldn't be more than %d chars`, curLine, Config.Line),
				`sizes.line`,
			)
			pos, curLine = forward2NewLine(file, pos)
		} else {
			pos -= position.Column - 1 // move backward to first column
			curLine = position.Line
		}
	}
}

func forward2NewLine(file *token.File, pos int) (int, int) {
	end := file.Base() + file.Size()
	if pos > end {
		return pos, -1
	}
	position := file.Position(token.Pos(pos))
	line := position.Line
	for curLine := line; line == curLine; line = position.Line {
		// it's safe to forward maxLine + 2. it won't skip lines that's too long.
		pos += Config.Line + 2
		if pos > end {
			return pos, -1
		}
		position = file.Position(token.Pos(pos))
	}
	pos -= position.Column - 1 // move backward to first column
	return pos, line
}

func CheckFunc(funct ast.Node, file *token.File) {
	if Config.Func <= 0 {
		return
	}
	position := file.Position(funct.Pos())
	if file.Position(funct.End()).Line-position.Line <= Config.Func {
		return
	}
	var name string
	if fun, ok := funct.(*ast.FuncDecl); ok {
		name = fun.Name.Name
	}
	problems.Add(position,
		fmt.Sprintf(`func %s shouldn't be more than %d lines`, name, Config.Func),
		`sizes.func`,
	)
}
