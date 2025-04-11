package rules

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	dirPkg "github.com/lovego/gospec/rules/objects/dir"
	pkgPkg "github.com/lovego/gospec/rules/objects/names/pkg"

	filePkg "github.com/lovego/gospec/rules/objects/file"
	funcPkg "github.com/lovego/gospec/rules/objects/func"
	structPkg "github.com/lovego/gospec/rules/objects/struct"

	constPkg "github.com/lovego/gospec/rules/objects/names/const"
	labelPkg "github.com/lovego/gospec/rules/objects/names/label"
	typePkg "github.com/lovego/gospec/rules/objects/names/type"
	varPkg "github.com/lovego/gospec/rules/objects/names/var"
)

var config = configT{
	Dir:        &dirPkg.Dir,
	Pkg:        &pkgPkg.Pkg,
	File:       &filePkg.File,
	TestFile:   &filePkg.TestFile,
	Func:       &funcPkg.Func,
	FuncInTest: &funcPkg.FuncInTest,
	Struct:     &structPkg.Struct,

	Const:      &constPkg.Const,
	LocalConst: &constPkg.LocalConst,
	Var:        &varPkg.Var,
	LocalVar:   &varPkg.LocalVar,
	Type:       &typePkg.Type,
	LocalType:  &typePkg.LocalType,
	Label:      &labelPkg.Label,
}

type configT struct {
	Targets    []string
	Dir        *dirPkg.Rule
	Pkg        *pkgPkg.Rule
	File       *filePkg.Rule
	TestFile   *filePkg.Rule
	Func       *funcPkg.Rule
	FuncInTest *funcPkg.Rule `yaml:"funcInTest"`
	Struct     *structPkg.Rule

	Const      *constPkg.Rule
	LocalConst *constPkg.Rule `yaml:"localConst"`
	Var        *varPkg.Rule
	LocalVar   *varPkg.Rule `yaml:"localVar"`
	Type       *typePkg.Rule
	LocalType  *typePkg.Rule `yaml:"localType"`
	Label      *labelPkg.Rule
}

// load config for rules
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

// specified dirs to run gospec
func Targets() []string {
	return config.Targets
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
	return testConfig(dir)
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
