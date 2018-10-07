package constpkg

import (
	"go/ast"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleCheck() {
	var src = `package example
const C =  3
func F() {
  const c, ConstName =  3, 4
  var   v =  5
}
`
	problems.Clear()
	w := walker.Parse("example.go", src)
	w.Walk(func(isLocal bool, node ast.Node) {
		Check(isLocal, node, w.FileSet)
	})
	problems.Render()
	// Output:
	// +-----------------+-----------------------------------------------------------+------------------+
	// |    position     |                          problem                          |       rule       |
	// +-----------------+-----------------------------------------------------------+------------------+
	// | example.go:4:12 | local const name ConstName should be lowerCamelCase style | localConst.style |
	// +-----------------+-----------------------------------------------------------+------------------+
}
