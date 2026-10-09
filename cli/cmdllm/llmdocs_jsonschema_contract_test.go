package cmdllm

// JSON schema contract for `gitmap llm-docs --format=json`. Pairs the
// runtime encoder (encodeLLMDocsJSON in llmdocsrender.go) with the
// published schema at 02-spec/08-json-schemas/llm-docs.schema.json so
// drift in either side fails the build.
//
// NOTE: this test lives in cmdllm (not cmd) because the encoder
// (encodeLLMDocsJSON) is unexported here. Schema helpers are local,
// minimal equivalents of the cmd package's schema-test
// infrastructure.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const llmDocsSchemaFilename = "llm-docs.schema.json"

// TestLLMDocsJSONSchema_TopLevelShape pins the root type (object) and
// verifies every contracted top-level key is declared as a property.
func TestLLMDocsJSONSchema_TopLevelShape(t *testing.T) {
	root := llmLoadSchemaFile(t, llmDocsSchemaFilename)
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
	root := llmLoadSchemaFile(t, llmDocsSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	var buf bytes.Buffer
	if err := encodeLLMDocsJSON(&buf, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}

	gotKeys := llmFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}

// TestLLMDocsJSONSchema_NestedCommandShape asserts the nested
// command-group + per-command property declarations are present.
func TestLLMDocsJSONSchema_NestedCommandShape(t *testing.T) {
	root := llmLoadSchemaFile(t, llmDocsSchemaFilename)
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

// --- local schema-test helpers (minimal equivalents of cmd's infra) ---

func llmPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("llmPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

// llmFindPublishedSchema walks up from the package dir to locate
// 02-spec/08-json-schemas/<filename>.
func llmFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(llmPackageDir())
	for i := 0; i < 8; i++ {
		candidate := filepath.Join(dir, "02-spec", "08-json-schemas", filename)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
	}

	t.Fatalf("published schema %s not found walking up from %s", filename, llmPackageDir())

	return ""
}

func llmLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(llmFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

// llmFirstObjectKeys returns the keys of the first JSON object in raw,
// in wire order.
func llmFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("llmFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("llmFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("llmFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("llmFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
