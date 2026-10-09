package cmdlist

// JSON contract tests for `gitmap list-versions --json`.
//
// list-versions emits an array of {version, source?, changelog?}.
// The contract covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: version, source, changelog.
//   - source/changelog are omitted when empty (legacy omitempty shape).
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdlist/ -run ListVersionsJSONContract
//
// NOTE: this test lives in cmdlist (not cmd) because the encoder and its
// type (versionEntry, encodeListVersionsJSON) are unexported here. The
// golden helpers below are local, minimal equivalents of the cmd
// package's golden-test infrastructure.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/goldenguard"
	"github.com/alimtvnetwork/gitmap-v28/cli/release"
)

// TestListVersionsJSONContract_EmptyIsArrayNotNull is the jq-compat guarantee.
func TestListVersionsJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	lvAssertGoldenBytesDeterministic(t, "list_versions_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeListVersionsJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalListVersionsEntries builds a deterministic two-row fixture:
// the first row exercises all fields; the second omits both optional
// fields so the omitempty contract is pinned.
func canonicalListVersionsEntries(t *testing.T) []versionEntry {
	t.Helper()
	v1, err := release.Parse("v5.72.0")
	if err != nil {
		t.Fatalf("parse v5.72.0: %v", err)
	}

	v2, err := release.Parse("v5.71.0")
	if err != nil {
		t.Fatalf("parse v5.71.0: %v", err)
	}

	return []versionEntry{
		{Version: v1, Source: "local", Notes: []string{"first note", "second note"}},
		{Version: v2},
	}
}

// TestListVersionsJSONContract_CanonicalRows pins the bytes of a
// canonical two-row sample, locking key order AND the omitempty
// behavior for source/changelog.
func TestListVersionsJSONContract_CanonicalRows(t *testing.T) {
	entries := canonicalListVersionsEntries(t)
	lvAssertGoldenBytesDeterministic(t, "list_versions_canonical.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeListVersionsJSON(&buf, entries)

		return buf.Bytes(), err
	})
}

// TestListVersionsJSONContract_KeyOrder asserts the first object's
// key order matches the schema registry declaration.
func TestListVersionsJSONContract_KeyOrder(t *testing.T) {
	entries := canonicalListVersionsEntries(t)
	var buf bytes.Buffer
	if err := encodeListVersionsJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}

	lvAssertSchemaKeysFirstObject(t, buf.Bytes(), "list-versions")
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---

func lvPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("lvPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

func lvResolveGoldenPath(name string) string {
	return filepath.Join(lvPackageDir(), "testdata", name)
}

func lvAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := lvResolveGoldenPath(name)
	trigger := os.Getenv("GITMAP_UPDATE_GOLDEN") == "1"
	if goldenguard.AllowUpdate(t, trigger) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}

		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("write golden %s: %v", path, err)
		}

		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v (run with GITMAP_UPDATE_GOLDEN=1 and GITMAP_ALLOW_GOLDEN_UPDATE=1 to create)", path, err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("golden mismatch for %s\n--- want\n%s\n--- got\n%s", name, string(want), string(got))
	}
}

func lvAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
	t.Helper()
	first, err := encode()
	if err != nil {
		t.Fatalf("%s: encode run 0: %v", name, err)
	}

	for i := 1; i < 3; i++ {
		got, err := encode()
		if err != nil {
			t.Fatalf("%s: encode run %d: %v", name, i, err)
		}

		if !bytes.Equal(got, first) {
			t.Fatalf("%s: determinism broken — run %d differs from run 0", name, i)
		}
	}

	lvAssertGoldenBytes(t, name, first)
}

// lvAssertSchemaKeysFirstObject checks the top-level key order of the
// first array element against testdata/schemas/<name>.v1.json.
func lvAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(lvPackageDir(), "testdata", "schemas", name+".v1.json")
	schemaRaw, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema %s: %v", schemaPath, err)
	}

	var schemaDoc struct {
		Keys []string `json:"keys"`
	}

	if err := json.Unmarshal(schemaRaw, &schemaDoc); err != nil {
		t.Fatalf("parse schema %s: %v", schemaPath, err)
	}

	got := lvFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

// lvFirstObjectKeys returns the keys of the first JSON object in raw
// (skipping a top-level array wrapper), in wire order.
func lvFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("lvFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("lvFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("lvFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("lvFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
