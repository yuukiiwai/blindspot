package core

import (
	"fmt"
	"sort"
	"strings"
)

// いつか使う
type NodeInterface interface {
	GetID() string
	Equals(other NodeInterface) bool
	GetResources() []string
}

// Node リソースの状態を表現
type Node []string

// NewNode 新しいNodeを作成
func NewNode(resources []string) Node {
	// 空文字列を除去
	filteredResources := make([]string, 0, len(resources))
	for _, resource := range resources {
		if resource != "" {
			filteredResources = append(filteredResources, resource)
		}
	}

	// リソースをコピーしてソート
	sortedResources := make([]string, len(filteredResources))
	copy(sortedResources, filteredResources)
	sort.Strings(sortedResources)
	return Node(sortedResources)
}

// GetID ノードの一意な識別子を生成
func (n *Node) GetID() string {
	if len(*n) == 0 {
		return "empty"
	}
	return strings.Join(*n, ",")
}

// Equals ノードが同じかどうかを判定
func (n *Node) Equals(other *Node) bool {
	return n.GetID() == other.GetID()
}

// String ノードの文字列表現
func (n *Node) String() string {
	return fmt.Sprintf("Node[%s]", strings.Join(*n, ", "))
}

// GetResources リソースのリストを取得
func (n *Node) GetResources() []string {
	return *n
}
