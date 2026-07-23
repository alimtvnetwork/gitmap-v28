package cmd

// JSON encoder for `gitmap temp-releaselist --json`.
//
// Migrated off json.MarshalIndent([]model.TempRelease) onto
// gitmap/stablejson so every key has compile-time-stable order rather than
// reflection-defined order.
//
// Schema: spec/08-json-schemas/temp-release-list.schema.json.

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/stablejson"
)

// tempReleaseList wire keys. Names + order are the contract.
const (
	trlKeyID             = "id"
	trlKeyBranch         = "branch"
	trlKeyVersionPrefix  = "versionPrefix"
	trlKeySequenceNumber = "sequenceNumber"
	trlKeyCommitSha      = "commit"
	trlKeyCommitMessage  = "commitMessage"
	trlKeyCreatedAt      = "createdAt"
)

// encodeTempReleaseListJSON writes temp-release records as stable JSON.
func encodeTempReleaseListJSON(w io.Writer, releases []model.TempRelease) error {
	return stablejson.WriteArray(w, buildTempReleaseListItems(releases))
}

// buildTempReleaseListItems is the single source of (field name, field
// order, value) for temp-release rows.
func buildTempReleaseListItems(releases []model.TempRelease) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(releases))
	for _, r := range releases {
		items = append(items, []stablejson.Field{
			{Key: trlKeyID, Value: r.ID},
			{Key: trlKeyBranch, Value: r.Branch},
			{Key: trlKeyVersionPrefix, Value: r.VersionPrefix},
			{Key: trlKeySequenceNumber, Value: r.SequenceNumber},
			{Key: trlKeyCommitSha, Value: r.CommitSha},
			{Key: trlKeyCommitMessage, Value: r.CommitMessage},
			{Key: trlKeyCreatedAt, Value: r.CreatedAt},
		})
	}

	return items
}
