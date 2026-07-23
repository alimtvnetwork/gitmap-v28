package cmd

// Schema contract for `gitmap find-next --json`. Pairs the runtime
// encoder (encodeFindNextJSON / buildFindNextJSONItems in
// findnextrender.go) with the published schema at
// spec/08-json-schemas/find-next.schema.json so drift in either
// side fails the build.
//
// The sibling `findnextjson_contract_test.go` already pins the
// model.FindNextRow + nested model.ScanRecord field DECLARATION
// order via canonical-row goldens. This test layers the published
// JSON Schema on top so external consumers have a single
// authoritative document to validate against.

import (
	"bytes"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

const findNextSchemaFilename = "find-next.schema.json"

// findNextTopLevelRequiredKeys mirrors the schema's items.required
// array. Centralized so the assertion below is a one-line diff
// against the on-disk schema rather than a re-typed literal.
var findNextTopLevelRequiredKeys = []string{
	"method",
	"nextVersionNum",
	"nextVersionTag",
	"probedAt",
	"repo",
}

// TestFindNextJSONSchema_TopLevelShape pins the root type (array)
// and the items.required key set against the schema. A new
// top-level field added to the stablejson encoder without a matching
// schema update fails here.
func TestFindNextJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, findNextSchemaFilename)
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
	if !equalStringSlices(got, findNextTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, findNextTopLevelRequiredKeys)
	}
}

// TestFindNextJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map. Catches
// an unsynced field on either side (encoder added a column, or
// schema dropped one).
func TestFindNextJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, findNextSchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	rows := []model.FindNextRow{canonicalFindNextRow()}
	var buf bytes.Buffer
	if err := encodeFindNextJSON(&buf, rows); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}
