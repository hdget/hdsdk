package mbtree

// GetNode obtain node from the tree
func (t *SafeMultiBranchTree) GetNode(id int64) *Node {
	v, exist := t.Nodes.Load(id)
	if !exist {
		return nil
	}
	n, ok := v.(*Node)
	if !ok || n == nil {
		return nil
	}
	return n
}

// GetParentNode 获取父节点, 如果没找到，返回nil
func (t *SafeMultiBranchTree) GetParentNode(id int64) *Node {
	node := t.GetNode(id)
	if node == nil {
		return nil
	}
	return t.GetNode(node.Pid)
}

// GetRootNode 获取根节点
func (t *SafeMultiBranchTree) GetRootNode() *Node {
	return t.GetNode(t.RootId)
}

// GetAncestorNode for a given id, get ancestor node object has distance to @nodeId
// @param: specified node id
// @param: distance value must > 1 and <= current level,
// the parent distance is 1
// the grandparent distance is 2
// the great-grandfather distance is 3
func (t *SafeMultiBranchTree) GetAncestorNode(id int64, distance int) *Node {
	// no ancestor for root node
	if id == t.RootId {
		return nil
	}

	// if specified node is not found, no ancestor too
	node := t.GetNode(id)
	if node == nil {
		return nil
	}

	// if distance lesser than 1, return nil
	if distance < 1 {
		return nil
	}

	// here myLevel larger than 1, as root node has been checked before
	// distance must not larger than current level
	myLevel := t.Level(id)
	if distance > myLevel {
		return nil
	}

	var found *Node
	count := 1
	descendant := node
	ascendant := t.GetNode(descendant.Pid)
	for {
		if count == distance {
			found = ascendant
			break
		}
		descendant = ascendant
		ascendant = t.GetNode(descendant.Pid)
		count += 1
	}
	return found
}

// GetLeafNodes get leaf nodes of the tree,
// if specified node id, then get leaf node for that node
func (t *SafeMultiBranchTree) GetLeafNodes(args ...int64) []*Node {
	leafNodes := make([]*Node, 0)
	id, exist := getIdFromArgs(args...)

	// 如果没有传id参数，直接返回所有叶子节点
	if !exist {
		t.Nodes.Range(func(k, v interface{}) bool {
			if node, ok := v.(*Node); ok && node.IsLeaf() {
				leafNodes = append(leafNodes, node)
			}
			return true
		})
		return leafNodes
	}

	for nid := range t.DepthFirstTraversal(id) {
		if node := t.GetNode(nid); node != nil && node.IsLeaf() {
			leafNodes = append(leafNodes, node)
		}
	}
	return leafNodes
}

// GetSiblingNodes return the siblings of given @id.
// If @nid is root or there are no siblings, nil is returned.
func (t *SafeMultiBranchTree) GetSiblingNodes(id int64) []*Node {
	// there is no sibling nodes for root node
	if id == t.RootId {
		return nil
	}

	parentNode := t.GetParentNode(id)
	if parentNode == nil {
		return nil
	}

	siblingNodes := make([]*Node, 0)
	// 将指定id的父节点的children中不等于自己的节点全部加入
	for _, childId := range parentNode.ChildIds {
		// 只考虑不等于自己的节点
		if childId != id {
			// 只考虑找到的有效的节点
			if childNode := t.GetNode(childId); childNode != nil {
				siblingNodes = append(siblingNodes, childNode)
			}
		}
	}
	return siblingNodes
}

// GetAllNodes 获取所有树节点slice
func (t *SafeMultiBranchTree) GetAllNodes() []*Node {
	allNodes := make([]*Node, 0)
	t.Nodes.Range(func(k, v interface{}) bool {
		if node, ok := v.(*Node); ok {
			allNodes = append(allNodes, node)
		}
		return true
	})
	return allNodes
}

// GetChildNodes return the children node slice of specified node id
// empty slice is returned if no corresponding node exist for specified node id
func (t *SafeMultiBranchTree) GetChildNodes(id int64) []*Node {
	node := t.GetNode(id)
	if node == nil {
		return nil
	}

	childNodes := make([]*Node, 0)
	for _, childId := range node.ChildIds {
		if childNode := t.GetNode(childId); childNode != nil {
			childNodes = append(childNodes, childNode)
		}
	}
	return childNodes
}

// GetDescendantNodes get all descendant nodes, including myself
func (t *SafeMultiBranchTree) GetDescendantNodes(id int64, args ...FilterFunc) []*Node {
	subtree := t.SubTree(id)
	if subtree == nil {
		return nil
	}

	nodes := make([]*Node, 0)
	filter := getFilterFromArgs(args...)
	subtree.Nodes.Range(func(k, v interface{}) bool {
		if node, ok := v.(*Node); ok && filter(node) {
			nodes = append(nodes, node)
		}
		return true
	})
	return nodes
}
