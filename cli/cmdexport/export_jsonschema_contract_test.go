package cmdexport

// JSON schema contract for `gitmap export`. Pairs the runtime encoder
// (encodeDatabaseExportJSON in exportrender.go) with the published
// schema at 02-spec/08-json-schemas/export.schema.json so drift in either
// side fails the build.
//
// NOTE: this test lives in cmdexport (not cmd) because the encoder
// (encodeDatabaseExportJSON) is unexported here. Schema helpers are
// local, minimal equivalents of the cmd package's schema-test
// infrastructure.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

const exportSchemaFilename = "export.schema.json"

// exportRequiredKeys mirrors the schema's required array (sorted for
// the slice comparison below).
var exportRequiredKeys = []string{
	"bookmarks", "exportedAt", "groups", "history", "releases", "repos", "version",
}

// TestExportJSONSchema_TopLevelShape pins the root type (object) and
// the required key set against the schema.
func TestExportJSONSchema_TopLevelShape(t *testing.T) {
	root := exLoadSchemaFile(t, exportSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}

	got := exStringSliceFromAny(root["required"])
	sort.Strings(got)
	if !exEqualStringSlices(got, exportRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, exportRequiredKeys)
	}
}

// TestExportJSONSchema_EncoderMatchesSchema runs the real stablejson
// encoder, then asserts every key in the emitted top-level object is
// declared in the schema's properties map.
func TestExportJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := exLoadSchemaFile(t, exportSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	var buf bytes.Buffer
	if err := encodeDatabaseExportJSON(&buf, model.DatabaseExport{}); err != nil {
		t.Fatalf("encode: %v", err)
	}

	gotKeys := exFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}

// --- local schema-test helpers (minimal equivalents of cmd's infra) ---

func exPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("exPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

// exFindPublishedSchema walks up from the package dir to locate
// 02-spec/08-json-schemas/<filename>.
func exFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(exPackageDir())
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

	t.Fatalf("published schema %s not found walking up from %s", filename, exPackageDir())

	return ""
}

func exLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(exFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

func exStringSliceFromAny(v any) []string {
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

func exEqualStringSlices(a, b []string) bool {
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

// exFirstObjectKeys returns the keys of the first JSON object in raw,
// in wire order.
func exFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("exFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("exFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("exFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("exFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
