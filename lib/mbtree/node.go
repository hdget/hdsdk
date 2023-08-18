package mbtree

import (
	"fmt"
)

type Node struct {
	// node id
	Id int64
	// node tag
	Tag string
	// parent id of the node
	Pid int64
	// child id list
	ChildIds []int64
	// node data
	Data interface{}
}

const defaultRootId = 0

// NewNode new node
func NewNode(id int64, data interface{}) *Node {
	return &Node{
		Id:       id,
		Pid:      0,
		ChildIds: make([]int64, 0),
		Data:     data,
		Tag:      buildNodeTag(id),
	}
}

// NewRootNode new root node
func NewRootNode(data interface{}) *Node {
	return NewNode(defaultRootId, data)
}

// IsLeaf return true if current node has no children
func (n *Node) IsLeaf() bool {
	return len(n.ChildIds) == 0
}

// HasChildren return true if current node has children
func (n *Node) HasChildren() bool {
	return !n.IsLeaf()
}

// updateChildren update the children list with different modes:
// addition (Node.ADD or Node.INSERT) and deletion (Node.DELETE).
func (n *Node) updateChildren(childId int64, action Action) {
	// ignore invalid childId
	if childId <= 0 {
		return
	}

	switch action {
	case ADD:
		n.ChildIds = append(n.ChildIds, childId)
	case DELETE:
		for i, id := range n.ChildIds {
			if id == childId {
				n.ChildIds = append(n.ChildIds[:i], n.ChildIds[i+1:]...)
			}
		}
	}
}

// updateParent Update the parent id
func (n *Node) updateParent(id int64) {
	n.Pid = id
}

// buildNodeTag 构造默认的节点tag
func buildNodeTag(id int64) string {
	return fmt.Sprintf("%d", id)
}
