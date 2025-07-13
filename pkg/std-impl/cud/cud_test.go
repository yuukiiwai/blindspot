package cud

import (
	"strings"
	"testing"
)

func TestCudYamlParser(t *testing.T) {
	yamlInput := `
start_resources:
  user_count: 0
  server_status: "stopped"

edge_rules:
  - name: start_server
    effect:
      - action: update
        resource:
          key: server_status
          value: "running"
    fire_condition: server_status == "stopped"
    block_condition: user_count > 100

  - name: add_user
    effect:
      - action: update
        resource:
          key: user_count
          value: user_count + 1
    fire_condition: server_status == "running"
    block_condition: user_count >= 50

  - name: create_log
    effect:
      - action: create
        resource:
          key: log_entry
          value: "server started"
    fire_condition: server_status == "running" && !has("log_entry")
    block_condition: ""

  - name: stop_server
    effect:
      - action: update
        resource:
          key: server_status
          value: "stopped"
      - action: delete
        resource:
          key: log_entry
          value: ""
    fire_condition: server_status == "running" && user_count == 0
    block_condition: ""
`

	parser, err := NewCudYamlParser()
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	firstResource, newNode, edgeRules, err := parser.Parse(yamlInput)
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	// 初期リソースの確認
	if firstResource == nil {
		t.Fatal("First resource is nil")
	}

	resources := firstResource.GetResources().(map[string]any)
	if resources["user_count"] != 0 {
		t.Errorf("Expected user_count to be 0, got %v", resources["user_count"])
	}
	if resources["server_status"] != "stopped" {
		t.Errorf("Expected server_status to be 'stopped', got %v", resources["server_status"])
	}

	// エッジルールの数を確認
	if len(edgeRules) != 4 {
		t.Errorf("Expected 4 edge rules, got %d", len(edgeRules))
	}

	// エッジルール名の確認
	expectedNames := []string{"start_server", "add_user", "create_log", "stop_server"}
	for i, rule := range edgeRules {
		if rule.Name != expectedNames[i] {
			t.Errorf("Expected rule name %s, got %s", expectedNames[i], rule.Name)
		}
	}

	// newNode関数のテスト
	testResources := map[string]any{"test": "value"}
	testNode := newNode(testResources)
	if testNode == nil {
		t.Fatal("newNode returned nil")
	}

	// ID生成のテスト
	id := firstResource.GetID()
	if id == "" {
		t.Error("GetID returned empty string")
	}

	// ResourcesStringのテスト
	resourceStrings := firstResource.GetResourcesString()
	if len(resourceStrings) == 0 {
		t.Error("GetResourcesString returned empty slice")
	}

	// 複数要素を含む文字列表現の確認
	resourceString := strings.Join(resourceStrings, " ")
	if !strings.Contains(resourceString, "user_count") || !strings.Contains(resourceString, "server_status") {
		t.Errorf("Resource string doesn't contain expected keys: %s", resourceString)
	}
}

func TestCudNodeOperations(t *testing.T) {
	// 空ノードのテスト
	emptyNode := newCudNode(nil)
	if emptyNode.GetID() != "empty" {
		t.Errorf("Expected empty node ID to be 'empty', got %s", emptyNode.GetID())
	}

	// 通常ノードのテスト
	resources := map[string]any{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}
	node := newCudNode(resources)

	// ID生成の一意性テスト
	id1 := node.GetID()
	id2 := node.GetID()
	if id1 != id2 {
		t.Error("GetID should return consistent results")
	}

	// 異なるリソースの異なるIDテスト
	resources2 := map[string]any{
		"key1": "different_value",
		"key2": 42,
		"key3": true,
	}
	node2 := newCudNode(resources2)
	if node.GetID() == node2.GetID() {
		t.Error("Different resources should have different IDs")
	}

	// Equalsテスト
	node3 := newCudNode(resources)
	if !node.Equals(node3) {
		t.Error("Nodes with same resources should be equal")
	}
	if node.Equals(node2) {
		t.Error("Nodes with different resources should not be equal")
	}

	// GetResourcesStringテスト
	resourceStrings := node.GetResourcesString()
	if len(resourceStrings) != 3 {
		t.Errorf("Expected 3 resource strings, got %d", len(resourceStrings))
	}
}
