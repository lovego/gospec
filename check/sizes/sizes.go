package sizes

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/lovego/spec/problems"
	"github.com/mattn/go-runewidth"
)

type ConfigT struct {
	Dir, File, TestFile, Line, Func int
}

var Config = ConfigT{Dir: 20, File: 300, TestFile: 600, Line: 100, Func: 30}

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
	isTestFile := strings.HasSuffix(file.Name(), "_test.go")
	if isTestFile && count <= Config.TestFile || count <= Config.File {
		return
	}
	count -= commentsLineCount(nil, f, file, src)
	if isTestFile && count <= Config.TestFile || count <= Config.File {
		return
	}
	problems.Add(
		token.Position{Filename: file.Name()}, fmt.Sprintf(
			`file %s size: %d lines, limits %d`, path.Base(file.Name()), count, Config.File,
		), `sizes.file`,
	)
}

func CheckLines(p string, f *ast.File, file *token.File, src []string) {
	if Config.Line <= 0 {
		return
	}
	for i, line := range src {
		if width := runewidth.StringWidth(line); width > Config.Line && !isComment(i+1, f, file, src) {
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
		if !blankPrefix(src[start.Line-1], start) {
			start.Line += 1
		}
		if blankSuffix(src[end.Line-1], end) {
			end.Line += 1
		}
		count += end.Line - start.Line
	}
	return count
}

func isComment(lineNum int, f *ast.File, file *token.File, src []string) bool {
	for _, cg := range f.Comments {
		start, end := file.Position(cg.Pos()), file.Position(cg.End())
		if lineNum < start.Line || lineNum > end.Line {
			continue
		}
		if lineNum == start.Line && !blankPrefix(src[start.Line-1], start) ||
			lineNum == end.Line && !blankSuffix(src[end.Line-1], end) {
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
