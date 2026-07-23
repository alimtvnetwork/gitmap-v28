package cmd

// JSON encoder for `gitmap stats --json`.
//
// Migrated off json.MarshalIndent(model.OverallStats) onto
// gitmap/stablejson so every key (top-level + nested per-command rows)
// has compile-time-stable order rather than reflection-defined order.
//
// The stats output is a single top-level object whose `commands` field
// is an array of per-command stat objects. The nested array is
// pre-rendered into json.RawMessage using COMPACT mode so it can be
// embedded without indentation-context mismatches; the top-level object
// itself is pretty-printed with 2-space indent. This yields valid,
// stable JSON where key order is the headline guarantee.
//
// Schema: spec/08-json-schemas/stats.schema.json.

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/stablejson"
)

// stats top-level wire keys. Names + order are the contract.
const (
	statsKeyTotalCommands   = "totalCommands"
	statsKeyUniqueCommands  = "uniqueCommands"
	statsKeyTotalSuccess    = "totalSuccess"
	statsKeyTotalFail       = "totalFail"
	statsKeyOverallFailRate = "overallFailRate"
	statsKeyAvgDurationMs   = "avgDurationMs"
	statsKeyCommands        = "commands"
)

// per-command stats wire keys.
const (
	statsCmdKeyCommand       = "command"
	statsCmdKeyTotalRuns     = "totalRuns"
	statsCmdKeySuccessCount  = "successCount"
	statsCmdKeyFailCount     = "failCount"
	statsCmdKeyFailRate      = "failRate"
	statsCmdKeyAvgDurationMs = "avgDurationMs"
	statsCmdKeyMinDurationMs = "minDurationMs"
	statsCmdKeyMaxDurationMs = "maxDurationMs"
	statsCmdKeyLastUsed      = "lastUsed"
)

// encodeStatsJSON writes overall + per-command stats as stable JSON.
func encodeStatsJSON(w io.Writer, overall model.OverallStats, commands []model.CommandStats) error {
	commandsRaw, err := renderStatsCommandsRaw(commands)
	if err != nil {
		return err
	}

	return stablejson.WriteObject(w, []stablejson.Field{
		{Key: statsKeyTotalCommands, Value: overall.TotalCommands},
		{Key: statsKeyUniqueCommands, Value: overall.UniqueCommands},
		{Key: statsKeyTotalSuccess, Value: overall.TotalSuccess},
		{Key: statsKeyTotalFail, Value: overall.TotalFail},
		{Key: statsKeyOverallFailRate, Value: overall.OverallFailRate},
		{Key: statsKeyAvgDurationMs, Value: overall.AvgDuration},
		{Key: statsKeyCommands, Value: commandsRaw},
	})
}

// renderStatsCommandsRaw pre-renders the per-command array in compact
// mode so it embeds cleanly as a top-level object value.
func renderStatsCommandsRaw(commands []model.CommandStats) (json.RawMessage, error) {
	var buf bytes.Buffer
	if err := stablejson.WriteArrayIndent(&buf, buildStatsCommandsItems(commands), ""); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})), nil
}

// buildStatsCommandsItems is the single source of (field name, field
// order, value) for per-command stats rows.
func buildStatsCommandsItems(commands []model.CommandStats) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(commands))
	for _, s := range commands {
		items = append(items, []stablejson.Field{
			{Key: statsCmdKeyCommand, Value: s.Command},
			{Key: statsCmdKeyTotalRuns, Value: s.TotalRuns},
			{Key: statsCmdKeySuccessCount, Value: s.SuccessCount},
			{Key: statsCmdKeyFailCount, Value: s.FailCount},
			{Key: statsCmdKeyFailRate, Value: s.FailRate},
			{Key: statsCmdKeyAvgDurationMs, Value: s.AvgDuration},
			{Key: statsCmdKeyMinDurationMs, Value: s.MinDuration},
			{Key: statsCmdKeyMaxDurationMs, Value: s.MaxDuration},
			{Key: statsCmdKeyLastUsed, Value: s.LastUsed},
		})
	}

	return items
}
