package cmd

// JSON schema contract for `gitmap llm-docs --format=json`. Pairs the
// runtime encoder (encodeLLMDocsJSON in llmdocsrender.go) with the
// published schema at spec/08-json-schemas/llm-docs.schema.json so
// drift in either side fails the build.

import (
	"bytes"
	"testing"
)

const llmDocsSchemaFilename = "llm-docs.schema.json"

// TestLLMDocsJSONSchema_TopLevelShape pins the root type (object) and
// verifies every contracted top-level key is declared as a property.
func TestLLMDocsJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, llmDocsSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}
	for _, want := range []string{
		"commands", "architecture", "flags", "conventions",
		"structure", "database", "installation", "patterns",
	} {
		if _, ok := props[want]; !ok {
			t.Errorf("schema missing top-level property %q", want)
		}
	}
}

// TestLLMDocsJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder with all sections enabled, then asserts every emitted
// top-level key is declared in the schema's properties map.
func TestLLMDocsJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, llmDocsSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	var buf bytes.Buffer
	if err := encodeLLMDocsJSON(&buf, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}

// TestLLMDocsJSONSchema_NestedCommandShape asserts the nested
// command-group + per-command property declarations are present.
func TestLLMDocsJSONSchema_NestedCommandShape(t *testing.T) {
	root := loadSchemaFile(t, llmDocsSchemaFilename)
	props, _ := root["properties"].(map[string]any)
	commands, ok := props["commands"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing commands property")
	}
	groupItems, ok := commands["items"].(map[string]any)
	if !ok {
		t.Fatalf("commands.items missing")
	}
	groupProps, ok := groupItems["properties"].(map[string]any)
	if !ok {
		t.Fatalf("group items.properties missing")
	}
	for _, k := range []string{"title", "commands"} {
		if _, ok := groupProps[k]; !ok {
			t.Errorf("group missing property %q", k)
		}
	}
	cmdsField, _ := groupProps["commands"].(map[string]any)
	cmdItems, _ := cmdsField["items"].(map[string]any)
	cmdProps, _ := cmdItems["properties"].(map[string]any)
	for _, k := range []string{"name", "alias", "description", "example"} {
		if _, ok := cmdProps[k]; !ok {
			t.Errorf("per-command missing property %q", k)
		}
	}
}
