package cmd

// JSON encoder for `gitmap amend audit` file output.
//
// Migrated off json.MarshalIndent onto gitmap/stablejson so key order
// becomes a compile-time decision rather than a reflection accident.
// Nested objects (previousAuthor, newAuthor) and the commits array are
// pre-rendered in compact mode and embedded as json.RawMessage so their
// key order is also stable.
//
// Schema: spec/08-json-schemas/amend-audit.schema.json.

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/stablejson"
)

// amend-audit wire keys. Names + order are the contract; reordering
// or renaming here is a consumer-facing break and the schema contract
// test fails on drift in either direction.
const (
	amendAuditKeyID             = "id"
	amendAuditKeyTimestamp      = "timestamp"
	amendAuditKeyBranch         = "branch"
	amendAuditKeyFromCommit     = "fromCommit"
	amendAuditKeyToCommit       = "toCommit"
	amendAuditKeyTotalCommits   = "totalCommits"
	amendAuditKeyPreviousAuthor = "previousAuthor"
	amendAuditKeyNewAuthor      = "newAuthor"
	amendAuditKeyMode           = "mode"
	amendAuditKeyForcePushed    = "forcePushed"
	amendAuditKeyCommits        = "commits"
)

// encodeAmendAuditJSON writes a single AmendmentRecord as stable
// JSON with 2-space indentation. Split out from CLI dispatch so
// contract tests can capture the bytes into a buffer instead of a
// file.
func encodeAmendAuditJSON(w io.Writer, record model.AmendmentRecord) error {
	prevAuthorRaw, err := renderAmendAuthorRaw(record.PreviousAuthor)
	if err != nil {
		return err
	}

	newAuthorRaw, err := renderAmendAuthorRaw(record.NewAuthor)
	if err != nil {
		return err
	}

	commitsRaw, err := renderCommitEntriesRaw(record.Commits)
	if err != nil {
		return err
	}

	return stablejson.WriteObject(w, []stablejson.Field{
		{Key: amendAuditKeyID, Value: record.ID},
		{Key: amendAuditKeyTimestamp, Value: record.Timestamp},
		{Key: amendAuditKeyBranch, Value: record.Branch},
		{Key: amendAuditKeyFromCommit, Value: record.FromCommit},
		{Key: amendAuditKeyToCommit, Value: record.ToCommit},
		{Key: amendAuditKeyTotalCommits, Value: record.TotalCommits},
		{Key: amendAuditKeyPreviousAuthor, Value: prevAuthorRaw},
		{Key: amendAuditKeyNewAuthor, Value: newAuthorRaw},
		{Key: amendAuditKeyMode, Value: record.Mode},
		{Key: amendAuditKeyForcePushed, Value: record.ForcePushed},
		{Key: amendAuditKeyCommits, Value: commitsRaw},
	})
}

// renderAmendAuthorRaw pre-renders an AmendAuthor object in compact
// mode so it embeds cleanly as a top-level object value.
func renderAmendAuthorRaw(a model.AmendAuthor) (json.RawMessage, error) {
	var buf bytes.Buffer
	if err := stablejson.WriteObjectIndent(&buf, []stablejson.Field{
		{Key: "name", Value: a.Name},
		{Key: "email", Value: a.Email},
	}, ""); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})), nil
}

// renderCommitEntriesRaw pre-renders the commits array in compact
// mode so it embeds cleanly as a top-level object value.
func renderCommitEntriesRaw(entries []model.CommitEntry) (json.RawMessage, error) {
	items := make([][]stablejson.Field, 0, len(entries))
	for _, e := range entries {
		items = append(items, []stablejson.Field{
			{Key: "sha", Value: e.SHA},
			{Key: "message", Value: e.Message},
		})
	}

	var buf bytes.Buffer
	if err := stablejson.WriteArrayIndent(&buf, items, ""); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})), nil
}
