package check

import (
	"os"

	"github.com/bughou-go/spec/c"
	"github.com/bughou-go/spec/check/names"
	"github.com/bughou-go/spec/check/sizes"
)

func Check(dir *c.Dir) {
	sizes.Check(dir)
	names.Check(dir)

	if c.ProblemsCount() > 0 {
		os.Exit(1)
	}
}
