package orders

import (
	"testing"
)

func TestGraph(t *testing.T) {
	g := NewGraph()
	g.AddNode("1", "1")
	g.AddNode("2", "2")
	g.AddNode("3", "3")
	g.AddNode("4", "4")
	g.AddEdges("1", map[string]interface{}{"2": "2"})
	g.AddEdges("2", map[string]interface{}{"4": "4"})
	g.AddEdges("3", map[string]interface{}{"2": "2"})
	nodes := g.TopoSort()
	expected := []string{"1", "3", "2", "4"}
	checkResult(expected, nodes, t)

	g.AddEdges("4", map[string]interface{}{"2": "2"})
	nodes = g.TopoSort()
	checkGraphResult(expected[:2], nodes, t)
}

func checkGraphResult(expected, result []string, t *testing.T) {
	if len(result) != len(expected) {
		t.Fatalf("expected:%v got:%v\n", expected, result)
	}
	for i, k := range expected {
		if k != result[i] {
			t.Fatalf("expected:%s got:%s  result:%v\n", expected[i], k, result)
		}
	}
}
