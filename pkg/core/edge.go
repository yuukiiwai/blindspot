package core

import (
	"fmt"
)

// Edge リソースの変化を表現
type Edge struct {
	from *Node
	to   *Node
	rule *EdgeRule
}

// NewEdge 新しいEdgeを作成
func NewEdge(from, to *Node, rule *EdgeRule) *Edge {
	return &Edge{
		from: from,
		to:   to,
		rule: rule,
	}
}

// String エッジの文字列表現
func (e *Edge) String() string {
	return fmt.Sprintf("%s --%s--> %s", e.from.GetID(), e.rule.Name, e.to.GetID())
}

// GetFrom 開始ノードを取得
func (e *Edge) GetFrom() *Node {
	return e.from
}

// GetTo 終了ノードを取得
func (e *Edge) GetTo() *Node {
	return e.to
}

// GetRule エッジのルールを取得
func (e *Edge) GetRule() *EdgeRule {
	return e.rule
}
