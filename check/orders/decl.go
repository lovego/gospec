package orders

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"

	"github.com/lovego/spec/problems"
)

var errorTable = map[token.Token]map[token.Token]string{
	token.CONST: {token.TYPE: " type should after const.\n",
		token.VAR: " var should after const\n"},
	token.VAR: {token.CONST: " const should before var\n",
		token.TYPE: " type should after var\n"},
	token.TYPE: {token.CONST: " const should before type\n",
		token.VAR: " var should berfore type\n"},
}

//1:const 2: var 3: type 4: function
func checkDeclOrder(sc *SymbolCollector) {
	srcFile := sc.srcFile
	declTable := sc.declTable
	funcTable := sc.funcTable
	newDeclTable := getExpectedOrderTable(declTable)
	for i, d := range newDeclTable {
		oldDecl := declTable[i]
		if d.Tok == declTable[i].Tok {
			continue
		}
		desc := fmt.Sprintf("file %s, %s %s", path.Base(srcFile.Name()),
			d.Tok, getDeclName(d))
		problems.Add(srcFile.Position(d.Pos()), desc, errorTable[d.Tok][oldDecl.Tok])
	}
	if len(funcTable) <= 0 {
		return
	}
	firstFunc := funcTable[0]
	for _, d := range declTable {
		pos := srcFile.Position(d.Pos())
		if pos.Line > firstFunc.end {
			desc := fmt.Sprintf("file %s, %s %s", path.Base(srcFile.Name()), d.Tok,
				getDeclName(d))
			problems.Add(pos, desc, "function should after const, var,type declaration")
		}
	}
	declTable = make([]*ast.GenDecl, 0)
}

func getExpectedOrderTable(declTable []*ast.GenDecl) []*ast.GenDecl {
	num := len(declTable)
	newDeclTable := make([]*ast.GenDecl, num)
	copy(newDeclTable, declTable)
	constIdx := 0
	typeIdx := num - 1
	for i := 0; i < num; i++ {
		switch newDeclTable[i].Tok {
		case token.CONST:
			if constIdx == i {
				break
			}
			newDeclTable[constIdx], newDeclTable[i] = newDeclTable[i],
				newDeclTable[constIdx]
			constIdx++
		case token.TYPE:
			if typeIdx <= i {
				break
			}
			newDeclTable[typeIdx], newDeclTable[i] = newDeclTable[i],
				newDeclTable[typeIdx]
			typeIdx--
		default:
		}
	}
	return newDeclTable
}
