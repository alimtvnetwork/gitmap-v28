package cmdtemprelease

// JSON schema contract for `gitmap temp-releaselist --json`. Pairs the
// runtime encoder (encodeTempReleaseListJSON / buildTempReleaseListItems
// in tempreleaselistrender.go) with the published schema at
// 02-spec/08-json-schemas/temp-release-list.schema.json so drift in either
// side fails the build.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
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
	root := trlLoadSchemaFile(t, tempReleaseListSchemaFilename)
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

	got := trlStringSliceFromAny(items["required"])
	sort.Strings(got)
	if !trlEqualStringSlices(got, tempReleaseListItemsRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, tempReleaseListItemsRequiredKeys)
	}
}

// TestTempReleaseListJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder, then asserts every key in the first emitted object is declared in
// the schema's items.properties map.
func TestTempReleaseListJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := trlLoadSchemaFile(t, tempReleaseListSchemaFilename)
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

	gotKeys := trlExtractFirstObjectKeyOrder(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}

// --- local schema helpers (minimal equivalents of cmd's infra) ---

func trlFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(trlPackageDir())
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
	t.Fatalf("schema %s not found", filename)
	return ""
}

func trlLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(trlFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	return s
}

func trlStringSliceFromAny(v any) []string {
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

func trlEqualStringSlices(a, b []string) bool {
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

func trlExtractFirstObjectKeyOrder(t *testing.T, raw []byte) []string {
	return trlFirstObjectKeys(t, raw)
}
