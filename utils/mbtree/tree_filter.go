package mbtree

// filterNodes traverse all nodes in the tree by using filter function
// @param f: is passed one node as an argument and that node is included if function returns true,
// @return []*Node
func (t *SafeMultiBranchTree) filterNodes(f FilterFunc) []*Node {
	filtered := make([]*Node, 0)
	t.Nodes.Range(func(k, v interface{}) bool {
		node := v.(*Node)
		// skip root node
		if node.Id != 0 && f(node) {
			filtered = append(filtered, node)
		}
		return true
	})
	return filtered
}

// filterNodesWithin filter nodes in specified id slice
func (t *SafeMultiBranchTree) filterNodesWithin(ids []int64, f FilterFunc) []*Node {
	filtered := make([]*Node, 0)
	for _, id := range ids {
		node := t.GetNode(id)
		if node == nil {
			continue
		}

		if f(node) {
			filtered = append(filtered, node)
		}
	}
	return filtered
}
