package varpkg

import (
	"go/ast"
	"go/token"

	namepkg "github.com/lovego/gospec/rules/name"
)

type Rule struct {
	thing, key   string
	namepkg.Rule `yaml:",inline"`
}

func (r Rule) checkDecl(value *ast.ValueSpec, fileSet *token.FileSet) {
	if value.Names[0].Obj.Kind != ast.Var {
		return
	}
	for _, ident := range value.Names {
		r.Exec(ident.Name, r.thing, r.key, fileSet.Position(ident.Pos()))
	}
}

func (r Rule) checkShortDecl(assign *ast.AssignStmt, fileSet *token.FileSet) {
	if assign.Tok != token.DEFINE {
		return
	}
	for _, exp := range assign.Lhs {
		if ident, ok := exp.(*ast.Ident); ok {
			r.Exec(ident.Name, r.thing, r.key, fileSet.Position(ident.Pos()))
		}
	}
}

func (r Rule) checkRangeStmt(rang *ast.RangeStmt, fileSet *token.FileSet) {
	if rang.Tok != token.DEFINE {
		return
	}
	if ident, ok := rang.Key.(*ast.Ident); ok {
		r.Exec(ident.Name, `range var`, r.key, fileSet.Position(ident.Pos()))
	}
	if ident, ok := rang.Value.(*ast.Ident); ok {
		r.Exec(ident.Name, `range var`, r.key, fileSet.Position(ident.Pos()))
	}
}

func (r Rule) checkFuncType(typ *ast.FuncType, fileSet *token.FileSet) {
	r.checkFieldList(`func param`, typ.Params, fileSet)
	r.checkFieldList(`func result`, typ.Results, fileSet)
}

func (r Rule) checkFieldList(thing string, fl *ast.FieldList, fileSet *token.FileSet) {
	if fl == nil {
		return
	}
	for _, f := range fl.List {
		for _, ident := range f.Names {
			r.Exec(ident.Name, thing, r.key, fileSet.Position(ident.Pos()))
		}
	}
}
