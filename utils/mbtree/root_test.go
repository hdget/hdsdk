package mbtree

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
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

var paths = map[int][]int{
	0:  []int{1, 2, 3, 4},
	2:  []int{5, 6, 7},
	3:  []int{8, 9},
	4:  []int{10, 11},
	6:  []int{12, 13, 14, 15},
	13: []int{16},
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
	tt := NewTree()

	tt.CreateRootNode(0, &NodeData{})

	pids := make([]int, 0)
	for k := range paths {
		pids = append(pids, k)
	}

	sort.Slice(pids, func(i, j int) bool { return pids[i] < pids[j] })

	for _, pid := range pids {
		for _, id := range paths[pid] {
			data := getData(id)
			if _, err := tt.CreateNode(int64(id), int64(pid), data); err != nil {
				log.Printf("create node error: %v", err)
			}
		}
	}

	return tt
}

func TestGetAncestor(t *testing.T) {
	tt := getTree()

	// 0->2->6->13->16
	node := tt.GetAncestorNode(16, 2)
	fmt.Println(node)
	assert.Equal(t, node.Id, int64(6))

	node = tt.GetAncestorNode(16, 0)
	assert.Equal(t, node.Id, int64(13))
}

func TestTree_CreateRootNode(t *testing.T) {
	tt := NewTree()
	tt.CreateRootNode(3, nil)
	assert.Equal(t, tt.RootId, int64(3))
}

func TestTree_Paths(t *testing.T) {
	tt := getTree()
	paths := tt.Paths()
	assert.ElementsMatch(t, paths, [][]int64{
		[]int64{0, 2, 6, 12}, []int64{0, 2, 6, 14}, []int64{0, 1}, []int64{0, 3, 8}, []int64{0, 3, 9},
		[]int64{0, 4, 11}, []int64{0, 2, 5}, []int64{0, 4, 10}, []int64{0, 2, 6, 13, 16}, []int64{0, 2, 7}, []int64{0, 2, 6, 15}})
}

func TestTree_Depth(t *testing.T) {
	tt := getTree()
	assert.Equal(t, tt.Depth(), 4)
	assert.Equal(t, tt.Depth(6), 2)
}

func TestTree_DepthFirstTraversal(t *testing.T) {
	tt := getTree()

	nids := make([]int64, 0)
	for nid := range tt.DepthFirstTraversal(2) {
		nids = append(nids, nid)
	}
	assert.ElementsMatch(t, nids, []int64{2, 5, 6, 7, 12, 13, 14, 15, 16})

	filterNids := make([]int64, 0)
	for nid := range tt.DepthFirstTraversal(2, func(node *Node) bool {
		nodeData := node.Data.(*NodeData)
		return nodeData.Level > 5
	}) {
		filterNids = append(filterNids, nid)
	}
	assert.ElementsMatch(t, filterNids, []int64{6, 7, 12, 13, 14, 15, 16})
}

func TestTree_DeleteNode(t *testing.T) {
	tt := getTree()

	// [0 2 6 14]
	// [0 2 6 13 16]
	// [0 2 6 12]
	//  [0 2 6 15]
	assert.Equal(t, len(tt.GetChildIds(2)), 3)
	assert.Equal(t, tt.GetNode(12).Pid, int64(6))

	err := tt.DeleteNode(6)
	if err != nil {
		t.Error("Error delete node: ", err)
	}

	assert.Equal(t, len(tt.GetChildIds(2)), 6)
	assert.Equal(t, tt.GetNode(12).Pid, int64(2))
}

func TestTree_GetAllNodes(t *testing.T) {
	tt := getTree()

	for _, node := range tt.GetAllNodes() {
		t.Log(node.Id)
	}
}

func TestTree_GetChildIds(t *testing.T) {
	tt := getTree()

	t.Log(tt.GetChildIds(2))
}

func TestTree_GetChildNodes(t *testing.T) {
	tt := getTree()
	childNodes := tt.GetChildNodes(2)
	for _, node := range childNodes {
		t.Log(node.Id)
	}
}

func TestTree_GetDescendantNodeIds(t *testing.T) {
	tt := getTree()
	for _, nid := range tt.GetDescendantIds(2) {
		t.Log(nid)
	}
}

func TestTree_GetDescendantNodes(t *testing.T) {
	tt := getTree()
	for _, node := range tt.GetDescendantNodes(2) {
		t.Log(node.Id)
	}
}

func TestTree_GetNode(t *testing.T) {
	tt := getTree()

	assert.Equal(t, tt.GetNode(2).Id, int64(2))
	assert.Nil(t, tt.GetNode(999))
}

func TestTree_GetParentNode(t *testing.T) {
	tt := getTree()
	assert.Equal(t, tt.GetParentNode(2).Id, int64(0))
	assert.Equal(t, tt.GetParentNode(16).Id, int64(13))
	assert.Nil(t, tt.GetParentNode(999))
}

func TestTree_GetParentId(t *testing.T) {
	tt := getTree()
	assert.Equal(t, tt.GetParentId(2), int64(0))
	assert.Equal(t, tt.GetParentId(16), int64(13))
	assert.Equal(t, tt.GetParentId(999), int64(-1))
}

func TestTree_GetRootNode(t *testing.T) {
	tt := getTree()

	rootNode := tt.GetRootNode()
	assert.Equal(t, rootNode.Id, tt.RootId)
}

