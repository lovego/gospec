package names

import (
	"fmt"
	"go/ast"
	"regexp"
)

type ConfigT struct {
	Dir, File, Pkg, Const, Type, Var, Func, Label configT
}
type configT struct {
	Style  string
	MaxLen int
}

var Config = ConfigT{
	Dir:  configT{Style: `lowercase`, MaxLen: 20},
	File: configT{Style: `lowercase`, MaxLen: 20},
	Pkg:  configT{Style: `lowercase`, MaxLen: 20},

	Type: configT{Style: `camelCase`, MaxLen: 20},
	Func: configT{Style: `camelCase`, MaxLen: 20},

	Const: configT{Style: `camelCase`, MaxLen: 20},
	Var:   configT{Style: `camelCase`, MaxLen: 20},
	Label: configT{Style: `camelCase`, MaxLen: 20},
	// LocalConst: configT{Style: `camelCase`, MaxLen: 10},
	// LocalVar:   configT{Style: `camelCase`, MaxLen: 10},
}

func getConfig(objKind ast.ObjKind) configT {
	switch objKind {
	case ast.Pkg:
		return Config.Pkg
	case ast.Con:
		return Config.Const
	case ast.Typ:
		return Config.Type
	case ast.Var:
		return Config.Var
	case ast.Fun:
		return Config.Func
	case ast.Lbl:
		return Config.Label
	default:
		return configT{}
	}
}

func checkName(name string, config configT) string {
	if name == `_` {
		return ``
	}
	styleRight := checkStyle(config.Style, name)
	sizeRight := len(name) <= config.MaxLen
	switch {
	case styleRight && sizeRight:
		return ``
	case !styleRight && sizeRight:
		return fmt.Sprintf(`should be %s style`, config.Style)
	case styleRight && !sizeRight:
		return fmt.Sprintf(`shouldn't be more than %d chars`, config.MaxLen)
	default:
		return fmt.Sprintf(`should be %s style and not more than %d chars`, config.Style, config.MaxLen)
	}
}

var lowercaseRegexp = regexp.MustCompile(`^(_?[a-z0-9]+)+$`)
var camelcaseRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func checkStyle(style, name string) bool {
	switch style {
	case `lowercase`:
		return lowercaseRegexp.MatchString(name)
	case `camelCase`:
		return camelcaseRegexp.MatchString(name)
	default:
		panic(fmt.Sprintf(`unknown style config: "%s".`, style))
	}
}
