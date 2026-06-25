package api

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSwaggerYAMLIsValid(t *testing.T) {
	raw, err := os.ReadFile("swagger.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	if err := yaml.Unmarshal(raw, &document); err != nil {
		t.Fatal(err)
	}
	if document["openapi"] != "3.0.3" {
		t.Fatalf("unexpected OpenAPI version: %v", document["openapi"])
	}
	paths, ok := document["paths"].(map[string]any)
	if !ok || len(paths) < 30 {
		t.Fatalf("expected complete API paths, got %d", len(paths))
	}
}
