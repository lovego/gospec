package rules

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/lovego/gospec/rules/names"
	"github.com/lovego/gospec/rules/sizes"
	"gopkg.in/yaml.v2"
)

func init() {
	loadConfig()
}

func loadConfig() {
	p := getConfigPath()
	if p == `` {
		return
	}
	if content, err := ioutil.ReadFile(p); err == nil {
		parseConfig(content)
	} else {
		panic(err)
	}
}

func parseConfig(content []byte) {
	var config = &struct {
		Sizes *sizes.RulesT
		Names *names.RulesT
	}{
		Sizes: &sizes.Rules,
		Names: &names.Rules,
	}
	if err := yaml.Unmarshal(content, config); err != nil {
		panic(err)
	}
}

func getConfigPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for ; dir != `/`; dir = path.Dir(dir) {
		if p := testConfig(dir); p != `` {
			return p
		}
	}
	if p := testConfig(dir); p != `` {
		return p
	}
	return ``
}

func testConfig(dir string) string {
	p := path.Join(dir, `.gospec.yml`)
	if _, err := os.Stat(p); err == nil {
		return p
	} else if os.IsNotExist(err) {
		return ``
	} else {
		panic(err)
	}
}
