package funcpkg

import (
	"go/ast"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleCheck() {
	problems.Clear()
	w := walker.New("func.go")
	w.Walk(func(isLocal bool, node ast.Node) {
		Check(false, node, w.FileSet)
	})

	w = walker.New("func_test.go")
	w.Walk(func(isLocal bool, node ast.Node) {
		Check(true, node, w.FileSet)
	})

	problems.Render()
	// Output:
}
