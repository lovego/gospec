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

var Config = ConfigT{Dir: 20, File: 300, TestFile: 600, Line: 100, Func: 30}

type ConfigT struct {
	Dir, File, TestFile, Line, Func int
}

type SizeWalker struct {
	srcFile *token.File
	funcs   []funcNode
	counts  []*int
	index   int
}

type funcNode struct {
	node  ast.Node
	begin int
	end   int
}

func NewWalker(srcFile *token.File) *SizeWalker {
	return &SizeWalker{
		srcFile: srcFile,
		index:   -1,
	}
}

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

func (sw *SizeWalker) getIndex(node ast.Node) int {
	pos := sw.srcFile.Position(node.Pos())
	for i := len(sw.funcs) - 1; i >= 0; i-- {
		fnNode := sw.funcs[i]
		if pos.Line > fnNode.begin && pos.Line < fnNode.end {
			return i
		}
	}
	return 0
}

func (sw *SizeWalker) addFunc(node ast.Node, body *ast.BlockStmt) {
	sw.index++
	begin := sw.srcFile.Position(body.Lbrace).Line
	end := sw.srcFile.Position(body.Rbrace).Line
	sw.funcs = append(sw.funcs, funcNode{node: node, begin: begin, end: end})
	sw.counts = append(sw.counts, new(int))
}

func (sw *SizeWalker) Visit(node ast.Node) ast.Node {
	switch node.(type) {
	case ast.Stmt:
		if _, ok := node.(*ast.BlockStmt); !ok {
			index := sw.getIndex(node)
			count := sw.counts[index]
			*count++
		}
	case *ast.FuncLit:
		v := node.(*ast.FuncLit)
		sw.addFunc(node, v.Body)
	case *ast.FuncDecl:
		v := node.(*ast.FuncDecl)
		sw.addFunc(node, v.Body)
	}
	return node
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

func (sw *SizeWalker) CheckFuncs() {
	if Config.Func <= 0 {
		return
	}
	for i, fn := range sw.funcs {
		count := sw.counts[i]
		if *count <= Config.Func {
			continue
		}
		var name string
		if funct, ok := fn.node.(*ast.FuncDecl); ok {
			name = funct.Name.Name
		}
		problems.Add(sw.srcFile.Position(fn.node.Pos()),
			fmt.Sprintf(`func %s size: %d statements, limits %d`, name, *count, Config.Func), `sizes.func`,
		)
	}
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
