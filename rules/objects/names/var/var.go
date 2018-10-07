package varpkg

import (
	"go/ast"
	"go/token"

	"github.com/lovego/gospec/rules/name"
)

var Rule = name.Rule{
	MaxLen: 40,
	Style:  "camelCase",
}
var LocalRule = name.Rule{
	MaxLen: 30,
	Style:  "lowerCamelCase",
}

func Check(local bool, node ast.Node, fileSet *token.FileSet) {
	switch n := node.(type) {
	case *ast.ValueSpec:
		checkVarDecl(local, n, fileSet)
		return
	case *ast.FuncType:
		checkFuncType(n, fileSet)
		return
	}
	if local {
		switch n := node.(type) {
		case *ast.AssignStmt:
			checkShortVarDecl(n, fileSet)
		case *ast.RangeStmt:
			checkRangeStmt(n, fileSet)
		}
	}
}

func checkVarDecl(local bool, vars *ast.ValueSpec, fileSet *token.FileSet) {
	if vars.Names[0].Obj.Kind != ast.Var {
		return
	}
	for _, ident := range vars.Names {
		if local {
			LocalRule.Exec(ident.Name, `local var`, `localVar`, fileSet.Position(ident.Pos()))
		} else {
			Rule.Exec(ident.Name, `var`, `var`, fileSet.Position(ident.Pos()))
		}
	}
}

func checkShortVarDecl(assign *ast.AssignStmt, fileSet *token.FileSet) {
	if assign.Tok != token.DEFINE {
		return
	}
	for _, exp := range assign.Lhs {
		if ident, ok := exp.(*ast.Ident); ok {
			LocalRule.Exec(ident.Name, `local var`, `localVar`, fileSet.Position(ident.Pos()))
		}
	}
}

func checkRangeStmt(rang *ast.RangeStmt, fileSet *token.FileSet) {
	if rang.Tok != token.DEFINE {
		return
	}
	if ident, ok := rang.Key.(*ast.Ident); ok {
		LocalRule.Exec(ident.Name, `range var`, `localVar`, fileSet.Position(ident.Pos()))
	}
	if ident, ok := rang.Value.(*ast.Ident); ok {
		LocalRule.Exec(ident.Name, `range var`, `localVar`, fileSet.Position(ident.Pos()))
	}
}

func checkFuncType(typ *ast.FuncType, fileSet *token.FileSet) {
	checkFieldList(`func param`, typ.Params, fileSet)
	checkFieldList(`func result`, typ.Results, fileSet)
}

func checkFieldList(thing string, fl *ast.FieldList, fileSet *token.FileSet) {
	if fl == nil {
		return
	}
	for _, f := range fl.List {
		for _, ident := range f.Names {
			LocalRule.Exec(ident.Name, thing, `localVar`, fileSet.Position(ident.Pos()))
		}
	}
}
