package mbtree

// Contains check if nid in tree or not
func (t *SafeMultiBranchTree) Contains(id int64) bool {
	_, exist := t.Nodes.Load(id)
	return exist
}

// IsAncestor check if the @ancestor is the preceding nodes of @grandchild
// @param: ancestorId, the ancestor node id
// @param: grandchild, the grandchild node id
// @return: true or false
func (t *SafeMultiBranchTree) IsAncestor(ancestorId int64, grandchildId int64) bool {
	for id := range t.RSearch(grandchildId) {
		if id == ancestorId {
			return true
		}
	}
	return false
}
