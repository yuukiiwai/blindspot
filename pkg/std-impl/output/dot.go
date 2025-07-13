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
	dot.WriteString("  rankdir=LR;\n")
	dot.WriteString("  node [shape=box];\n\n")

	// ノードの出力
	for _, node := range generator.GetNodes() {
		nodeID := getDotNodeID(node)
		label := getDotNodeLabel(node)
		dot.WriteString(fmt.Sprintf("  %s [label=\"%s\"];\n", nodeID, label))
	}

	dot.WriteString("\n")

	// エッジの出力
	for _, edge := range generator.GetEdges() {
		fromID := getDotNodeID(edge.GetFrom())
		toID := getDotNodeID(edge.GetTo())
		edgeLabel := edge.GetRule().GetName()
		dot.WriteString(fmt.Sprintf("  %s -> %s [label=\"%s\"];\n", fromID, toID, edgeLabel))
	}

	dot.WriteString("}\n")

	return dot.String(), nil
}

// getDotNodeID ノードのDOT IDを生成
func getDotNodeID(node *core.Node) string {
	id := (*node).GetID()

	// 明示的にemptyの場合
	if id == "empty" {
		return "empty"
	}

	// 空文字列の場合（これは起こるべきではないが、安全のため）
	if id == "" {
		return "empty"
	}

	// 通常のIDをDOTで使える形式に変換
	// カンマをアンダースコアに、スペースをアンダースコアに置換
	result := strings.ReplaceAll(id, ",", "_")
	result = strings.ReplaceAll(result, " ", "_")

	// その他の特殊文字も安全な文字に置換
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, ".", "_")

	return result
}

// getDotNodeLabel ノードのDOT表示名を生成
func getDotNodeLabel(node *core.Node) string {
	resources := (*node).GetResourcesString()

	return strings.Join(resources, "\\n")
}
