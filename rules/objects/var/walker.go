package names

import (
	"go/ast"
	"go/token"
	"strings"
)

type walker struct {
	local   bool
	fileSet *token.FileSet
}

func (w walker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		w.checkFuncDecl(n)
	case *ast.FuncType:
		w.checkFuncType(n)
	case *ast.StructType:
		w.checkStruct(n)
	case *ast.GenDecl:
		w.checkGenDecl(n)
	case *ast.InterfaceType:
		w.checkInterface(n)
	case *ast.AssignStmt:
		w.checkShortVarDecl(n)
	case *ast.RangeStmt:
		w.checkRangeStmt(n)
	case *ast.BlockStmt:
		if !w.local {
			return walker{local: true, fileSet: w.fileSet}
		}
	}
	return w
}

func (w walker) checkStruct(strut *ast.StructType) {
	for _, f := range strut.Fields.List {
		for _, ident := range f.Names {
			Rules.StructField.Exec(`struct field`, ident.Name, w.fileSet.Position(ident.Pos()))
		}
	}
}

func (w walker) checkGenDecl(decl *ast.GenDecl) {
	if decl.Tok == token.IMPORT {
		return
	}
	for _, s := range decl.Specs {
		switch spec := s.(type) {
		case *ast.TypeSpec:
			ident := spec.Name
			checkIdent(``, ident, w.local, w.fileSet)
		case *ast.ValueSpec:
			for _, ident := range spec.Names {
				checkIdent(``, ident, w.local, w.fileSet)
			}
		}
	}
}

func (w walker) checkFuncType(fun *ast.FuncType) {
	checkFieldList(`func param`, fun.Params, true, w.fileSet)
	checkFieldList(`func result`, fun.Results, true, w.fileSet)
}

func (w walker) checkInterface(ifc *ast.InterfaceType) {
	checkFieldList(`interface method`, ifc.Methods, w.local, w.fileSet)
}

func (w walker) checkShortVarDecl(assign *ast.AssignStmt) {
	if assign.Tok != token.DEFINE {
		return
	}
	for _, exp := range assign.Lhs {
		if ident, ok := exp.(*ast.Ident); ok {
			checkIdent(``, ident, true, w.fileSet)
		}
	}
}

func (w walker) checkRangeStmt(rang *ast.RangeStmt) {
	if rang.Tok != token.DEFINE {
		return
	}
	if ident, ok := rang.Key.(*ast.Ident); ok {
		checkIdent(`range var`, ident, true, w.fileSet)
	}
	if ident, ok := rang.Value.(*ast.Ident); ok {
		checkIdent(`range var`, ident, true, w.fileSet)
	}
}

func checkFieldList(thing string, fl *ast.FieldList, local bool, fileSet *token.FileSet) {
	if fl == nil {
		return
	}
	for _, f := range fl.List {
		for _, ident := range f.Names {
			checkIdent(thing, ident, local, fileSet)
		}
	}
}
