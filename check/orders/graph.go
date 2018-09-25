package orders

import (
	"fmt"
)

type node struct {
	name  string
	next  *node
	value interface{}
	index int
}

type Graph struct {
	nodes []*node
}

func NewGraph() *Graph {
	return &Graph{}
}

func (pg *Graph) find(id string) *node {
	for _, n := range pg.nodes {
		if n.name == id {
			return n
		}
	}
	return nil
}

func (pg *Graph) AddNode(id string, value interface{}) {
	n := pg.find(id)
	if n != nil {
		return
	}
	num := len(pg.nodes)
	pg.nodes = append(pg.nodes, &node{name: id, value: value, index: num})
}

func (pg *Graph) AddEdges(startNodeId string, nodes map[string]interface{}) {
	n := pg.find(startNodeId)
	if n == nil {
		return
	}
	for _, v := range nodes {
		id := v.(string)
		tmpNode := &node{name: id, value: v}
		tmpNode.next = n.next
		n.next = tmpNode
	}
}

func (pg *Graph) findIncomeEdge(idx int, dst *node,
	sortedNodes map[string]bool) *node {
	for j, pn := range pg.nodes {
		if idx == j {
			continue
		}
		next := pn.next
		for next != nil {
			isSorted := sortedNodes[pn.name]
			if next.name == dst.name && !isSorted {
				return nil
			}
			next = next.next
		}
	}
	return dst
}

func (pg *Graph) getUnSorted() map[string]bool {
	sortedNodes := make(map[string]bool, len(pg.nodes))
	for _, node := range pg.nodes {
		sortedNodes[node.name] = false
	}
	return sortedNodes
}

//topology sort
func (pg *Graph) TopoSort() []string {
	ids := make([]string, 0, len(pg.nodes))
	sortedNodes := pg.getUnSorted()
	for j := 0; j < len(pg.nodes); j++ {
		for i, n := range pg.nodes {
			isSorted := sortedNodes[n.name]
			if isSorted {
				continue
			}
			in := pg.findIncomeEdge(i, n, sortedNodes)
			if in == nil {
				continue
			}
			ids = append(ids, in.name)
			sortedNodes[n.name] = true
		}
	}
	return ids
}

func (pg *Graph) Dump() {
	for _, n := range pg.nodes {
		fmt.Printf("node:%s\n\t", n.name)
		next := n.next
		for next != nil {
			fmt.Printf("[%s]-->", next.name)
			next = next.next
		}
		fmt.Printf("\n")
	}
}
