package cmdbookmark

// JSON schema contract for `gitmap bookmark list --json`. Pairs the
// runtime encoder (encodeBookmarkListJSON / buildBookmarkListJSONItems
// in bookmarklistrender.go) with the published schema at
// 02-spec/08-json-schemas/bookmark-list.schema.json so drift in either
// side fails the build.
//
// NOTE: this test lives in cmdbookmark (not cmd) because the encoder
// (encodeBookmarkListJSON) is unexported here. Schema helpers are shared
// with bookmarklistjson_contract_test.go in this package.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
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
	root := bmLoadSchemaFile(t, bookmarkListSchemaFilename)
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

	got := bmStringSliceFromAny(items["required"])
	sort.Strings(got)
	if !bmEqualStringSlices(got, bookmarkListTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, bookmarkListTopLevelRequiredKeys)
	}
}

// TestBookmarkListJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map.
func TestBookmarkListJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := bmLoadSchemaFile(t, bookmarkListSchemaFilename)
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

	gotKeys := bmFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}

// --- local schema-test helpers (minimal equivalents of cmd's infra) ---

// bmFindPublishedSchema walks up from the package dir to locate
// 02-spec/08-json-schemas/<filename>.
func bmFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(bmPackageDir())
	for i := 0; i < 8; i++ {
		candidate := filepath.Join(dir, "02-spec", "08-json-schemas", filename)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
	}

	t.Fatalf("published schema %s not found walking up from %s", filename, bmPackageDir())

	return ""
}

func bmLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(bmFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

func bmStringSliceFromAny(v any) []string {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}

	out := make([]string, 0, len(arr))
	for _, e := range arr {
		s, ok := e.(string)
		if !ok {
			return nil
		}

		out = append(out, s)
	}

	return out
}

func bmEqualStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
