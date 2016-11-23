package names

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type ConfigT struct {
	Dir, File, Pkg                         configT
	Func, Const, Type, Var                 configT
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

	Func:  configT{Style: `camelCase`, MaxLen: 20},
	Const: configT{Style: `camelCase`, MaxLen: 20},
	Type:  configT{Style: `camelCase`, MaxLen: 20},
	Var:   configT{Style: `camelCase`, MaxLen: 20},

	LocalConst: configT{Style: `lowerCamelCase`, MaxLen: 15},
	LocalType:  configT{Style: `lowerCamelCase`, MaxLen: 15},
	LocalVar:   configT{Style: `lowerCamelCase`, MaxLen: 15},
	Label:      configT{Style: `lowerCamelCase`, MaxLen: 15},
}

var configValue = reflect.Indirect(reflect.ValueOf(&Config))

func getConfig(kind string, local bool) (configT, string) {
	switch kind {
	case `package`:
		return Config.Pkg, `names.pkg`
	case `func`:
		return Config.Func, `names.func`
	case `label`:
		return Config.Label, `names.label`
	default:
		key := kind
		if local {
			key = `local` + capitalize(key)
		}
		value := configValue.FieldByName(capitalize(key))
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

var lowercaseUnderline = regexp.MustCompile(`^(_?[a-z0-9]+)+$`)
var camelcase = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
var lowerCamelCase = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`)
var hasUppercase = regexp.MustCompile(`[A-Z]+`)

func checkStyle(name, style string, loose bool) bool {
	switch style {
	case `lower_case`:
		return loose && !hasUppercase.MatchString(name) || lowercaseUnderline.MatchString(name)
	case `camelCase`:
		return loose && strings.IndexByte(name, '_') < 0 || camelcase.MatchString(name)
	case `lowerCamelCase`:
		return loose && strings.IndexByte(name, '_') < 0 && (name[0] < 'a' || name[0] > 'z') ||
			lowerCamelCase.MatchString(name)
	default:
		panic(fmt.Sprintf(`unknown style config: "%s".`, style))
	}
}
