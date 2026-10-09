package cmdbookmark

// JSON contract tests for `gitmap bookmark list --json`.
//
// bookmark-list emits an array of model.BookmarkRecord. The contract
// covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: id, name, command, args, flags, createdAt.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdbookmark/ -run BookmarkListJSONContract
//
// NOTE: this test lives in cmdbookmark (not cmd) because the encoder
// (encodeBookmarkListJSON) is unexported here. The golden helpers below
// are local, minimal equivalents of the cmd package's golden-test
// infrastructure.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/goldenguard"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// TestBookmarkListJSONContract_EmptyIsArrayNotNull is the jq-compat guarantee.
func TestBookmarkListJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	bmAssertGoldenBytesDeterministic(t, "bookmark_list_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeBookmarkListJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalBookmarkRecord builds a deterministic single row.
func canonicalBookmarkRecord() model.BookmarkRecord {
	return model.BookmarkRecord{
		ID:        7,
		Name:      "my-bookmark",
		Command:   "scan",
		Args:      "--depth=3",
		Flags:     "--verbose",
		CreatedAt: "2025-01-01T12:00:00Z",
	}
}

// TestBookmarkListJSONContract_CanonicalRow_KeyOrder asserts the key
// order of the emitted object matches the schema declaration.
func TestBookmarkListJSONContract_CanonicalRow_KeyOrder(t *testing.T) {
	records := []model.BookmarkRecord{canonicalBookmarkRecord()}
	var buf bytes.Buffer
	if err := encodeBookmarkListJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}

	bmAssertSchemaKeysFirstObject(t, buf.Bytes(), "bookmark-list")
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---

func bmPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("bmPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

func bmResolveGoldenPath(name string) string {
	return filepath.Join(bmPackageDir(), "testdata", name)
}

func bmAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := bmResolveGoldenPath(name)
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

func bmAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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

	bmAssertGoldenBytes(t, name, first)
}

// bmAssertSchemaKeysFirstObject checks the key order of the first array
// element against testdata/schemas/<name>.v1.json.
func bmAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(bmPackageDir(), "testdata", "schemas", name+".v1.json")
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

	got := bmFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

// bmFirstObjectKeys returns the keys of the first JSON object in raw
// (skipping a top-level array wrapper), in wire order.
func bmFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("bmFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("bmFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("bmFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("bmFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
