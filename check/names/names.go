package names

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"regexp"
	"strings"

	"github.com/bughou-go/spec/problems"
)

type ConfigT struct {
	Dir, File, Pkg, Func, Const, Var, LocalConst, LocalVar configT
}
type configT struct {
	Style  string
	MaxLen int
}

var Config = ConfigT{
	Dir:        configT{Style: `lowercase`, MaxLen: 20},
	Pkg:        configT{Style: `lowercase`, MaxLen: 20},
	File:       configT{Style: `lowercase`, MaxLen: 20},
	Func:       configT{Style: `camelCase`, MaxLen: 20},
	Const:      configT{Style: `camelCase`, MaxLen: 20},
	Var:        configT{Style: `camelCase`, MaxLen: 20},
	LocalConst: configT{Style: `camelCase`, MaxLen: 10},
	LocalVar:   configT{Style: `camelCase`, MaxLen: 10},
}

func CheckDir(p string) {
	if p == `.` || p == `..` || p == `/` {
		return
	}
	name := path.Base(p)
	desc := checkName(name, Config.Dir)
	if desc != `` {
		problems.Add(token.Position{Filename: p}, fmt.Sprintf(`dir %s %s`, name, desc), `names.dir`)
	}
}

func CheckPkg(pkg *ast.Package, fset *token.FileSet) {
	desc := checkName(pkg.Name, Config.Pkg)
	if desc != `` {
		var f *ast.File
		for _, file := range pkg.Files {
			f = file
			break
		}
		problems.Add(fset.Position(f.Name.Pos()), fmt.Sprintf(`package %s %s`, pkg.Name, desc), `names.pkg`)
	}
}

func CheckFile(p string) {
	name := path.Base(p)
	desc := checkName(strings.TrimSuffix(name, `.go`), Config.File)
	if desc != `` {
		problems.Add(token.Position{Filename: p}, fmt.Sprintf(`file %s %s`, name, desc), `names.file`)
	}
}

func CheckFunc(fun *ast.FuncDecl, file *token.File) {
	desc := checkName(fun.Name.Name, Config.Func)
	if desc != `` {
		problems.Add(file.Position(fun.Pos()), fmt.Sprintf(`func %s %s`, fun.Name.Name, desc), `names.func`)
	}
}

func CheckVar(n ast.Node, file *token.File) {
}

func CheckConst(n ast.Node, file *token.File) {
}

func checkName(name string, config configT) string {
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
