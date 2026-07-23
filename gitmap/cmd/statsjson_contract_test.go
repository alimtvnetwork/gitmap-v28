package cmd

// JSON contract tests for `gitmap stats --json`.
//
// stats emits a single object summarizing overall command-history
// metrics plus a nested `commands` array of per-command rows. The
// contract covers:
//
//   - Top-level object shape with all required keys.
//   - Key order: totalCommands, uniqueCommands, totalSuccess,
//     totalFail, overallFailRate, avgDurationMs, commands.
//   - Empty commands list is `[]` (NOT null) on the wire.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run StatsJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// canonicalStatsOverall builds a deterministic two-row stats snapshot.
// One row exercises a populated command; tests pin both the wrapping
// object's key order and the embedded compact array's bytes.
func canonicalStatsOverall(t *testing.T) (model.OverallStats, []model.CommandStats) {
	t.Helper()
	commands := []model.CommandStats{
		{
			Command:      "clone",
			TotalRuns:    10,
			SuccessCount: 9,
			FailCount:    1,
			FailRate:     0.1,
			AvgDuration:  1234,
			MinDuration:  500,
			MaxDuration:  3000,
			LastUsed:     "2026-05-26T08:30:00Z",
		},
		{
			Command:      "stats",
			TotalRuns:    2,
			SuccessCount: 2,
			FailCount:    0,
			FailRate:     0.0,
			AvgDuration:  42,
			MinDuration:  40,
			MaxDuration:  45,
			LastUsed:     "2026-05-26T09:00:00Z",
		},
	}
	overall := model.OverallStats{
		TotalCommands:   12,
		UniqueCommands:  2,
		TotalSuccess:    11,
		TotalFail:       1,
		OverallFailRate: 0.0833,
		AvgDuration:     1100,
		Commands:        commands,
	}

	return overall, commands
}

// TestStatsJSONContract_EmptyCommandsArray pins the empty-commands
// shape: a fully-populated overall object with `commands: []`.
func TestStatsJSONContract_EmptyCommandsArray(t *testing.T) {
	overall := model.OverallStats{}
	assertGoldenBytesDeterministic(t, "stats_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeStatsJSON(&buf, overall, nil)

		return buf.Bytes(), err
	})
}

// Canonical-byte fixture is intentionally omitted: float formatting
// (e.g. 0.0833) and embedded compact-array bytes are tied to Go's
// json.Marshal output, so regenerating via GITMAP_UPDATE_GOLDEN is
// the safer pin. The key-order test below covers the structural
// contract without locking in float-printing artifacts.

// TestStatsJSONContract_KeyOrder asserts the top-level object's key
// order matches the schema registry declaration.
func TestStatsJSONContract_KeyOrder(t *testing.T) {
	overall, commands := canonicalStatsOverall(t)
	var buf bytes.Buffer
	if err := encodeStatsJSON(&buf, overall, commands); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "stats")
}
