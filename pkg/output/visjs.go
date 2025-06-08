package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuukiiwai/blindspot/pkg/core"
)

// VisjsFormatter Vis.js形式の出力フォーマッター
type VisjsFormatter struct{}

// NewVisjsFormatter 新しいVisjsFormatterを作成
func NewVisjsFormatter() *VisjsFormatter {
	return &VisjsFormatter{}
}

// Format ステートマシンをVis.js形式で出力
func (f *VisjsFormatter) Format(generator *core.Generator) (string, error) {
	// ノードの生成
	nodes := make([]map[string]interface{}, 0)
	for id, node := range generator.GetNodes() {
		nodes = append(nodes, map[string]interface{}{
			"id":    id,
			"label": strings.Join(node.GetResources(), "\n"),
		})
	}

	// エッジの生成
	edges := make([]map[string]interface{}, 0)
	for _, edge := range generator.GetEdges() {
		edges = append(edges, map[string]interface{}{
			"from":  edge.GetFrom().GetID(),
			"to":    edge.GetTo().GetID(),
			"label": edge.GetRule().GetName(),
		})
	}

	// JSONデータの生成
	data := map[string]interface{}{
		"nodes": nodes,
		"edges": edges,
	}

	// JSONに変換
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON変換エラー: %w", err)
	}

	return string(jsonData), nil
}
