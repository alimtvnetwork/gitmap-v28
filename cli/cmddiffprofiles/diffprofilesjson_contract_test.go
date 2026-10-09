package cmddiffprofiles

// JSON contract tests for `gitmap diff-profiles --json`.
//
// The contract covers key order of the emitted single object:
//   profileA, profileB, onlyInA, onlyInB, different, same.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmddiffprofiles/ -run DiffProfilesJSONContract
//
// NOTE: this test lives in cmddiffprofiles (not cmd) because the encoder
// and its types (dpResult, dpDiff, encodeDiffProfilesJSON) are unexported
// here. The golden helpers below are local, minimal equivalents of the
// cmd package's golden-test infrastructure.

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

// canonicalDPResult builds a deterministic diff-profiles result.
func canonicalDPResult() dpResult {
	return dpResult{
		onlyInA: []model.ScanRecord{
			{RepoName: "repo-a", AbsolutePath: "/home/user/a"},
		},
		onlyInB: []model.ScanRecord{
			{RepoName: "repo-b", AbsolutePath: "/home/user/b"},
		},
		different: []dpDiff{
			{
				Name:  "repo-diff",
				PathA: "/home/user/diff-a",
				PathB: "/home/user/diff-b",
				ModeA: "https",
				ModeB: "ssh",
			},
		},
		same: []model.ScanRecord{
			{RepoName: "repo-same", AbsolutePath: "/home/user/same"},
		},
	}
}

// TestDiffProfilesJSONContract_CanonicalRecord_KeyOrder asserts
// the key order of the emitted object matches the schema declaration.
func TestDiffProfilesJSONContract_CanonicalRecord_KeyOrder(t *testing.T) {
	encode := func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeDiffProfilesJSON(&buf, "alpha", "beta", canonicalDPResult())

		return buf.Bytes(), err
	}

	dpAssertGoldenBytesDeterministic(t, "diff_profiles_canonical.json", encode)
	raw, _ := encode()
	dpAssertSchemaKeysFirstObject(t, raw, "diff-profiles")
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---

func dpPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("dpPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

func dpResolveGoldenPath(name string) string {
	return filepath.Join(dpPackageDir(), "testdata", name)
}

func dpAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := dpResolveGoldenPath(name)
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

func dpAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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

	dpAssertGoldenBytes(t, name, first)
}

// dpAssertSchemaKeysFirstObject checks the top-level key order of the
// emitted object against testdata/schemas/<name>.v1.json.
func dpAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(dpPackageDir(), "testdata", "schemas", name+".v1.json")
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

	got := dpFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

// dpFirstObjectKeys returns the top-level keys of the first JSON object
// in raw, in wire order.
func dpFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("dpFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("dpFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("dpFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("dpFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
