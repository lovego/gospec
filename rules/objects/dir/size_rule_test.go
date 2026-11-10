package dirpkg

import (
	"fmt"

	"github.com/lovego/gospec/problems"
)

func Example_sizeRuleCheck() {
	problems.Clear()
	sizeRule{MaxEntries: 3}.check(".", "dir.size")
	problems.Render()

	// Output:
	// +----------+---------------------------------+---------------------+
	// | position |             problem             |        rule         |
	// +----------+---------------------------------+---------------------+
	// | .        | dir . size: 5 entries, limit: 3 | dir.size.maxEntries |
	// +----------+---------------------------------+---------------------+
}

func Example_entriesCount() {
	fmt.Println(entriesCount("."))
	// Output: 5
}
