package input

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/yuukiiwai/blindspot/pkg/core"
)

// InputData は入力JSONの構造を定義
type InputData struct {
	FirstResources []string `json:"first-resources"`
	Rules          []Rule   `json:"rules"`
}

// Rule はルールの構造を定義
type Rule struct {
	Name           string `json:"name"`
	Effect         string `json:"effect"`
	FireCondition  string `json:"fireCondition"`
	BlockCondition string `json:"blockCondition"`
}

// YaegiParser はYaegiを使用してGoのコードを評価するパーサー
type YaegiParser struct {
	interpreter *interp.Interpreter
}

// NewYaegiParser は新しいYaegiParserを作成する
func NewYaegiParser() (*YaegiParser, error) {
	i := interp.New(interp.Options{})
	if err := i.Use(stdlib.Symbols); err != nil {
		return nil, fmt.Errorf("failed to use stdlib: %w", err)
	}

	// coreパッケージの型を登録
	nodeValue := reflect.ValueOf(core.Node{})
	if err := i.Use(map[string]map[string]reflect.Value{
		"github.com/yuukiiwai/blindspot/pkg/core": {
			"Node":         nodeValue,
			"NewNode":      reflect.ValueOf(core.NewNode),
			"GetResources": nodeValue.MethodByName("GetResources"),
		},
	}); err != nil {
		return nil, fmt.Errorf("failed to use core package: %w", err)
	}

	// ヘルパー関数を定義
	helperCode := `
		func contains(slice []string, item string) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		}

		func filter(slice []string, f func(string) bool) []string {
			var result []string
			for _, s := range slice {
				if f(s) {
					result = append(result, s)
				}
			}
			return result
		}
	`
	if _, err := i.Eval(helperCode); err != nil {
		return nil, fmt.Errorf("failed to eval helper functions: %w", err)
	}

	return &YaegiParser{
		interpreter: i,
	}, nil
}

// Parse は入力ファイルをパースして、開始リソースとルールを返す
func (p *YaegiParser) Parse(input string) ([]string, []*core.EdgeRule, error) {
	var data InputData
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return nil, nil, fmt.Errorf("JSONのパースに失敗: %w", err)
	}

	if len(data.FirstResources) == 0 {
		return nil, nil, fmt.Errorf("開始リソースが指定されていません")
	}

	if len(data.Rules) == 0 {
		return nil, nil, fmt.Errorf("ルールが指定されていません")
	}

	edgeRules := make([]*core.EdgeRule, 0, len(data.Rules))
	for _, rule := range data.Rules {
		// エフェクトのコンパイル
		effectCode := fmt.Sprintf(`func(node *github.com/yuukiiwai/blindspot/pkg/core.Node) *github.com/yuukiiwai/blindspot/pkg/core.Node {
			resources := node.GetResources()
			%s
			return github.com/yuukiiwai/blindspot/pkg/core.NewNode(resources)
		}`, rule.Effect)

		// 発火条件のコンパイル
		fireCode := fmt.Sprintf(`func(node *github.com/yuukiiwai/blindspot/pkg/core.Node) bool {
			resources := node.GetResources()
			result := %s
			return result
		}`, rule.FireCondition)

		// ブロック条件のコンパイル
		blockCode := fmt.Sprintf(`func(node *github.com/yuukiiwai/blindspot/pkg/core.Node) bool {
			resources := node.GetResources()
			result := %s
			return result
		}`, rule.BlockCondition)

		log.Printf("[DEBUG] ルール: %s", rule.Name)
		log.Printf("[DEBUG] エフェクト: %s", effectCode)
		log.Printf("[DEBUG] 発火条件: %s", fireCode)
		log.Printf("[DEBUG] ブロック条件: %s", blockCode)

		edgeRules = append(edgeRules, core.NewEdgeRule(
			rule.Name,
			func(node *core.Node) *core.Node {
				v, err := p.interpreter.Eval(effectCode)
				if err != nil {
					log.Printf("[ERROR] エフェクトの評価に失敗: %v", err)
					return node
				}
				f := v.Interface().(func(*core.Node) *core.Node)
				return f(node)
			},
			func(node *core.Node) bool {
				v, err := p.interpreter.Eval(fireCode)
				if err != nil {
					log.Printf("[ERROR] 発火条件の評価に失敗: %v", err)
					return false
				}
				f := v.Interface().(func(*core.Node) bool)
				return f(node)
			},
			func(node *core.Node) bool {
				v, err := p.interpreter.Eval(blockCode)
				if err != nil {
					log.Printf("[ERROR] ブロック条件の評価に失敗: %v", err)
					return false
				}
				f := v.Interface().(func(*core.Node) bool)
				return f(node)
			},
		))
	}

	return data.FirstResources, edgeRules, nil
}
