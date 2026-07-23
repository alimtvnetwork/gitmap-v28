package cmd

// JSON encoder for `gitmap list-versions --json`.
//
// Migrated off json.MarshalIndent(lvJSONEntry, ...) onto
// gitmap/stablejson so key order becomes a compile-time decision
// rather than a reflection accident. Optional fields (source,
// changelog) are conditionally appended so they remain effectively
// omitempty — absent rather than null/empty in the wire output.
//
// Schema: spec/08-json-schemas/list-versions.schema.json.

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/stablejson"
)

// list-versions wire keys. Names + order are the contract.
const (
	listVersionsKeyVersion   = "version"
	listVersionsKeySource    = "source"
	listVersionsKeyChangelog = "changelog"
)

// encodeListVersionsJSON writes entries as a stablejson 2-space-indented
// array. Empty input emits `[]\n`.
func encodeListVersionsJSON(w io.Writer, entries []versionEntry) error {
	return stablejson.WriteArray(w, buildListVersionsJSONItems(entries))
}

// buildListVersionsJSONItems is the single source of (field name,
// field order, value) for list-versions. `source` and `changelog`
// are emitted only when non-empty, preserving the legacy omitempty
// shape that downstream consumers depend on.
func buildListVersionsJSONItems(entries []versionEntry) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(entries))
	for _, e := range entries {
		fields := []stablejson.Field{
			{Key: listVersionsKeyVersion, Value: e.Version.String()},
		}
		if e.Source != "" {
			fields = append(fields, stablejson.Field{Key: listVersionsKeySource, Value: e.Source})
		}
		if len(e.Notes) > 0 {
			fields = append(fields, stablejson.Field{Key: listVersionsKeyChangelog, Value: e.Notes})
		}
		items = append(items, fields)
	}

	return items
}
