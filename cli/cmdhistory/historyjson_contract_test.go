package cmdhistory

// JSON contract tests for `gitmap history --json`.
//
// history emits an array of model.CommandHistoryRecord. The contract
// covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: id, command, alias, args, flags, startedAt,
//     finishedAt, durationMs, exitCode, summary, repoCount, createdAt.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdhistory/ -run HistoryJSONContract
//
// NOTE: this test lives in cmdhistory (not cmd) because the encoder
// (encodeHistoryJSON) is unexported here. The golden helpers below
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

// TestHistoryJSONContract_EmptyIsArrayNotNull is the jq-compat
// guarantee: zero rows must encode as `[]\n` even when the input
// slice is nil.
func TestHistoryJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	hiAssertGoldenBytesDeterministic(t, "history_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeHistoryJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalHistoryRecord builds a deterministic single record whose
// every field is a fixed value, so the golden file's bytes are stable
// across machines and time. Used for the byte-exact + key-order
// tests below.
func canonicalHistoryRecord() model.CommandHistoryRecord {
	return model.CommandHistoryRecord{
		ID:         42,
		Command:    "scan",
		Alias:      "s",
		Args:       "all",
		Flags:      "--depth=3",
		StartedAt:  "2025-01-01T12:00:00Z",
		FinishedAt: "2025-01-01T12:00:05Z",
		DurationMs: 5123,
		ExitCode:   0,
		Summary:    "Scanned 7 repos, 2 dirty",
		RepoCount:  7,
		CreatedAt:  "2025-01-01T12:00:00Z",
	}
}

// TestHistoryJSONContract_CanonicalRow_KeyOrders asserts the
// key order of the emitted object matches the schema declaration.
// Structural-only (no byte-exact golden for the populated row) so
// the test stays robust against future numeric formatting changes
// in encoding/json or value-shape tweaks.
func TestHistoryJSONContract_CanonicalRow_KeyOrders(t *testing.T) {
	records := []model.CommandHistoryRecord{canonicalHistoryRecord()}
	var buf bytes.Buffer
	if err := encodeHistoryJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}

	hiAssertSchemaKeysFirstObject(t, buf.Bytes(), "history")
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---

func hiPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("hiPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

func hiResolveGoldenPath(name string) string {
	return filepath.Join(hiPackageDir(), "testdata", name)
}

func hiAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := hiResolveGoldenPath(name)
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

func hiAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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

	hiAssertGoldenBytes(t, name, first)
}

// hiAssertSchemaKeysFirstObject checks the key order of the first array
// element against testdata/schemas/<name>.v1.json.
func hiAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(hiPackageDir(), "testdata", "schemas", name+".v1.json")
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

	got := hiFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

// hiFirstObjectKeys returns the keys of the first JSON object in raw
// (skipping a top-level array wrapper), in wire order.
func hiFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("hiFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("hiFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("hiFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("hiFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
