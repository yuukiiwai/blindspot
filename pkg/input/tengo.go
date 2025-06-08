package input

import (
	"github.com/d5/tengo/v2"
	"github.com/yuukiiwai/blindspot/pkg/core"
)

type tengoParser struct {
	script *tengo.Script
}

func NewTengoParser() (Parser, error) {
	// tengoにTypeを追加する
	s := tengo.NewScript([]byte(""))
	s.Add("Node", core.Node{})
	s.Add("EdgeRule", core.EdgeRule{})
	s.Add("Edge", core.Edge{})
	return &tengoParser{s}, nil
}

func (p *tengoParser) Parse(input string) ([]string, []*core.EdgeRule, error) {
	return nil, nil, nil
}
