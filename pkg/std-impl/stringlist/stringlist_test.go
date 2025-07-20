// stringlist実装のテスト
package stringlist

import (
	"testing"

	"github.com/yuukiiwai/blindspot/pkg/core"
	"github.com/yuukiiwai/blindspot/pkg/std-impl/output"
)

var expectedOutput = `graph TD
    empty["empty"]
    a["a"]
    a_b["a<br/>b"]

    empty -->|create_a| a
    a -->|create_b_from_a| a_b
    a -->|delete_a| empty
    a_b -->|delete_b| a
`

func TestStringListNode(t *testing.T) {
	exampleContent := `
	{
		"start_resources": [],
		"edge_rules": [
			{
				"name": "create_a",
				"action": "create",
				"rule": ["a"],
				"fire_condition": [],
				"block_condition": []
			},
			{
				"name": "create_b_from_a",
				"action": "create",
				"rule": ["b"],
				"fire_condition": ["a"],
				"block_condition": ["b"]
			},
			{
				"name": "delete_b",
				"action": "delete",
				"rule": ["b"],
				"fire_condition": ["b"],
				"block_condition": []
			},
			{
				"name": "delete_a",
				"action": "delete",
				"rule": ["a"],
				"fire_condition": ["a"],
				"block_condition": ["b"]
			}
		]
	}
	`
	parser, err := NewRuledJsonParser()
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	firstResource, newNode, edgeRules, err := parser.Parse(exampleContent)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	generator := core.NewGenerator(
		newNode,
		firstResource,
		edgeRules,
		nil,
	)
	generator.Generate()
	mermaid := output.NewMermaidFormatter()
	outputResult, err := mermaid.Format(generator)
	if err != nil {
		t.Fatalf("failed to format: %v", err)
	}
	if outputResult != expectedOutput {
		t.Errorf("expected %s, but got %s", expectedOutput, outputResult)
	}
}
