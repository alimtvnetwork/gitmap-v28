package cmdlist

// JSON schema contract for `gitmap list-versions --json`. Pairs the
// runtime encoder (encodeListVersionsJSON / buildListVersionsJSONItems
// in listversionsrender.go) with the published schema at
// 02-spec/08-json-schemas/list-versions.schema.json so drift in either
// side fails the build.
//
// NOTE: this test lives in cmdlist (not cmd) because the encoder
// (encodeListVersionsJSON) is unexported here. Schema helpers are
// shared with listversionsjson_contract_test.go in this package.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
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
	root := lvLoadSchemaFile(t, listVersionsSchemaFilename)
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

	got := lvStringSliceFromAny(items["required"])
	sort.Strings(got)
	if !lvEqualStringSlices(got, listVersionsItemsRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, listVersionsItemsRequiredKeys)
	}
}

// TestListVersionsJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map.
func TestListVersionsJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := lvLoadSchemaFile(t, listVersionsSchemaFilename)
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

	gotKeys := lvFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}

// --- local schema-test helpers (minimal equivalents of cmd's infra) ---
// lvPackageDir and lvFirstObjectKeys are defined in
// listversionsjson_contract_test.go.

// lvFindPublishedSchema walks up from the package dir to locate
// 02-spec/08-json-schemas/<filename>.
func lvFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(lvPackageDir())
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

	t.Fatalf("published schema %s not found walking up from %s", filename, lvPackageDir())

	return ""
}

func lvLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(lvFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

func lvStringSliceFromAny(v any) []string {
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

func lvEqualStringSlices(a, b []string) bool {
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
