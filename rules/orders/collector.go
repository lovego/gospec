package orders

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"strings"
)

type SymbolCollector struct {
	srcFile       *token.File
	importedTable []string
	funcTable     []*funcInfo
	declFuncTable map[string]bool
	declTable     []*ast.GenDecl
}

func NewCollector(srcFile *token.File) *SymbolCollector {
	return &SymbolCollector{
		srcFile:       srcFile,
		declFuncTable: make(map[string]bool),
	}
}

func (sc *SymbolCollector) collectImported(node ast.Node) {
	v, _ := node.(*ast.ImportSpec)
	name := strings.TrimSuffix(strings.TrimPrefix(v.Path.Value, "\""), "\"")
	name = path.Base(name)
	sc.importedTable = append(sc.importedTable, name)
}

func (sc *SymbolCollector) collectDecls(node ast.Node) {
	if sc.isLocal(node) {
		return
	}
	declTable := sc.declTable
	decl, _ := node.(*ast.GenDecl)
	if decl.Tok == token.IMPORT {
		return
	}
	sc.declTable = append(declTable, decl)
}

func (sc *SymbolCollector) isLocal(node ast.Node) bool {
	num := len(sc.funcTable)
	if num <= 0 {
		return false
	}
	lastFunc := sc.funcTable[num-1]
	position := sc.srcFile.Position(node.Pos())
	if position.Line > lastFunc.begin && position.Line < lastFunc.end {
		return true
	}
	return false
}

func (sc *SymbolCollector) collectFuncDecl(node ast.Node) {
	v, _ := node.(*ast.FuncDecl)
	f := &funcInfo{}
	getReceiver(f, v)
	f.funcName = v.Name.Name
	completeName := getCompleteName(f)
	if _, ok := sc.declFuncTable[completeName]; ok {
		return
	}
	sc.declFuncTable[completeName] = true
	f.called = make(map[string]map[string]bool)
	f.calledFuncs = make([]string, 0)
	f.funcDecl = v
	getFuncPos(v, f, sc.srcFile)
	sc.funcTable = append(sc.funcTable, f)
}

func (sc *SymbolCollector) collectCallee(f *funcInfo) {
	if isInternalFunc(f) || sc.isImported(f.packageName) {
		return
	}
	//TODO: figure out how this happen
	/*
		if f.packageName == "" && f.funcName == "" {
			return
		}
	*/
	if f.funcName == "" {
		return
	}

	for _, fn := range sc.funcTable {
		if fn.begin < f.begin && fn.end > f.begin {
			if ok := addFuncName(fn.called, f); !ok {
				return
			}
			fn.calledFuncs = append(fn.calledFuncs, getCompleteName(f))
			break
		}
	}
}

func isInternalFunc(f *funcInfo) bool {
	if f.packageName != "" {
		return false
	}
	internalFuncTable := []string{
		"len",
		"append",
		"new",
		"make",
		"copy",
		"delete",
	}
	for _, name := range internalFuncTable {
		if name == f.funcName {
			return true
		}
	}
	return false
}

func (sc *SymbolCollector) isImported(packageName string) bool {
	for _, name := range sc.importedTable {
		if name == packageName {
			return true
		}
	}
	return false
}

func getDeclName(decl *ast.GenDecl) string {
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			return getIdentName(s.Name)
		case *ast.ValueSpec:
			return getIdentName(s.Names[0])
		}
	}
	return ""
}

func getIdentName(ident *ast.Ident) string {
	if ident == nil {
		return ""
	}
	return ident.Name
}

func getFuncPos(v *ast.FuncDecl, f *funcInfo, file *token.File) {
	startLine := v.Body.Lbrace
	endLine := v.Body.Rbrace
	f.begin = file.Position(startLine).Line
	f.end = file.Position(endLine).Line
}

func getIdentPos(v *ast.Ident, f *funcInfo, file *token.File) {
	pos := file.Position(v.NamePos).Line
	f.begin = pos
}

func newFunc(f *funcInfo, packName string, ident *ast.Ident, file *token.File) {
	f.funcName = getIdentName(ident)
	getIdentPos(ident, f, file)
}

func getReceiver(fn *funcInfo, v *ast.FuncDecl) {
	if v.Recv == nil {
		return
	}
	for _, field := range v.Recv.List {
		for _, ident := range field.Names {
			fn.packageName = ident.Name
			return
		}
	}
}

func getCompleteName(f *funcInfo) string {
	return fmt.Sprintf("%s:%s", f.packageName, f.funcName)
}

func addFuncName(table map[string]map[string]bool, f *funcInfo) bool {
	calleeTable, ok := table[f.packageName]
	if !ok {
		calleeTable = make(map[string]bool)
		table[f.packageName] = calleeTable
	}
	_, ok = calleeTable[f.funcName]
	if ok {
		return false
	}
	calleeTable[f.funcName] = true
	return true
}

func (sc *SymbolCollector) getValidCallee(declFunc *funcInfo) {
	calledTable := declFunc.called
	for receiver, names := range calledTable {
		validCallees := sc.filterCallee(receiver, names)
		if len(validCallees) == 0 {
			delete(calledTable, receiver)
		}
		calledTable[receiver] = validCallees
	}
	updateCalled(declFunc)
}

func (sc *SymbolCollector) filterInvalidCallee() {
	for _, f := range sc.funcTable {
		sc.getValidCallee(f)
	}
}

func (sc *SymbolCollector) filterCallee(receiver string, callees map[string]bool) map[string]bool {
	validCallees := make(map[string]bool)
	for calleeName, _ := range callees {
		for _, fn := range sc.funcTable {
			if receiver == fn.packageName && fn.funcName == calleeName {
				validCallees[calleeName] = true
			}
		}
	}
	return validCallees
}

func updateCalled(fn *funcInfo) {
	calledFuncs := make([]string, 0)
	for _, completeName := range fn.calledFuncs {
		names := strings.Split(completeName, ":")
		packageName := names[0]
		callees, ok := fn.called[packageName]
		if !ok {
			continue
		}
		if _, ok := callees[names[1]]; !ok {
			continue
		}
		calledFuncs = append(calledFuncs, fmt.Sprintf("%s:%s", packageName, names[1]))
	}
	fn.calledFuncs = calledFuncs
}

func (sc *SymbolCollector) getCallees(funcName string) []string {
	sc.filterInvalidCallee()
	for _, f := range sc.funcTable {
		if f.funcName == funcName {
			return f.calledFuncs
		}
	}
	return nil
}

func (sc *SymbolCollector) getFuncNames() []string {
	names := make([]string, 0, len(sc.funcTable))
	for _, f := range sc.funcTable {
		names = append(names, f.funcName)
	}
	return names
}

func (sc *SymbolCollector) getDeclNames() []string {
	names := make([]string, 0, len(sc.declTable))
	for _, d := range sc.declTable {
		names = append(names, getDeclName(d))
	}
	return names
}
