package funcpkg

import (
	"go/ast"

	"github.com/lovego/gospec/problems"
	namepkg "github.com/lovego/gospec/rules/name"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleRule_Check() {
	var src = `package example
func (receiver *Type) F1(a, b int, c bool)  () {
}

var f = func(a, b int, c bool) {
}


type I1 interface {
  M1()
  M2()(a int, b bool)
  I2
}

type I2 interface{}
`
	problems.Clear()
	rule := Rule{
		key: "func",
		Name: namepkg.Rule{
			MaxLen: 30, Style: "camelCase",
		},
		Size: sizeRule{
			MaxParams: 2, MaxResults: 1, MaxStatements: 30,
		},
	}

	w := walker.Parse("example.go", src)
	w.Walk(func(isLocal bool, node ast.Node) {
		rule.Check(node, w.FileSet)
	})
	problems.Render()
	// Output:
	// +-----------------+-----------------------------------------------+----------------------+
	// |    position     |                    problem                    |         rule         |
	// +-----------------+-----------------------------------------------+----------------------+
	// | example.go:2:25 | func F1 params size: 3, limit: 2              | func.size.maxParams  |
	// | example.go:5:13 | literal func params size: 3, limit: 2         | func.size.maxParams  |
	// | example.go:11:7 | interface method M2 results size: 2, limit: 1 | func.size.maxResults |
	// +-----------------+-----------------------------------------------+----------------------+
}

func ExampleFuncInTest() {
	var src = `package example
func ExampleA_too_long_too_long_too_long_too_long_too_long () {
}

`
	problems.Clear()

	w := walker.Parse("example_test.go", src)
	w.Walk(func(isLocal bool, node ast.Node) {
		FuncInTest.Check(node, w.FileSet)
	})
	problems.Render()
	// Output:
	// +---------------------+------------------------------------------------------------------------------------------+------------------------+
	// |      position       |                                         problem                                          |          rule          |
	// +---------------------+------------------------------------------------------------------------------------------+------------------------+
	// | example_test.go:2:6 | func name ExampleA_too_long_too_long_too_long_too_long_too_long 53 chars long, limit: 50 | funcInTest.name.maxLen |
	// +---------------------+------------------------------------------------------------------------------------------+------------------------+
}
