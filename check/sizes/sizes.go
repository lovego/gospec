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
	lastLine := 1
	for pos <= end {
		var line int
		line, pos = forwardAMaxLine(file, pos)
		if line == lastLine {
			c.Problem(token.Position{Filename: file.Name(), Line: line}, descs.Line, `sizes.line`)
		}
		lastLine = line
	}
}

// move forward a max line
func forwardAMaxLine(file *token.File, pos int) (int, int) {
	pos += Config.Line + 1 // Config.Line doesn't contain newline, so plus 1.
	position := file.Position(token.Pos(pos))
	pos -= position.Column - 1 // move to first colunn
	return position.Line, pos
}

func checkFuncSize(f *ast.File, fset *token.FileSet) {
}
