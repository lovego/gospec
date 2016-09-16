package check

import (
	"github.com/bughou-go/spec/check/names"
	"github.com/bughou-go/spec/check/sizes"
)

type configStruct struct {
	Sizes sizes.Config
	Names names.Config
}

var config = configStruct{
	Sizes: sizes.DefaultConfig,
	Names: names.DefaultConfig,
}

func init() {
	// parse config from file
}
