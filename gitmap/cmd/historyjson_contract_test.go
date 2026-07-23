package cmd

// JSON contract tests for `gitmap history --json`.
//
// history emits an array of model.CommandHistoryRecord. The contract
// covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: id, command, alias, args, flags, startedAt,
//     finishedAt, durationMs, exitCode, summary, repoCount, createdAt.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run HistoryJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// TestHistoryJSONContract_EmptyIsArrayNotNull is the jq-compat
// guarantee: zero rows must encode as `[]\n` even when the input
// slice is nil.
func TestHistoryJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "history_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeHistoryJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalHistoryRecord builds a deterministic single record whose
// every field is a fixed value, so the golden file's bytes are stable
// across machines and time. Used for the byte-exact + key-order
// tests below.
func canonicalHistoryRecord() model.CommandHistoryRecord {
	return model.CommandHistoryRecord{
		ID:         42,
		Command:    "scan",
		Alias:      "s",
		Args:       "all",
		Flags:      "--depth=3",
		StartedAt:  "2025-01-01T12:00:00Z",
		FinishedAt: "2025-01-01T12:00:05Z",
		DurationMs: 5123,
		ExitCode:   0,
		Summary:    "Scanned 7 repos, 2 dirty",
		RepoCount:  7,
		CreatedAt:  "2025-01-01T12:00:00Z",
	}
}

// TestHistoryJSONContract_CanonicalRow_KeyOrders asserts the
// key order of the emitted object matches the schema declaration.
// Structural-only (no byte-exact golden for the populated row) so
// the test stays robust against future numeric formatting changes
// in encoding/json or value-shape tweaks.
func TestHistoryJSONContract_CanonicalRow_KeyOrders(t *testing.T) {
	records := []model.CommandHistoryRecord{canonicalHistoryRecord()}
	var buf bytes.Buffer
	if err := encodeHistoryJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "history")
}
