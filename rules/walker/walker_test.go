package walker

import (
	"fmt"
	"go/ast"
)

func ExampleNew() {
	New("walker_test.go")

	// Output:
}
func ExampleWalker_Walk() {
	w := Parse("example.go", `
package testdata

import "fmt"

func Hello() {
	fmt.Println("hello")
}
`)
	w.Walk(func(isLocal bool, node ast.Node) {
		if isLocal {
			fmt.Printf("local %T\n", node)
		} else {
			fmt.Printf("%T\n", node)
		}
	})

	// Output:
	// *ast.File
	// *ast.Ident
	// *ast.GenDecl
	// *ast.ImportSpec
	// *ast.BasicLit
	// *ast.FuncDecl
	// *ast.Ident
	// *ast.FuncType
	// *ast.FieldList
	// *ast.BlockStmt
	// local *ast.ExprStmt
	// local *ast.CallExpr
	// local *ast.SelectorExpr
	// local *ast.Ident
	// local *ast.Ident
	// local *ast.BasicLit
}
