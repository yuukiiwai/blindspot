package cud

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yuukiiwai/blindspot/pkg/core"
)

// CudNode キーバリューペアのリソースの状態を表現
type CudNode map[string]any

// GetID ノードの一意な識別子を生成
func (n CudNode) GetID() string {
	if len(n) == 0 {
		return "empty"
	}

	// キーでソートして一意性を保証
	keys := make([]string, 0, len(n))
	for k := range n {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		// 値をJSON形式で文字列化して一意性を保証
		valueBytes, err := json.Marshal(n[k])
		if err != nil {
			panic(fmt.Sprintf("failed to marshal value for key %s: %v", k, err))
		}
		parts = append(parts, fmt.Sprintf("%s:%s", k, string(valueBytes)))
	}

	return strings.Join(parts, ",")
}

// Equals ノードが同じかどうかを判定
func (n CudNode) Equals(other core.Node) bool {
	return n.GetID() == other.GetID()
}

// GetResources リソースのマップを取得
func (n CudNode) GetResources() any {
	return map[string]any(n)
}

// GetResourcesString リソースのリストを文字列で取得
func (n CudNode) GetResourcesString() []string {
	if len(n) == 0 {
		return []string{"empty"}
	}

	// キーでソートして出力の一貫性を保証
	keys := make([]string, 0, len(n))
	for k := range n {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var result []string
	for _, k := range keys {
		// 値をJSON形式で表示
		valueBytes, err := json.Marshal(n[k])
		if err != nil {
			result = append(result, fmt.Sprintf("%s:<marshal_error>", k))
		} else {
			result = append(result, fmt.Sprintf("%s:%s", k, string(valueBytes)))
		}
	}

	return result
}

func newCudNode(resources map[string]any) CudNode {
	// nilチェックとコピー作成
	if resources == nil {
		return CudNode{}
	}

	result := make(CudNode, len(resources))
	for k, v := range resources {
		if k != "" { // 空キーは除外
			result[k] = v
		}
	}

	return result
}
