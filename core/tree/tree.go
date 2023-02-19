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

func (node *PatternTreeNode) Add(graphNode *PatternTreeNode) {
	node.Nodes = append(node.Nodes, graphNode)
}

func (node *PatternTreeNode) AddOwner(owner *signature.PatternOwner) {
	node.Owners = append(node.Owners, owner)
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

func (tree *PatternTree) Add(node *PatternTreeNode) {
	tree.Nodes = append(tree.Nodes, node)
}

func buildTree(owners []*signature.PatternOwner) *PatternTree {
	graph := &PatternTree{Nodes: make([]*PatternTreeNode, 0)}

	for _, owner := range owners {
		pattern := owner.Pattern
		currentNode := graph.FindByEdge(pattern)

		if currentNode == nil {
			currentNode = NewFromEdge(pattern)
			graph.Add(currentNode)
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

	return graph
}

func New(owners ...*signature.PatternOwner) *PatternTree {
	return buildTree(owners)
}
