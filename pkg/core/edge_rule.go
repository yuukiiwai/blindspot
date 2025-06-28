package core

import "fmt"

/*
EdgeRule エッジのルールを定義

	Name: ルールの名前
	Effect: エッジ発火後のNodeを返す関数
	FireCondition: ルールが発火する条件に合致した場合にtrueを返す関数
	BlockCondition: ルールがブロックされる条件に合致した場合にtrueを返す関数(前提として、FireConditionがtrueの場合に評価する)

EffectやFireCondition, BlockConditionは処理中に型が違う場合panicを起こしたほうが良い。

つまり発火条件は、FireConditionがtrueの上で、BlockConditionがfalseの場合に発火する。
*/
type EdgeRule struct {
	Name           string
	Effect         func(*Node) *Node
	FireCondition  func(*Node) bool
	BlockCondition func(*Node) bool
}

// NewEdgeRule 新しいEdgeRuleを作成
func NewEdgeRule(
	name string,
	effect func(*Node) *Node,
	fireCondition func(*Node) bool,
	blockCondition func(*Node) bool,
) (*EdgeRule, error) {
	if effect == nil {
		return nil, fmt.Errorf("effect function cannot be nil")
	}
	if fireCondition == nil {
		return nil, fmt.Errorf("fireCondition function cannot be nil")
	}
	if blockCondition == nil {
		return nil, fmt.Errorf("blockCondition function cannot be nil")
	}
	return &EdgeRule{
		Name:           name,
		Effect:         effect,
		FireCondition:  fireCondition,
		BlockCondition: blockCondition,
	}, nil
}

// GetName ルール名を取得
func (r *EdgeRule) GetName() string {
	return r.Name
}

// GetEffect エフェクト関数を取得
func (r *EdgeRule) GetEffect() func(*Node) *Node {
	return r.Effect
}

// GetFireCondition 発火条件関数を取得
func (r *EdgeRule) GetFireCondition() func(*Node) bool {
	return r.FireCondition
}

// GetBlockCondition ブロック条件関数を取得
func (r *EdgeRule) GetBlockCondition() func(*Node) bool {
	return r.BlockCondition
}
