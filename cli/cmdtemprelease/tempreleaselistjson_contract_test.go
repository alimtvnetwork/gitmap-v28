package cmdtemprelease

// JSON contract tests for `gitmap temp-releaselist --json`.
//
// temp-releaselist emits an array of temporary release branch records.
// The contract covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: id, branch, versionPrefix, sequenceNumber,
//     commit, commitMessage, createdAt.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run TempReleaseListJSONContract

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// canonicalTempReleaseList builds a deterministic two-row fixture.
func canonicalTempReleaseList(t *testing.T) []model.TempRelease {
	t.Helper()

	return []model.TempRelease{
		{
			ID:             1,
			Branch:         "release/v5.76.0-temp-3",
			VersionPrefix:  "v5.76",
			SequenceNumber: 3,
			CommitSha:      "abc123def456",
			CommitMessage:  "feat: add version-history stablejson encoder",
			CreatedAt:      "2026-05-26T10:00:00Z",
		},
		{
			ID:             2,
			Branch:         "release/v5.76.0-temp-4",
			VersionPrefix:  "v5.76",
			SequenceNumber: 4,
			CommitSha:      "def789abc012",
			CommitMessage:  "fix: correct edge case in history rewrite",
			CreatedAt:      "2026-05-26T11:30:00Z",
		},
	}
}

// TestTempReleaseListJSONContract_EmptyIsArrayNotNull is the jq-compat guarantee.
func TestTempReleaseListJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	trlAssertGoldenBytesDeterministic(t, "temp_release_list_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeTempReleaseListJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// TestTempReleaseListJSONContract_CanonicalRows pins the bytes of a
// canonical two-row sample, locking key order.
func TestTempReleaseListJSONContract_CanonicalRows(t *testing.T) {
	releases := canonicalTempReleaseList(t)
	trlAssertGoldenBytesDeterministic(t, "temp_release_list_canonical.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeTempReleaseListJSON(&buf, releases)

		return buf.Bytes(), err
	})
}

// TestTempReleaseListJSONContract_KeyOrder asserts the first object's
// key order matches the schema registry declaration.
func TestTempReleaseListJSONContract_KeyOrder(t *testing.T) {
	releases := canonicalTempReleaseList(t)
	var buf bytes.Buffer
	if err := encodeTempReleaseListJSON(&buf, releases); err != nil {
		t.Fatalf("encode: %v", err)
	}

	trlAssertSchemaKeysFirstObject(t, buf.Bytes(), "temp-release-list")
}

// --- local test helpers (minimal equivalents of cmd's infra) ---

func trlPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("trlPackageDir: runtime.Caller failed")
	}
	return filepath.Dir(file)
}

func trlAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join(trlPackageDir(), "testdata", name)
	if os.Getenv("GITMAP_UPDATE_GOLDEN") == "1" && os.Getenv("GITMAP_ALLOW_GOLDEN_UPDATE") == "1" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
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

func trlAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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
	trlAssertGoldenBytes(t, name, first)
}

func trlAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(trlPackageDir(), "testdata", "schemas", name+".v1.json")
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
	got := trlFirstObjectKeys(t, raw)
	if len(got) != len(doc.Keys) {
		t.Fatalf("key count: got %v, want %v", got, doc.Keys)
	}
	for i := range got {
		if got[i] != doc.Keys[i] {
			t.Fatalf("key order: got %v, want %v", got, doc.Keys)
		}
	}
}

func trlFirstObjectKeys(t *testing.T, raw []byte) []string {
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
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("keys: %v", err)
		}
		k, _ := tok.(string)
		keys = append(keys, k)
		var skip any
		_ = dec.Decode(&skip)
	}
	return keys
}
