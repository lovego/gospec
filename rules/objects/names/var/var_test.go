package varpkg

import (
	"go/ast"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleCheck() {
	var src = `package example
var V =  1
func F(arg1, Arg2 int) {
  const c =  2
  var V int
  slice := []int{}
  for k, v := range slice {
    k = 5
  }
  for k, v = range slice {
  }
}
`
	problems.Clear()
	w := walker.Parse("example.go", src)
	w.Walk(func(isLocal bool, node ast.Node) {
		Check(isLocal, node, w.FileSet)
	})
	problems.Render()

	// Output:
	// +-----------------+-----------------------------------------------------+----------------+
	// |    position     |                       problem                       |      rule      |
	// +-----------------+-----------------------------------------------------+----------------+
	// | example.go:3:14 | func param name Arg2 should be lowerCamelCase style | localVar.style |
	// | example.go:5:7  | local var name V should be lowerCamelCase style     | localVar.style |
	// +-----------------+-----------------------------------------------------+----------------+
}
