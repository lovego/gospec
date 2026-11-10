package rules

import (
	"fmt"
	"strings"

	"github.com/lovego/deep"
)

func Example_checkDefaultConfig() {
	configInCode := config
	config = configT{}
	LoadConfig()
	configInYaml := config

	fmt.Println(strings.Join(deep.Equal(configInYaml, configInCode), "\n"))

	// Output:
}
