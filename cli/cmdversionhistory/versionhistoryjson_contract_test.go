package cmdversionhistory

// JSON contract tests for `gitmap version-history --json`.
//
// version-history emits an array of version transition records.
// The contract covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: fromVersionTag, fromVersionNum, toVersionTag,
//     toVersionNum, flattenedPath?, createdAt?, id, repoId.
//   - flattenedPath and createdAt are omitted when empty (legacy
//     omitempty shape).
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run VersionHistoryJSONContract

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// canonicalVersionHistoryRecords builds a deterministic two-row fixture:
// the first row exercises all fields; the second omits both optional
// fields so the omitempty contract is pinned.
func canonicalVersionHistoryRecords(t *testing.T) []model.RepoVersionHistoryRecord {
	t.Helper()

	return []model.RepoVersionHistoryRecord{
		{
			ID:             1,
			RepoID:         42,
			FromVersionTag: "v5.40.0",
			FromVersionNum: 50400,
			ToVersionTag:   "v5.41.0",
			ToVersionNum:   50410,
			FlattenedPath:  "/home/user/repos/gitmap-v28",
			CreatedAt:      "2026-05-20T14:30:00Z",
		},
		{
			ID:             2,
			RepoID:         42,
			FromVersionTag: "v5.41.0",
			FromVersionNum: 50410,
			ToVersionTag:   "v5.42.0",
			ToVersionNum:   50420,
		},
	}
}

// TestVersionHistoryJSONContract_EmptyIsArrayNotNull is the jq-compat guarantee.
func TestVersionHistoryJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	vhAssertGoldenBytesDeterministic(t, "version_history_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeVersionHistoryJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// TestVersionHistoryJSONContract_CanonicalRows pins the bytes of a
// canonical two-row sample, locking key order AND the omitempty
// behavior for flattenedPath/createdAt.
func TestVersionHistoryJSONContract_CanonicalRows(t *testing.T) {
	records := canonicalVersionHistoryRecords(t)
	vhAssertGoldenBytesDeterministic(t, "version_history_canonical.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeVersionHistoryJSON(&buf, records)

		return buf.Bytes(), err
	})
}

// TestVersionHistoryJSONContract_KeyOrder asserts the first object's
// key order matches the schema registry declaration.
func TestVersionHistoryJSONContract_KeyOrder(t *testing.T) {
	records := canonicalVersionHistoryRecords(t)
	var buf bytes.Buffer
	if err := encodeVersionHistoryJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}

	vhAssertSchemaKeysFirstObject(t, buf.Bytes(), "version-history")
}

// --- local test helpers (minimal equivalents of cmd's infra) ---

func vhPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("vhPackageDir")
	}
	return filepath.Dir(file)
}

func vhAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join(vhPackageDir(), "testdata", name)
	if os.Getenv("GITMAP_UPDATE_GOLDEN") == "1" && os.Getenv("GITMAP_ALLOW_GOLDEN_UPDATE") == "1" {
		_ = os.MkdirAll(filepath.Dir(path), 0o755)
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("golden mismatch for %s", name)
	}
}

func vhAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
	t.Helper()
	first, err := encode()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	for i := 1; i < 3; i++ {
		got, err := encode()
		if err != nil {
			t.Fatalf("encode run %d: %v", i, err)
		}
		if !bytes.Equal(got, first) {
			t.Fatalf("determinism broken for %s", name)
		}
	}
	vhAssertGoldenBytes(t, name, first)
}

func vhAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(vhPackageDir(), "testdata", "schemas", name+".v1.json")
	schemaRaw, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var doc struct {
		Keys []string `json:"keys"`
	}
	if err := json.Unmarshal(schemaRaw, &doc); err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	got := vhFirstObjectKeys(t, raw)
	if len(got) != len(doc.Keys) {
		t.Fatalf("key count: got %v, want %v", got, doc.Keys)
	}
	for i := range got {
		if got[i] != doc.Keys[i] {
			t.Fatalf("key order: got %v, want %v", got, doc.Keys)
		}
	}
}

func vhFirstObjectKeys(t *testing.T, raw []byte) []string {
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
