package cmd

// JSON encoder for `gitmap version-history --json`.
//
// Migrated off json.MarshalIndent([]model.RepoVersionHistoryRecord) onto
// gitmap/stablejson so every key has compile-time-stable order rather than
// reflection-defined order.
//
// The version-history output is a top-level array of objects. Optional
// fields (`flattenedPath`, `createdAt`) are omitted from the wire when
// empty to preserve the legacy omitempty shape.
//
// Schema: spec/08-json-schemas/version-history.schema.json.

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/stablejson"
)

// versionHistory wire keys. Names + order are the contract.
const (
	vhKeyFromVersionTag = "fromVersionTag"
	vhKeyFromVersionNum = "fromVersionNum"
	vhKeyToVersionTag   = "toVersionTag"
	vhKeyToVersionNum   = "toVersionNum"
	vhKeyFlattenedPath  = "flattenedPath"
	vhKeyCreatedAt      = "createdAt"
	vhKeyID             = "id"
	vhKeyRepoID         = "repoId"
)

// encodeVersionHistoryJSON writes version-history records as stable JSON.
func encodeVersionHistoryJSON(w io.Writer, records []model.RepoVersionHistoryRecord) error {
	return stablejson.WriteArray(w, buildVersionHistoryItems(records))
}

// buildVersionHistoryItems is the single source of (field name, field
// order, value) for version-history rows.
func buildVersionHistoryItems(records []model.RepoVersionHistoryRecord) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(records))
	for _, r := range records {
		row := []stablejson.Field{
			{Key: vhKeyFromVersionTag, Value: r.FromVersionTag},
			{Key: vhKeyFromVersionNum, Value: r.FromVersionNum},
			{Key: vhKeyToVersionTag, Value: r.ToVersionTag},
			{Key: vhKeyToVersionNum, Value: r.ToVersionNum},
		}
		if r.FlattenedPath != "" {
			row = append(row, stablejson.Field{Key: vhKeyFlattenedPath, Value: r.FlattenedPath})
		}
		if r.CreatedAt != "" {
			row = append(row, stablejson.Field{Key: vhKeyCreatedAt, Value: r.CreatedAt})
		}
		row = append(row,
			stablejson.Field{Key: vhKeyID, Value: r.ID},
			stablejson.Field{Key: vhKeyRepoID, Value: r.RepoID},
		)
		items = append(items, row)
	}

	return items
}
