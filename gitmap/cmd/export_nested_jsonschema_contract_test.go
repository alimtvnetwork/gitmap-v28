package cmd

// Per-record schema contract for `gitmap export` nested arrays.
// Builds a canonical fixture with one record per nested array, runs
// the live encoder, then asserts every emitted key on each record is
// declared in the corresponding `properties` map of the published
// schema. Catches struct-tag drift in either the model or the schema.

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// canonicalNonEmptyExport builds a deterministic export with one row
// per nested array — used by per-record drift assertions.
func canonicalNonEmptyExport() model.DatabaseExport {
	return model.DatabaseExport{
		Version:    "1",
		ExportedAt: "2026-05-26T12:00:00Z",
		Repos:      []model.ScanRecord{{ID: 1, Slug: "demo", RepoName: "demo"}},
		Groups: []model.GroupExport{{
			Group:     model.Group{ID: 1, Name: "default"},
			RepoSlugs: []string{"demo"},
		}},
		Releases:  []model.ReleaseRecord{{ID: 1, RepoID: 1, Version: "1.0.0"}},
		History:   []model.CommandHistoryRecord{{ID: 1, Command: "scan", StartedAt: "2026-05-26T12:00:00Z"}},
		Bookmarks: []model.BookmarkRecord{{ID: 1, Name: "fav", Command: "scan"}},
	}
}

// TestExportJSONSchema_NestedRecordKeysSubsetOfProperties asserts that
// every key emitted by the encoder for each nested array is declared
// in that array's `items.properties` map.
func TestExportJSONSchema_NestedRecordKeysSubsetOfProperties(t *testing.T) {
	root := loadSchemaFile(t, exportSchemaFilename)
	topProps, _ := root["properties"].(map[string]any)

	var buf bytes.Buffer
	if err := encodeDatabaseExportJSON(&buf, canonicalNonEmptyExport()); err != nil {
		t.Fatalf("encode: %v", err)
	}

	var top map[string]json.RawMessage
	if err := json.Unmarshal(buf.Bytes(), &top); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	for _, arrayKey := range []string{"repos", "groups", "releases", "history", "bookmarks"} {
		assertNestedRecordKeys(t, topProps, arrayKey, top[arrayKey])
	}
}

// assertNestedRecordKeys walks one nested array's first record and
// verifies every emitted key is in the schema's items.properties map.
func assertNestedRecordKeys(t *testing.T, topProps map[string]any, arrayKey string, raw json.RawMessage) {
	t.Helper()
	arrayProp, _ := topProps[arrayKey].(map[string]any)
	items, _ := arrayProp["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("%s: items.properties missing", arrayKey)
	}

	keys := readEveryObjectKeys(t, raw)
	if len(keys) == 0 {
		t.Fatalf("%s: expected at least one record", arrayKey)
	}
	for i, row := range keys {
		for _, key := range row {
			if _, allowed := props[key]; !allowed {
				t.Errorf("%s[%d]: encoder emitted %q not declared in schema", arrayKey, i, key)
			}
		}
	}
}
