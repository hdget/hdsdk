package mbtree

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"sort"
	"testing"
)

type NodeData struct {
	Id       int64  `db,json:"id"`
	Name     string `db,json:"name"`
	Nickname string `db,json:"nickname"`
	Mobile   string `db,json:"mobile"`
	Level    int    `db,json:"level"`
}

var parentChildren = map[int][]int{
	0:  []int{1, 2, 3, 4},
	2:  []int{5, 6, 7},
	3:  []int{8, 9},
	4:  []int{10, 11},
	6:  []int{12, 13, 14, 15},
	13: []int{16},
}

var tree *SafeMultiBranchTree

func TestMain(m *testing.M) {
	tree = getTree()
	m.Run()
	os.Exit(0)
}

func TestNewTree(t *testing.T) {
	rootNode := NewRootNode(nil)
	tt := NewTree(rootNode)
	assert.Equal(t, tt.RootId, int64(0))
}

func TestSize(t *testing.T) {
	assert.Equal(t, tree.Size(), 17)

	rootNode := NewRootNode(nil)
	tt := NewTree(rootNode)
	assert.Equal(t, tt.Size(), 1)
}

func TestLevel(t *testing.T) {
	assert.Equal(t, tree.Level(0), 0)
	assert.Equal(t, tree.Level(15), 3)
	assert.Equal(t, tree.Level(13), 3)
	assert.Equal(t, tree.Level(16), 4)
}

func TestDepth(t *testing.T) {
	assert.Equal(t, tree.Depth(), 4)
	assert.Equal(t, tree.Depth(6), 2)
}

func TestPaths(t *testing.T) {
	allPaths := tree.AllPaths()
	assert.ElementsMatch(t, allPaths, [][]int64{
		[]int64{0, 1},
		[]int64{0, 2, 5}, []int64{0, 2, 6, 12}, []int64{0, 2, 6, 13, 16}, []int64{0, 2, 6, 14}, []int64{0, 2, 6, 15}, []int64{0, 2, 7},
		[]int64{0, 3, 8}, []int64{0, 3, 9},
		[]int64{0, 4, 10}, []int64{0, 4, 11},
	})
}

func TestSubTree(t *testing.T) {
	subtree1 := tree.SubTree(2)
	t.Logf("Subtree rootId now is: %d", subtree1.RootId)
	assert.ElementsMatch(t, subtree1.AllPaths(), [][]int64{
		[]int64{2, 6, 12}, []int64{2, 6, 14}, []int64{2, 5}, []int64{2, 6, 13, 16}, []int64{2, 7}, []int64{2, 6, 15}})

	subtree2 := tree.SubTree(0)
	t.Logf("Subtree rootId now is: %d", subtree2.RootId)
	assert.ElementsMatch(t, subtree2.AllPaths(), [][]int64{
		[]int64{0, 2, 6, 12}, []int64{0, 2, 6, 14}, []int64{0, 1}, []int64{0, 3, 8}, []int64{0, 3, 9},
		[]int64{0, 4, 11}, []int64{0, 2, 5}, []int64{0, 4, 10}, []int64{0, 2, 6, 13, 16}, []int64{0, 2, 7}, []int64{0, 2, 6, 15}})

	subtree3 := tree.SubTree(3)
	t.Logf("Subtree rootId now is: %d", subtree3.RootId)
	assert.ElementsMatch(t, subtree3.AllPaths(), [][]int64{[]int64{3, 8}, []int64{3, 9}})

	subtree4 := tree.SubTree(1)
	t.Logf("Subtree rootId now is: %d", subtree4.RootId)
	assert.ElementsMatch(t, subtree4.AllPaths(), [][]int64{})
}

func TestIsAncestor(t *testing.T) {
	assert.Equal(t, tree.IsAncestor(0, 2), true)
	assert.Equal(t, tree.IsAncestor(0, 15), true)
	assert.Equal(t, tree.IsAncestor(6, 16), true)
	assert.Equal(t, tree.IsAncestor(1, 13), false)
}

func TestContains(t *testing.T) {
	assert.Equal(t, tree.Contains(-1), false)
	assert.Equal(t, tree.Contains(3), true)
	assert.Equal(t, tree.Contains(18), false)
}

func TestFilterNodes(t *testing.T) {
	filter := func(n *Node) bool {
		nodeData := n.Data.(*NodeData)
		return nodeData.Level > 10
	}

	filteredIds := make([]int64, 0)
	for _, node := range tree.filterNodes(filter) {
		filteredIds = append(filteredIds, node.Id)
	}
	assert.ElementsMatch(t, filteredIds, []int64{11, 12, 13, 14, 15, 16})
}

