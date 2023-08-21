// Package mbtree
// 树结构相关方法
package mbtree

import (
	"sync"
)

// SafeMultiBranchTree 并发安全的多叉树
type SafeMultiBranchTree struct {
	sync.Mutex
	// 根节点ID
	RootId int64
	// 并发安全的map
	Nodes sync.Map
	// the maximum tree level
	// RSearch() use this value to traverse back how many levels
	MaxTreeLevel int
}

type Action int

const (
	ADD Action = iota
	DELETE
)

type TraversalMode int

const (
	// DEPTH 深度优先遍历算法
	DEPTH TraversalMode = iota
	// WIDTH 广度优先遍历算法
	WIDTH
)

// defaultTreeLevel the maximum tree level RSearch() will traverse back
const defaultTreeLevel = 100

// FilterFunc filter function to traverse and filter nodes
// pass a node as an argument,
// if function returns true then the node is matched and will be kept
type FilterFunc func(*Node) bool

func NewTree(rootNode *Node, args ...int) *SafeMultiBranchTree {
	if rootNode == nil {
		return nil
	}

	// max level
	maxLevel := defaultTreeLevel
	if len(args) > 0 {
		maxLevel = args[0]
		if maxLevel < defaultTreeLevel {
			maxLevel = defaultTreeLevel
		}
	}

	// new tree
	t := &SafeMultiBranchTree{
		RootId:       rootNode.Id,
		MaxTreeLevel: maxLevel,
	}
	// save root node
	t.Nodes.Store(rootNode.Id, rootNode)
	return t
}

func getFilterFromArgs(args ...FilterFunc) func(*Node) bool {
	var filter func(*Node) bool
	if len(args) > 0 {
		filter = args[0]
	} else {
		return filterAlwaysTrue
	}
	return filter
}

func getIdFromArgs(args ...int64) (int64, bool) {
	if len(args) > 0 {
		return args[0], true
	}
	return -1, false
}

// 定义一个总是返回true的filter函数
func filterAlwaysTrue(_ *Node) bool {
	return true
}
