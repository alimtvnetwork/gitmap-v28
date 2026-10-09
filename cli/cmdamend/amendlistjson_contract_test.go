package cmdamend

// JSON contract tests for `gitmap amend list --json`.
//
// amend-list emits an array of store.AmendmentRow. The contract
// covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: ID, Branch, FromCommit, ToCommit, TotalCommits,
//     PreviousName, PreviousEmail, NewName, NewEmail, Mode,
//     ForcePushed, CreatedAt.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdamend/ -run AmendListJSONContract
//
// NOTE: this test lives in cmdamend (not cmd) because the encoder
// (encodeAmendListJSON) is unexported here. The golden helpers below
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
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TestAmendListJSONContract_EmptyIsArrayNotNull is the jq-compat
// guarantee.
func TestAmendListJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	amAssertGoldenBytesDeterministic(t, "amend_list_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeAmendListJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalAmendmentRow builds a deterministic single row.
func canonicalAmendmentRow() store.AmendmentRow {
	return store.AmendmentRow{
		ID:            7,
		Branch:        "main",
		FromCommit:    "a1b2c3d",
		ToCommit:      "e4f5g6h",
		TotalCommits:  3,
		PreviousName:  "Alice Old",
		PreviousEmail: "alice.old@example.com",
		NewName:       "Alice New",
		NewEmail:      "alice.new@example.com",
		Mode:          "rewrite",
		ForcePushed:   1,
		CreatedAt:     "2025-01-01T12:00:00Z",
	}
}

// TestAmendListJSONContract_CanonicalRow_KeyOrders asserts the key
// order of the emitted object matches the schema declaration.
func TestAmendListJSONContract_CanonicalRow_KeyOrders(t *testing.T) {
	rows := []store.AmendmentRow{canonicalAmendmentRow()}
	var buf bytes.Buffer
	if err := encodeAmendListJSON(&buf, rows); err != nil {
		t.Fatalf("encode: %v", err)
	}

	amAssertSchemaKeysFirstObject(t, buf.Bytes(), "amend-list")
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---

func amPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("amPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

func amResolveGoldenPath(name string) string {
	return filepath.Join(amPackageDir(), "testdata", name)
}

func amAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := amResolveGoldenPath(name)
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

func amAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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

	amAssertGoldenBytes(t, name, first)
}

// amAssertSchemaKeysFirstObject checks the key order of the first array
// element against testdata/schemas/<name>.v1.json.
func amAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(amPackageDir(), "testdata", "schemas", name+".v1.json")
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

	got := amFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

// amFirstObjectKeys returns the keys of the first JSON object in raw
// (skipping a top-level array wrapper), in wire order.
func amFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("amFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("amFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("amFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("amFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
