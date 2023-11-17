package mbtree

import (
	"github.com/hdget/hdsdk/hdutils"
)

// GetParentId 获取父节点的ID, 如果没找到，返回-1
func (t *SafeMultiBranchTree) GetParentId(id int64) int64 {
	node := t.GetNode(id)
	if node == nil {
		return -1
	}
	return node.Pid
}

// GetChildIds return the children nid list of specified id.
// empty slice is returned if nid does not exist
func (t *SafeMultiBranchTree) GetChildIds(id int64) []int64 {
	node := t.GetNode(id)
	if node == nil {
		return nil
	}
	return node.ChildIds
}

// GetPaths 获取某个指定Node的路径，从根到叶节点
func (t *SafeMultiBranchTree) GetPaths(id int64) [][]int64 {
	paths := make([][]int64, 0)
	for _, leafNode := range t.GetLeafNodes() {
		pathIds := make([]int64, 0)
		if t.IsAncestor(id, leafNode.Id) {
			for id := range t.RSearch(leafNode.Id) {
				pathIds = append(pathIds, id)
			}
		}
		// 倒序
		if len(pathIds) > 0 {
			reversedPathIds := hdutils.ReverseInt64Slice(pathIds)
			paths = append(paths, reversedPathIds)
		}
	}
	return paths
}

// GetDescendantIds get all descendant node ids, including myself
func (t *SafeMultiBranchTree) GetDescendantIds(id int64, args ...FilterFunc) []int64 {
	subtree := t.SubTree(id)
	if subtree == nil {
		return nil
	}

	nodeIds := make([]int64, 0)
	filter := getFilterFromArgs(args...)
	subtree.Nodes.Range(func(k, v interface{}) bool {
		id, ok := k.(int64)
		if !ok {
			return true
		}

		node, ok := v.(*Node)
		if !ok {
			return true
		}

		if filter(node) {
			nodeIds = append(nodeIds, id)
		}
		return true
	})
	return nodeIds
}
