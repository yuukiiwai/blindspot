// pkg/core/node.go
package core

type Node interface {
	GetID() string
	Equals(other Node) bool
	GetResources() any
	GetResourcesString() string
}
