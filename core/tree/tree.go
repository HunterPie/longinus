package tree

import (
	"github.com/HunterPie/Longinus/core/signature"
)

type PatternTreeNode struct {
	Owners     []*signature.PatternOwner
	Value      uint8
	IsWildcard bool
	Nodes      []*PatternTreeNode
}

func (node *PatternTreeNode) IsEqualToEdge(edge *signature.PatternEdge) bool {
	return node.IsWildcard == edge.IsWildcard && node.Value == edge.Value
}

func (node *PatternTreeNode) IsEqualToByte(byte uint8) bool {
	return node.IsWildcard || node.Value == byte
}

func (node *PatternTreeNode) Add(graphNode *PatternTreeNode) {
	node.Nodes = append(node.Nodes, graphNode)
}

func (node *PatternTreeNode) AddOwner(owner *signature.PatternOwner) {
	node.Owners = append(node.Owners, owner)
}

func (node *PatternTreeNode) hasOwners() bool {
	return len(node.Owners) > 0
}

func (node *PatternTreeNode) findRecursively(bytearray []uint8, owners []*signature.PatternOwner) []*signature.PatternOwner {
	if len(bytearray) == 0 {
		return owners
	}

	currentByte := bytearray[0]
	if !node.IsEqualToByte(currentByte) {
		return owners
	}

	if node.hasOwners() {
		owners = append(owners, node.Owners...)
	}

	for _, nextNode := range node.Nodes {
		owners = nextNode.findRecursively(bytearray[1:], owners)
	}

	return owners
}

func (node *PatternTreeNode) Find(bytearray []uint8) []*signature.PatternOwner {
	return node.findRecursively(bytearray, make([]*signature.PatternOwner, 0))
}

func FindNodeByEdge(nodes []*PatternTreeNode, edge *signature.PatternEdge) *PatternTreeNode {
	if edge == nil {
		return nil
	}

	for _, node := range nodes {
		if node.IsEqualToEdge(edge) {
			return node
		}
	}

	return nil
}

func NewFromEdge(edge *signature.PatternEdge) *PatternTreeNode {
	return &PatternTreeNode{
		Value:      edge.Value,
		IsWildcard: edge.IsWildcard,
		Owners:     make([]*signature.PatternOwner, 0),
		Nodes:      make([]*PatternTreeNode, 0),
	}
}

type PatternTree struct {
	Nodes []*PatternTreeNode
}

func (tree *PatternTree) FindByEdge(edge *signature.PatternEdge) *PatternTreeNode {
	for _, node := range tree.Nodes {
		if node.IsEqualToEdge(edge) {
			return node
		}
	}

	return nil
}

func (tree *PatternTree) FindPattern(bytearray []uint8) []*signature.PatternOwner {
	found := make([]*signature.PatternOwner, 0)

	for _, node := range tree.Nodes {
		found = append(found, node.Find(bytearray)...)
	}

	return found
}

func (tree *PatternTree) Add(node *PatternTreeNode) {
	tree.Nodes = append(tree.Nodes, node)
}

func buildTree(owners []*signature.PatternOwner) *PatternTree {
	tree := &PatternTree{Nodes: make([]*PatternTreeNode, 0)}

	for _, owner := range owners {
		pattern := owner.Pattern
		currentNode := tree.FindByEdge(pattern)

		if currentNode == nil {
			currentNode = NewFromEdge(pattern)
			tree.Add(currentNode)
		}

		pattern = pattern.Next

		for pattern != nil {
			nextPattern := FindNodeByEdge(currentNode.Nodes, pattern)

			if nextPattern == nil {
				nextPattern = NewFromEdge(pattern)
				currentNode.Add(nextPattern)
			}

			currentNode = nextPattern
			pattern = pattern.Next
		}

		currentNode.AddOwner(owner)
	}

	return tree
}

func New(owners ...*signature.PatternOwner) *PatternTree {
	return buildTree(owners)
}
