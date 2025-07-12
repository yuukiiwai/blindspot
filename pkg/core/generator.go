package core

import (
	"log/slog"
	"sort"
)

// Generator ステートマシン生成器
type Generator struct {
	newNode        func(resources any) Node
	startResources Node
	edgeRules      []*EdgeRule
	nodes          map[string]Node // インターフェースを使用
	edges          []*Edge
	processedNodes map[string]bool
}

// NewGenerator 新しいジェネレーターを作成
func NewGenerator(
	newNode func(resources any) Node,
	startResources Node,
	edgeRules []*EdgeRule,
) *Generator {
	return &Generator{
		newNode:        newNode,
		startResources: startResources,
		edgeRules:      edgeRules,
		nodes:          make(map[string]Node),
		edges:          make([]*Edge, 0),
		processedNodes: make(map[string]bool),
	}
}

// Generate ステートマシンを生成
func (g *Generator) Generate() error {
	startNode := g.newNode(g.startResources.GetResources())
	startNodePtr := g.addOrGetNode(&startNode)
	queue := []*Node{startNodePtr}

	slog.Debug("[START] 開始ノード", "resources", (*startNodePtr).GetResources(), "id", (*startNodePtr).GetID())

	iterationCount := 0
	for len(queue) > 0 {
		iterationCount++
		if iterationCount > 1000 { // 安全のための上限設定
			slog.Error("[ERROR] 1000回以上の反復でループ検出。強制終了します。")
			break
		}

		currentNode := queue[0]
		queue = queue[1:]
		nodeID := (*currentNode).GetID()

		slog.Debug("[ITERATION]", "count", iterationCount, "resources", (*currentNode).GetResources(), "id", nodeID, "queueSize", len(queue))

		if g.processedNodes[nodeID] {
			slog.Debug("[SKIP] すでに処理済み", "id", nodeID)
			continue
		}

		g.processedNodes[nodeID] = true
		newEdges := g.generateEdgesFromNode(currentNode)

		slog.Debug("[EDGES]", "count", len(newEdges))

		for _, edge := range newEdges {
			g.edges = append(g.edges, edge)
			targetNodeID := (*edge.GetTo()).GetID()
			slog.Debug("[EDGE_ADD]", "edge", edge.String())
			if !g.processedNodes[targetNodeID] {
				queue = append(queue, edge.GetTo())
				slog.Debug("[QUEUE_ADD] キューに追加", "resources", (*edge.GetTo()).GetResources(), "id", (*edge.GetTo()).GetID())
			} else {
				slog.Debug("[QUEUE_SKIP] すでに処理済みのためキューに追加しない", "resources", (*edge.GetTo()).GetResources(), "id", (*edge.GetTo()).GetID())
			}
		}

		slog.Debug("[QUEUE_STATUS]", "size", len(queue))
	}

	slog.Debug("[COMPLETE]", "iterations", iterationCount, "nodes", len(g.nodes), "edges", len(g.edges))

	// デバッグ: すべてのノードを出力
	slog.Debug("[DEBUG] 生成されたノード一覧")
	for id, node := range g.nodes {
		slog.Debug("ノード情報", "id", id, "resources", node.GetResources())
	}

	return nil
}

// GetNodes 生成されたノードを取得
// 設計資料として出力される際の一貫性と可読性のため、ノードIDでソートして返す
func (g *Generator) GetNodes() []*Node {
	// ノードIDのスライスを作成してソート
	var ids []string
	for id := range g.nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	
	// ソート済みのIDに基づいてノードスライスを構築
	var nodes []*Node
	for _, id := range ids {
		node := g.nodes[id]
		nodes = append(nodes, &node)
	}
	return nodes
}

// GetEdges 生成されたエッジを取得
func (g *Generator) GetEdges() []*Edge {
	return g.edges
}

// addOrGetNode ノードを追加または取得
func (g *Generator) addOrGetNode(node *Node) *Node {
	id := (*node).GetID()
	slog.Debug("[NODE_DEBUG]", "resources", (*node).GetResources(), "id", id)
	if existingNode, exists := g.nodes[id]; exists {
		slog.Debug("[NODE_REUSE] 既存ノードを再利用", "id", id)
		return &existingNode
	}
	g.nodes[id] = *node
	slog.Debug("[NODE_CREATE] 新しいノードを作成", "id", id, "resources", (*node).GetResources())
	return node
}

// generateEdgesFromNode 指定されたノードから適用可能なエッジを生成
func (g *Generator) generateEdgesFromNode(node *Node) []*Edge {
	var edges []*Edge

	for _, rule := range g.edgeRules {
		fire := rule.GetFireCondition()(node)
		block := rule.GetBlockCondition()(node)
		slog.Debug("[CHECK]", "resources", (*node).GetResources(), "rule", rule.GetName(), "fire", fire, "block", block)
		if fire && !block {
			newNode := rule.GetEffect()(node)
			slog.Debug("[EFFECT]", "resources", (*node).GetResources(), "rule", rule.GetName(), "newResources", (*newNode).GetResources(), "newId", (*newNode).GetID())
			actualNewNode := g.addOrGetNode(newNode)
			edge := NewEdge(node, actualNewNode, rule)
			edges = append(edges, edge)
		}
	}

	return edges
}
