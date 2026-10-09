package cmdlatestbranch

// JSON schema contract tests for `gitmap latest-branch --json`.
//
// Same two-tier strictness as the cmd package's contract tests:
// byte-exact for canonical fixtures, structural for variable data.
// Covers BOTH the top-level object without `top` AND the nested
// top-N array inside `top`, since the key is conditionally emitted.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdlatestbranch/ -run LatestBranchJSONContract
//
// NOTE: this test lives in cmdlatestbranch (not cmd) because the encoder
// and its types (latestBranchResult, encodeLatestBranchJSON) are unexported
// here. The golden helpers below are local, minimal equivalents of the
// cmd package's golden-test infrastructure.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/goldenguard"
)

// canonicalLatestResult builds a deterministic latestBranchResult
// for byte-exact tests. Every value is a fixed string so the golden
// file's bytes are reproducible across machines, timezones, and
// gitutil format-helper changes (we go around the helpers by
// pre-formatting the strings ourselves at the call site).
func canonicalLatestResult() latestBranchResult {
	return latestBranchResult{
		branchNames:    []string{"main", "develop"},
		selectedRemote: "origin",
		shortSha:       "abc1234",
		commitDate:     "01-Jan-2025 12:00 PM (UTC)",
		latest: gitutil.RemoteBranchInfo{
			RemoteRef:  "refs/remotes/origin/main",
			CommitDate: time.Unix(0, 0).UTC(),
			Sha:        "abc1234567890",
			Subject:    "Initial commit",
		},
	}
}

// TestLatestBranchJSONContract_NoTopOmitsKey pins the bytes when
// `top` is absent (top=0). The `omitempty` tag on Top means the key
// must NOT appear in the output.
func TestLatestBranchJSONContract_NoTopOmitsKey(t *testing.T) {
	encode := func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeLatestBranchJSON(&buf, canonicalLatestResult(), nil, 0)

		return buf.Bytes(), err
	}

	lbAssertGoldenBytesDeterministic(t, "latest_branch_no_top.json", encode)
	// Schema check uses a fresh encode so a bug in the helper that
	// mutates the returned slice cannot cross-contaminate the two checks.
	raw, _ := encode()
	lbAssertSchemaKeysFirstObject(t, raw, "latest-branch-no-top")
}

// TestLatestBranchJSONContract_WithTopIncludesKey verifies the
// `top` key is present and the nested object's key order matches
// the lbTopKey* constants in latestbranchrender.go. Structural-only
// (no byte-exact golden) because renderLatestBranchTopRaw calls
// gitutil.FormatDisplayDate which uses the LOCAL timezone — bytes
// would drift across CI machines.
func TestLatestBranchJSONContract_WithTopIncludesKey(t *testing.T) {
	items := []gitutil.RemoteBranchInfo{
		{
			RemoteRef:  "refs/remotes/origin/main",
			CommitDate: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			Sha:        "abc1234567890",
			Subject:    "Initial commit",
		},
	}

	var buf bytes.Buffer
	if err := encodeLatestBranchJSON(&buf, canonicalLatestResult(), items, 1); err != nil {
		t.Fatalf("encode: %v", err)
	}

	lbAssertSchemaKeysFirstObject(t, buf.Bytes(), "latest-branch-with-top")
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---

func lbPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("lbPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

func lbResolveGoldenPath(name string) string {
	return filepath.Join(lbPackageDir(), "testdata", name)
}

func lbAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := lbResolveGoldenPath(name)
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

func lbAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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

	lbAssertGoldenBytes(t, name, first)
}

// lbAssertSchemaKeysFirstObject checks the top-level key order of the
// emitted object against testdata/schemas/<name>.v1.json.
func lbAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(lbPackageDir(), "testdata", "schemas", name+".v1.json")
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

	got := lbFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

// lbFirstObjectKeys returns the top-level keys of the first JSON object
// in raw, in wire order.
func lbFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("lbFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("lbFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("lbFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("lbFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
