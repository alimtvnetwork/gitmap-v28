// Package cmd — visibilityundoflags_test.go: unit coverage for
// parseVisUndoArgs (the `vu` / `vr` flag parser shipped in v6.7-v6.9),
// matchesFromResults (adapter feeding the shared audit pipeline),
// and bulkExitCode (collapsed-tally → exit code matrix). These guard
// against silent demotion of --force to a no-op, mis-routing of
// --run <id>, and drift in the bulk exit-code contract.
package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestParseUndoArgsDefaults(t *testing.T) {
	f := parseVisUndoArgs(nil)
	if f.Verbose || f.DryRun || f.Force || f.RunID != 0 {
		t.Fatalf("expected zero-valued undoFlags, got %+v", f)
	}
}

func TestParseUndoArgsAllFlags(t *testing.T) {
	f := parseVisUndoArgs([]string{"--verbose", "--dry-run", "--force", "--run", "42"})
	if !f.Verbose || !f.DryRun || !f.Force || f.RunID != 42 {
		t.Fatalf("expected all flags set + RunID=42, got %+v", f)
	}
}

func TestParseUndoArgsForceWithoutOtherFlags(t *testing.T) {
	f := parseVisUndoArgs([]string{"--force"})
	if !f.Force || f.DryRun || f.Verbose || f.RunID != 0 {
		t.Fatalf("expected only Force=true, got %+v", f)
	}
}

func TestParseUndoArgsIgnoresUnknownTokens(t *testing.T) {
	f := parseVisUndoArgs([]string{"garbage", "--force", "extra"})
	if !f.Force {
		t.Fatalf("expected Force=true despite unknown tokens, got %+v", f)
	}
}

func TestMatchesFromResultsPreservesNameAndPattern(t *testing.T) {
	in := []model.MakeAllVisibilityResultRecord{
		{RepoName: "alpha", MatchedPattern: "p1"},
		{RepoName: "beta", MatchedPattern: "p2"},
	}
	got := matchesFromResults(in)
	if len(got) != 2 || got[0].RepoName != "alpha" || got[0].MatchedPattern != "p1" {
		t.Fatalf("unexpected matches[0]: %+v", got)
	}
	if got[1].RepoName != "beta" || got[1].MatchedPattern != "p2" {
		t.Fatalf("unexpected matches[1]: %+v", got)
	}
}

func TestBulkExitCodeMatrix(t *testing.T) {
	cases := []struct {
		name            string
		changed, failed int
		want            int
	}{
		{"all-ok", 5, 0, constants.ExitVisOK},
		{"all-failed", 0, 3, constants.ExitVisAuthFailed},
		{"mixed", 2, 1, constants.ExitVisBulkPartial},
		{"zero-zero", 0, 0, constants.ExitVisOK},
	}
	for _, c := range cases {
		if got := bulkExitCode(c.changed, c.failed); got != c.want {
			t.Errorf("%s: bulkExitCode(%d,%d)=%d, want %d",
				c.name, c.changed, c.failed, got, c.want)
		}
	}
}
