package output

import (
	"fmt"
	"strings"

	"github.com/yuukiiwai/blindspot/pkg/core"
)

// MermaidFormatter Mermaid形式の出力フォーマッター
type MermaidFormatter struct{}

// NewMermaidFormatter 新しいMermaidFormatterを作成
func NewMermaidFormatter() *MermaidFormatter {
	return &MermaidFormatter{}
}

// Format ステートマシンをMermaid形式で出力
func (f *MermaidFormatter) Format(generator *core.Generator) (string, error) {
	var mermaid strings.Builder
	mermaid.WriteString("graph TD\n")

	// ノードの出力
	for _, node := range generator.GetNodes() {
		nodeID := getMermaidNodeID(node)
		label := getMermaidNodeLabel(node)
		mermaid.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", nodeID, label))
	}

	mermaid.WriteString("\n")

	// エッジの出力
	for _, edge := range generator.GetEdges() {
		fromID := getMermaidNodeID(edge.GetFrom())
		toID := getMermaidNodeID(edge.GetTo())
		edgeLabel := edge.GetRule().GetName()
		mermaid.WriteString(fmt.Sprintf("    %s -->|%s| %s\n", fromID, edgeLabel, toID))
	}

	return mermaid.String(), nil
}

// getMermaidNodeID ノードのMermaid IDを生成
func getMermaidNodeID(node *core.Node) string {
	id := (*node).GetID()

	// 明示的にemptyの場合
	if id == "empty" {
		return "empty"
	}

	// 空文字列の場合（これは起こるべきではないが、安全のため）
	if id == "" {
		return "empty"
	}

	// 通常のIDをMermaidで使える形式に変換
	// カンマをアンダースコアに、スペースをアンダースコアに置換
	result := strings.ReplaceAll(id, ",", "_")
	result = strings.ReplaceAll(result, " ", "_")

	// その他の特殊文字も安全な文字に置換
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, ".", "_")

	return result
}

// getMermaidNodeLabel ノードのMermaid表示名を生成
func getMermaidNodeLabel(node *core.Node) string {
	resources := (*node).GetResourcesString()

	return resources
}
