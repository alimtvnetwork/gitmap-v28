package cmdamend

// JSON contract tests for `gitmap amend audit` file output.
//
// The contract covers key order of the emitted single object:
//   id, timestamp, branch, fromCommit, toCommit, totalCommits,
//   previousAuthor, newAuthor, mode, forcePushed, commits.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdamend/ -run AmendAuditJSONContract
//
// NOTE: this test lives in cmdamend (not cmd) because the encoder
// (encodeAmendAuditJSON) is unexported here. The golden helpers below
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

// canonicalAmendmentRecord builds a deterministic single audit record.
func canonicalAmendmentRecord() model.AmendmentRecord {
	return model.AmendmentRecord{
		ID:           42,
		Timestamp:    "2025-01-01T12:00:00Z",
		Branch:       "main",
		FromCommit:   "abc123",
		ToCommit:     "def456",
		TotalCommits: 2,
		PreviousAuthor: model.AmendAuthor{
			Name:  "Old Name",
			Email: "old@example.com",
		},
		NewAuthor: model.AmendAuthor{
			Name:  "New Name",
			Email: "new@example.com",
		},
		Mode:          "all",
		IsForcePushed: true,
		Commits: []model.CommitEntry{
			{SHA: "abc123", Message: "First"},
			{SHA: "def456", Message: "Second"},
		},
	}
}

// TestAmendAuditJSONContract_CanonicalRecord_KeyOrder asserts the
// key order of the emitted object matches the schema declaration.
func TestAmendAuditJSONContract_CanonicalRecord_KeyOrder(t *testing.T) {
	encode := func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeAmendAuditJSON(&buf, canonicalAmendmentRecord())

		return buf.Bytes(), err
	}

	aaAssertGoldenBytesDeterministic(t, "amend_audit_canonical.json", encode)
	raw, _ := encode()
	aaAssertSchemaKeysFirstObject(t, raw, "amend-audit")
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---

func aaPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("aaPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

func aaResolveGoldenPath(name string) string {
	return filepath.Join(aaPackageDir(), "testdata", name)
}

func aaAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := aaResolveGoldenPath(name)
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

func aaAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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

	aaAssertGoldenBytes(t, name, first)
}

// aaAssertSchemaKeysFirstObject checks the top-level key order of the
// emitted object against testdata/schemas/<name>.v1.json.
func aaAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(aaPackageDir(), "testdata", "schemas", name+".v1.json")
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

	got := aaFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

// aaFirstObjectKeys returns the top-level keys of the first JSON object
// in raw, in wire order.
func aaFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("aaFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("aaFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("aaFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("aaFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
