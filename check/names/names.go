package names

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"strings"

	"github.com/bughou-go/spec/problems"
)

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

func CheckGenDecl(decl *ast.GenDecl, file *token.File) {
	if decl.Tok == token.IMPORT {
		return
	}
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			CheckIdent(s.Name, file, ``)
		case *ast.ValueSpec:
			for _, ident := range s.Names {
				CheckIdent(ident, file, ``)
			}
		}
	}
}

func CheckFunc(n ast.Node, file *token.File) {
	switch fun := n.(type) {
	case *ast.FuncDecl:
		CheckIdent(fun.Name, file, ``)
		checkFieldList(fun.Recv, file, `func receiver`)
		checkFieldList(fun.Type.Params, file, `func param`)
		checkFieldList(fun.Type.Results, file, `func result`)
	case *ast.FuncLit:
		checkFieldList(fun.Type.Params, file, `func param`)
		checkFieldList(fun.Type.Results, file, `func result`)
	}
}

func CheckShortVarDecl(as *ast.AssignStmt, file *token.File) {
	if as.Tok != token.DEFINE {
		return
	}
	for _, exp := range as.Lhs {
		if ident, ok := exp.(*ast.Ident); ok {
			CheckIdent(ident, file, ``)
		}
	}
}

func CheckRangeStmt(rg *ast.RangeStmt, file *token.File) {
	if rg.Tok != token.DEFINE {
		return
	}
	if ident, ok := rg.Key.(*ast.Ident); ok {
		CheckIdent(ident, file, `range var`)
	}
	if ident, ok := rg.Value.(*ast.Ident); ok {
		CheckIdent(ident, file, `range var`)
	}
}

func CheckInterface(itfc *ast.InterfaceType, file *token.File) {
	checkFieldList(itfc.Methods, file, `interface method`)
}

func CheckStruct(st *ast.StructType, file *token.File) {
	checkFieldList(st.Fields, file, `struct field`)
}

func CheckIdent(ident *ast.Ident, file *token.File, thing string) {
	if ident == nil || ident.Obj == nil {
		return
	}
	objKind := ident.Obj.Kind
	conf := getConfig(objKind)
	if conf.Style == `` {
		return
	}
	desc := checkName(ident.Name, conf)
	if desc == `` {
		return
	}
	kind := objKind.String()
	if thing == `` {
		thing = kind
	}
	problems.Add(file.Position(ident.Pos()),
		fmt.Sprintf(`%s %s %s`, thing, ident.Name, desc), `names.`+kind,
	)
}

func checkFieldList(fl *ast.FieldList, file *token.File, thing string) {
	if fl == nil {
		return
	}
	for _, f := range fl.List {
		for _, ident := range f.Names {
			CheckIdent(ident, file, thing)
		}
		if ft, ok := f.Type.(*ast.FuncType); ok {
			checkFieldList(ft.Params, file, thing+` param`)
			checkFieldList(ft.Results, file, thing+` result`)
		}
	}
}
