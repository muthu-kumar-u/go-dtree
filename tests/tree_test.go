package tests

import (
	"testing"

	dtree "github.com/muthu-kumar-u/go-dtree/tree"
)

func TestCreateTreeFromjson_ValidInput(t *testing.T) {
	jsonData := `{
		"root": {
			"node_id": 1,
			"predict": "ShouldBuy",
			"branches": [
				{
					"condition": {"type": "bool", "value": true},
					"outcome": {"nextNode": {"node_id": 2, "predict": "CheckBudget", "branches": []}}
				}
			]
		}
	}`

	tree, err := dtree.CreateTreeFromjson([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to create tree from JSON: %v", err)
	}

	if tree.Root == nil || tree.Root.Predict != "ShouldBuy" {
		t.Errorf("Unexpected root node: %+v", tree.Root)
	}
}

func TestCreateTreeFromjson_InvalidJSON(t *testing.T) {
	invalidJsonData := `{ "root": { "node_id": 1, "predict": "ShouldBuy", "branches": [ { "condition": { "type": "bool", "value": true }, "outcome": { "nextNode": } } ] } }`

	_, err := dtree.CreateTreeFromjson([]byte(invalidJsonData))
	if err == nil {
		t.Fatal("Expected error for invalid JSON but got nil")
	}
}

func TestCreateTreeFromjson_MissingRootNode(t *testing.T) {
	jsonData := `{ "branches": [] }`

	tree, err := dtree.CreateTreeFromjson([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected error for missing root node but got nil")
	}

	if tree != nil {
		t.Error("Tree should be nil when root node is missing")
	}
}

func TestCreateTreeFromjson_MissingConditionOrOutcome(t *testing.T) {
	jsonData := `{
		"root": {
			"node_id": 1,
			"predict": "ShouldBuy",
			"branches": [
				{
					"condition": null,
					"outcome": {"nextNode": {"node_id": 2, "predict": "CheckBudget", "branches": []}}
				}
			]
		}
	}`

	_, err := dtree.CreateTreeFromjson([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected error for missing condition but got nil")
	}
}
