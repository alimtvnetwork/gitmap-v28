package cmd

// JSON encoder for `gitmap probe --json`.
//
// Migrated off json.Encoder onto gitmap/stablejson so key order
// becomes a compile-time decision rather than a reflection accident.
// Schema: spec/08-json-schemas/probe-report.schema.json.

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/stablejson"
)

// probe-report wire keys. Names + order are the contract; reordering
// or renaming here is a consumer-facing break and the schema
// contract test fails on drift in either direction.
const (
	probeKeyRepoID         = "repoId"
	probeKeySlug           = "slug"
	probeKeyAbsolutePath   = "absolutePath"
	probeKeyNextVersionTag = "nextVersionTag"
	probeKeyNextVersionNum = "nextVersionNum"
	probeKeyMethod         = "method"
	probeKeyIsAvailable    = "isAvailable"
	probeKeyError          = "error"
)

// encodeProbeJSON writes entries as a stablejson 2-space-indented
// array. Empty input emits `[]\n` so `jq length` works without a
// special case. Split out from CLI dispatch so contract tests can
// capture the bytes into a buffer instead of stdout.
func encodeProbeJSON(w io.Writer, entries []probeJSONEntry) error {
	return stablejson.WriteArray(w, buildProbeJSONItems(entries))
}

// buildProbeJSONItems is the single source of (field name, field
// order, value) for probe-report. Centralized so a future column
// rename/reorder is one diff and the contract test catches schema
// drift in the same PR.
func buildProbeJSONItems(entries []probeJSONEntry) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(entries))
	for _, e := range entries {
		items = append(items, []stablejson.Field{
			{Key: probeKeyRepoID, Value: e.RepoID},
			{Key: probeKeySlug, Value: e.Slug},
			{Key: probeKeyAbsolutePath, Value: e.AbsolutePath},
			{Key: probeKeyNextVersionTag, Value: e.NextVersionTag},
			{Key: probeKeyNextVersionNum, Value: e.NextVersionNum},
			{Key: probeKeyMethod, Value: e.Method},
			{Key: probeKeyIsAvailable, Value: e.IsAvailable},
			{Key: probeKeyError, Value: e.Error},
		})
	}

	return items
}
