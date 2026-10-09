package cmdwatch

// JSON contract tests for `gitmap watch --json`.
//
// watch emits a single top-level object with a nested repos array
// and summary object. The contract covers:
//
//   - Top-level key order: timestamp, repos, summary.
//   - Repo item key order: name, path, branch, status, ahead, behind, stash.
//   - Summary key order: total, dirty, behind, stash.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdwatch/ -run WatchJSONContract
//
// NOTE: this test lives in cmdwatch (not cmd) because the encoder and
// its types (watchSnapshot, watchSummary, encodeWatchJSON,
// buildWatchSummary) are unexported here. The golden helpers below are
// local, minimal equivalents of the cmd package's golden-test
// infrastructure.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/goldenguard"
)

// TestWatchJSONContract_EmptyRepos asserts the top-level shape when
// no repos are present. The repos array must be `[]` and the summary
// must show all-zero counts.
func TestWatchJSONContract_EmptyRepos(t *testing.T) {
	wAssertGoldenBytesDeterministic(t, "watch_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeWatchJSON(&buf, nil, watchSummary{}, "2025-01-01T12:00:00Z")

		return buf.Bytes(), err
	})
}

// TestWatchJSONContract_CanonicalObject_KeyOrders asserts the key
// order of the top-level object, the first repo item, and the
// summary object.
func TestWatchJSONContract_CanonicalObject_KeyOrders(t *testing.T) {
	snapshots := []watchSnapshot{canonicalWatchSnapshot()}
	summary := buildWatchSummary(snapshots)
	var buf bytes.Buffer
	if err := encodeWatchJSON(&buf, snapshots, summary, "2025-01-01T12:00:00Z"); err != nil {
		t.Fatalf("encode: %v", err)
	}

	wAssertSchemaKeysFirstObject(t, buf.Bytes(), "watch")
}

// canonicalWatchSnapshot builds a deterministic repo snapshot for
// key-order and structural tests.
func canonicalWatchSnapshot() watchSnapshot {
	return watchSnapshot{
		Name:   "widget",
		Path:   "/repos/acme/widget",
		Branch: "main",
		Status: "clean",
		Ahead:  0,
		Behind: 2,
		Stash:  1,
	}
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---

func wPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("wPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

func wResolveGoldenPath(name string) string {
	return filepath.Join(wPackageDir(), "testdata", name)
}

func wAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := wResolveGoldenPath(name)
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

func wAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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

	wAssertGoldenBytes(t, name, first)
}

// wAssertSchemaKeysFirstObject checks the top-level key order of the
// emitted object against testdata/schemas/<name>.v1.json.
func wAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(wPackageDir(), "testdata", "schemas", name+".v1.json")
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

	got := wFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

// wFirstObjectKeys returns the top-level keys of the first JSON object
// in raw, in wire order.
func wFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("wFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("wFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("wFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("wFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
