package rules

import (
	"io/ioutil"
	"os"
	"path"

	dirPkg "github.com/lovego/gospec/rules/objects/dir"
	filePkg "github.com/lovego/gospec/rules/objects/file"
	funPkg "github.com/lovego/gospec/rules/objects/fun"
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
		Dir        *dirPkg.RuleT
		File       *filePkg.RuleT
		Func       *funPkg.RuleT
		FuncInTest *funPkg.RuleT
	}{
		Dir:        &dirPkg.Rule,
		File:       &filePkg.Rule,
		Func:       &funPkg.Rule,
		FuncInTest: &funPkg.RuleInTest,
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
