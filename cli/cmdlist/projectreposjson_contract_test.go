package cmdlist

// JSON contract tests for `gitmap <type>-repos --json`.
//
// project-repos emits an array of model.DetectedProject. The contract
// covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: id, repoId, repoName, projectTypeId, projectType,
//     projectName, absolutePath, repoPath, relativePath,
//     primaryIndicator, detectedAt.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdlist/ -run ProjectReposJSONContract
//
// NOTE: this test lives in cmdlist (not cmd) because the encoder
// (encodeProjectReposJSON) is unexported here. The golden helpers
// below are local, minimal equivalents of the cmd package's
// golden-test infrastructure.

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

// TestProjectReposJSONContract_EmptyIsArrayNotNull is the jq-compat guarantee.
func TestProjectReposJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	pjAssertGoldenBytesDeterministic(t, "project_repos_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeProjectReposJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalDetectedProject builds a deterministic single row.
func canonicalDetectedProject() model.DetectedProject {
	return model.DetectedProject{
		ID:               42,
		RepoID:           7,
		RepoName:         "gitmap-v28",
		ProjectTypeID:    1,
		ProjectType:      "go",
		ProjectName:      "gitmap",
		AbsolutePath:     "/home/user/code/gitmap-v28/gitmap",
		RepoPath:         "/home/user/code/gitmap-v28",
		RelativePath:     "gitmap",
		PrimaryIndicator: "go.mod",
		DetectedAt:       "2025-01-01T12:00:00Z",
	}
}

// TestProjectReposJSONContract_CanonicalRow_KeyOrder asserts the key
// order of the emitted object matches the schema declaration.
func TestProjectReposJSONContract_CanonicalRow_KeyOrder(t *testing.T) {
	projects := []model.DetectedProject{canonicalDetectedProject()}
	var buf bytes.Buffer
	if err := encodeProjectReposJSON(&buf, projects); err != nil {
		t.Fatalf("encode: %v", err)
	}

	pjAssertSchemaKeysFirstObject(t, buf.Bytes(), "project-repos")
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---

func pjPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("pjPackageDir: runtime.Caller failed")
	}

	return filepath.Dir(file)
}

func pjResolveGoldenPath(name string) string {
	return filepath.Join(pjPackageDir(), "testdata", name)
}

func pjAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := pjResolveGoldenPath(name)
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

func pjAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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

	pjAssertGoldenBytes(t, name, first)
}

// pjAssertSchemaKeysFirstObject checks the key order of the first array
// element against testdata/schemas/<name>.v1.json.
func pjAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(pjPackageDir(), "testdata", "schemas", name+".v1.json")
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

	got := pjFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

// pjFirstObjectKeys returns the keys of the first JSON object in raw
// (skipping a top-level array wrapper), in wire order.
func pjFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("pjFirstObjectKeys: no object found: %v", err)
		}

		if delim, ok := tok.(json.Delim); ok && delim == '{' {
			break
		}
	}

	var keys []string

	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("pjFirstObjectKeys: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("pjFirstObjectKeys: expected string key, got %T", tok)
		}

		keys = append(keys, key)

		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("pjFirstObjectKeys: skip value: %v", err)
		}
	}

	return keys
}
