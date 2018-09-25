package names

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"strings"

	"github.com/lovego/gospec/problems"
)

type NameWalker struct {
	srcFile *token.File
	fn      ast.Node //stack
	begin   int
	end     int
}

func NewNameWalker(srcFile *token.File) *NameWalker {
	return &NameWalker{
		srcFile: srcFile,
	}
}

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

func (nw *NameWalker) Visit(node ast.Node) ast.Node {
	srcFile := nw.srcFile
	switch n := node.(type) {
	case *ast.FuncDecl:
		v := node.(*ast.FuncDecl)
		nw.begin = nw.srcFile.Position(v.Body.Lbrace).Line
		nw.end = nw.srcFile.Position(v.Body.Rbrace).Line
		checkFuncDecl(n, srcFile)
	case *ast.AssignStmt:
		checkShortVarDecl(n, srcFile)
	case *ast.RangeStmt:
		checkRangeStmt(n, srcFile)
	case *ast.GenDecl:
		v := node.(*ast.GenDecl)
		checkGenDecl(n, nw.isLocal(v.TokPos), srcFile)
	// type define
	case *ast.StructType:
		v := node.(*ast.StructType)
		checkStruct(n, nw.isLocal(v.Struct), srcFile)
	case *ast.InterfaceType:
		v := node.(*ast.InterfaceType)
		checkInterface(n, nw.isLocal(v.Interface), srcFile)
	// func literal
	case *ast.FuncLit:
		checkFuncLit(n, srcFile)
	}
	return node
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

func checkFuncDecl(fun *ast.FuncDecl, file *token.File) {
	CheckIdent(fun.Name, false, file, ``)
	checkFieldList(fun.Recv, true, file, `func receiver`)
	checkFieldList(fun.Type.Params, true, file, `func param`)
	checkFieldList(fun.Type.Results, true, file, `func result`)
	checkParamNum(fun.Name.Name, fun.Type.Params, `param`,
		file.Position(fun.Pos()))
	checkResultNum("", fun.Type.Results, `result`, file.Position(fun.Pos()))
}

func checkParamNum(name string, list *ast.FieldList, kind string,
	pos token.Position) {
	checkNum(name, list, kind, pos)
}

func checkResultNum(name string, list *ast.FieldList, kind string,
	pos token.Position) {
	checkNum(name, list, kind, pos)
}

func checkNum(name string, list *ast.FieldList, kind string,
	pos token.Position) {
	if list == nil {
		return
	}
	num := list.NumFields()
	upperLimit := getFuncConfig(kind)
	if num <= upperLimit || upperLimit == 0 {
		return
	}
	desc := fmt.Sprintf("func:%s %s number:%d beyond limit:%d",
		name, kind, num, upperLimit)
	problems.Add(pos, desc, `names.file`)
}

func checkShortVarDecl(as *ast.AssignStmt, file *token.File) {
	if as.Tok != token.DEFINE {
		return
	}
	for _, exp := range as.Lhs {
		if ident, ok := exp.(*ast.Ident); ok {
			CheckIdent(ident, true, file, ``)
		}
	}
}

func checkRangeStmt(rg *ast.RangeStmt, file *token.File) {
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
func checkGenDecl(decl *ast.GenDecl, local bool, file *token.File) {
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
func (nw *NameWalker) isLocal(pos token.Pos) bool {
	position := nw.srcFile.Position(pos)
	if position.Line > nw.begin && position.Line < nw.end {
		return true
	}
	return false
}

func checkStruct(st *ast.StructType, local bool, file *token.File) {
	checkFieldList(st.Fields, local, file, `struct field`)
}
func checkInterface(itfc *ast.InterfaceType, local bool, file *token.File) {
	checkFieldList(itfc.Methods, local, file, `interface method`)
}
func checkFuncLit(fun *ast.FuncLit, file *token.File) {
	checkFieldList(fun.Type.Params, true, file, `func param`)
	checkFieldList(fun.Type.Results, true, file, `func result`)
	checkParamNum("", fun.Type.Params, `param`,
		file.Position(fun.Pos()))
	checkResultNum("", fun.Type.Params, `param`,
		file.Position(fun.Pos()))

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

	if desc := checkName(ident.Name, cfg, false); desc != `` {
		problems.Add(file.Position(ident.Pos()),
			fmt.Sprintf(`%s name %s %s`, getThing(thing, local, kind), ident.Name, desc), rule,
		)
	}
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

func getThing(thing string, local bool, kind string) string {
	if thing == `` {
		if local {
			thing = `local `
		}
		thing += kind
	}
	return thing
}
