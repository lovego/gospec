package main

import (
	"fmt"
	"go/ast"
	"go/parser"
)

type walker struct {
}

func (w walker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Ident:
		fmt.Printf("%+v %+v %#v\n", *n, *n.Obj, n.Obj.Decl)
	}
	return w
}

func ExampleIdent() {
	src := `
  struct {
    A string
    B struct {
      C int64
    }
  }
  `
	expr, err := parser.ParseExpr(src)
	if err != nil {
		fmt.Println(err)
	}

	ast.Walk(walker{}, expr)

	// Output:
}
