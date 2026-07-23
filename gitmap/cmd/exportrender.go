package cmd

// JSON encoder for `gitmap export`.
//
// Migrated off json.MarshalIndent(model.DatabaseExport) onto
// gitmap/stablejson so the top-level key order is a compile-time
// decision rather than struct-field-tag-defined. Nested arrays
// (`repos`, `groups`, `releases`, `history`, `bookmarks`) are
// pre-rendered with json.MarshalIndent — their per-record key order
// is already deterministic (Go struct field declaration order) and
// migrating every nested record type is out of scope for this
// release; only the top-level shape is contractually pinned here.
//
// Schema: spec/08-json-schemas/export.schema.json.

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/stablejson"
)

// export top-level wire keys. Names + order are the contract.
const (
	exportKeyVersion    = "version"
	exportKeyExportedAt = "exportedAt"
	exportKeyRepos      = "repos"
	exportKeyGroups     = "groups"
	exportKeyReleases   = "releases"
	exportKeyHistory    = "history"
	exportKeyBookmarks  = "bookmarks"
)

// encodeDatabaseExportJSON writes the database export as stable JSON.
func encodeDatabaseExportJSON(w io.Writer, e model.DatabaseExport) error {
	fields := []stablejson.Field{
		{Key: exportKeyVersion, Value: e.Version},
		{Key: exportKeyExportedAt, Value: e.ExportedAt},
	}

	for _, entry := range buildExportArrayFields(e) {
		raw, err := renderExportArrayRaw(entry.value)
		if err != nil {
			return err
		}
		fields = append(fields, stablejson.Field{Key: entry.key, Value: raw})
	}

	return stablejson.WriteObject(w, fields)
}

// exportArrayField pairs a wire key with its slice value so the
// top-level iteration is table-driven.
type exportArrayField struct {
	key   string
	value any
}

// buildExportArrayFields returns the five nested-array fields in
// contractual order.
func buildExportArrayFields(e model.DatabaseExport) []exportArrayField {
	return []exportArrayField{
		{exportKeyRepos, e.Repos},
		{exportKeyGroups, e.Groups},
		{exportKeyReleases, e.Releases},
		{exportKeyHistory, e.History},
		{exportKeyBookmarks, e.Bookmarks},
	}
}

// renderExportArrayRaw pre-renders a nested slice as a JSON array so
// it can be embedded as a stablejson Field value. Empty slices emit
// `[]` (never `null`).
func renderExportArrayRaw(slice any) (json.RawMessage, error) {
	if isEmptyExportSlice(slice) {
		return json.RawMessage("[]"), nil
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("  ", "  ")
	if err := enc.Encode(slice); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimRight(buf.Bytes(), "\n")), nil
}

// isEmptyExportSlice returns true for nil or zero-length slice values
// across the five known DatabaseExport array fields.
func isEmptyExportSlice(slice any) bool {
	switch v := slice.(type) {
	case []model.ScanRecord:
		return len(v) == 0
	case []model.GroupExport:
		return len(v) == 0
	case []model.ReleaseRecord:
		return len(v) == 0
	case []model.CommandHistoryRecord:
		return len(v) == 0
	case []model.BookmarkRecord:
		return len(v) == 0
	}

	return false
}
