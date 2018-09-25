package orders

import (
	"go/ast"
)

const (
	stateNone = iota
	stateCallBegin
	stateSelector
	statePackage
	stateFunc
	stateCallEnd
)

type MyVisiter struct {
	state    int
	lastFunc *funcInfo
}

func (mv *MyVisiter) Visit(node ast.Node, sc *SymbolCollector) {
	switch node.(type) {
	case *ast.ImportSpec:
		sc.collectImported(node)
	case *ast.FuncDecl:
		sc.collectFuncDecl(node)
	case *ast.GenDecl:
		sc.collectDecls(node)
	case *ast.CallExpr:
		mv.state = stateCallBegin
		mv.lastFunc = &funcInfo{}

	case *ast.SelectorExpr:
		switch mv.state {
		case stateCallBegin:
			mv.state = statePackage
		default:
			mv.state = stateNone
		}

	case *ast.Ident:
		ident, _ := node.(*ast.Ident)
		mv.handleIdent(ident, sc)
	default:
	}
}

func (mv *MyVisiter) handleIdent(node *ast.Ident, sc *SymbolCollector) {
	f := mv.lastFunc
	switch mv.state {
	case stateCallBegin:
		mv.state = stateNone
		newFunc(f, "", node, sc.srcFile)
		sc.collectCallee(f)
	case statePackage:
		f.packageName = getIdentName(node)
		mv.state = stateFunc

	case stateFunc:
		f.funcName = getIdentName(node)
		mv.state = stateCallEnd
		getIdentPos(node, f, sc.srcFile)
		sc.collectCallee(f)
	default:
	}
}
