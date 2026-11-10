package name

import (
	"fmt"
	"go/token"

	"github.com/lovego/gospec/problems"
)

var testPos = token.Position{Filename: "test.go", Line: 15, Column: 9}

func ExampleRule_empty() {
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

func ExampleRule_MaxLen_style() {
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

func ExampleRule_multiStyle() {
	problems.Clear()
	rule := Rule{MaxLen: 20, Style: "lower_case,lower-case"}
	rule.Exec(`dir_name`, `dir`, `dir.name`, testPos)
	rule.Exec(`dir-name`, `dir`, `dir.name`, testPos)
	Rule{MaxLen: 20, Style: "lower_case, lower-case"}.Exec(`dirName`, `dir`, `dir.name`, testPos)
	problems.Render()
	// Output:
	// +--------------+---------------------------------------------------------+----------------+
	// |   position   |                         problem                         |      rule      |
	// +--------------+---------------------------------------------------------+----------------+
	// | test.go:15:9 | dir name dirName should be lower_case, lower-case style | dir.name.style |
	// +--------------+---------------------------------------------------------+----------------+
}

func ExampleRule_unknownStyle() {
	check := func(style string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		Rule{MaxLen: 20, Style: style}.Exec(`var_name`, `var`, `var.name`, testPos)
	}
	check("camel-Case")
	check("lower_case,camel-Case")
	// Output:
	// unknown style config: "camel-Case".
	// unknown style config: "camel-Case".
}
