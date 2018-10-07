package filepkg

import (
	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleSizeRule_check() {
	problems.Clear()
	w := walker.New("size_rule_test.go")
	sizeRule{MaxLines: 20, MaxLineWidth: 80}.check(w.SrcFile, "size_rule_test.go", "file.size", w.AstFile, w.FileSet)
	problems.Render()

	// Output:
	// +----------------------+--------------------------------------------------+------------------------+
	// |       position       |                     problem                      |          rule          |
	// +----------------------+--------------------------------------------------+------------------------+
	// | size_rule_test.go    | file size_rule_test.go size: 21 lines, limit: 20 | file.size.maxLines     |
	// | size_rule_test.go:11 | line 11 width: 114, limit: 80                    | file.size.maxLineWidth |
	// +----------------------+--------------------------------------------------+------------------------+
}
