package funcpkg

import (
	"go/ast"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleSizeRule_checkType() {
	var src = `package example
func F(a, b int, c bool) (a int, b bool) {
}

`
	problems.Clear()
	w := walker.Parse("example.go", src)
	w.Walk(func(isLocal bool, node ast.Node) {
		if typ, ok := node.(*ast.FuncType); ok {
			sizeRule{MaxParams: 2, MaxResults: 1}.checkType(typ, "func F", "func", w.FileSet)
		}
	})
	problems.Render()
	// Output:
	// +-----------------+----------------------------------+-----------------+
	// |    position     |             problem              |      rule       |
	// +-----------------+----------------------------------+-----------------+
	// | example.go:2:7  | func F params size: 3, limit: 2  | func.maxParams  |
	// | example.go:2:26 | func F results size: 2, limit: 1 | func.maxResults |
	// +-----------------+----------------------------------+-----------------+

}

func ExampleSizeRule_checkBody() {
	var src = `package example
func F() {
  type  t  int64  // 1
  const a = 4     // 2
  var   b = 5     // 3
label:            // 4
  if c := 3; c > 0  { // 5, 6 (if, assign)
    c = a + b     // 7
  }
}
`

	problems.Clear()
	w := walker.Parse("example.go", src)
	w.Walk(func(isLocal bool, node ast.Node) {
		if body, ok := node.(*ast.BlockStmt); ok {
			sizeRule{MaxStatements: 6}.checkBody(body, "func F", "func", w.FileSet)
		}
	})
	problems.Render()
	// Output:
	// +-----------------+------------------------------------------+--------------------+
	// |    position     |                 problem                  |        rule        |
	// +-----------------+------------------------------------------+--------------------+
	// | example.go:2:10 | func F body size: 7 statements, limit: 6 | func.maxStatements |
	// +-----------------+------------------------------------------+--------------------+
}
