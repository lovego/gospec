package rules

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func ExampleConfig_default() {
	out, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(out))

	// Output:
	// dir:
	//   name:
	//     maxLen: 20
	//     style: lower_case
	//   size:
	//     maxEntries: 20
	// pkg:
	//   maxLen: 20
	//   style: lower_case
	// file:
	//   name:
	//     maxLen: 20
	//     style: lower_case
	//   size:
	//     maxLines: 300
	//     maxLineWidth: 120
	// testfile:
	//   name:
	//     maxLen: 20
	//     style: lower_case
	//   size:
	//     maxLines: 300
	//     maxLineWidth: 120
	// func:
	//   name:
	//     maxLen: 30
	//     style: camelCase
	//   size:
	//     maxParams: 5
	//     maxResults: 3
	//     maxStatements: 30
	// funcInTest:
	//   name:
	//     maxLen: 50
	//     style: camelCase
	//   size:
	//     maxParams: 5
	//     maxResults: 3
	//     maxStatements: 30
	// struct:
	//   name:
	//     maxLen: 30
	//     style: camelCase
	//   size:
	//     maxFields: 100
	// const:
	//   maxLen: 30
	//   style: camelCase
	// localConst:
	//   maxLen: 20
	//   style: lowerCamelCase
	// var:
	//   maxLen: 40
	//   style: camelCase
	// localVar:
	//   maxLen: 30
	//   style: lowerCamelCase
	// type:
	//   maxLen: 30
	//   style: camelCase
	// localType:
	//   maxLen: 20
	//   style: lowerCamelCase
	// label:
	//   maxLen: 30
	//   style: lowerCamelCase
}
