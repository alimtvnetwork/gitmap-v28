package cmd

// JSON schema contract for `gitmap watch --json`. Pairs the runtime
// encoder (encodeWatchJSON / renderWatchReposRaw / renderWatchSummaryRaw
// in watchrender.go) with the published schema at
// spec/08-json-schemas/watch.schema.json so drift in either side fails
// the build.

import (
	"bytes"
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
	root := loadSchemaFile(t, watchSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}
	got := stringSliceFromAny(root["required"])
	sort.Strings(got)
	if !equalStringSlices(got, watchTopLevelRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, watchTopLevelRequiredKeys)
	}
}

// TestWatchJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder and asserts every top-level key is declared in the schema's
// properties map.
func TestWatchJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := loadSchemaFile(t, watchSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	var buf bytes.Buffer
	if err := encodeWatchJSON(&buf, nil, watchSummary{}, "2025-01-01T12:00:00Z"); err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotKeys := readFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}

// TestWatchJSONSchema_RepoItemShape pins the items.required key set
// for the nested repo objects.
func TestWatchJSONSchema_RepoItemShape(t *testing.T) {
	root := loadSchemaFile(t, watchSchemaFilename)
	repos, _ := root["properties"].(map[string]any)["repos"].(map[string]any)
	items, _ := repos["items"].(map[string]any)
	got := stringSliceFromAny(items["required"])
	sort.Strings(got)
	want := []string{"ahead", "behind", "branch", "name", "path", "stash", "status"}
	if !equalStringSlices(got, want) {
		t.Fatalf("repo items.required = %v, want %v", got, want)
	}
}

// TestWatchJSONSchema_SummaryShape pins the summary.required key set.
func TestWatchJSONSchema_SummaryShape(t *testing.T) {
	root := loadSchemaFile(t, watchSchemaFilename)
	summary, _ := root["properties"].(map[string]any)["summary"].(map[string]any)
	got := stringSliceFromAny(summary["required"])
	sort.Strings(got)
	want := []string{"behind", "dirty", "stash", "total"}
	if !equalStringSlices(got, want) {
		t.Fatalf("summary.required = %v, want %v", got, want)
	}
}
