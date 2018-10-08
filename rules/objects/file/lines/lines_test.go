package lines

import (
	"fmt"
	"sort"

	"github.com/lovego/gospec/rules/walker"
)

func ExampleLines() {
	w := walker.New("lines.go")
	New(w.SrcFile, w.AstFile, w.FileSet)
	// Output:
}

var testWalker = walker.Parse("example.go", `package example
// 45678

/* 45 */var a = 8
var a=3 //

/*

*/
`)

func ExampleLines_2() {
	lines := New(testWalker.SrcFile, testWalker.AstFile, testWalker.FileSet)

	fmt.Println(lines.Num())
	fmt.Println(lines.Get(9))
	fmt.Println(lines.IsComment(2), lines.IsComment(4))
	fmt.Println(mapKeys(lines.comments))

	// Output:
	// 9
	// */
	// true false
	// [2 7 8 9]
}

func ExampleComment_position() {
	for _, commentGroup := range testWalker.AstFile.Comments {
		start := testWalker.FileSet.Position(commentGroup.Pos())
		end := testWalker.FileSet.Position(commentGroup.End())
		fmt.Printf("%d:%d ~ %d:%d\n", start.Line, start.Column, end.Line, end.Column)
	}
	// Output:
	// 2:1 ~ 2:9
	// 4:1 ~ 4:9
	// 5:9 ~ 5:11
	// 7:1 ~ 9:3
}

func mapKeys(m map[int]bool) (keys []int) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return
}
