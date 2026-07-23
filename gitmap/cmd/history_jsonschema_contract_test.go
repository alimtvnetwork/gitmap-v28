package cmd

// JSON schema contract for `gitmap history --json`. Pairs the
// runtime encoder (encodeHistoryJSON / buildHistoryJSONItems in
// historyrender.go) with the published schema at
// spec/08-json-schemas/history.schema.json so drift in either side
// fails the build.

import (
	"bytes"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

const historySchemaFilename = "history.schema.json"

// historyTopLevelRequiredKeys mirrors the schema's items.required
// array. Centralized so the assertion below is a one-line diff
// against the on-disk schema rather than a re-typed literal.
var historyTopLevelRequiredKeys = []string{
	"alias",
	"args",
	"command",
	"createdAt",
	"durationMs",
	"exitCode",
	"finishedAt",
	"flags",
	"id",
	"repoCount",
	"startedAt",
	"summary",
}

// TestHistoryJSONSchema_TopLevelShape pins the root type (array)
// and the items.required key set against the schema. A new
// top-level field added to the stablejson encoder without a matching
// schema update fails here.
func TestHistoryJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, historySchemaFilename)
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
	if !equalStringSlices(got, historyTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, historyTopLevelRequiredKeys)
	}
}

// TestHistoryJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map. Catches
// an unsynced field on either side (encoder added a column, or
// schema dropped one).
func TestHistoryJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, historySchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	records := []model.CommandHistoryRecord{canonicalHistoryRecord()}
	var buf bytes.Buffer
	if err := encodeHistoryJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}
