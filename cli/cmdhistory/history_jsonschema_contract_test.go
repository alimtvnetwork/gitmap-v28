package cmdhistory

// JSON schema contract for `gitmap history --json`. Pairs the
// runtime encoder (encodeHistoryJSON / buildHistoryJSONItems in
// historyrender.go) with the published schema at
// 02-spec/08-json-schemas/history.schema.json so drift in either side
// fails the build.
//
// NOTE: this test lives in cmdhistory (not cmd) because the encoder
// (encodeHistoryJSON) is unexported here. Schema helpers are local,
// minimal equivalents of the cmd package's schema-test
// infrastructure.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

const historySchemaFilename = "history.schema.json"

// historyTopLevelRequiredKeys mirrors the schema's items.required
// array. Centralized so the assertion below is a one-line diff
// against the on-disk schema rather than a re-typed literal.
var historyTopLevelRequiredKeys = []string{
	"alias",
	"args",
	"command",
	"createdAt",
	"durationMs",
	"exitCode",
	"finishedAt",
	"flags",
	"id",
	"repoCount",
	"startedAt",
	"summary",
}

// TestHistoryJSONSchema_TopLevelShape pins the root type (array)
// and the items.required key set against the schema. A new
// top-level field added to the stablejson encoder without a matching
// schema update fails here.
func TestHistoryJSONSchema_TopLevelShape(t *testing.T) {
	root := hiLoadSchemaFile(t, historySchemaFilename)
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

	got := hiStringSliceFromAny(items["required"])
	sort.Strings(got)
	if !hiEqualStringSlices(got, historyTopLevelRequiredKeys) {
		t.Fatalf("items.required = %v, want %v", got, historyTopLevelRequiredKeys)
	}
}

// TestHistoryJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the first emitted
// object is declared in the schema's items.properties map. Catches
// an unsynced field on either side (encoder added a column, or
// schema dropped one).
func TestHistoryJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := hiLoadSchemaFile(t, historySchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing items.properties object")
	}

	records := []model.CommandHistoryRecord{canonicalHistoryRecord()}
	var buf bytes.Buffer
	if err := encodeHistoryJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}

	gotKeys := hiFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema items.properties", key)
		}
	}
}

// --- local schema-test helpers (minimal equivalents of cmd's infra) ---
// hiPackageDir is defined in historyjson_contract_test.go.

// hiFindPublishedSchema walks up from the package dir to locate
// 02-spec/08-json-schemas/<filename>.
func hiFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(hiPackageDir())
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

	t.Fatalf("published schema %s not found walking up from %s", filename, hiPackageDir())

	return ""
}

func hiLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(hiFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

func hiStringSliceFromAny(v any) []string {
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

func hiEqualStringSlices(a, b []string) bool {
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
