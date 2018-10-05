package sizes

import (
	"fmt"
	"go/parser"
	"go/token"
)

func ExampleStmtsNum() {
	var src = `package example

func function() {
  type  t  int64  // 1
  const a = 4     // 2
  var   b = 5     // 3
label:            // 4
  if c := 3; c > 0  { // 5, 6 (if, assign)
    c = a + b     // 7
  }
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, `example.go`, src, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(stmtsNum(file))

	// Output: 7
}