func TestFilterNodesWithin(t *testing.T) {
	filter := func(n *Node) bool {
		nodeData := n.Data.(*NodeData)
		return nodeData.Level > 10
	}

	ids := []int64{1, 5, 10, 13, 15, 16}
	filteredIds := make([]int64, 0)
	for _, node := range tree.filterNodesWithin(ids, filter) {
		filteredIds = append(filteredIds, node.Id)
	}
	assert.ElementsMatch(t, filteredIds, []int64{13, 15, 16})
}

func TestGetParentId(t *testing.T) {
	assert.Equal(t, tree.GetParentId(-2), int64(-1))
	assert.Equal(t, tree.GetParentId(0), int64(0))
	assert.Equal(t, tree.GetParentId(2), int64(0))
	assert.Equal(t, tree.GetParentId(16), int64(13))
	assert.Equal(t, tree.GetParentId(999), int64(-1))
}

func TestGetChildIds(t *testing.T) {
	childIds1 := tree.GetChildIds(2)
	assert.ElementsMatch(t, childIds1, []int64{5, 6, 7})

	childIds2 := tree.GetChildIds(-1)
	assert.ElementsMatch(t, childIds2, nil)

	childIds3 := tree.GetChildIds(0)
	assert.ElementsMatch(t, childIds3, []int64{1, 2, 3, 4})
}

func TestGetPaths(t *testing.T) {
	paths1 := tree.GetPaths(2)
	assert.ElementsMatch(t, paths1, [][]int64{
		[]int64{0, 2, 6, 12}, []int64{0, 2, 6, 14}, []int64{0, 2, 5}, []int64{0, 2, 6, 13, 16}, []int64{0, 2, 7}, []int64{0, 2, 6, 15}})

	paths2 := tree.GetPaths(-1)
	assert.ElementsMatch(t, paths2, nil)

	paths3 := tree.GetPaths(0)
	assert.ElementsMatch(t, paths3, [][]int64{
		[]int64{0, 2, 6, 12}, []int64{0, 2, 6, 14}, []int64{0, 1}, []int64{0, 3, 8}, []int64{0, 3, 9},
		[]int64{0, 4, 11}, []int64{0, 2, 5}, []int64{0, 4, 10}, []int64{0, 2, 6, 13, 16}, []int64{0, 2, 7}, []int64{0, 2, 6, 15}})
}

func TestGetDescendantIds(t *testing.T) {
	descendantIds := make([]int64, 0)
	descendantIds = append(descendantIds, tree.GetDescendantIds(2)...)

	assert.ElementsMatch(t, descendantIds,
		[]int64{2, 6, 12, 14, 5, 13, 16, 7, 15},
	)
}

func TestGetNode(t *testing.T) {
	assert.Equal(t, tree.GetNode(2).Id, int64(2))
	assert.Equal(t, tree.GetNode(0).Id, int64(0))
	assert.Nil(t, tree.GetNode(999))
	assert.Nil(t, tree.GetNode(-1))
}

func TestGetParentNode(t *testing.T) {
	assert.Equal(t, tree.GetParentNode(2).Id, int64(0))
	assert.Equal(t, tree.GetParentNode(16).Id, int64(13))
	assert.Nil(t, tree.GetParentNode(999))
}

func TestGetRootNode(t *testing.T) {
	assert.Equal(t, tree.GetRootNode().Id, int64(0))

	subtree1 := tree.SubTree(2)
	assert.Equal(t, subtree1.GetRootNode().Id, int64(2))
}

func TestGetAncestorNode(t *testing.T) {
	node0 := tree.GetAncestorNode(16, 0)
	assert.Equal(t, node0.Id, int64(13))

	// 0->2->6->13->16
	node1 := tree.GetAncestorNode(16, 1)
	assert.Equal(t, node1.Id, int64(13))

	node2 := tree.GetAncestorNode(16, 2)
	assert.Equal(t, node2.Id, int64(6))

	node3 := tree.GetAncestorNode(16, 3)
	assert.Equal(t, node3.Id, int64(2))

	node4 := tree.GetAncestorNode(16, 4)
	assert.Equal(t, node4.Id, int64(0))

	node5 := tree.GetAncestorNode(16, 5)
	assert.Equal(t, node5.Id, int64(0))

	node6 := tree.GetAncestorNode(0, 4)
	assert.Nil(t, node6)

	node7 := tree.GetAncestorNode(2, 3)
	assert.Equal(t, node7.Id, int64(0))
}

