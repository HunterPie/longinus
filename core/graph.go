package core

type PatternGraphNode struct {
	Owners     []*PatternOwner
	Value      uint8
	IsWildcard bool
	Nodes      []*PatternGraphNode
}

func (node *PatternGraphNode) IsEqualToEdge(edge *PatternEdge) bool {
	return node.IsWildcard == edge.IsWildcard && node.Value == edge.Value
}

func (node *PatternGraphNode) Add(graphNode *PatternGraphNode) {
	node.Nodes = append(node.Nodes, graphNode)
}

func (node *PatternGraphNode) AddOwner(owner *PatternOwner) {
	node.Owners = append(node.Owners, owner)
}

func FindNodeByEdge(nodes []*PatternGraphNode, edge *PatternEdge) *PatternGraphNode {
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

func NewFromEdge(edge *PatternEdge) *PatternGraphNode {
	return &PatternGraphNode{
		Value:      edge.Value,
		IsWildcard: edge.IsWildcard,
		Owners:     make([]*PatternOwner, 0),
		Nodes:      make([]*PatternGraphNode, 0),
	}
}

type PatternGraph struct {
	Nodes []*PatternGraphNode
}

func (graph *PatternGraph) FindByEdge(edge *PatternEdge) *PatternGraphNode {
	for _, node := range graph.Nodes {
		if node.IsEqualToEdge(edge) {
			return node
		}
	}

	return nil
}

func (graph *PatternGraph) Add(node *PatternGraphNode) {
	graph.Nodes = append(graph.Nodes, node)
}

func buildGraph(owners []*PatternOwner) *PatternGraph {
	graph := &PatternGraph{Nodes: make([]*PatternGraphNode, 0)}

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

func NewGraph(owners ...*PatternOwner) *PatternGraph {
	return buildGraph(owners)
}
