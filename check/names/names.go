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
	Dir, File, Pkg, Const, Type, Var, Func, Label configT
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

func CheckFunc(n ast.Node, file *token.File) {
	switch fun := n.(type) {
	case *ast.FuncDecl:
		CheckIdent(fun.Name, file)
		checkFieldList(fun.Type.Recv, file, `func receiver`)
		checkFieldList(fun.Type.Params, file, `func param`)
		checkFieldList(fun.Type.Results, file, `func result`)
	case *ast.FuncLit:
		checkFieldList(fun.Type.Params, file, `func param`)
		checkFieldList(fun.Type.Results, file, `func result`)
	}
}

func checkFieldList(fl *ast.FieldList, file *token.File, kind string) {
	if fl == nil {
		return
	}
	for _, f := range fl.List {
		for _, ident := range f.Names {
			desc := checkName(ident.Name, Config.Var)
			if desc == `` {
				continue
			}
			problems.Add(file.Position(ident.Pos()),
				fmt.Sprintf(`%s %s %s`, kind, ident.Name, desc), `names.var`,
			)
		}
	}
}

func CheckShortVarDecl(as *ast.AssignStmt, file *token.File) {
	if as.Tok != token.DEFINE {
		return
	}
	for _, exp := range as.Lhs {
		if ident, ok := exp.(*ast.Ident); ok {
			names.CheckIdent(ident, file)
		}
	}
}

func CheckGenDecl(decl *ast.GenDecl, file *token.File) {
	if decl.Tok == token.IMPORT {
		return true
	}
	for _, spec := range v.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			CheckIdent(s.Name, file)
		case *ast.ValueSpec:
			for _, ident := range s.Names {
				CheckIdent(ident, file)
			}
		}
	}
}

func CheckIdent(ident *ast.Ident, file *token.File) {
	var desc string
	kind := ident.Obj.ObjKind
	switch kind {
	case ast.Con:
		desc = checkName(ident.Name, Config.Const)
	case ast.Typ:
		desc = checkName(ident.Name, Config.Type)
	case ast.Var:
		desc = checkName(ident.Name, Config.Var)
	case ast.Fun:
		desc = checkName(ident.Name, Config.Func)
	case ast.Label:
		desc = checkName(ident.Name, Config.Label)
	}
	if desc != `` {
		problems.Add(file.Position(ident.Pos()),
			fmt.Sprintf(`%s %s %s`, kind, ident.Name, desc), `names.`+kind.String(),
		)
	}
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