func TestGetLeafNodes(t *testing.T) {
	leafIds1 := make([]int64, 0)
	for _, leafNode := range tree.GetLeafNodes() {
		leafIds1 = append(leafIds1, leafNode.Id)
	}
	assert.ElementsMatch(t, leafIds1, []int64{12, 14, 1, 8, 9, 11, 5, 10, 16, 7, 15})

	leafIds2 := make([]int64, 0)
	for _, leafNode := range tree.GetLeafNodes(2) {
		leafIds2 = append(leafIds2, leafNode.Id)
	}
	assert.ElementsMatch(t, leafIds2, []int64{5, 12, 16, 14, 15, 7})

	leafIds3 := make([]int64, 0)
	for _, leafNode := range tree.GetLeafNodes(0) {
		leafIds3 = append(leafIds3, leafNode.Id)
	}
	assert.ElementsMatch(t, leafIds3, []int64{12, 14, 1, 8, 9, 11, 5, 10, 16, 7, 15})

	leafIds4 := make([]int64, 0)
	for _, leafNode := range tree.GetLeafNodes(11) {
		leafIds4 = append(leafIds4, leafNode.Id)
	}
	assert.Equal(t, leafIds4, []int64{11})

	leafIds5 := make([]int64, 0)
	for _, leafNode := range tree.GetLeafNodes(18) {
		leafIds5 = append(leafIds5, leafNode.Id)
	}
	assert.Equal(t, len(leafIds5), 0)
}

func TestGetSiblingNodes(t *testing.T) {
	siblingNodes := tree.GetSiblingNodes(13)
	siblingIds := make([]int64, 0)
	for _, node := range siblingNodes {
		siblingIds = append(siblingIds, node.Id)
	}
	assert.ElementsMatch(t, siblingIds, []int64{12, 15, 14})
}

