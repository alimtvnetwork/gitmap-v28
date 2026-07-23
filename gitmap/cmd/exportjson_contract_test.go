package cmd

// JSON contract tests for `gitmap export`.
//
// export emits a single top-level object whose seven keys appear in
// contractual order: version, exportedAt, repos, groups, releases,
// history, bookmarks. The five nested arrays are always present (`[]`
// when empty, never `null`).
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run ExportJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
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
	assertGoldenBytesDeterministic(t, "export_empty.json", func() ([]byte, error) {
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
	assertSchemaKeysFirstObject(t, buf.Bytes(), "export")
}
