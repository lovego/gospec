package names

import (
	"go/ast"
	"go/token"
)

type walker struct {
	local   bool
	fileSet *token.FileSet
}

func (w walker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		w.checkFuncDecl(n, w.file)
	case *ast.AssignStmt:
		w.checkShortVarDecl(n, w.file)
	case *ast.RangeStmt:
		w.checkRangeStmt(n, w.file)
	case *ast.StructType:
		w.checkStruct(n)
	case *ast.BlockStmt:
		if !w.local {
			return walker{local: true, fileSet: w.fileSet}
		}
	}
	return w
}

func (w walker) checkGenDecl(decl *ast.GenDecl, local bool) {
	if decl.Tok == token.IMPORT {
		return
	}
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			checkIdent(s.Name, local, file, ``)
		case *ast.ValueSpec:
			for _, ident := range s.Names {
				checkIdent(ident, local, file, ``)
			}
		}
	}
}

func (w walker) checkFuncDecl(fun *ast.FuncDecl) {
	checkIdent(n, w.local, w.fileSet)
	checkFieldList(fun.Recv, true, file, `func receiver`)
	checkFieldList(fun.Type.Params, true, file, `func param`)
	checkFieldList(fun.Type.Results, true, file, `func result`)
}

func (w walker) checkFuncLit(fun *ast.FuncLit) {
	checkFieldList(fun.Type.Params, true, file, `func param`)
	checkFieldList(fun.Type.Results, true, file, `func result`)
}

func (w walker) checkShortVarDecl(as *ast.AssignStmt) {
	if as.Tok != token.DEFINE {
		return
	}
	for _, exp := range as.Lhs {
		if ident, ok := exp.(*ast.Ident); ok {
			checkIdent(ident, true, file, ``)
		}
	}
}

func (w walker) checkRangeStmt(rg *ast.RangeStmt) {
	if rg.Tok != token.DEFINE {
		return
	}
	if ident, ok := rg.Key.(*ast.Ident); ok {
		checkIdent(ident, true, file, `range var`)
	}
	if ident, ok := rg.Value.(*ast.Ident); ok {
		checkIdent(ident, true, file, `range var`)
	}
}

func (w walker) checkInterface(itfc *ast.InterfaceType, local bool) {
	checkFieldList(itfc.Methods, local, file, `interface method`)
}

func (w walker) checkStruct(st *ast.StructType, local bool) {
	checkFieldList(st.Fields, local, file, `struct field`)
}

func (w walker) checkFieldList(fl *ast.FieldList, local bool, thing string) {
	if fl == nil {
		return
	}
	for _, f := range fl.List {
		for _, ident := range f.Names {
			checkIdent(ident, local, file, thing)
		}
		if ft, ok := f.Type.(*ast.FuncType); ok {
			checkFieldList(ft.Params, true, file, thing+` param`)
			checkFieldList(ft.Results, true, file, thing+` result`)
		}
	}
}

func (w walker) checkStruct(st *ast.StructType) {
	for _, f := range st.Fields.List {
		for _, ident := range f.Names {
			Rules.StructField.exec(ident.Name, w.fileSet.Position(ident.Pos()))
		}
	}
}

func (w walker) checkIdent(ident *ast.Ident, local bool, thing string) {
	if ident == nil || ident.Obj == nil {
		return
	}
	kind := ident.Obj.Kind.String()
	cfg, rule := getConfig(kind, local, file.Name())
	if cfg.Style == `` {
		return
	}

	if desc := checkName(ident.Name, cfg, false); desc != `` {
		problems.Add(file.Position(ident.Pos()),
			fmt.Sprintf(`%s name %s %s`, getThing(thing, local, kind), ident.Name, desc), rule,
		)
	}
}
