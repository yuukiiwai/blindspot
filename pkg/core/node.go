// pkg/core/node.go
package core

type Node interface {
	GetID() string                // ノードを識別するためのIDを返す
	Equals(other Node) bool       // ノードが同じかどうかを判定
	GetResources() any            // ノード実装を返す
	GetResourcesString() []string // ノードを表現するときに1行で表現するものを1要素として返す
}
