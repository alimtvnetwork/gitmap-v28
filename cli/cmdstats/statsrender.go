package cmdstats

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
// Schema: 02-spec/08-json-schemas/stats.schema.json.

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstats/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonx"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// buildStatsOverallFields packages top-level stats into stablejson fields.
func buildStatsOverallFields(overall model.OverallStats, commandsRaw json.RawMessage) []jsonx.Field {
	return []jsonx.Field{
		{Key: constants.WireKeyTotalCommands, Value: overall.TotalCommands},
		{Key: constants.WireKeyUniqueCommands, Value: overall.UniqueCommands},
		{Key: constants.WireKeyTotalSuccess, Value: overall.TotalSuccess},
		{Key: constants.WireKeyTotalFail, Value: overall.TotalFail},
		{Key: constants.WireKeyOverallFailRate, Value: overall.OverallFailRate},
		{Key: constants.WireKeyAvgDurationMs, Value: overall.AvgDuration},
		{Key: constants.WireKeyCommands, Value: commandsRaw},
	}
}

// encodeStatsJSON writes overall + per-command stats as stable JSON.
func encodeStatsJSON(w io.Writer, overall model.OverallStats, commands []model.CommandStats) error {
	commandsRaw, err := renderStatsCommandsRaw(commands)
	if err != nil {
		return err
	}

	return jsonx.WriteObject(w, buildStatsOverallFields(overall, commandsRaw))
}

// renderStatsCommandsRaw pre-renders the per-command array in compact
// mode so it embeds cleanly as a top-level object value.
func renderStatsCommandsRaw(commands []model.CommandStats) (json.RawMessage, error) {
	var buf bytes.Buffer
	if err := jsonx.WriteArrayIndent(&buf, buildStatsCommandsItems(commands), ""); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})), nil
}

// buildSingleCommandFields maps a CommandStats record to its wire fields.
func buildSingleCommandFields(s model.CommandStats) []jsonx.Field {
	return []jsonx.Field{
		{Key: constants.WireKeyCmdCommand, Value: s.Command},
		{Key: constants.WireKeyCmdTotalRuns, Value: s.TotalRuns},
		{Key: constants.WireKeyCmdSuccessCount, Value: s.SuccessCount},
		{Key: constants.WireKeyCmdFailCount, Value: s.FailCount},
		{Key: constants.WireKeyCmdFailRate, Value: s.FailRate},
		{Key: constants.WireKeyCmdAvgDurationMs, Value: s.AvgDuration},
		{Key: constants.WireKeyCmdMinDurationMs, Value: s.MinDuration},
		{Key: constants.WireKeyCmdMaxDurationMs, Value: s.MaxDuration},
		{Key: constants.WireKeyCmdLastUsed, Value: s.LastUsed},
	}
}

// buildStatsCommandsItems is the single source of (field name, field
// order, value) for per-command stats rows.
func buildStatsCommandsItems(commands []model.CommandStats) [][]jsonx.Field {
	items := make([][]jsonx.Field, 0, len(commands))
	for _, s := range commands {
		items = append(items, buildSingleCommandFields(s))
	}

	return items
}
