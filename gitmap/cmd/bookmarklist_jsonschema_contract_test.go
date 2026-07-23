package cmd

// JSON schema contract for `gitmap bookmark list --json`. Pairs the
// runtime encoder (encodeBookmarkListJSON / buildBookmarkListJSONItems
// in bookmarklistrender.go) with the published schema at
// spec/08-json-schemas/bookmark-list.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

const bookmarkListSchemaFilename = "bookmark-list.schema.json"

// bookmarkListTopLevelRequiredKeys mirrors the schema's items.required array.
var bookmarkListTopLevelRequiredKeys = []string{
	"command",
	"id",
	"name",
}

// TestBookmarkListJSONSchema_TopLevelShape pins the root type (array)
// and the items.required key set against the schema.
func TestBookmarkListJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, bookmarkListSchemaFilename)
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
	if !equalStringSlices(got, bookmarkListTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, bookmarkListTopLevelRequiredKeys)
	}
}

// TestBookmarkListJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map.
func TestBookmarkListJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, bookmarkListSchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	records := []model.BookmarkRecord{canonicalBookmarkRecord()}
	var buf bytes.Buffer
	if err := encodeBookmarkListJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}
