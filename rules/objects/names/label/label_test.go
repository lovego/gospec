package labelpkg

import (
	"go/ast"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleCheck() {
	var src = `package example
func F() {
  Label:
}
`
	problems.Clear()
	w := walker.Parse("example.go", src)
	w.Walk(func(isLocal bool, node ast.Node) {
		Check(node, w.FileSet)
	})
	problems.Render()
	// Output:
	// +----------------+-------------------------------------------------+-------------+
	// |    position    |                     problem                     |    rule     |
	// +----------------+-------------------------------------------------+-------------+
	// | example.go:3:3 | label name Label should be lowerCamelCase style | label.style |
	// +----------------+-------------------------------------------------+-------------+
}
