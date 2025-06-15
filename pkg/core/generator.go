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
	startNode := NewNode(g.startResources)
	startNodePtr := g.addOrGetNode(&startNode)
	queue := []*Node{startNodePtr}

	log.Printf("[START] 開始ノード: %v (ID: %s)", startNodePtr.GetResources(), startNodePtr.GetID())

	iterationCount := 0
	for len(queue) > 0 {
		iterationCount++
		if iterationCount > 100 { // 安全のための上限設定
			log.Printf("[ERROR] 100回以上の反復でループ検出。強制終了します。")
			break
		}

		currentNode := queue[0]
		queue = queue[1:]
		nodeID := currentNode.GetID()

		log.Printf("[ITERATION %d] 処理中ノード: %v (ID: %s), キューサイズ: %d", iterationCount, currentNode.GetResources(), nodeID, len(queue))

		if g.processedNodes[nodeID] {
			log.Printf("[SKIP] すでに処理済み: %s", nodeID)
			continue
		}

		g.processedNodes[nodeID] = true
		newEdges := g.generateEdgesFromNode(currentNode)

		log.Printf("[EDGES] %d個のエッジを生成", len(newEdges))

		for _, edge := range newEdges {
			g.edges = append(g.edges, edge)
			targetNodeID := edge.GetTo().GetID()
			log.Printf("[EDGE_ADD] %s", edge.String())
			if !g.processedNodes[targetNodeID] {
				queue = append(queue, edge.GetTo())
				log.Printf("[QUEUE_ADD] キューに追加: %v (ID: %s)", edge.GetTo().GetResources(), edge.GetTo().GetID())
			} else {
				log.Printf("[QUEUE_SKIP] すでに処理済みのためキューに追加しない: %v (ID: %s)", edge.GetTo().GetResources(), edge.GetTo().GetID())
			}
		}

		log.Printf("[QUEUE_STATUS] 現在のキューサイズ: %d", len(queue))
	}

	log.Printf("[COMPLETE] 総反復回数: %d, 総ノード数: %d, 総エッジ数: %d", iterationCount, len(g.nodes), len(g.edges))

	// デバッグ: すべてのノードを出力
	log.Printf("[DEBUG] 生成されたノード一覧:")
	for id, node := range g.nodes {
		log.Printf("  ID: '%s' -> Resources: %v", id, node.GetResources())
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
	log.Printf("[NODE_DEBUG] addOrGetNode: resources=%v, id='%s'", node.GetResources(), id)
	if existingNode, exists := g.nodes[id]; exists {
		log.Printf("[NODE_REUSE] 既存ノードを再利用: ID='%s'", id)
		return existingNode
	}
	g.nodes[id] = node
	log.Printf("[NODE_CREATE] 新しいノードを作成: ID='%s', resources=%v", id, node.GetResources())
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
			log.Printf("[EFFECT] ノード: %v, ルール: %s, 新リソース: %v, 新ID: '%s'", node.GetResources(), rule.GetName(), newNode.GetResources(), newNode.GetID())
			actualNewNode := g.addOrGetNode(newNode)
			edge := NewEdge(node, actualNewNode, rule)
			edges = append(edges, edge)
		}
	}

	return edges
}
