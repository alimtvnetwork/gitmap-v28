package cmdwatch

// JSON schema contract for `gitmap watch --json`. Pairs the runtime
// encoder (encodeWatchJSON / renderWatchReposRaw / renderWatchSummaryRaw
// in watchrender.go) with the published schema at
// 02-spec/08-json-schemas/watch.schema.json so drift in either side fails
// the build.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

const watchSchemaFilename = "watch.schema.json"

// watchTopLevelRequiredKeys mirrors the schema's required array.
var watchTopLevelRequiredKeys = []string{
	"repos",
	"summary",
	"timestamp",
}

// TestWatchJSONSchema_TopLevelShape pins the root type (object) and
// the required key set against the schema.
func TestWatchJSONSchema_TopLevelShape(t *testing.T) {
	root := wLoadSchemaFile(t, watchSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}

	got := wStringSliceFromAny(root["required"])
	sort.Strings(got)
	if !wEqualStringSlices(got, watchTopLevelRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, watchTopLevelRequiredKeys)
	}
}

// TestWatchJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder and asserts every top-level key is declared in the schema's
// properties map.
func TestWatchJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := wLoadSchemaFile(t, watchSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	var buf bytes.Buffer
	if err := encodeWatchJSON(&buf, nil, watchSummary{}, "2025-01-01T12:00:00Z"); err != nil {
		t.Fatalf("encode: %v", err)
	}

	gotKeys := wReadFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}

// TestWatchJSONSchema_RepoItemShape pins the items.required key set
// for the nested repo objects.
func TestWatchJSONSchema_RepoItemShape(t *testing.T) {
	root := wLoadSchemaFile(t, watchSchemaFilename)
	repos, _ := root["properties"].(map[string]any)["repos"].(map[string]any)
	items, _ := repos["items"].(map[string]any)
	got := wStringSliceFromAny(items["required"])
	sort.Strings(got)
	want := []string{"ahead", "behind", "branch", "name", "path", "stash", "status"}
	if !wEqualStringSlices(got, want) {
		t.Fatalf("repo items.required = %v, want %v", got, want)
	}
}

// TestWatchJSONSchema_SummaryShape pins the summary.required key set.
func TestWatchJSONSchema_SummaryShape(t *testing.T) {
	root := wLoadSchemaFile(t, watchSchemaFilename)
	summary, _ := root["properties"].(map[string]any)["summary"].(map[string]any)
	got := wStringSliceFromAny(summary["required"])
	sort.Strings(got)
	want := []string{"behind", "dirty", "stash", "total"}
	if !wEqualStringSlices(got, want) {
		t.Fatalf("summary.required = %v, want %v", got, want)
	}
}

// --- local schema helpers (minimal equivalents of cmd's infra) ---

func wFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(wPackageDir())
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

func wLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(wFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	return s
}

func wStringSliceFromAny(v any) []string {
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

func wEqualStringSlices(a, b []string) bool {
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

func wReadFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))
	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("no object: %v", err)
		}
		if d, ok := tok.(json.Delim); ok && d == '{' {
			break
		}
	}
	var keys []string
	for dec.More() {
		tok, _ := dec.Token()
		k, _ := tok.(string)
		keys = append(keys, k)
		var skip any
		_ = dec.Decode(&skip)
	}
	return keys
}
