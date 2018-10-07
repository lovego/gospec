package typepkg

import (
	"go/ast"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleCheck() {
	var src = `package example
type T int
func F() {
  type TypeName int
}
`
	problems.Clear()
	w := walker.Parse("example.go", src)
	w.Walk(func(isLocal bool, node ast.Node) {
		Check(isLocal, node, w.FileSet)
	})
	problems.Render()

	// Output:
	// +----------------+---------------------------------------------------------+-----------------+
	// |    position    |                         problem                         |      rule       |
	// +----------------+---------------------------------------------------------+-----------------+
	// | example.go:4:8 | local type name TypeName should be lowerCamelCase style | localType.style |
	// +----------------+---------------------------------------------------------+-----------------+
}
