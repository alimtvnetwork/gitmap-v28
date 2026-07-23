package cmd

// JSON schema contract for `gitmap export`. Pairs the runtime encoder
// (encodeDatabaseExportJSON in exportrender.go) with the published
// schema at spec/08-json-schemas/export.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

const exportSchemaFilename = "export.schema.json"

// exportRequiredKeys mirrors the schema's required array (sorted for
// the slice comparison below).
var exportRequiredKeys = []string{
	"bookmarks", "exportedAt", "groups", "history", "releases", "repos", "version",
}

// TestExportJSONSchema_TopLevelShape pins the root type (object) and
// the required key set against the schema.
func TestExportJSONSchema_TopLevelShape(t *testing.T) {
	root := loadSchemaFile(t, exportSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}
	got := stringSliceFromAny(root["required"])
	sort.Strings(got)
	if !equalStringSlices(got, exportRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, exportRequiredKeys)
	}
}

// TestExportJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder, then asserts every key in the emitted top-level object is
// declared in the schema's properties map.
func TestExportJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, exportSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	var buf bytes.Buffer
	if err := encodeDatabaseExportJSON(&buf, model.DatabaseExport{}); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := extractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}
