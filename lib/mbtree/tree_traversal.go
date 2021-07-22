package mbtree

// RSearch traverse back from the nid node to its ancestors (until root)
// @param: it supports pass one filter func to check the node
func (t *SafeMultiBranchTree) RSearch(id int64, args ...FilterFunc) chan int64 {
	// IMPORTANT: as non-buffered channel will be blocked if receiver quit
	// here use buffered-channel, the channel buffer number is SafeMultiBranchTree.MaxTreeLevel
	// default is 100, which should be enough to accommodate as much nodes as possible
	c := make(chan int64, t.MaxTreeLevel)

	go func() {
		// if iterate to root node, then need quit
		if id == t.RootId {
			close(c)
			return
		}

		node := t.GetNode(id)
		if node == nil {
			close(c)
			return
		}

		// 初始化current
		currentId := id
		currentNode := node

		filter := getFilterFromArgs(args...)
		for {
			// 如果遍历到根节点，需要停止
			if currentId == t.RootId {
				c <- t.RootId
				break
			}

			// 如果没有过滤函数，直接返回当前值
			if filter(currentNode) {
				c <- currentId
			}

			// 如果没有到根节点，遍历到根节点
			currentId = currentNode.Pid
			currentNode = t.GetNode(currentId)
			// 如果没有找到节点也需要停止
			if currentNode == nil {
				break
			}
		} // END for

		// 关闭channel
		close(c)
	}()

	return c
}

// DepthFirstTraversal 深度优先遍历
func (t *SafeMultiBranchTree) DepthFirstTraversal(nid int64, args ...FilterFunc) chan int64 {
	return t.traversal(nid, DEPTH, args...)
}

// WidthFirstTraversal 广度优先遍历
func (t *SafeMultiBranchTree) WidthFirstTraversal(nid int64, args ...FilterFunc) chan int64 {
	return t.traversal(nid, WIDTH, args...)
}

// traversal traverse the tree (or a subtree) with optional node filtering and sorting.
func (t *SafeMultiBranchTree) traversal(nid int64, mode TraversalMode, args ...FilterFunc) chan int64 {
	c := make(chan int64)
	go func() {
		node := t.GetNode(nid)
		if node == nil {
			close(c)
			return
		}

		filter := getFilterFromArgs(args...)

		// 满足条件的当前节点的id返回
		if filter(node) {
			c <- node.Id
		}

		// 过滤子节点
		queue := t.filterNodesWithin(node.ChildIds, filter)
		for {
			// 如果无要遍历的节点了，break
			if len(queue) == 0 {
				break
			}
			// 将需遍历的子节点的第一个取出来
			c <- queue[0].Id
			// 过滤该节点的子节点
			expansion := t.filterNodesWithin(queue[0].ChildIds, filter)
			// 如果深度优先，重新组装queue
			if mode == DEPTH {
				queue = append(expansion, queue[1:]...)
			} else if mode == WIDTH {
				queue = append(queue[1:], expansion...)
			}
		}
		close(c)
	}()

	return c
}
