package filepkg

import (
	"github.com/lovego/gospec/problems"
	"github.com/lovego/gospec/rules/walker"
)

func ExampleSizeRule_check() {
	problems.Clear()
	w := walker.Parse("example.go", `package example
var a = 3
var b = "long long line"
// comment
// long long long long long comment
 `)
	sizeRule{MaxLineWidth: 20, MaxCommentLineWidth: 30, MaxLines: 3}.check(
		w.SrcFile, w.AstFile, w.FileSet, "size_rule_test.go", "file.size",
	)
	problems.Render()

	// Output:
	// +---------------------+------------------------------------------------+-------------------------------+
	// |      position       |                    problem                     |             rule              |
	// +---------------------+------------------------------------------------+-------------------------------+
	// | size_rule_test.go   | file size_rule_test.go size: 6 lines, limit: 3 | file.size.maxLines            |
	// | size_rule_test.go:3 | line 3 width: 24, limit: 20                    | file.size.maxLineWidth        |
	// | size_rule_test.go:5 | line 5 width: 35, limit: 30                    | file.size.maxCommentLineWidth |
	// +---------------------+------------------------------------------------+-------------------------------+
}
