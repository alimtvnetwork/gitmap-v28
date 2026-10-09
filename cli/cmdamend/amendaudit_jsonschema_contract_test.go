package cmdamend

// JSON schema contract for `gitmap amend audit` file output. Pairs
// the runtime encoder (encodeAmendAuditJSON / buildAuditRecord in
// amendaudit.go + amendauditrender.go) with the published schema at
// 02-spec/08-json-schemas/amend-audit.schema.json so drift in either
// side fails the build.
//
// NOTE: this test lives in cmdamend (not cmd) because the encoder
// (encodeAmendAuditJSON) is unexported here. The schema helpers below
// are local, minimal equivalents of the cmd package's schema-test
// infrastructure.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

const amendAuditSchemaFilename = "amend-audit.schema.json"

// amendAuditTopLevelRequiredKeys mirrors the schema's required array.
var amendAuditTopLevelRequiredKeys = []string{
	"branch",
	"commits",
	"forcePushed",
	"fromCommit",
	"id",
	"mode",
	"newAuthor",
	"previousAuthor",
	"timestamp",
	"toCommit",
	"totalCommits",
}

// TestAmendAuditJSONSchema_TopLevelShape pins the root type (object)
// and the required key set against the schema.
func TestAmendAuditJSONSchema_TopLevelShape(t *testing.T) {
	root := aaLoadSchemaFile(t, amendAuditSchemaFilename)
	if root["type"] != "object" {
		t.Fatalf("top-level type = %v, want object", root["type"])
	}

	got := aaStringSliceFromAny(root["required"])
	sort.Strings(got)
	if !aaEqualStringSlices(got, amendAuditTopLevelRequiredKeys) {
		t.Fatalf("required = %v, want %v", got, amendAuditTopLevelRequiredKeys)
	}
}

// TestAmendAuditJSONSchema_EncoderMatchesSchema runs the real
// stablejson encoder, then asserts every key in the emitted object
// is declared in the schema's properties map.
func TestAmendAuditJSONSchema_EncoderMatchesSchema(t *testing.T) {
	root := aaLoadSchemaFile(t, amendAuditSchemaFilename)
	props, ok := root["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties object")
	}

	record := canonicalAmendmentRecord()
	var buf bytes.Buffer
	if err := encodeAmendAuditJSON(&buf, record); err != nil {
		t.Fatalf("encode: %v", err)
	}

	gotKeys := aaFirstObjectKeys(t, buf.Bytes())
	for _, key := range gotKeys {
		if _, allowed := props[key]; !allowed {
			t.Errorf("encoder emitted %q not declared in schema properties", key)
		}
	}
}

// --- local schema-test helpers (minimal equivalents of cmd's infra) ---

// aaFindPublishedSchema walks up from the package dir to locate
// 02-spec/08-json-schemas/<filename>.
func aaFindPublishedSchema(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(aaPackageDir())
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

	t.Fatalf("published schema %s not found walking up from %s", filename, aaPackageDir())

	return ""
}

func aaLoadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(aaFindPublishedSchema(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

func aaStringSliceFromAny(v any) []string {
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

func aaEqualStringSlices(a, b []string) bool {
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
