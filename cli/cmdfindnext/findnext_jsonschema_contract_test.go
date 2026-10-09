package cmdfindnext

// Schema contract for `gitmap find-next --json`. Pairs the runtime
// encoder (encodeFindNextJSON / buildFindNextJSONItems in
// findnextrender.go) with the published schema at
// 02-spec/08-json-schemas/find-next.schema.json so drift in either
// side fails the build.
//
// The sibling `findnextjson_contract_test.go` already pins the
// model.FindNextRow + nested model.ScanRecord field DECLARATION
// order via canonical-row goldens. This test layers the published
// JSON Schema on top so external consumers have a single
// authoritative document to validate against.
//
// NOTE: this test lives in cmdfindnext (not cmd) because the encoder
// (encodeFindNextJSON) is unexported here. Schema helpers are shared
// with findnextjson_contract_test.go in this package.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
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
	root := fnLoadSchemaFile(t, findNextSchemaFilename)
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

	got := fnStringSliceFromAny(items["required"])
	sort.Strings(got)
	if !fnEqualStringSlices(got, findNextTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, findNextTopLevelRequiredKeys)
	}
}

// TestFindNextJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map. Catches
// an unsynced field on either side (encoder added a column, or
// schema dropped one).
func TestFindNextJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := fnLoadSchemaFile(t, findNextSchemaFilename)
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

	gotKeys := fnFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}

// --- local schema-test helpers (minimal equivalents of cmd's infra) ---
// fnPackageDir is defined in findnextjson_contract_test.go.

// fnFindPublishedSchema walks up from the package dir to locate
// 02-spec/08-json-schemas/<filename>.
func fnFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(fnPackageDir())
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

	t.Fatalf("published schema %s not found walking up from %s", filename, fnPackageDir())

	return ""
}

func fnLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(fnFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

func fnStringSliceFromAny(v any) []string {
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

func fnEqualStringSlices(a, b []string) bool {
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
