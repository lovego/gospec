package dirpkg

import (
	"github.com/lovego/gospec/problems"
)

func ExampleCheck() {
	problems.Clear()
	Check("")
	Check("../dir")
	Check(".")
	Check("..")
	problems.Render()
	// Output:
}
