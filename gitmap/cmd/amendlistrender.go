package cmd

// JSON encoder for `gitmap amend list --json`.
//
// Migrated off json.MarshalIndent onto gitmap/stablejson so key order
// becomes a compile-time decision rather than a reflection accident.
//
// AmendmentRow has no json struct tags, so the legacy encoder emits
// PascalCase keys matching the Go field names. The wire keys below
// preserve that shape for backward compatibility.
//
// Schema: spec/08-json-schemas/amend-list.schema.json.

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/stablejson"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// amend-list wire keys. Names + order are the contract.
const (
	amendListKeyID            = "ID"
	amendListKeyBranch        = "Branch"
	amendListKeyFromCommit    = "FromCommit"
	amendListKeyToCommit      = "ToCommit"
	amendListKeyTotalCommits  = "TotalCommits"
	amendListKeyPreviousName  = "PreviousName"
	amendListKeyPreviousEmail = "PreviousEmail"
	amendListKeyNewName       = "NewName"
	amendListKeyNewEmail      = "NewEmail"
	amendListKeyMode          = "Mode"
	amendListKeyForcePushed   = "ForcePushed"
	amendListKeyCreatedAt     = "CreatedAt"
)

// encodeAmendListJSON writes rows as a stablejson 2-space-indented
// array. Empty input emits `[]\n`.
func encodeAmendListJSON(w io.Writer, rows []store.AmendmentRow) error {
	return stablejson.WriteArray(w, buildAmendListJSONItems(rows))
}

// buildAmendListJSONItems is the single source of (field name, field
// order, value) for amend-list.
func buildAmendListJSONItems(rows []store.AmendmentRow) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(rows))
	for _, r := range rows {
		items = append(items, []stablejson.Field{
			{Key: amendListKeyID, Value: r.ID},
			{Key: amendListKeyBranch, Value: r.Branch},
			{Key: amendListKeyFromCommit, Value: r.FromCommit},
			{Key: amendListKeyToCommit, Value: r.ToCommit},
			{Key: amendListKeyTotalCommits, Value: r.TotalCommits},
			{Key: amendListKeyPreviousName, Value: r.PreviousName},
			{Key: amendListKeyPreviousEmail, Value: r.PreviousEmail},
			{Key: amendListKeyNewName, Value: r.NewName},
			{Key: amendListKeyNewEmail, Value: r.NewEmail},
			{Key: amendListKeyMode, Value: r.Mode},
			{Key: amendListKeyForcePushed, Value: r.ForcePushed},
			{Key: amendListKeyCreatedAt, Value: r.CreatedAt},
		})
	}

	return items
}
