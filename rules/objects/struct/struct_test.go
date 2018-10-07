package structpkg

import (
	"go/ast"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleCheck() {
	var src = `package example
type T struct {
  Name string
  A, B, C int
}
`
	problems.Clear()
	w := walker.Parse("example.go", src)
	w.Walk(func(isLocal bool, node ast.Node) {
		Check(node, w.FileSet)
	})
	problems.Render()
	// Output:
}
