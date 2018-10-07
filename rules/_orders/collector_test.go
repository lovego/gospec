package orders

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"
)

type testWalker struct {
	checker *FileChecker
}

const (
	noValue = iota
)

var testChecker = createChecker()

func (tw testWalker) Visit(node ast.Node) ast.Visitor {
	tw.checker.Visit(node)
	return tw
}

func createChecker() *FileChecker {
	fileName := "collector_test.go"
	fileContent, _ := ioutil.ReadFile(fileName)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, fileContent,
		parser.ParseComments)
	if err != nil {
		panic(err)
	}

	checker := NewFileChecker(fset.File(file.Package))
	walker := testWalker{checker: checker}
	ast.Walk(walker, file)
	return checker
}

func TestCollectImported(t *testing.T) {
	expected := []string{
		"ast",
		"parser",
		"token",
		"ioutil",
		"testing",
	}
	checkCollectorResult(expected, testChecker.sc.importedTable, t)
}

func checkCollectorResult(expected, result []string, t *testing.T) {
	if len(expected) != len(result) {
		t.Fatalf("expected:%v got:%v\n", expected, result)
	}
	for i, p := range expected {
		if p != result[i] {
			t.Fatalf("expected:%s get:%s\n", p, result[i])
		}
	}
}

func TestCollectDecl(t *testing.T) {
	var expectedDecl = []string{
		"testWalker",
		"noValue",
		"testChecker",
	}
	checkCollectorResult(expectedDecl, testChecker.sc.getDeclNames(), t)
	var expectedFuncDecl = []string{
		"Visit",
		"createChecker",
		"TestCollectImported",
		"checkCollectorResult",
		"TestCollectDecl",
		"TestCollectCallee",
	}
	checkCollectorResult(expectedFuncDecl, testChecker.sc.getFuncNames(), t)
}

func TestCollectCallee(t *testing.T) {
	callees := map[string][]string{
		"createChecker":   nil,
		"TestCollectDecl": []string{":checkCollectorResult"},
	}
	for fn, names := range callees {
		checkCollectorResult(names, testChecker.sc.getCallees(fn), t)
	}
}
