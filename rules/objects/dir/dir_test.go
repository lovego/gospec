package dir

import (
	"fmt"

	"github.com/lovego/gospec/problems"
)

func ExampleCheck() {
	defer func(old uint) {
		Rule.Size.MaxEntries = old
	}(Rule.Size.MaxEntries)

	Rule.Size.MaxEntries = 1

	problems.Clear()
	Check("../dir")
	problems.Render()

	// Output:
	// +----------+-----------------------------------+---------------------+
	// | position |              problem              |        rule         |
	// +----------+-----------------------------------+---------------------+
	// | ../dir   | dir dir size: 2 entries, limit: 1 | dir.size.maxEntries |
	// +----------+-----------------------------------+---------------------+
}

func ExampleCheck_emtpty() {
	problems.Clear()
	Check("../dir")
	Check(".")
	Check("..")
	problems.Render()
	// Output:
}

func ExampleEntriesCount() {
	fmt.Println(entriesCount("."))
	// Output: 2
}
