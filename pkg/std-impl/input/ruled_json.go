package input

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/yuukiiwai/blindspot/pkg/core"
)

type RuledJson struct {
	StartResources []string `json:"start_resources"`
	EdgeRules      []struct {
		Name           string   `json:"name"`            // ルール名
		Action         string   `json:"action"`          // create, update, delete
		Rule           []string `json:"rule"`            // create, deleteは対象の文字列1つ, updateは[0]が既存, [1]が新しいもの
		FireCondition  []string `json:"fire_condition"`  // ある物を指定して、その物がある場合に発火する
		BlockCondition []string `json:"block_condition"` // Fireがtrueのときに評価する。ある物を指定して、その物がある場合にブロックする
	} `json:"edge_rules"`
}

func createFireConditionFunc(conditions []string) func(*core.Node) bool {
	return func(n *core.Node) bool {
		// 条件が空の場合は、ノードのリソースも空の場合のみtrueを返す
		if len(conditions) == 0 {
			return len(n.GetResources()) == 0
		}

		resources := n.GetResources()
		for _, condition := range conditions {
			if slices.Contains(resources, condition) {
				return true
			}
		}
		return false
	}
}

func createBlockConditionFunc(conditions []string) func(*core.Node) bool {
	return func(n *core.Node) bool {
		// 条件が空の場合は常にfalseを返す（ブロックしない）
		if len(conditions) == 0 {
			return false
		}

		resources := n.GetResources()
		for _, condition := range conditions {
			if slices.Contains(resources, condition) {
				return true
			}
		}
		return false
	}
}

func NewRuledJsonParser() (core.Parser, error) {
	return &RuledJson{}, nil
}

func (r *RuledJson) Parse(input string) ([]string, []*core.EdgeRule, error) {
	var ruledJson RuledJson
	err := json.Unmarshal([]byte(input), &ruledJson)
	if err != nil {
		return nil, nil, err
	}
	var edgeRules []*core.EdgeRule
	for _, rule := range ruledJson.EdgeRules {
		// クロージャー内で使用するためにルールをコピー
		currentRule := rule
		fireCondition := createFireConditionFunc(currentRule.FireCondition)
		blockCondition := createBlockConditionFunc(currentRule.BlockCondition)

		if currentRule.Action == "create" {
			edgeRules = append(edgeRules, core.NewEdgeRule(
				currentRule.Name,
				func(n *core.Node) *core.Node {
					currentResources := n.GetResources()
					newResources := make([]string, len(currentResources)+1)
					copy(newResources, currentResources)
					newResources[len(currentResources)] = currentRule.Rule[0]
					newNode := core.NewNode(newResources)
					return &newNode
				},
				fireCondition,
				blockCondition,
			))
		} else if currentRule.Action == "update" {
			edgeRules = append(edgeRules, core.NewEdgeRule(
				currentRule.Name,
				func(n *core.Node) *core.Node {
					currentResources := n.GetResources()
					newResources := make([]string, len(currentResources))
					copy(newResources, currentResources)

					targetIndex := slices.Index(newResources, currentRule.Rule[0])
					if targetIndex == -1 {
						panic(fmt.Sprintf("rule: %v, rule.Rule[0] %s not found in %v", currentRule, currentRule.Rule[0], currentResources))
					}
					newResources[targetIndex] = currentRule.Rule[1]
					newNode := core.NewNode(newResources)
					return &newNode
				},
				fireCondition,
				blockCondition,
			))
		} else if currentRule.Action == "delete" {
			edgeRules = append(edgeRules, core.NewEdgeRule(
				currentRule.Name,
				func(n *core.Node) *core.Node {
					currentResources := n.GetResources()
					newResources := make([]string, 0, len(currentResources))

					// 削除対象以外の要素をコピー
					for _, resource := range currentResources {
						if resource != currentRule.Rule[0] {
							newResources = append(newResources, resource)
						}
					}

					newNode := core.NewNode(newResources)
					return &newNode
				},
				fireCondition,
				blockCondition,
			))
		} else {
			panic(fmt.Sprintf("the rule action is not defined: %v, action: %s", currentRule, currentRule.Action))
		}
	}
	return ruledJson.StartResources, edgeRules, nil
}
