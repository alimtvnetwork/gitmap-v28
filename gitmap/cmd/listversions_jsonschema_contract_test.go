package cmd

// JSON schema contract for `gitmap list-versions --json`. Pairs the
// runtime encoder (encodeListVersionsJSON / buildListVersionsJSONItems
// in listversionsrender.go) with the published schema at
// spec/08-json-schemas/list-versions.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"sort"
	"testing"
)

const listVersionsSchemaFilename = "list-versions.schema.json"

// listVersionsItemsRequiredKeys mirrors the schema's items.required array.
var listVersionsItemsRequiredKeys = []string{
	"version",
}

// TestListVersionsJSONSchema_TopLevelShape pins the root type (array)
// and the items.required key set against the schema.
func TestListVersionsJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, listVersionsSchemaFilename)
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
	if !equalStringSlices(got, listVersionsItemsRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, listVersionsItemsRequiredKeys)
	}
}

// TestListVersionsJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map.
func TestListVersionsJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, listVersionsSchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	entries := canonicalListVersionsEntries(t)
	var buf bytes.Buffer
	if err := encodeListVersionsJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}
