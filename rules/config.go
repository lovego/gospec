package rules

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/lovego/gospec/rules/names"
	"github.com/lovego/gospec/rules/sizes"
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
	if err := json.Unmarshal(content, config); err != nil {
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
	const file = `gospec.json`
	p := path.Join(dir, file)
	if _, err := os.Stat(p); err == nil {
		return p
	} else if os.IsNotExist(err) {
		return ``
	} else {
		panic(err)
	}
}
