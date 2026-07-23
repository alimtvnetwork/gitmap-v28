package cmd

// JSON schema contract for `gitmap version-history --json`. Pairs the
// runtime encoder (encodeVersionHistoryJSON / buildVersionHistoryItems
// in versionhistoryrender.go) with the published schema at
// spec/08-json-schemas/version-history.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"sort"
	"testing"
)

const versionHistorySchemaFilename = "version-history.schema.json"

// versionHistoryItemsRequiredKeys mirrors the schema's items.required array.
var versionHistoryItemsRequiredKeys = []string{
	"fromVersionNum",
	"fromVersionTag",
	"id",
	"repoId",
	"toVersionNum",
	"toVersionTag",
}

// TestVersionHistoryJSONSchema_TopLevelShape pins the root type (array)
// and the items.required key set against the schema.
func TestVersionHistoryJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, versionHistorySchemaFilename)
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
	if !equalStringSlices(got, versionHistoryItemsRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, versionHistoryItemsRequiredKeys)
	}
}

// TestVersionHistoryJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder, then asserts every key in the first emitted object is declared in
// the schema's items.properties map.
func TestVersionHistoryJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, versionHistorySchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	records := canonicalVersionHistoryRecords(t)
	var buf bytes.Buffer
	if err := encodeVersionHistoryJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}
