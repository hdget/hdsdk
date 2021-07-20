package mbtree

import (
	"github.com/pkg/errors"
)

// CreateNode create a child node for given @pid node.
func (t *SafeMultiBranchTree) CreateNode(id int64, pid int64, data interface{}) (*Node, error) {
	node := NewNode(id, data)
	if err := t.mountTo(node, pid); err != nil {
		return nil, err
	}
	return node, nil
}

// MoveNode move node @id to be a child of parent @destPid
func (t *SafeMultiBranchTree) MoveNode(id int64, destPid int64) error {
	t.Lock()
	defer t.Unlock()

	node := t.GetNode(id)
	if node == nil {
		return errors.Wrapf(ErrNodeNotFound, "MoveNode() nodeId: %d", id)
	}

	sourceParentNode := t.GetNode(node.Pid)
	if sourceParentNode == nil {
		return errors.Wrapf(ErrNodeNotFound, "MoveNode() sourcePid: %d", node.Pid)
	}

	destParentNode := t.GetNode(destPid)
	if destParentNode == nil {
		return errors.Wrapf(ErrNodeNotFound, "MoveNode() sourcePid: %d", destPid)
	}

	// if destination parent node is not root node,
	// node should not be the ancestor of destination node
	if destPid != t.RootId && t.IsAncestor(id, destPid) {
		return ErrInvalidSourceDest
	}

	{
		// 1. delete id from source parent node's children
		sourceParentNode.updateChildren(id, DELETE)
		// 2. add id to destination parent node's children
		destParentNode.updateChildren(id, ADD)
		// 3. change current node's parent to destPid
		node.updateParent(destPid)
	}
	return nil
}

// DeleteNode delete a node by linking past it.
// For example, if we have `a -> b -> ccc` and delete node b,
// we are left with `a -> ccc`.
func (t *SafeMultiBranchTree) DeleteNode(id int64) error {
	// 同一时刻只允许一个goroutine进入
	t.Lock()
	defer t.Unlock()

	node := t.GetNode(id)
	if node == nil {
		return errors.Wrapf(ErrNodeNotFound, "DeleteNode() nodeId: %d", id)
	}

	if id == t.RootId {
		return errors.Wrap(ErrDeleteRootForbidden, "DeleteNode()")
	}

	// get parent node
	parentNode := t.GetNode(node.Pid)
	if parentNode == nil {
		return errors.Wrapf(ErrNodeNotFound, "DeleteNode() parentId: %d", node.Pid)
	}

	{
		// 1. let all my children points to the parent node
		for _, childNode := range t.GetChildNodes(id) {
			childNode.updateParent(parentNode.Id)
		}
		// 2. add my children to parent's children
		parentNode.ChildIds = append(parentNode.ChildIds, node.ChildIds...)
		// 3. delete me from parent's children
		parentNode.updateChildren(id, DELETE)
		// 4. delete from the map
		t.Nodes.Delete(id)
	}

	return nil
}

// mountTo 将node挂载在@pid节点下
func (t *SafeMultiBranchTree) mountTo(node *Node, pid int64) error {
	t.Lock()
	defer t.Unlock()

	if node == nil {
		return errors.Wrap(ErrInvalidNode, "mountTo()")
	}

	// node to be mounted should not has the same id with root's id
	if node.Id == t.RootId {
		return errors.Wrap(ErrInvalidNodeId, "mountTo()")
	}

	// check if node to be mounted already exists or not
	if exist := t.Contains(node.Id); exist {
		return errors.Wrapf(ErrNodeAlreadyExist, "mountTo() nodeId: %d", node.Id)
	}

	// check if pid already exists or not
	parentNode := t.GetNode(pid)
	if parentNode == nil {
		return errors.Wrapf(ErrNodeNotFound, "mountTo() parentId: %d", pid)
	}

	// 具体操作
	{
		// 1. 添加到节点map中
		t.Nodes.Store(node.Id, node)
		// 2. 更新当前节点的父节点
		node.updateParent(pid)
		// 3. 维护父节点下的children
		parentNode.updateChildren(node.Id, ADD)
	}

	return nil
}
