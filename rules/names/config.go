package names

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type ConfigT struct {
	Dir, File, Pkg                         configT
	Const, Type, Var, Func, FuncInTest     configT
	LocalConst, LocalType, LocalVar, Label configT
}
type configT struct {
	Style  string
	MaxLen int
}

var Config = ConfigT{
	Dir:  configT{Style: `lower_case`, MaxLen: 20},
	File: configT{Style: `lower_case`, MaxLen: 20},
	Pkg:  configT{Style: `lower_case`, MaxLen: 20},

	Type:       configT{Style: `camelCase`, MaxLen: 30},
	Const:      configT{Style: `camelCase`, MaxLen: 30},
	Var:        configT{Style: `camelCase`, MaxLen: 40},
	Func:       configT{Style: `camelCase`, MaxLen: 30},
	FuncInTest: configT{Style: `camelCaseInTest`, MaxLen: 50},

	LocalConst: configT{Style: `lowerCamelCase`, MaxLen: 20},
	LocalType:  configT{Style: `lowerCamelCase`, MaxLen: 20},
	LocalVar:   configT{Style: `lowerCamelCase`, MaxLen: 30},
	Label:      configT{Style: `lowerCamelCase`, MaxLen: 20},
}

type FuncConfigT struct {
	Param  int
	Result int
}

var FuncConfig = FuncConfigT{
	Param:  5,
	Result: 3,
}

var configValue = reflect.ValueOf(&Config).Elem()

func getFuncConfig(kind string) int {
	switch kind {
	case `param`:
		return FuncConfig.Param
	case `result`:
		return FuncConfig.Result
	}
	return 0
}

func getConfig(kind string, local bool, fileName string) (configT, string) {
	switch kind {
	case `package`:
		return Config.Pkg, `names.pkg`
	case `func`:
		if strings.HasSuffix(fileName, "_test.go") {
			return Config.FuncInTest, `names.funcInTest`
		} else {
			return Config.Func, `names.func`
		}
	case `label`:
		return Config.Label, `names.label`
	default:
		key := capitalize(kind)
		if local {
			key = `Local` + key
		}
		value := configValue.FieldByName(key)
		return value.Interface().(configT), `names.` + key
	}
}

func capitalize(s string) string {
	b := []byte(s)
	b[0] -= 0x20
	return string(b)
}

func checkName(name string, config configT, loose bool) string {
	if name == `_` {
		return ``
	}
	desc := []string{}
	if len(name) > config.MaxLen {
		desc = append(desc, fmt.Sprintf(`%d chars long, limits %d`, len(name), config.MaxLen))
	}
	if !checkStyle(name, config.Style, loose) {
		desc = append(desc, fmt.Sprintf(`should be %s style`, config.Style))
	}
	return strings.Join(desc, ` and `)
}

var lowercaseUnderscore = regexp.MustCompile(`^(_?[a-z0-9]+)+$`)
var lowercaseDash = regexp.MustCompile(`^(-?[a-z0-9]+)+$`)
var camelcase = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
var lowerCamelCase = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`)
var hasUppercase = regexp.MustCompile(`[A-Z]+`)
var hasUppercaseOrUnderscore = regexp.MustCompile(`[A-Z_]+`)
var exampleTestCase = regexp.MustCompile(`^Example[_A-Z]`)

func checkStyle(name, style string, loose bool) bool {
	switch style {
	case `lower_case`:
		return loose && !hasUppercase.MatchString(name) || lowercaseUnderscore.MatchString(name)
	case `lower-case`:
		return loose && !hasUppercaseOrUnderscore.MatchString(name) ||
			lowercaseDash.MatchString(name)
	case `camelCase`:
		return loose && strings.IndexByte(name, '_') < 0 || camelcase.MatchString(name)
	case `camelCaseInTest`:
		return loose && strings.IndexByte(name, '_') < 0 || camelcase.MatchString(name) ||
			exampleTestCase.MatchString(name)
	case `lowerCamelCase`:
		return loose && strings.IndexByte(name, '_') < 0 && (name[0] < 'a' || name[0] > 'z') ||
			lowerCamelCase.MatchString(name)
	default:
		panic(fmt.Sprintf(`unknown style config: "%s".`, style))
	}
}