package structpkg

import (
	"go/ast"

	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleRule_Check() {
	var src = `package example
type T struct {
  Name string
  A, B, C int
}
`
	problems.Clear()
	structCopy := Struct
	structCopy.Size.MaxFields = 3
	w := walker.Parse("example.go", src)
	w.Walk(func(isLocal bool, node ast.Node) {
		switch strut := node.(type) {
		case *ast.StructType:
			structCopy.check(strut, w.FileSet)
		}
	})
	problems.Render()
	// Output:
	// +----------------+---------------------------------+-----------------------+
	// |    position    |             problem             |         rule          |
	// +----------------+---------------------------------+-----------------------+
	// | example.go:2:8 | struct size: 4 fields, limit: 3 | struct.size.maxFields |
	// +----------------+---------------------------------+-----------------------+
}
