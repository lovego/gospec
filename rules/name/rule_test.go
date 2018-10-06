package name

import (
	"fmt"
	"go/token"

	"github.com/lovego/gospec/problems"
)

var testPos = token.Position{Filename: "test.go", Line: 15, Column: 9}

func ExampleRule_Empty() {
	problems.Clear()
	Rule{}.Exec(`_`, `var`, `var.name`, testPos)
	Rule{MaxLen: 20, Style: "lower-case"}.Exec(`dir-name`, `dir`, `dir.name`, testPos)
	Rule{MaxLen: 20, Style: "lowerCamelCase"}.Exec(`varName`, `var`, `var.name`, testPos)
	problems.Render()
	// Output:
}

func ExampleRule_MaxLen() {
	problems.Clear()
	Rule{MaxLen: 6}.Exec(`varName`, `var`, `var.name`, testPos)
	problems.Render()
	// Output:
	// +--------------+-----------------------------------------+-----------------+
	// |   position   |                 problem                 |      rule       |
	// +--------------+-----------------------------------------+-----------------+
	// | test.go:15:9 | var name varName 7 chars long, limit: 6 | var.name.maxLen |
	// +--------------+-----------------------------------------+-----------------+
}

func ExampleRule_Style() {
	problems.Clear()
	Rule{MaxLen: 20, Style: "lower_case"}.Exec(`varName`, `var`, `var.name`, testPos)
	problems.Render()
	// Output:
	// +--------------+---------------------------------------------+----------------+
	// |   position   |                   problem                   |      rule      |
	// +--------------+---------------------------------------------+----------------+
	// | test.go:15:9 | var name varName should be lower_case style | var.name.style |
	// +--------------+---------------------------------------------+----------------+
}

func ExampleRule_MaxLen_Style() {
	problems.Clear()
	Rule{MaxLen: 7, Style: "camelCase"}.Exec(`var_name`, `var`, `var.name`, testPos)
	problems.Render()
	// Output:
	// +--------------+------------------------------------------------------------------------+--------------------------+
	// |   position   |                                problem                                 |           rule           |
	// +--------------+------------------------------------------------------------------------+--------------------------+
	// | test.go:15:9 | var name var_name 8 chars long, limit: 7 and should be camelCase style | var.name.{maxLen, style} |
	// +--------------+------------------------------------------------------------------------+--------------------------+
}

func ExampleRule_UnknownStyle() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	problems.Clear()
	Rule{MaxLen: 20, Style: "camel-Case"}.Exec(`var_name`, `var`, `var.name`, testPos)
	problems.Render()
	// Output: unknown style config: "camel-Case".
}
