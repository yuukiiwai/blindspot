package node

import (
	"sort"
	"strings"

	"github.com/yuukiiwai/blindspot/pkg/core"
)

// Node リソースの状態を表現
type StringListNode []string

// NewNode 新しいNodeを作成

// GetID ノードの一意な識別子を生成
func (n StringListNode) GetID() string {
	if len(n) == 0 {
		return "empty"
	}
	return strings.Join(n, ",")
}

// Equals ノードが同じかどうかを判定
func (n StringListNode) Equals(other core.Node) bool {
	return n.GetID() == other.GetID()
}

// GetResources リソースのリストを取得
func (n StringListNode) GetResources() any {
	return []string(n)
}

// GetResourcesString リソースのリストを文字列で取得
func (n StringListNode) GetResourcesString() string {
	if len(n) == 0 {
		return "empty"
	}
	return strings.Join(n, "<br/>")
}

func NewStringListNode(resources []string) StringListNode {
	// 空文字列を除去してソート
	filtered := make([]string, 0, len(resources))
	for _, resource := range resources {
		if resource != "" {
			filtered = append(filtered, resource)
		}
	}

	sorted := make([]string, len(filtered))
	copy(sorted, filtered)
	sort.Strings(sorted)

	return StringListNode(sorted)
}
