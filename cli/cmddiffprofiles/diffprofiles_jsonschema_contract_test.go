package cmddiffprofiles

// JSON schema contract for `gitmap diff-profiles --json`. Pairs
// the runtime encoder (encodeDiffProfilesJSON in diffprofilesrender.go)
// with the published schema at
// 02-spec/08-json-schemas/diff-profiles.schema.json so drift in either
// side fails the build.
//
// NOTE: this test lives in cmddiffprofiles (not cmd) because the
// encoder (encodeDiffProfilesJSON) is unexported here. Schema helpers
// are shared with diffprofilesjson_contract_test.go in this package.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

const diffProfilesSchemaFilename = "diff-profiles.schema.json"

// diffProfilesTopLevelRequiredKeys mirrors the schema's required array.
var diffProfilesTopLevelRequiredKeys = []string{
	"different",
	"onlyInA",
	"onlyInB",
	"profileA",
	"profileB",
	"same",
}

// TestDiffProfilesJSONSchema_TopLevelShape pins the root type
// (object) and the required key set against the schema.
func TestDiffProfilesJSONSchema_TopLevelShape(t *testing.T) {
	root := dpLoadSchemaFile(t, diffProfilesSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}

	got := dpStringSliceFromAny(root["required"])
	sort.Strings(got)
	if !dpEqualStringSlices(got, diffProfilesTopLevelRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, diffProfilesTopLevelRequiredKeys)
	}
}

// TestDiffProfilesJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the emitted object
// is declared in the schema's properties map.
func TestDiffProfilesJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := dpLoadSchemaFile(t, diffProfilesSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	result := canonicalDPResult()
	var buf bytes.Buffer
	if err := encodeDiffProfilesJSON(&buf, "alpha", "beta", result); err != nil {
		t.Fatalf("encode: %v", err)
	}

	gotKeys := dpFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}

// --- local schema-test helpers (minimal equivalents of cmd's infra) ---

// dpFindPublishedSchema walks up from the package dir to locate
// 02-spec/08-json-schemas/<filename>.
func dpFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(dpPackageDir())
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

	t.Fatalf("published schema %s not found walking up from %s", filename, dpPackageDir())

	return ""
}

func dpLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(dpFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

func dpStringSliceFromAny(v any) []string {
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

func dpEqualStringSlices(a, b []string) bool {
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
