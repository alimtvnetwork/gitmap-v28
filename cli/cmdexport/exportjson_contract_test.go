package cmdexport

// JSON contract tests for `gitmap export`.
//
// export emits a single top-level object whose seven keys appear in
// contractual order: version, exportedAt, repos, groups, releases,
// history, bookmarks. The five nested arrays are always present (`[]`
// when empty, never `null`).
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdexport/ -run ExportJSONContract
//
// NOTE: this test lives in cmdexport (not cmd) because the encoder
// (encodeDatabaseExportJSON) is unexported here. The golden helpers
// below are local, minimal equivalents of the cmd package's
// golden-test infrastructure.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/goldenguard"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// canonicalEmptyExport builds a deterministic empty-database export.
func canonicalEmptyExport() model.DatabaseExport {
	return model.DatabaseExport{
		Version:    "1",
		ExportedAt: "2026-05-26T12:00:00Z",
	}
}

// TestExportJSONContract_EmptyArraysNotNull pins the empty-database
// shape so downstream `jq '.repos | length'` consumers never see `null`.
func TestExportJSONContract_EmptyArraysNotNull(t *testing.T) {
	export := canonicalEmptyExport()
	exAssertGoldenBytesDeterministic(t, "export_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeDatabaseExportJSON(&buf, export)

		return buf.Bytes(), err
	})
}

// TestExportJSONContract_TopLevelKeyOrder asserts the top-level key
// order matches the schema registry declaration.
func TestExportJSONContract_TopLevelKeyOrder(t *testing.T) {
	export := canonicalEmptyExport()
	var buf bytes.Buffer
	if err := encodeDatabaseExportJSON(&buf, export); err != nil {
		t.Fatalf("encode: %v", err)
	}

	exAssertSchemaKeysFirstObject(t, buf.Bytes(), "export")
}

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---
// exPackageDir is defined in export_jsonschema_contract_test.go.

func exResolveGoldenPath(name string) string {
	return filepath.Join(exPackageDir(), "testdata", name)
}

func exAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := exResolveGoldenPath(name)
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

func exAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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

	exAssertGoldenBytes(t, name, first)
}

// exAssertSchemaKeysFirstObject checks the top-level key order of the
// emitted object against testdata/schemas/<name>.v1.json.
func exAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(exPackageDir(), "testdata", "schemas", name+".v1.json")
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

	got := exFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}
