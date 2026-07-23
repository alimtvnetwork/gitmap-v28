package cmd

// JSON encoder for `gitmap watch --json`.
//
// Migrated off json.MarshalIndent(map[string]any) onto gitmap/stablejson
// so every key (top-level and nested) has compile-time-stable order.
//
// The watch output is a single top-level object containing an array
// (`repos`) and a nested object (`summary`). Both nested shapes are
// pre-rendered into json.RawMessage using COMPACT mode (zero indent)
// so they can be embedded without indentation-context mismatches.
// The top-level object itself is pretty-printed with 2-space indent.
// This yields valid, stable JSON where key order is the headline
// guarantee.
//
// Schema: spec/08-json-schemas/watch.schema.json.

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/stablejson"
)

// watch top-level wire keys. Names + order are the contract.
const (
	watchKeyTimestamp = "timestamp"
	watchKeyRepos     = "repos"
	watchKeySummary   = "summary"
)

// watch summary wire keys.
const (
	watchSummaryKeyTotal  = "total"
	watchSummaryKeyDirty  = "dirty"
	watchSummaryKeyBehind = "behind"
	watchSummaryKeyStash  = "stash"
)

// watch repo snapshot wire keys.
const (
	watchRepoKeyName   = "name"
	watchRepoKeyPath   = "path"
	watchRepoKeyBranch = "branch"
	watchRepoKeyStatus = "status"
	watchRepoKeyAhead  = "ahead"
	watchRepoKeyBehind = "behind"
	watchRepoKeyStash  = "stash"
)

// encodeWatchJSON writes a single watch snapshot as stable JSON.
// `timestamp` is injected so tests can pass a fixed value.
func encodeWatchJSON(w io.Writer, snapshots []watchSnapshot, summary watchSummary, timestamp string) error {
	reposRaw, err := renderWatchReposRaw(snapshots)
	if err != nil {
		return err
	}

	summaryRaw, err := renderWatchSummaryRaw(summary)
	if err != nil {
		return err
	}

	return stablejson.WriteObject(w, []stablejson.Field{
		{Key: watchKeyTimestamp, Value: timestamp},
		{Key: watchKeyRepos, Value: reposRaw},
		{Key: watchKeySummary, Value: summaryRaw},
	})
}

// renderWatchReposRaw pre-renders the repos array in compact mode
// so it embeds cleanly as a top-level object value.
func renderWatchReposRaw(snapshots []watchSnapshot) (json.RawMessage, error) {
	var buf bytes.Buffer
	if err := stablejson.WriteArrayIndent(&buf, buildWatchSnapshotItems(snapshots), ""); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})), nil
}

// renderWatchSummaryRaw pre-renders the summary object in compact
// mode so it embeds cleanly as a top-level object value.
func renderWatchSummaryRaw(summary watchSummary) (json.RawMessage, error) {
	var buf bytes.Buffer
	if err := stablejson.WriteObjectIndent(&buf, buildWatchSummaryFields(summary), ""); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})), nil
}

// buildWatchSnapshotItems is the single source of (field name, field
// order, value) for watch repo snapshots.
func buildWatchSnapshotItems(snapshots []watchSnapshot) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(snapshots))
	for _, s := range snapshots {
		items = append(items, []stablejson.Field{
			{Key: watchRepoKeyName, Value: s.Name},
			{Key: watchRepoKeyPath, Value: s.Path},
			{Key: watchRepoKeyBranch, Value: s.Branch},
			{Key: watchRepoKeyStatus, Value: s.Status},
			{Key: watchRepoKeyAhead, Value: s.Ahead},
			{Key: watchRepoKeyBehind, Value: s.Behind},
			{Key: watchRepoKeyStash, Value: s.Stash},
		})
	}

	return items
}

// buildWatchSummaryFields is the single source of (field name, field
// order, value) for the watch summary.
func buildWatchSummaryFields(summary watchSummary) []stablejson.Field {
	return []stablejson.Field{
		{Key: watchSummaryKeyTotal, Value: summary.Total},
		{Key: watchSummaryKeyDirty, Value: summary.Dirty},
		{Key: watchSummaryKeyBehind, Value: summary.Behind},
		{Key: watchSummaryKeyStash, Value: summary.Stash},
	}
}
