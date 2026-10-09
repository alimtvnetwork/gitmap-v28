package cmdstats

// JSON schema contract for `gitmap stats --json`. Pairs the runtime
// encoder (encodeStatsJSON / build*Items in statsrender.go) with the
// published schema at 02-spec/08-json-schemas/stats.schema.json so drift
// in either side fails the build.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

const statsSchemaFilename = "stats.schema.json"

// statsTopLevelRequiredKeys mirrors the schema's top-level required array.
var statsTopLevelRequiredKeys = []string{
	"avgDurationMs",
	"commands",
	"overallFailRate",
	"totalCommands",
	"totalFail",
	"totalSuccess",
	"uniqueCommands",
}

// TestStatsJSONSchema_TopLevelShape pins the root type (object) and
// the required-keys set against the schema.
func TestStatsJSONSchema_TopLevelShape(t *testing.T) {
	root := stLoadSchemaFile(t, statsSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}

	got := stStringSliceFromAny(root["required"])
	sort.Strings(got)
	if !stEqualStringSlices(got, statsTopLevelRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, statsTopLevelRequiredKeys)
	}
}

// TestStatsJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder, then asserts every key in the emitted top-level object is
// declared in the schema's properties map.
func TestStatsJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := stLoadSchemaFile(t, statsSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	overall, commands := canonicalStatsOverall(t)
	var buf bytes.Buffer
	if err := encodeStatsJSON(&buf, overall, commands); err != nil {
		t.Fatalf("encode: %v", err)
	}

	gotKeys := stReadFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}

// --- local schema helpers (minimal equivalents of cmd's infra) ---

func stFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(stPackageDir())
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

func stLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(stFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	return s
}

func stStringSliceFromAny(v any) []string {
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

func stEqualStringSlices(a, b []string) bool {
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

func stReadFirstObjectKeys(t *testing.T, raw []byte) []string {
	return stFirstObjectKeys(t, raw)
}
