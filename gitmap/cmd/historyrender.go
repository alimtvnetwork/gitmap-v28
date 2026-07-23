package cmd

// JSON encoder for `gitmap history --json`.
//
// Migrated off json.MarshalIndent onto gitmap/stablejson so key order
// becomes a compile-time decision rather than a reflection accident.
// Schema: spec/08-json-schemas/history.schema.json.

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/stablejson"
)

// history wire keys. Names + order are the contract; reordering or
// renaming here is a consumer-facing break and the schema contract
// test fails on drift in either direction.
const (
	historyKeyID         = "id"
	historyKeyCommand    = "command"
	historyKeyAlias      = "alias"
	historyKeyArgs       = "args"
	historyKeyFlags      = "flags"
	historyKeyStartedAt  = "startedAt"
	historyKeyFinishedAt = "finishedAt"
	historyKeyDurationMs = "durationMs"
	historyKeyExitCode   = "exitCode"
	historyKeySummary    = "summary"
	historyKeyRepoCount  = "repoCount"
	historyKeyCreatedAt  = "createdAt"
)

// encodeHistoryJSON writes records as a stablejson 2-space-indented
// array. Empty input emits `[]\n` so `jq length` works without a
// special case. Split out from CLI dispatch so contract tests can
// capture the bytes into a buffer instead of stdout.
func encodeHistoryJSON(w io.Writer, records []model.CommandHistoryRecord) error {
	return stablejson.WriteArray(w, buildHistoryJSONItems(records))
}

// buildHistoryJSONItems is the single source of (field name, field
// order, value) for history. Centralized so a future column
// rename/reorder is one diff and the contract test catches schema
// drift in the same PR.
func buildHistoryJSONItems(records []model.CommandHistoryRecord) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(records))
	for _, r := range records {
		items = append(items, []stablejson.Field{
			{Key: historyKeyID, Value: r.ID},
			{Key: historyKeyCommand, Value: r.Command},
			{Key: historyKeyAlias, Value: r.Alias},
			{Key: historyKeyArgs, Value: r.Args},
			{Key: historyKeyFlags, Value: r.Flags},
			{Key: historyKeyStartedAt, Value: r.StartedAt},
			{Key: historyKeyFinishedAt, Value: r.FinishedAt},
			{Key: historyKeyDurationMs, Value: r.DurationMs},
			{Key: historyKeyExitCode, Value: r.ExitCode},
			{Key: historyKeySummary, Value: r.Summary},
			{Key: historyKeyRepoCount, Value: r.RepoCount},
			{Key: historyKeyCreatedAt, Value: r.CreatedAt},
		})
	}

	return items
}
