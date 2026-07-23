package cmd

import (
	"encoding/json"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestRenderUndoJSONShapeStable(t *testing.T) {
	s := undoJSONSummary{
		Command: constants.CmdVisibilityUndo, RunID: 42, SourceRun: 7,
		Provider: constants.ProviderGitHub, Owner: "acme",
		Matched: 3, Changed: 2, Skipped: 1, Failed: 0, ExitCode: 0,
	}
	got, err := renderUndoJSON(s)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	var round undoJSONSummary
	if err := json.Unmarshal(got, &round); err != nil {
		t.Fatalf("round-trip: %v / %s", err, got)
	}
	if round != s {
		t.Fatalf("round-trip diff:\n got %+v\nwant %+v", round, s)
	}
}

func TestRenderUndoJSONZeroValuesEmitted(t *testing.T) {
	got, err := renderUndoJSON(undoJSONSummary{Command: constants.CmdVisibilityRedo})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	for _, key := range []string{
		`"runId":0`, `"sourceRunId":0`, `"matched":0`, `"changed":0`,
		`"skipped":0`, `"failed":0`, `"exitCode":0`,
	} {
		if !containsJSONFragment(string(got), key) {
			t.Fatalf("missing zero-key %q in %s", key, got)
		}
	}
}

func containsJSONFragment(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}

	return false
}
