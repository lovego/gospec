package sizes

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path"
	"strings"

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
	count := entriesCount(dir)
	if count <= Config.Dir {
		return
	}
	problems.Add(
		token.Position{Filename: dir}, fmt.Sprintf(
			`dir %s size: %d entries, limits %d`, path.Base(dir), count, Config.Dir,
		), `sizes.dir`,
	)
}

func CheckFile(f *ast.File, file *token.File, src []string) {
	if Config.File <= 0 {
		return
	}
	count := file.LineCount()
	if count <= Config.File {
		return
	}
	if count -= commentsLineCount(nil, f, file, src); count <= Config.File {
		return
	}
	problems.Add(
		token.Position{Filename: file.Name()}, fmt.Sprintf(
			`file %s size: %d lines, limits %d`, path.Base(file.Name()), count, Config.File,
		), `sizes.file`,
	)
}

func CheckLines(p string, src []string) {
	if Config.Line <= 0 {
		return
	}
	for i, line := range src {
		if width := runewidth.StringWidth(line); width > Config.Line {
			problems.Add(token.Position{Filename: p, Line: i + 1}, fmt.Sprintf(
				`line %d size: %d chars wide, limits %d`, i+1, width, Config.Line), `sizes.line`,
			)
		}
	}
}

func CheckFunc(fun ast.Node, f *ast.File, file *token.File, src []string) {
	if Config.Func <= 0 {
		return
	}
	w := &stmtWalker{}
	ast.Walk(w, fun)
	if w.count <= Config.Func {
		return
	}
	var name string
	if funct, ok := fun.(*ast.FuncDecl); ok {
		name = funct.Name.Name
	}
	problems.Add(file.Position(fun.Pos()),
		fmt.Sprintf(`func %s size: %d statements, limits %d`, name, w.count, Config.Func), `sizes.func`,
	)
}

type stmtWalker struct {
	count int
}

func (w *stmtWalker) Visit(node ast.Node) ast.Visitor {
	if stmt, ok := node.(ast.Stmt); ok {
		if _, ok := stmt.(*ast.BlockStmt); !ok {
			// fmt.Printf("%T\n", stmt)
			w.count++
		}
	}
	return w
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

func commentsLineCount(context ast.Node, f *ast.File, file *token.File, src []string) (count int) {
	for _, cg := range f.Comments {
		if context != nil && (cg.Pos() < context.Pos() || cg.End() > context.End()) {
			continue
		}
		start, end := file.Position(cg.Pos()), file.Position(cg.End())
		// non blank prefix
		if line := src[start.Line-1]; !isWhitespace(line[:start.Column-1]) {
			start.Line += 1
		}
		// blank suffix
		if line := src[end.Line-1]; isWhitespace(line[end.Column-1:]) {
			end.Line += 1
		}
		count += end.Line - start.Line
	}
	return count
}

func isWhitespace(s string) bool {
	return len(s) == 0 || len(strings.TrimSpace(s)) == 0
}
