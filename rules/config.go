package rules

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/lovego/gospec/rules/name"
	"gopkg.in/yaml.v2"

	dirPkg "github.com/lovego/gospec/rules/objects/dir"
	filePkg "github.com/lovego/gospec/rules/objects/file"
	funcPkg "github.com/lovego/gospec/rules/objects/func"
	structPkg "github.com/lovego/gospec/rules/objects/struct"

	constPkg "github.com/lovego/gospec/rules/objects/names/const"
	labelPkg "github.com/lovego/gospec/rules/objects/names/label"
	pkgPkg "github.com/lovego/gospec/rules/objects/names/pkg"
	typePkg "github.com/lovego/gospec/rules/objects/names/type"
	varPkg "github.com/lovego/gospec/rules/objects/names/var"
)

var config = configT{
	Dir:        &dirPkg.Rule,
	File:       &filePkg.Rule,
	TestFile:   &filePkg.TestFileRule,
	Func:       &funcPkg.Rule,
	FuncInTest: &funcPkg.RuleInTest,
	Struct:     &structPkg.Rule,

	Pkg:        &pkgPkg.Rule,
	Const:      &constPkg.Rule,
	LocalConst: &constPkg.LocalRule,
	Var:        &varPkg.Rule,
	LocalVar:   &varPkg.LocalRule,
	Type:       &typePkg.Rule,
	LocalType:  &typePkg.LocalRule,
	Label:      &labelPkg.Rule,
}

type configT struct {
	Dir        *dirPkg.RuleT
	Pkg        *name.Rule
	File       *filePkg.RuleT
	TestFile   *filePkg.RuleT
	Func       *funcPkg.RuleT
	FuncInTest *funcPkg.RuleT `yaml:"funcInTest"`
	Struct     *structPkg.RuleT

	Const      *name.Rule
	LocalConst *name.Rule `yaml:"localConst"`
	Var        *name.Rule
	LocalVar   *name.Rule `yaml:"localVar"`
	Type       *name.Rule
	LocalType  *name.Rule `yaml:"localType"`
	Label      *name.Rule
}

func LoadConfig() {
	p := getConfigPath()
	if p == `` {
		return
	}
	content, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(content, &config); err != nil {
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
