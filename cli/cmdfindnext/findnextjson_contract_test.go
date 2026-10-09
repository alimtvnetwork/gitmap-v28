package cmdfindnext

// JSON schema contract tests for `gitmap find-next --json`.
//
// find-next emits an array of model.FindNextRow, each containing an
// embedded model.ScanRecord under the `repo` key. The contract
// covers BOTH layers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - FindNextRow key order: repo, nextVersionTag, nextVersionNum,
//     method, probedAt.
//   - ScanRecord key order (nested under `repo`): id, slug,
//     repoName, httpsUrl, sshUrl, branch, branchSource,
//     relativePath, absolutePath, cloneInstruction, notes, depth.
//
// ScanRecord is large (12 fields) and shared with several other CLI
// outputs, so a rename/reorder there would silently ripple into
// multiple consumers. Pinning it here adds one tripwire that catches
// every such change.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdfindnext/ -run FindNextJSONContract
//
// NOTE: this test lives in cmdfindnext (not cmd) because the encoder
// (encodeFindNextJSON) is unexported here. The golden helpers below
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

// TestFindNextJSONContract_EmptyIsArrayNotNull is the jq-compat
// guarantee: zero rows must encode as `[]\n` even when the input
// slice is nil.
func TestFindNextJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	fnAssertGoldenBytesDeterministic(t, "find_next_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeFindNextJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalFindNextRow builds a deterministic single row whose every
// field is a fixed value, so the golden file's bytes are stable
// across machines and time. Used for the byte-exact + key-order
// tests below.
func canonicalFindNextRow() model.FindNextRow {
	return model.FindNextRow{
		Repo: model.ScanRecord{
			ID:               42,
			Slug:             "acme/widget",
			RepoName:         "widget",
			HTTPSUrl:         "https://github.com/acme/widget.git",
			SSHUrl:           "git@github.com:acme/widget.git",
			Branch:           "main",
			BranchSource:     "remote-head",
			RelativePath:     "acme/widget",
			AbsolutePath:     "/repos/acme/widget",
			CloneInstruction: "git clone https://github.com/acme/widget.git",
			Notes:            "",
		},
		NextVersionTag: "v1.2.3",
		NextVersionNum: 123,
		Method:         "tag-probe",
		ProbedAt:       "2025-01-01T12:00:00Z",
	}
}

// TestFindNextJSONContract_CanonicalRow_KeyOrders asserts the
// FindNextRow-level AND nested ScanRecord-level key orders.
// Structural-only (no byte-exact golden for the populated row) so
// the test stays robust against future numeric formatting changes
// in encoding/json or value-shape tweaks.
func TestFindNextJSONContract_CanonicalRow_KeyOrders(t *testing.T) {
	rows := []model.FindNextRow{canonicalFindNextRow()}
	var buf bytes.Buffer
	if err := encodeFindNextJSON(&buf, rows); err != nil {
		t.Fatalf("encode: %v", err)
	}

	fnAssertSchemaKeysFirstObject(t, buf.Bytes(), "find-next")
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---

func fnPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("fnPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

func fnResolveGoldenPath(name string) string {
	return filepath.Join(fnPackageDir(), "testdata", name)
}

func fnAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := fnResolveGoldenPath(name)
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

func fnAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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

	fnAssertGoldenBytes(t, name, first)
}

// fnAssertSchemaKeysFirstObject checks the key order of the first array
// element against testdata/schemas/<name>.v1.json.
func fnAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(fnPackageDir(), "testdata", "schemas", name+".v1.json")
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

	got := fnFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

// fnFirstObjectKeys returns the keys of the first JSON object in raw
// (skipping a top-level array wrapper), in wire order.
func fnFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("fnFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("fnFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("fnFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("fnFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