func TestGetAllNodes(t *testing.T) {
	ids := make([]int64, 0)
	for _, node := range tree.GetAllNodes() {
		ids = append(ids, node.Id)
	}
	assert.ElementsMatch(t, ids, []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
}

func TestGetChildNodes(t *testing.T) {
	ids1 := make([]int64, 0)
	for _, node := range tree.GetChildNodes(2) {
		ids1 = append(ids1, node.Id)
	}
	assert.ElementsMatch(t, ids1, []int64{5, 6, 7})

	ids2 := make([]int64, 0)
	for _, node := range tree.GetChildNodes(11) {
		ids2 = append(ids2, node.Id)
	}
	assert.ElementsMatch(t, ids2, []int64{})
}

func TestGetDescendantNodes(t *testing.T) {
	ids1 := make([]int64, 0)
	for _, node := range tree.GetDescendantNodes(2) {
		ids1 = append(ids1, node.Id)
	}
	assert.ElementsMatch(t, ids1, []int64{2, 5, 6, 7, 12, 13, 14, 15, 16})

	ids2 := make([]int64, 0)
	for _, node := range tree.GetDescendantNodes(11) {
		ids2 = append(ids2, node.Id)
	}
	assert.ElementsMatch(t, ids2, []int64{11})
}

func TestNodeIsLeaf(t *testing.T) {
	assert.Equal(t, tree.GetNode(0).IsLeaf(), false)
	assert.Equal(t, tree.GetNode(13).IsLeaf(), false)
	assert.Equal(t, tree.GetNode(9).IsLeaf(), true)

	n := tree.GetNode(16)
	fmt.Println(n.ChildIds)
	assert.Equal(t, tree.GetNode(15).IsLeaf(), true)
}

func TestRSearch(t *testing.T) {
	ids1 := make([]int64, 0)
	for id := range tree.RSearch(0) {
		ids1 = append(ids1, id)
	}
	assert.ElementsMatch(t, ids1, []int64{})

	ids2 := make([]int64, 0)
	for id := range tree.RSearch(45) {
		ids2 = append(ids2, id)
	}
	assert.ElementsMatch(t, ids2, []int64{})

	ids3 := make([]int64, 0)
	for id := range tree.RSearch(12) {
		ids3 = append(ids3, id)
	}
	assert.ElementsMatch(t, ids3, []int64{12, 6, 2, 0})
}

func TestDepthFirstTraversal(t *testing.T) {
	nids := make([]int64, 0)
	for nid := range tree.DepthFirstTraversal(2) {
		nids = append(nids, nid)
	}
	assert.Equal(t, nids, []int64{2, 5, 6, 12, 13, 16, 14, 15, 7})

	filterNids := make([]int64, 0)
	for nid := range tree.DepthFirstTraversal(2, func(node *Node) bool {
		nodeData := node.Data.(*NodeData)
		return nodeData.Level > 5
	}) {
		filterNids = append(filterNids, nid)
	}
	assert.Equal(t, filterNids, []int64{6, 12, 13, 16, 14, 15, 7})
}

func TestWidthFirstTraversal(t *testing.T) {
	ids := make([]int64, 0)
	for nid := range tree.WidthFirstTraversal(2) {
		ids = append(ids, nid)
	}
	assert.Equal(t, ids, []int64{2, 5, 6, 7, 12, 13, 14, 15, 16})

	filterIds := make([]int64, 0)
	for nid := range tree.WidthFirstTraversal(2, func(node *Node) bool {
		nodeData := node.Data.(*NodeData)
		return nodeData.Level > 5
	}) {
		filterIds = append(filterIds, nid)
	}
	assert.Equal(t, filterIds, []int64{6, 7, 12, 13, 14, 15, 16})
}

func TestNodeHasChildren(t *testing.T) {
	node := tree.GetNode(2)
	assert.Equal(t, node.HasChildren(), true)

	node = tree.GetNode(16)
	assert.Equal(t, node.HasChildren(), false)
}

func TestCreateNode(t *testing.T) {
	tt := getTree()
	_, err := tt.CreateNode(16, 15, nil)
	assert.Error(t, err, ErrNodeAlreadyExist)

	node, _ := tt.CreateNode(17, 15, nil)
	assert.Equal(t, node.Pid, int64(15))
}

func TestMoveNode(t *testing.T) {
	tt := getTree()

	// [0 2 6 14]
	// [0 2 6 13 16]
	// [0 2 6 12]
	//  [0 2 6 15]
	childIds2 := tt.GetChildIds(2)
	assert.ElementsMatch(t, childIds2, []int64{5, 6, 7})

	childIds3 := tt.GetChildIds(3)
	assert.ElementsMatch(t, childIds3, []int64{8, 9})

	err := tt.MoveNode(6, 3)
	if err != nil {
		t.Error("Error move node: ", err)
	}

	childIds2 = tt.GetChildIds(2)
	childIds3 = tt.GetChildIds(3)
	assert.ElementsMatch(t, childIds2, []int64{5, 7})
	assert.ElementsMatch(t, childIds3, []int64{6, 8, 9})
}

func TestDeleteNode(t *testing.T) {
	tt := getTree()
	// [0 2 6 14]
	// [0 2 6 13 16]
	// [0 2 6 12]
	//  [0 2 6 15]
	assert.ElementsMatch(t, tt.GetChildIds(2), []int64{5, 6, 7})
	assert.Equal(t, tt.GetNode(12).Pid, int64(6))

	err := tt.DeleteNode(6)
	if err != nil {
		t.Error("Error delete node: ", err)
	}

	assert.ElementsMatch(t, tt.GetChildIds(2), []int64{5, 12, 13, 14, 15, 7})
	assert.Equal(t, tt.GetNode(12).Pid, int64(2))
}

func TestAddNode(t *testing.T) {
	tt := getTree()

	newNode := NewNode(17, nil)
	err := tt.mountTo(newNode, 0)
	if err != nil {
		t.Error("Error add node: ", err)
	}

	getNode := tt.GetNode(17)
	assert.Equal(t, newNode, getNode)
}

func BenchmarkTree(b *testing.B) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tree)
	if err != nil {
		b.Error("Error encode: ", err)
	}

	data := result.Bytes()

	var ttt SafeMultiBranchTree
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err = decoder.Decode(&ttt)
	if err != nil {
		b.Error("Error decode: ", err)
	}
}

func getData(id int) *NodeData {
	data := make(map[int]*NodeData)
	for i := 1; i <= 16; i++ {
		data[i] = &NodeData{
			Name:     fmt.Sprintf("name%d", i),
			Mobile:   fmt.Sprintf("1350731000%d", i),
			Nickname: fmt.Sprintf("nick%d", i),
			Level:    i,
		}
	}
	return data[id]
}

func getTree() *SafeMultiBranchTree {
	rootNode := NewRootNode(&NodeData{})
	tt := NewTree(rootNode)

	pids := make([]int, 0)
	for k := range parentChildren {
		pids = append(pids, k)
	}

	sort.Slice(pids, func(i, j int) bool { return pids[i] < pids[j] })

	for _, pid := range pids {
		for _, id := range parentChildren[pid] {
			data := getData(id)
			if _, err := tt.CreateNode(int64(id), int64(pid), data); err != nil {
				log.Printf("create node error: %v", err)
			}
		}
	}

	return tt
}
