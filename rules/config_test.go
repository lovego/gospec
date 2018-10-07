package rules

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/lovego/deep"
	"gopkg.in/yaml.v2"
)

func ExampleCheckDefaultConfig() {
	content, err := ioutil.ReadFile("../gospec.yml")
	if err != nil {
		panic(err)
	}

	var defaultConfigInYaml configT
	if err := yaml.Unmarshal(content, &defaultConfigInYaml); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(strings.Join(deep.Equal(defaultConfigInYaml, config), "\n"))

	// Output:
}
