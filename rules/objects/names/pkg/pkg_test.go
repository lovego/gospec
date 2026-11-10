package pkgpkg

import (
	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleChecker_Check() {
	var src = `package Example
 `
	problems.Clear()
	w := walker.Parse("example.go", src)
	checker := NewChecker()
	checker.Check(w.AstFile.Name, w.FileSet)
	checker.Check(w.AstFile.Name, w.FileSet)
	problems.Render()
	// Output:
	// +----------------+-------------------------------------------------+-----------+
	// |    position    |                     problem                     |   rule    |
	// +----------------+-------------------------------------------------+-----------+
	// | example.go:1:9 | package name Example should be lower_case style | pkg.style |
	// +----------------+-------------------------------------------------+-----------+
}
