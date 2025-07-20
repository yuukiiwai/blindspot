package cud

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/yuukiiwai/blindspot/pkg/core"
	"gopkg.in/yaml.v3"
)

type CudYaml struct {
	StartResources map[string]any `yaml:"start_resources"`
	EdgeRules      []struct {
		Name   string `yaml:"name"`
		Effect []struct {
			Action   string `yaml:"action"` // create, update, delete
			Resource struct {
				Key   string `yaml:"key"`
				Value any    `yaml:"value"`
			} `yaml:"resource"`
		} `yaml:"effect"`
		FireCondition  string `yaml:"fire_condition"`  // expr-lang expression
		BlockCondition string `yaml:"block_condition"` // expr-lang expression
	} `yaml:"edge_rules"`
}

func createFireConditionFunc(conditionExpr string) func(*core.Node) bool {
	if conditionExpr == "" {
		return func(n *core.Node) bool {
			resources, ok := (*n).GetResources().(map[string]any)
			if !ok {
				panic(fmt.Sprintf("node resources is not map[string]any: %v", n))
			}
			return len(resources) == 0
		}
	}

	program, err := expr.Compile(conditionExpr, expr.AllowUndefinedVariables())
	if err != nil {
		panic(fmt.Sprintf("failed to compile fire condition expression: %s, error: %v", conditionExpr, err))
	}

	return func(n *core.Node) bool {
		resources, ok := (*n).GetResources().(map[string]any)
		if !ok {
			panic(fmt.Sprintf("node resources is not map[string]any: %v", n))
		}

		result, err := expr.Run(program, resources)
		if err != nil {
			panic(fmt.Sprintf("failed to evaluate fire condition: %s, error: %v", conditionExpr, err))
		}

		boolResult, ok := result.(bool)
		if !ok {
			panic(fmt.Sprintf("fire condition expression must return bool, got: %T", result))
		}

		return boolResult
	}
}

func createBlockConditionFunc(conditionExpr string) func(*core.Node) bool {
	if conditionExpr == "" {
		return func(n *core.Node) bool {
			return false
		}
	}

	program, err := expr.Compile(conditionExpr, expr.AllowUndefinedVariables())
	if err != nil {
		panic(fmt.Sprintf("failed to compile block condition expression: %s, error: %v", conditionExpr, err))
	}

	return func(n *core.Node) bool {
		resources, ok := (*n).GetResources().(map[string]any)
		if !ok {
			panic(fmt.Sprintf("node resources is not map[string]any: %v", n))
		}

		result, err := expr.Run(program, resources)
		if err != nil {
			panic(fmt.Sprintf("failed to evaluate block condition: %s, error: %v", conditionExpr, err))
		}

		boolResult, ok := result.(bool)
		if !ok {
			panic(fmt.Sprintf("block condition expression must return bool, got: %T", result))
		}

		return boolResult
	}
}

func NewCudYamlParser() (core.Parser, error) {
	return &CudYaml{}, nil
}

func (c *CudYaml) Parse(input string) (
	firstResource core.Node,
	newNode func(any) core.Node,
	edgeRules []*core.EdgeRule,
	err error,
) {
	var cudYaml CudYaml
	err = yaml.Unmarshal([]byte(input), &cudYaml)
	if err != nil {
		return nil, nil, nil, err
	}

	newNode = func(resources any) core.Node {
		return newCudNode(resources.(map[string]any))
	}

	for _, rule := range cudYaml.EdgeRules {
		// クロージャー内で使用するためにルールをコピー
		currentRule := rule
		fireCondition := createFireConditionFunc(currentRule.FireCondition)
		blockCondition := createBlockConditionFunc(currentRule.BlockCondition)

		// effect配列の処理
		if len(currentRule.Effect) == 0 {
			return nil, nil, nil, fmt.Errorf("effect cannot be empty for rule: %s", currentRule.Name)
		}

		edgeRule, err := core.NewEdgeRule(
			currentRule.Name,
			func(n *core.Node) *core.Node {
				currentResources, ok := (*n).GetResources().(map[string]any)
				if !ok {
					panic(fmt.Sprintf("node resources is not map[string]any: %v", n))
				}
				newResources := make(map[string]any)
				// 既存のリソースをコピー
				for k, v := range currentResources {
					newResources[k] = v
				}

				// 各effectを順番に適用
				for _, effect := range currentRule.Effect {
					if effect.Resource.Key == "" {
						panic(fmt.Sprintf("resource key cannot be empty for rule: %s", currentRule.Name))
					}

					switch effect.Action {
					case "create":
						if effect.Resource.Value == nil {
							panic(fmt.Sprintf("resource value cannot be nil for create action in rule: %s", currentRule.Name))
						}
						newResources[effect.Resource.Key] = effect.Resource.Value

					case "update":
						if _, exists := newResources[effect.Resource.Key]; !exists {
							panic(fmt.Sprintf("rule: %v, key %s not found in current resources %v", currentRule, effect.Resource.Key, newResources))
						}
						if effect.Resource.Value == nil {
							panic(fmt.Sprintf("resource value cannot be nil for update action in rule: %s", currentRule.Name))
						}
						newResources[effect.Resource.Key] = effect.Resource.Value

					case "delete":
						delete(newResources, effect.Resource.Key)

					default:
						panic(fmt.Sprintf("unknown action: %s in rule: %s", effect.Action, currentRule.Name))
					}
				}

				newNode := newNode(newResources)
				return &newNode
			},
			fireCondition,
			blockCondition,
		)
		if err != nil {
			return nil, nil, nil, err
		}
		edgeRules = append(edgeRules, edgeRule)
	}

	return newNode(cudYaml.StartResources), newNode, edgeRules, nil
}
