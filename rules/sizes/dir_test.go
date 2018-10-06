package sizes

import (
	"fmt"

	"github.com/lovego/gospec/problems"
)

func ExampleCheckDir() {
	defer func(int old) {
		Rules.Dir = old
	}(Rules.Dir)
	Rules.Dir = 5

	problems.Clear()
	checkDir(".")
	problems.Render()

	// Output:
	// +----------+---------------------------------+-----------+
	// | position |             problem             |   rule    |
	// +----------+---------------------------------+-----------+
	// |        . | dir . size: 7 entries, limit: 5 | sizes.dir |
	// +----------+---------------------------------+-----------+
}

func ExampleEntriesCount() {
	fmt.Println(entriesCount("."))

	// Output: 7
}
