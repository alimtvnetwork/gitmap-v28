package cmdprobe

// JSON contract tests for `gitmap probe --json`.
//
// probe emits an array of probeJSONEntry. The contract covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: repoId, slug, absolutePath, nextVersionTag,
//     nextVersionNum, method, isAvailable, error.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdprobe/ -run ProbeJSONContract
//
// NOTE: this test lives in cmdprobe (not cmd) because the encoder and
// its type (probeJSONEntry, encodeProbeJSON) are unexported here. The
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
)

// TestProbeJSONContract_EmptyIsArrayNotNull is the jq-compat
// guarantee: zero rows must encode as `[]\n` even when the input
// slice is nil.
func TestProbeJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	prAssertGoldenBytesDeterministic(t, "probe_report_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeProbeJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalProbeEntry builds a deterministic single entry whose
// every field is a fixed value, so the golden file's bytes are
// stable across machines and time.
func canonicalProbeEntry() probeJSONEntry {
	return probeJSONEntry{
		RepoID:         42,
		Slug:           "acme/widget",
		AbsolutePath:   "/repos/acme/widget",
		NextVersionTag: "v1.2.3",
		NextVersionNum: 123,
		Method:         "tag-probe",
		IsAvailable:    true,
		Error:          "",
	}
}

// TestProbeJSONContract_CanonicalRow_KeyOrders asserts the key order
// of the emitted object matches the schema declaration.
// Structural-only (no byte-exact golden for the populated row) so
// the test stays robust against future numeric formatting changes
// in encoding/json or value-shape tweaks.
func TestProbeJSONContract_CanonicalRow_KeyOrders(t *testing.T) {
	entries := []probeJSONEntry{canonicalProbeEntry()}
	var buf bytes.Buffer
	if err := encodeProbeJSON(&buf, entries); err != nil {
		t.Fatalf("encode: %v", err)
	}

	prAssertSchemaKeysFirstObject(t, buf.Bytes(), "probe-report")
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---

func prPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("prPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

func prResolveGoldenPath(name string) string {
	return filepath.Join(prPackageDir(), "testdata", name)
}

func prAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := prResolveGoldenPath(name)
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

func prAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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

	prAssertGoldenBytes(t, name, first)
}

// prAssertSchemaKeysFirstObject checks the key order of the first array
// element against testdata/schemas/<name>.v1.json.
func prAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(prPackageDir(), "testdata", "schemas", name+".v1.json")
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

	got := prFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

// prFirstObjectKeys returns the keys of the first JSON object in raw
// (skipping a top-level array wrapper), in wire order.
func prFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("prFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("prFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("prFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("prFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
