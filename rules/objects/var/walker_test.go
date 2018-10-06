package names

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"

	"bou.ke/monkey"
)

func patchRuleExec() {
	var r *rule
	monkey.PatchInstanceMethod(reflect.TypeOf(r), "Exec",
		func(r *rule, thing string, name string, pos token.Position) {
			if thing == "" {
				thing = r.Desc
			}
			fmt.Printf("%-15s %-18s %-15s %v\n", r.Key, thing, name, pos)
		},
	)
}

func ExampleWalker() {
	patchRuleExec()

	var src = `package example

const constant = 3
var  variable int
type Type struct {
  field func(argument int) (result int)
}

func function() {
  const localConstant = 4
  var   localVariable = 5
  type  localType  int64
label:
  localVariable = localConstant + constant
}

func (receiver *Type) method() {
}

type Interface interface {
  method()
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, `example.go`, src, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast.Walk(walker{fileSet: fset}, file)

	// Output:
	// const           const              constant        example.go:3:7
	// var             var                variable        example.go:4:6
	// type            type               Type            example.go:5:6
	// structField     struct field       field           example.go:6:3
	// localVar        func param         argument        example.go:6:14
	// localVar        func result        result          example.go:6:29
	// func            func               function        example.go:9:6
	// localConst      local const        localConstant   example.go:10:9
	// localVar        local var          localVariable   example.go:11:9
	// localType       local type         localType       example.go:12:9
	// localVar        func receiver      receiver        example.go:17:7
	// func            method             method          example.go:17:23
	// type            type               Interface       example.go:20:6
	// func            interface method   method          example.go:21:3
}

func ExampleWalker_test() {
	patchRuleExec()

	var src = `package example

func function() {
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, `example_test.go`, src, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast.Walk(walker{fileSet: fset}, file)

	// Output:
	// funcInTest      func               function        example_test.go:3:6
}
