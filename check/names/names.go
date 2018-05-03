package names

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"strings"

	"github.com/lovego/spec/problems"
)

func CheckDir(p string) {
	if p == `.` || p == `..` || p == `/` {
		return
	}
	name := path.Base(p)
	desc := checkName(name, Config.Dir, false)
	if desc != `` {
		problems.Add(token.Position{Filename: p}, fmt.Sprintf(`dir name %s %s`, name, desc), `names.dir`)
	}
}

func CheckFile(p string) {
	name := path.Base(p)
	namePart := strings.TrimSuffix(name, `.go`)
	namePart = strings.TrimSuffix(namePart, `_test`)
	desc := checkName(namePart, Config.File, false)
	if desc != `` {
		problems.Add(token.Position{Filename: p},
			fmt.Sprintf(`file name %s %s`, name, desc), `names.file`,
		)
	}
}

func CheckGenDecl(decl *ast.GenDecl, local bool, file *token.File) {
	if decl.Tok == token.IMPORT {
		return
	}
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			CheckIdent(s.Name, local, file, ``)
		case *ast.ValueSpec:
			for _, ident := range s.Names {
				CheckIdent(ident, local, file, ``)
			}
		}
	}
}

func CheckFuncDecl(fun *ast.FuncDecl, file *token.File) {
	CheckIdent(fun.Name, false, file, ``)
	checkFieldList(fun.Recv, true, file, `func receiver`)
	checkFieldList(fun.Type.Params, true, file, `func param`)
	checkFieldList(fun.Type.Results, true, file, `func result`)
}

func CheckFuncLit(fun *ast.FuncLit, file *token.File) {
	checkFieldList(fun.Type.Params, true, file, `func param`)
	checkFieldList(fun.Type.Results, true, file, `func result`)
}

func CheckShortVarDecl(as *ast.AssignStmt, file *token.File) {
	if as.Tok != token.DEFINE {
		return
	}
	for _, exp := range as.Lhs {
		if ident, ok := exp.(*ast.Ident); ok {
			CheckIdent(ident, true, file, ``)
		}
	}
}

func CheckRangeStmt(rg *ast.RangeStmt, file *token.File) {
	if rg.Tok != token.DEFINE {
		return
	}
	if ident, ok := rg.Key.(*ast.Ident); ok {
		CheckIdent(ident, true, file, `range var`)
	}
	if ident, ok := rg.Value.(*ast.Ident); ok {
		CheckIdent(ident, true, file, `range var`)
	}
}

func CheckInterface(itfc *ast.InterfaceType, local bool, file *token.File) {
	checkFieldList(itfc.Methods, local, file, `interface method`)
}

func CheckStruct(st *ast.StructType, local bool, file *token.File) {
	checkFieldList(st.Fields, local, file, `struct field`)
}

func CheckIdent(ident *ast.Ident, local bool, file *token.File, thing string) {
	if ident == nil || ident.Obj == nil {
		return
	}
	kind := ident.Obj.Kind.String()
	cfg, rule := getConfig(kind, local, file.Name())
	if cfg.Style == `` {
		return
	}
	if desc := checkName(ident.Name, cfg, true); desc != `` {
		problems.Add(file.Position(ident.Pos()),
			fmt.Sprintf(`%s name %s %s`, getThing(thing, local, kind), ident.Name, desc), rule,
		)
	}
}

func getThing(thing string, local bool, kind string) string {
	if thing == `` {
		if local {
			thing = `local `
		}
		thing += kind
	}
	return thing
}

func checkFieldList(fl *ast.FieldList, local bool, file *token.File, thing string) {
	if fl == nil {
		return
	}
	for _, f := range fl.List {
		for _, ident := range f.Names {
			CheckIdent(ident, local, file, thing)
		}
		if ft, ok := f.Type.(*ast.FuncType); ok {
			checkFieldList(ft.Params, true, file, thing+` param`)
			checkFieldList(ft.Results, true, file, thing+` result`)
		}
	}
}
