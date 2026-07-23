package cmd

// JSON schema contract for `gitmap amend list --json`. Pairs the
// runtime encoder (encodeAmendListJSON / buildAmendListJSONItems in
// amendlistrender.go) with the published schema at
// spec/08-json-schemas/amend-list.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

const amendListSchemaFilename = "amend-list.schema.json"

// amendListTopLevelRequiredKeys mirrors the schema's items.required
// array.
var amendListTopLevelRequiredKeys = []string{
	"Branch",
	"CreatedAt",
	"ForcePushed",
	"FromCommit",
	"ID",
	"Mode",
	"NewEmail",
	"NewName",
	"PreviousEmail",
	"PreviousName",
	"ToCommit",
	"TotalCommits",
}

// TestAmendListJSONSchema_TopLevelShape pins the root type (array)
// and the items.required key set against the schema.
func TestAmendListJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, amendListSchemaFilename)
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
	if !equalStringSlices(got, amendListTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, amendListTopLevelRequiredKeys)
	}
}

// TestAmendListJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map.
func TestAmendListJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, amendListSchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	rows := []store.AmendmentRow{canonicalAmendmentRow()}
	var buf bytes.Buffer
	if err := encodeAmendListJSON(&buf, rows); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}
