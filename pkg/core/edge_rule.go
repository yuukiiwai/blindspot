package core

/*
EdgeRule エッジのルールを定義

	Name: ルールの名前
	Effect: エッジ発火後のNodeを返す関数
	FireCondition: ルールが発火する条件に合致した場合にtrueを返す関数
	BlockCondition: ルールがブロックされる条件に合致した場合にtrueを返す関数(前提として、FireConditionがtrueの場合に評価する)

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
) *EdgeRule {
	return &EdgeRule{
		Name:           name,
		Effect:         effect,
		FireCondition:  fireCondition,
		BlockCondition: blockCondition,
	}
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
