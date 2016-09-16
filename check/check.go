package check

import (
	"github.com/bughou-go/spec/check/names"
	"github.com/bughou-go/spec/check/sizes"
	"github.com/bughou-go/spec/d"
)

func Check(dir *d.Dir) {
	names.Check(dir)
	sizes.Check(dir)
}
