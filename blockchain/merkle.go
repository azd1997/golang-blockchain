package blockchain

import "crypto/sha256"

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	//对于叶节点（交易）来说，直接哈希
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		//对于其他节点，需将下边左右两个节点哈希拼起来再哈希
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	node.Left = left
	node.Right = right

	return &node

}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	//MerkleTree是根据叶节点也就是交易列表建立的
	//如果叶节点为奇数，则需要最后一个复制本身来成对

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	//获取节点集合,此时nodes只有Data
	for _, dat := range data {
		node := NewMerkleNode(nil, nil, dat)
		nodes = append(nodes, *node)
	}

	//由MerkleNodes构建MerkleTree
	for i := 0; i < len(data)/2; i++ {
		var level []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			level = append(level, *node)
		}

		//将nodes填充指针
		nodes = level
	}

	tree := MerkleTree{&nodes[0]}

	return &tree

}
