package core

import (
	"fmt"
	"sort"
	"strings"
)

// Node リソースの状態を表現
type Node struct {
	resources []string
}

// NewNode 新しいNodeを作成
func NewNode(resources []string) *Node {
	// リソースをコピーしてソート
	sortedResources := make([]string, len(resources))
	copy(sortedResources, resources)
	sort.Strings(sortedResources)
	return &Node{resources: sortedResources}
}

// GetID ノードの一意な識別子を生成
func (n *Node) GetID() string {
	if len(n.resources) == 0 {
		return "empty"
	}
	return strings.Join(n.resources, ",")
}

// Equals ノードが同じかどうかを判定
func (n *Node) Equals(other *Node) bool {
	return n.GetID() == other.GetID()
}

// String ノードの文字列表現
func (n *Node) String() string {
	return fmt.Sprintf("Node[%s]", strings.Join(n.resources, ", "))
}

// GetResources リソースのリストを取得
func (n *Node) GetResources() []string {
	return n.resources
}
