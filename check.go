package spec

import (
	"github.com/bughou-go/spec/specs"
	"github.com/bughou-go/spec/specs/names"
	"github.com/bughou-go/spec/specs/sizes"
)

func check(dir *spec.Dir) {
	names.Check(dir)
	sizes.Check(dir)
}
