package sizes

import (
	"go/ast"
	"go/token"
	"os"

	"github.com/bughou-go/spec/c"
)

func Check(dir *c.Dir) {
	checkDirSize(dir.Path)
	for _, pkg := range dir.Pkgs {
		for _, f := range pkg.Files {
			file := dir.Fset.File(f.Package)
			checkFileSize(file)
			checkLineSize(file)
			checkFuncSize(f, dir.Fset)
		}
	}
}

func checkDirSize(dir string) {
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
		c.Problem(token.Position{Filename: dir}, descs.Dir, `sizes.dir`)
	}
}

func checkFileSize(file *token.File) {
	if Config.File <= 0 {
		return
	}
	if file.LineCount() > Config.File {
		c.Problem(token.Position{Filename: file.Name()}, descs.File, `sizes.file`)
	}
}

func checkLineSize(file *token.File) {
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
			c.Problem(token.Position{Filename: file.Name(), Line: curLine}, descs.Line, `sizes.line`)
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

func checkFuncSize(f *ast.File, fset *token.FileSet) {
}
