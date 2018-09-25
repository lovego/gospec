package orders

import (
	"fmt"
	"testing"
)

func createTree() *Tree {
	tree := NewTree()
	tree.Add("A", []string{"B", "C", "D"})
	return tree
}

func TestAddNode(t *testing.T) {
	tree := createTree()
	nodes := tree.Traverse()
	expected := []string{"A", "B", "C", "D"}
	checkResult(expected, nodes, t)
}

func checkResult(expected, result []string, t *testing.T) {
	if len(expected) != len(result) {
		t.Fatalf("expected:%s got:%s \n", expected, result)
	}
	for i, r := range expected {
		if r != result[i] {
			t.Fatalf("expected:%s got:%s   got nodes:%v\n", expected[i], r, result)
		}
	}
}

func TestMerge(t *testing.T) {
	tree := createTree()
	tree2 := NewTree()
	tree2.Add("B", []string{"f", "D", "k", "z"})
	tree.Merge(tree2)
	nodes := tree.Traverse()
	expected := []string{"A", "B", "C", "D", "f", "k", "z"}
	checkResult(expected, nodes, t)

	tree = createTree()
	tree2 = NewTree()
	tree2.Add("B", []string{"f", "g", "k", "D"})
	tree.Merge(tree2)
	nodes = tree.Traverse()
	expected = []string{"A", "B", "C", "D", "f", "g", "k"}
	checkResult(expected, nodes, t)

	tree = createTree()
	tree2 = NewTree()
	tree2.Add("B", []string{"D", "g", "k", "f"})
	tree.Merge(tree2)
	nodes = tree.Traverse()
	expected = []string{"A", "B", "C", "D", "g", "k", "f"}
	checkResult(expected, nodes, t)

	//case2: merge twice
	tree2 = NewTree()
	tree2.Add("k", []string{"x", "z"})
	tree.Merge(tree2)
	nodes = tree.Traverse()
	expected = []string{"A", "B", "C", "D", "g", "k", "f", "x", "z"}
	checkResult(expected, nodes, t)

	//case 3:
	fmt.Printf("case 3======================\n")
	tree = NewTree()
	tree.Add("A", []string{"B", "C"})

	tree2 = NewTree()
	tree2.Add("B", []string{"D"})
	tree.Merge(tree2.Clone())

	tree3 := NewTree()
	tree3.Add("C", []string{"B"})
	tree3.Merge(tree2.Clone())

	tree.Merge(tree3.Clone())
	nodes = tree.Traverse()
	expected = []string{"A", "B", "C", "D"}
	checkResult(expected, nodes, t)

	//case 4:
	tree = createTree()
	tree2 = NewTree()
	tree2.Add("B", []string{"D"})
	tree.Merge(tree2)
	nodes = tree.Traverse()
	expected = []string{"A", "B", "C", "D"}
	checkResult(expected, nodes, t)
}

func TestClone(t *testing.T) {
	tree := NewTree()
	tree.Add("A", []string{"B", "C", "D"})

	clonedTree := tree.Clone()
	nodes := clonedTree.Traverse()
	expected := []string{"A", "B", "C", "D"}
	checkResult(expected, nodes, t)

	tree2 := NewTree()
	tree2.Add("B", []string{"f", "g", "k", "D"})
	clonedTree.Merge(tree2)
	nodes = clonedTree.Traverse()
	expected = []string{"A", "B", "C", "D", "f", "g", "k"}
	checkResult(expected, nodes, t)

	expected = []string{"A", "B", "C", "D"}
	nodes = tree.Traverse()
	checkResult(expected, nodes, t)
}
