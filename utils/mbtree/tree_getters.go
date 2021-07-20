package mbtree

// SubTree return a shallow COPY of subtree with nid being the new root.
func (t *SafeMultiBranchTree) SubTree(id int64) *SafeMultiBranchTree {
	if !t.Contains(id) {
		return nil
	}

	st := NewTree()
	st.RootId = id
	for nid := range t.DepthFirstTraversal(id) {
		if node := t.GetNode(nid); node != nil {
			st.Nodes.Store(nid, node)
		}
	}
	return st
}
