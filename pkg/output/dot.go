package output

import (
	"fmt"
	"strings"

	"github.com/yuukiiwai/blindspot/pkg/core"
)

// DotFormatter Graphviz（DOT）形式の出力フォーマッター
type DotFormatter struct{}

// NewDotFormatter 新しいDotFormatterを作成
func NewDotFormatter() *DotFormatter {
	return &DotFormatter{}
}

// Format ステートマシンをDOT形式で出力
func (f *DotFormatter) Format(generator *core.Generator) (string, error) {
	var dot strings.Builder
	dot.WriteString("digraph G {\n")
	dot.WriteString("    rankdir=LR;\n")
	dot.WriteString("    node [shape=box];\n\n")

	// ノードの出力
	for _, node := range generator.GetNodes() {
		nodeID := getDotNodeID(node)
		label := getDotNodeLabel(node)
		dot.WriteString(fmt.Sprintf("    %s [label=\"%s\"];\n", nodeID, label))
	}

	dot.WriteString("\n")

	// エッジの出力
	for _, edge := range generator.GetEdges() {
		fromID := getDotNodeID(edge.GetFrom())
		toID := getDotNodeID(edge.GetTo())
		edgeLabel := edge.GetRule().GetName()
		dot.WriteString(fmt.Sprintf("    %s -> %s [label=\"%s\"];\n", fromID, toID, edgeLabel))
	}

	dot.WriteString("}\n")
	return dot.String(), nil
}

// getDotNodeID ノードのDOT IDを生成
func getDotNodeID(node *core.Node) string {
	id := node.GetID()
	if id == "empty" {
		return "empty"
	}
	return strings.ReplaceAll(strings.ReplaceAll(id, ",", "_"), " ", "_")
}

// getDotNodeLabel ノードのDOT表示名を生成
func getDotNodeLabel(node *core.Node) string {
	if len(node.GetResources()) == 0 {
		return "empty"
	}
	return strings.Join(node.GetResources(), "\\n")
}
