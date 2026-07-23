package cmd

// JSON-schema test helpers shared by all
// `*_jsonschema_contract_test.go` files. Pulled out of
// startuplist_jsonschema_contract_test.go so that file stays under
// the 200-line budget AND so future schemas (see
// spec/08-json-schemas/_TODO.md) can reuse the same primitives
// without copy-paste.
//
// Deliberately small surface area: only file-locator + schema-load
// + propertyOrder-extraction live here. Generic primitives like
// `equalStringSlices`, `expectDelim`, and `collectObjectKeys`
// already exist in jsonsnapshot_helpers_test.go and are reused.
//
// All helpers are test-only (suffix `_test.go`) so they don't bloat
// the production binary.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// findSchemaFile resolves a schema filename to an absolute path by
// walking up from the test's CWD (Go sets it to the package dir,
// i.e. gitmap/cmd) until it finds a `spec/08-json-schemas/`
// sibling. Same idiom used by other gitmap contract tests that need
// to read project-relative fixtures.
func findSchemaFile(t *testing.T, filename string) string {
	t.Helper()
	dir := filepath.Dir(cmdPackageDir())
	for i := 0; i < 8; i++ {
		candidate := filepath.Join(dir, "spec", "08-json-schemas", filename)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("could not locate spec/08-json-schemas/%s walking up from %s", filename, dir)

	return ""
}

// loadSchemaFile reads + parses a schema file into a generic map so
// the test can read both standard JSON Schema fields AND our
// `propertyOrder` extension without binding to a struct.
func loadSchemaFile(t *testing.T, filename string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(findSchemaFile(t, filename))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var s map[string]any
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	return s
}

// stringSliceFromAny converts a JSON-unmarshaled []any into
// []string. Returns nil on any non-string element so the caller's
// equality check fails loudly rather than silently coercing.
func stringSliceFromAny(v any) []string {
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

// extractFirstObjectKeyOrder uses json.Decoder's token stream to
// recover the literal on-the-wire key order of the first object
// inside a top-level array. Standard json.Unmarshal into
// map[string]any would lose ordering (Go maps are unordered); the
// Token API preserves it because it walks the raw bytes
// left-to-right. Reuses `expectDelim` and `collectObjectKeys` from
// jsonsnapshot_helpers_test.go to keep the helper surface small.
func extractFirstObjectKeyOrder(t *testing.T, data []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(data))
	if err := skipUntilFirstObjectStart(dec); err != nil {
		t.Fatalf("opening first object: %v", err)
	}

	return collectObjectKeys(t, dec)
}
