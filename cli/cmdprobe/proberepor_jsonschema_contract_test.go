package cmdprobe

// JSON schema contract for `gitmap probe --json`. Pairs the runtime
// encoder (encodeProbeJSON / buildProbeJSONItems in proberender.go)
// with the published schema at
// 02-spec/08-json-schemas/probe-report.schema.json so drift in either
// side fails the build.
//
// NOTE: this test lives in cmdprobe (not cmd) because the encoder
// (encodeProbeJSON) is unexported here. Schema helpers are shared
// with probereporjson_contract_test.go in this package.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

const probeSchemaFilename = "probe-report.schema.json"

// probeTopLevelRequiredKeys mirrors the schema's items.required
// array. Centralized so the assertion below is a one-line diff
// against the on-disk schema rather than a re-typed literal.
var probeTopLevelRequiredKeys = []string{
	"absolutePath",
	"isAvailable",
	"method",
	"nextVersionNum",
	"nextVersionTag",
	"repoId",
	"slug",
}

// TestProbeJSONSchema_TopLevelShape pins the root type (array)
// and the items.required key set against the schema.
func TestProbeJSONSchema_TopLevelShape(t *testing.T) {
	root := prLoadSchemaFile(t, probeSchemaFilename)
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

	got := prStringSliceFromAny(items["required"])
	sort.Strings(got)
	if !prEqualStringSlices(got, probeTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, probeTopLevelRequiredKeys)
	}
}

// TestProbeJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder, then asserts every key in the first emitted object is
// declared in the schema's items.properties map.
func TestProbeJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := prLoadSchemaFile(t, probeSchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	entries := []probeJSONEntry{canonicalProbeEntry()}
	var buf bytes.Buffer
	if err := encodeProbeJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}

	gotKeys := prFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}

// --- local schema-test helpers (minimal equivalents of cmd's infra) ---
// prPackageDir and prFirstObjectKeys are defined in
// probereporjson_contract_test.go.

// prFindPublishedSchema walks up from the package dir to locate
// 02-spec/08-json-schemas/<filename>.
func prFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(prPackageDir())
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

	t.Fatalf("published schema %s not found walking up from %s", filename, prPackageDir())

	return ""
}

func prLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(prFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

func prStringSliceFromAny(v any) []string {
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

func prEqualStringSlices(a, b []string) bool {
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
