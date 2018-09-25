package orders

import (
//"fmt"
)

type treeNode struct {
	child    **treeNode
	slibling **treeNode
	parent   *treeNode
	value    string
	height   int
}

type Tree struct {
	root *treeNode
	keys map[string]**treeNode
	size int
}

func NewTree() *Tree {
	return &Tree{
		keys: make(map[string]**treeNode),
	}
}

func (t *Tree) Size() int {
	return t.size
}
func (t *Tree) Clone() *Tree {
	if t.root == nil {
		return nil
	}
	newTree := NewTree()
	next := newIter(t.root)
	newRoot := &(newTree.root)
	newNode := &newRoot
	for node := next(); node != nil; node = next() {
		tmp := cloneNode(node)
		newTree.keys[tmp.value] = &tmp
		*newNode = &tmp
		if tmp.child != nil {
			newNode = &(tmp.child)
			continue
		}
		newNode = &(tmp.slibling)
	}
	newTree.root = *newRoot
	newTree.size = t.size
	return newTree
}

func (t *Tree) canMerge(src *Tree) bool {
	srcRoot := src.root
	if src == nil || srcRoot == nil {
		return false
	}
	if srcRoot.value == t.root.value {
		return false
	}
	_, ok := t.keys[srcRoot.value]
	if !ok {
		return false
	}
	return true
}

func (t *Tree) Merge(src *Tree) bool {
	srcRoot := src.root
	if !t.canMerge(src) {
		return false
	}
	if src.size == 1 {
		return true
	}
	n := t.keys[srcRoot.value]
	(*n).child = srcRoot.child
	for k, v := range src.keys {
		if k == srcRoot.value {
			continue
		}
		if _, ok := t.keys[k]; !ok {
			t.keys[k] = v
			continue
		}
		//find same node. and remove it
		if (*v).slibling == nil {
			*v = nil
			continue
		}
		*v = *((*v).slibling)
		//src.keys[k] = v
	}
	return true
}

//Keys must be unique
func (t *Tree) Add(rootKey string, keys []string) {
	t.root = &treeNode{value: rootKey}
	t.keys[rootKey] = &(t.root)
	t.size++
	t.addChildren(keys)
}

func (t *Tree) addChildren(keys []string) {
	num := len(keys)
	if num == 0 {
		return
	}
	next := &(t.root.child)
	for _, key := range keys {
		if _, ok := t.keys[key]; ok {
			continue
		}
		node := &treeNode{value: key, parent: t.root}
		*next = &node
		next = &((**next).slibling)
		t.size++
		t.keys[key] = &node
	}
}

//root must be tree's root node
func newIter(root *treeNode) func() *treeNode {
	queue := []*treeNode{root}
	return func() *treeNode {
		if len(queue) == 0 {
			return nil
		}
		node := queue[0]
		queue = append(queue[:0], queue[1:]...)
		child := node.child
		for child != nil && *child != nil {
			(*child).height = node.height + 1
			queue = append(queue, *child)
			child = (*child).slibling
		}
		return node
	}
}

func (t *Tree) Traverse() []string {
	if t.root == nil {
		return nil
	}
	keys := make([]string, 0, t.size)
	next := newIter(t.root)
	for node := next(); node != nil; node = next() {
		keys = append(keys, node.value)
	}
	return keys
}

func cloneNode(srcNode *treeNode) *treeNode {
	if srcNode == nil {
		return nil
	}
	node := new(treeNode)
	*node = *srcNode
	return node
}