func TestTree_MoveNode(t *testing.T) {
	tt := getTree()

	// [0 2 6 14]
	// [0 2 6 13 16]
	// [0 2 6 12]
	//  [0 2 6 15]
	childIds2 := tt.GetChildIds(2)
	childIds3 := tt.GetChildIds(3)
	t.Log("Before move, node '2' has childs:", childIds2)
	t.Log("Before move, node '3' has childs:", childIds3)

	err := tt.MoveNode(6, 3)
	if err != nil {
		t.Error("Error move node: ", err)
	}

	childIds2 = tt.GetChildIds(2)
	childIds3 = tt.GetChildIds(3)

	t.Log("After move, node '2' has childs:", childIds2)
	t.Log("After move, node '3' has childs:", childIds3)
}

func TestNewNode(t *testing.T) {

}

func TestNewTree(t *testing.T) {
	tt := NewTree()
	t.Log("RootId is: ", tt.RootId)
}

func TestTree_Siblings(t *testing.T) {
	tt := getTree()

	sibilings := tt.GetSiblingNodes(13)
	sibilingIds := make([]int64, 0)
	for _, node := range sibilings {
		sibilingIds = append(sibilingIds, node.Id)
	}
	assert.ElementsMatch(t, sibilingIds, []int64{12, 15, 14})
}

func TestTree_Size(t *testing.T) {
	tt := getTree()
	assert.Equal(t, tt.Size(), 17)
}

func TestTree_Traversal(t *testing.T) {

}

func TestTree_WidthFirstTraversal(t *testing.T) {

}

func TestTree_CreateNode(t *testing.T) {
	tt := getTree()

	node, err := tt.CreateNode(16, 15, nil)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(node)
	}

	node, err = tt.CreateNode(17, 15, nil)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(node)
	}
}

func TestNode_HasChildren(t *testing.T) {
	tt := getTree()

	node := tt.GetNode(2)
	assert.Equal(t, node.HasChildren(), true)

	node = tt.GetNode(16)
	assert.Equal(t, node.HasChildren(), false)
}

func TestNode_IsLeaf(t *testing.T) {
	tt := getTree()

	node := tt.GetNode(2)
	assert.Equal(t, node.IsLeaf(), false)

	node = tt.GetNode(16)
	assert.Equal(t, node.IsLeaf(), true)
}

func TestTree_AddNode(t *testing.T) {
	tt := getTree()

	newNode := NewNode(17, nil)
	err := tt.mountTo(newNode, 0)
	if err != nil {
		t.Error("Error add node: ", err)
	}
	log.Println(newNode.Pid)

	getNode := tt.GetNode(17)
	assert.Equal(t, newNode, getNode)

	log.Println(getNode.Pid)
}

func TestTree_FilterNodes(t *testing.T) {
	tt := getTree()

	filter := func(n *Node) bool {
		nodeData := n.Data.(*NodeData)
		return nodeData.Level > 10
	}

	for _, node := range tt.filterNodes(filter) {
		nodeData := node.Data.(*NodeData)
		t.Logf("DbServerIndex: %d, Level: %d", node.Id, nodeData.Level)
	}
}

func TestTree_FilterNodesByIds(t *testing.T) {
	tt := getTree()

	ids := []int64{1, 5, 10, 13, 15, 16}
	filter := func(n *Node) bool {
		nodeData := n.Data.(*NodeData)
		return nodeData.Level > 10
	}

	for _, node := range tt.filterNodesInIds(ids, filter) {
		nodeData := node.Data.(*NodeData)
		t.Logf("DbServerIndex: %d, Level: %d", node.Id, nodeData.Level)
	}
}

func TestTree_Contains(t *testing.T) {
	tt := getTree()
	assert.Equal(t, tt.Contains(-1), false)
	assert.Equal(t, tt.Contains(3), true)
	assert.Equal(t, tt.Contains(18), false)
}

func TestTree_IsAncestor(t *testing.T) {
	tt := getTree()
	assert.Equal(t, tt.IsAncestor(0, 2), true)
	assert.Equal(t, tt.IsAncestor(0, 15), true)
	assert.Equal(t, tt.IsAncestor(6, 16), true)
	assert.Equal(t, tt.IsAncestor(1, 13), false)
}

func TestTree_SubTree(t *testing.T) {
	tt := getTree()

	subtree := tt.SubTree(2)
	t.Logf("Subtree rootId now is: %d", subtree.RootId)
	//for nid := range subtree.Nodes {
	//	t.Log(nid)
	//}
	subtree.Nodes.Range(func(k, v interface{}) bool {
		nid := k.(int64)
		t.Log(nid)
		return true
	})
}

func TestTree_RSearch(t *testing.T) {
	tt := getTree()

	for id := range tt.RSearch(0) {
		log.Printf("RSearch => %d\n", id)
	}

	for id := range tt.RSearch(45) {
		log.Printf("RSearch => %d\n", id)
	}

	for id := range tt.RSearch(12) {
		log.Printf("RSearch => %d\n", id)
	}

	for id := range tt.traversal(0, DEPTH) {
		log.Printf("DepthFirstTraversal => %d\n", id)
	}

	for id := range tt.traversal(0, WIDTH) {
		log.Printf("WidthFirstTraversal => %d\n", id)
	}

}

func TestTree_Leaves(t *testing.T) {
	tt := getTree()
	t.Log(tt.GetLeafNodes())
}

func TestTree_Level(t *testing.T) {
	tt := getTree()

	assert.Equal(t, tt.Level(0), 0)
	assert.Equal(t, tt.Level(15), 3)
	assert.Equal(t, tt.Level(13), 3)
	assert.Equal(t, tt.Level(16), 4)
}

func BenchmarkTree_xxx(b *testing.B) {
	tt := getTree()

	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tt)
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
