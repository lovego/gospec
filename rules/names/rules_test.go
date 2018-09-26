package names

import (
	"fmt"
	"reflect"
)

func ExampleRules_init() {
	value := reflect.ValueOf(&Rules).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.FieldByIndex([]int{i})
		fmt.Printf("%-11v %v\n", field.FieldByName("Key"), field.FieldByName("Desc"))
	}
	// Output:
	// dir         dir
	// file        file
	// pkg         pkg
	// func        func
	// funcInTest  func
	// label       label
	// const       const
	// var         var
	// type        type
	// localConst  local const
	// localVar    local var
	// localType   local type
	// structField struct field
}

func ExampleLowercaseFirstChar() {
	fmt.Println(lowercaseFirstChar("A"))
	fmt.Println(lowercaseFirstChar("Abc"))
	fmt.Println(lowercaseFirstChar("AbcDef"))
	// Output:
	// a
	// abc
	// abcDef
}
