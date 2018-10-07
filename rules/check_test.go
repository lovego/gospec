package rules

import (
	"github.com/lovego/gospec/problems"
)

func ExampleCheck() {
	problems.Clear()
	Check(".", []string{"check_test.go"})
	problems.Render()
	// Output:
}
