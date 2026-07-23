package cmd

// JSON schema contract for `gitmap temp-releaselist --json`. Pairs the
// runtime encoder (encodeTempReleaseListJSON / buildTempReleaseListItems
// in tempreleaselistrender.go) with the published schema at
// spec/08-json-schemas/temp-release-list.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"sort"
	"testing"
)

const tempReleaseListSchemaFilename = "temp-release-list.schema.json"

// tempReleaseListItemsRequiredKeys mirrors the schema's items.required array.
var tempReleaseListItemsRequiredKeys = []string{
	"branch",
	"commit",
	"commitMessage",
	"createdAt",
	"id",
	"sequenceNumber",
	"versionPrefix",
}

// TestTempReleaseListJSONSchema_TopLevelShape pins the root type (array)
// and the items.required key set against the schema.
func TestTempReleaseListJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, tempReleaseListSchemaFilename)
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
	if !equalStringSlices(got, tempReleaseListItemsRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, tempReleaseListItemsRequiredKeys)
	}
}

// TestTempReleaseListJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder, then asserts every key in the first emitted object is declared in
// the schema's items.properties map.
func TestTempReleaseListJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, tempReleaseListSchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	releases := canonicalTempReleaseList(t)
	var buf bytes.Buffer
	if err := encodeTempReleaseListJSON(&buf, releases); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}
