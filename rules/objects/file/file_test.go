package filepkg

import (
	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleCheck() {
	problems.Clear()
	w := walker.New("file.go")
	Check(false, "file.go", w.SrcFile, w.AstFile, w.FileSet)

	w = walker.New("file_test.go")
	Check(true, "file_test.go", w.SrcFile, w.AstFile, w.FileSet)

	problems.Render()
	// Output:
}
