package orders

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"

	"github.com/lovego/gospec/problems"
)

type OrdersConfigT struct {
	FuncCheck bool
	VarCheck  bool
}

var Config = OrdersConfigT{
	FuncCheck: true,
	VarCheck:  true,
}

type FileChecker struct {
	myVisiter MyVisiter
	sc        *SymbolCollector
}

func NewFileChecker(srcFile *token.File) *FileChecker {
	return &FileChecker{
		sc: NewCollector(srcFile),
	}
	return nil
}

func (fc *FileChecker) Visit(node ast.Node) ast.Node {
	fc.myVisiter.Visit(node, fc.sc)
	return node
}

func (fc *FileChecker) Check() {
	sc := fc.sc
	if Config.VarCheck {
		checkDeclOrder(sc)
	}
	if !Config.FuncCheck {
		return
	}
	sc.filterInvalidCallee()
	checkFuncOrder(sc.funcTable, sc.srcFile)
}

//only check direct call relationship
func checkFuncOrder(funcTable []*funcInfo, srcFile *token.File) {
	expFunOrderTable := getExpectedFuncOrder(funcTable)
	for _, funcs := range expFunOrderTable {
		checkByCallingOrder(srcFile, funcs, funcTable)
	}
}

func checkByCallingOrder(file *token.File, funcs []string, funcTable []*funcInfo) {
	for i := 1; i < len(funcs); i++ {
		prevFun := lookupFunc(funcs[i-1], funcTable)
		fn := lookupFunc(funcs[i], funcTable)
		if fn == nil {
			continue
		}
		desc := fmt.Sprintf("file %s, function %s should after function %s",
			path.Base(file.Name()), fn.funcName, prevFun.funcName)
		if fn.end < prevFun.begin {
			problems.Add(file.Position(fn.funcDecl.Pos()),
				desc, "function declare disorder")
		}
	}
}

func getExpectedFuncOrder(funcTable []*funcInfo) [][]string {
	removeCircle(funcTable)
	callTrees := createCallTree(funcTable)
	callTrees = filterCallTree(callTrees)
	expFunOrderTable := make([][]string, 0)
	for _, calleeTree := range callTrees {
		expFunOrderTable = append(expFunOrderTable, calleeTree.Traverse())
	}
	return expFunOrderTable
}

func lookupFunc(name string, funcTable []*funcInfo) *funcInfo {
	for _, fn := range funcTable {
		completeName := getCompleteName(fn)
		if completeName == name {
			return fn
		}
	}
	return nil
}

func removeCircle(funcTable []*funcInfo) {
	graph := createGraph(funcTable)
	validFuns := graph.TopoSort()
	for i, fn := range funcTable {
		name := getCompleteName(fn)
		for _, v := range validFuns {
			if v == name {
				goto next
			}
		}
		funcTable = append(funcTable[:i], funcTable[i+1:]...)
	next:
	}
}

func createGraph(funcTable []*funcInfo) *Graph {
	graph := NewGraph()
	for _, fn := range funcTable {
		name := getCompleteName(fn)
		nodes := make(map[string]interface{})
		for _, callee := range fn.calledFuncs {
			nodes[callee] = callee
		}
		graph.AddNode(name, name)
		graph.AddEdges(name, nodes)
	}
	return graph
}

func createCallTree(funcs []*funcInfo) []*Tree {
	forest := make([]*Tree, 0, len(funcs))
	for _, fn := range funcs {
		tree := NewTree()
		rootId := fmt.Sprintf("%s:%s", fn.packageName, fn.funcName)
		tree.Add(rootId, fn.calledFuncs)
		forest = append(forest, tree)
	}
	return mergeTrees(forest)
}

func filterCallTree(trees []*Tree) []*Tree {
	forest := make([]*Tree, 0)
	for _, tree := range trees {
		if tree.Size() <= 2 {
			continue
		}
		forest = append(forest, tree)
	}
	return forest
}

//Two different call tree can merge to each other when A call B and B call A
func mergeTrees(forest []*Tree) []*Tree {
	treeNum := len(forest)
	trees := make([]*Tree, 0)
	for i := 0; i < treeNum; i++ {
		tree := forest[i]
		if num := tryMerge(i, tree, forest); num > 0 {
			continue
		}
		trees = append(trees, tree)
	}
	return trees
}

func tryMerge(me int, tree *Tree, forest []*Tree) int {
	//clonedTree := tree.Clone()
	mergedNum := 0
	treeNum := len(forest)
	for j := 0; j < treeNum; j++ {
		//TODO: Merge should avoid change the tree's internal state
		//fix: the clonedTree's interal state is changed
		//if j == me || !forest[j].Merge(clonedTree) {
		if j == me || !forest[j].Merge(tree.Clone()) {
			continue
		}
		mergedNum++
	}
	return mergedNum
}
