// Package cmd — visibilityundojson.go: step-37 `--json` summary
// renderer for `vu` / `vr`. Pure, stdlib-only, no side effects so
// the contract is table-testable independent of the apply loop.
//
// Wire-format mirrors the v5.43.0+ JSON contract used by `--json`
// on every other gitmap CLI: stable key order, lowerCamel field names,
// integer counters, ISO-8601 timestamps from the caller.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §undo-redo.
package cmd

import (
	"encoding/json"
)

// undoJSONSummary is the wire shape emitted under `vu --json` / `vr --json`.
// All counters are populated even when zero (downstream JSON parsers
// hate missing keys). `Command` is "visibility-undo" or "visibility-redo".
type undoJSONSummary struct {
	Command   string `json:"command"`
	RunID     int64  `json:"runId"`
	SourceRun int64  `json:"sourceRunId"`
	Provider  string `json:"provider"`
	Owner     string `json:"owner"`
	Matched   int    `json:"matched"`
	Changed   int    `json:"changed"`
	Skipped   int    `json:"skipped"`
	Failed    int    `json:"failed"`
	ExitCode  int    `json:"exitCode"`
}

// renderUndoJSON returns the canonical JSON bytes (no trailing newline).
// Determinism: encoding/json preserves struct field declaration order
// so the byte output is stable for identical inputs — important for
// golden tests added later.
func renderUndoJSON(s undoJSONSummary) ([]byte, error) {
	return json.Marshal(s)
}
