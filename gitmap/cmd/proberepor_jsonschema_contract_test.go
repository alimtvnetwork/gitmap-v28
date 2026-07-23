package cmd

// JSON schema contract for `gitmap probe --json`. Pairs the runtime
// encoder (encodeProbeJSON / buildProbeJSONItems in proberender.go)
// with the published schema at
// spec/08-json-schemas/probe-report.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"sort"
	"testing"
)

const probeSchemaFilename = "probe-report.schema.json"

// probeTopLevelRequiredKeys mirrors the schema's items.required
// array. Centralized so the assertion below is a one-line diff
// against the on-disk schema rather than a re-typed literal.
var probeTopLevelRequiredKeys = []string{
	"absolutePath",
	"isAvailable",
	"method",
	"nextVersionNum",
	"nextVersionTag",
	"repoId",
	"slug",
}

// TestProbeJSONSchema_TopLevelShape pins the root type (array)
// and the items.required key set against the schema.
func TestProbeJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, probeSchemaFilename)
	if root["type"] != "array" {
		t.Fatalf("top-level type = %v, want array", root["type"])
	}
	items, ok := root["items"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items object")
	}
	if items["type"] != "object" {
		t.Fatalf("items.type = %v, want object", items["type"])
	}
	got := stringSliceFromAny(items["required"])
	sort.Strings(got)
	if !equalStringSlices(got, probeTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, probeTopLevelRequiredKeys)
	}
}

// TestProbeJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder, then asserts every key in the first emitted object is
// declared in the schema's items.properties map.
func TestProbeJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, probeSchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	entries := []probeJSONEntry{canonicalProbeEntry()}
	var buf bytes.Buffer
	if err := encodeProbeJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}
