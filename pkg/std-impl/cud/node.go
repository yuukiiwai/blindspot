package cud

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/yuukiiwai/blindspot/pkg/core"
)

// CudNode キーバリューペアのリソースの状態を表現
type CudNode map[string]any

// GetID ノードの一意な識別子を生成
func (n CudNode) GetID() string {
	var marshaled []byte
	var err error
	if len(n) == 0 {
		marshaled = []byte("")
	} else {
		marshaled, err = json.Marshal(n)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal node: %v", err))
		}
	}
	return fmt.Sprintf("%x", md5.Sum(marshaled))
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
