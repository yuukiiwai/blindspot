package core

import (
	"log"
)

// Generator ステートマシン生成器
type Generator struct {
	startResources []string
	edgeRules      []*EdgeRule
	nodes          map[string]*Node
	edges          []*Edge
	processedNodes map[string]bool
}

// NewGenerator 新しいジェネレーターを作成
func NewGenerator(startResources []string, edgeRules []*EdgeRule) *Generator {
	return &Generator{
		startResources: startResources,
		edgeRules:      edgeRules,
		nodes:          make(map[string]*Node),
		edges:          make([]*Edge, 0),
		processedNodes: make(map[string]bool),
	}
}

// Generate ステートマシンを生成
func (g *Generator) Generate() error {
	startNode := g.addOrGetNode(NewNode(g.startResources))
	queue := []*Node{startNode}

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]
		nodeID := currentNode.GetID()

		if g.processedNodes[nodeID] {
			continue
		}

		g.processedNodes[nodeID] = true
		newEdges := g.generateEdgesFromNode(currentNode)

		for _, edge := range newEdges {
			g.edges = append(g.edges, edge)
			targetNodeID := edge.GetTo().GetID()
			if !g.processedNodes[targetNodeID] {
				queue = append(queue, edge.GetTo())
			}
		}
	}

	return nil
}

// GetNodes 生成されたノードを取得
func (g *Generator) GetNodes() map[string]*Node {
	return g.nodes
}

// GetEdges 生成されたエッジを取得
func (g *Generator) GetEdges() []*Edge {
	return g.edges
}

// addOrGetNode ノードを追加または取得
func (g *Generator) addOrGetNode(node *Node) *Node {
	id := node.GetID()
	if existingNode, exists := g.nodes[id]; exists {
		return existingNode
	}
	g.nodes[id] = node
	return node
}

// generateEdgesFromNode 指定されたノードから適用可能なエッジを生成
func (g *Generator) generateEdgesFromNode(node *Node) []*Edge {
	var edges []*Edge

	for _, rule := range g.edgeRules {
		fire := rule.GetFireCondition()(node)
		block := rule.GetBlockCondition()(node)
		log.Printf("[CHECK] ノード: %v, ルール: %s, fireCondition: %v, blockCondition: %v", node.GetResources(), rule.GetName(), fire, block)
		if fire && !block {
			newNode := rule.GetEffect()(node)
			log.Printf("[EFFECT] ノード: %v, ルール: %s, 新リソース: %v", node.GetResources(), rule.GetName(), newNode.GetResources())
			actualNewNode := g.addOrGetNode(newNode)
			edge := NewEdge(node, actualNewNode, rule)
			edges = append(edges, edge)
		}
	}

	return edges
}
