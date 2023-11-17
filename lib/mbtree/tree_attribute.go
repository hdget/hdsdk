package mbtree

import "github.com/hdget/hdsdk/hdutils"

// Size get how many nodes in the tree
func (t *SafeMultiBranchTree) Size() int {
	count := 0
	t.Nodes.Range(func(k, v interface{}) bool {
		count += 1
		return true
	})
	return count
}

// Level get specified node's level
// The level is an integer starting with 0 at the root.
// In other words, the root lives at level 0
// @id: node id
// @args: it can pass filter function to calculate level passing exclusive nodes.
func (t *SafeMultiBranchTree) Level(id int64, args ...FilterFunc) int {
	// if it is root node, just return 0
	if id == t.RootId {
		return 0
	}

	// begin from 1, if we can reverse back and not reach root then add 1
	count := 0
	for range t.RSearch(id, args...) {
		count += 1
	}
	return count - 1
}

// Depth get the maximum level of the tree or the level of the given node
// if specified the node id, then it get the depth of that given node
// if not specified the node id, it get the maximum level of the tree
// @param:  args it could be @nodeId
// @return: the depth of the tree or the depth of the @nodeId
func (t *SafeMultiBranchTree) Depth(args ...int64) int {
	depth := 0
	id, exist := getIdFromArgs(args...)
	// If not specified id, then get maximum level of the tree
	if !exist {
		leaves := t.GetLeafNodes()
		for _, leafNode := range leaves {
			level := t.Level(leafNode.Id)
			if level > depth {
				depth = level
			}
		}
		return depth
	}

	// if specified id, then get level of the given node
	if !t.Contains(id) {
		return 0
	}
	return t.Level(id)
}

// AllPaths use this function to get the identifiers allowing to go from the root nodes to each leaf.
// @return: a list of list of identifiers, root being not omitted.
// For example:
//
//	Harry
//	|___ Bill
//	|___ Jane
//	|    |___ Diane
//	|         |___ George
//	|              |___ Jill
//	|         |___ Mary
//	|    |___ Mark
//
// Expected result:
//
//	 [['harry', 'jane', 'diane', 'mary'],
//		['harry', 'jane', 'mark'],
//		['harry', 'jane', 'diane', 'george', 'jill'],
//		['harry', 'bill']]
func (t *SafeMultiBranchTree) AllPaths() [][]int64 {
	paths := make([][]int64, 0)
	for _, leafNode := range t.GetLeafNodes() {
		pathIds := make([]int64, 0)
		for id := range t.RSearch(leafNode.Id) {
			pathIds = append(pathIds, id)
		}
		// 倒序
		if len(pathIds) > 0 {
			reversedPathIds := hdutils.ReverseInt64Slice(pathIds)
			paths = append(paths, reversedPathIds)
		}
	}
	return paths
}

// SubTree return a shallow COPY of subtree with nid being the new root.
func (t *SafeMultiBranchTree) SubTree(id int64) *SafeMultiBranchTree {
	newRootNode := t.GetNode(id)
	if newRootNode == nil {
		return nil
	}

	st := NewTree(newRootNode)
	for nid := range t.DepthFirstTraversal(id) {
		if node := t.GetNode(nid); node != nil {
			st.Nodes.Store(nid, node)
		}
	}
	return st
}
