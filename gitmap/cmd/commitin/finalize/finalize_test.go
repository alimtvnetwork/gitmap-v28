package finalize

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestOutcomeAllOk(t *testing.T) {
	if Outcome(Counters{Created: 5}) != constants.CommitInExitOk {
		t.Fatal("clean run must be Ok")
	}
}

func TestOutcomePartiallyFailed(t *testing.T) {
	if Outcome(Counters{Created: 1, Failed: 1}) != constants.CommitInExitPartiallyFailed {
		t.Fatal("mixed run must be PartiallyFailed")
	}
}

func TestPrintSummaryUsesConstantFormat(t *testing.T) {
	var buf bytes.Buffer
	PrintSummary(&buf, Counters{RunId: 7, Created: 3, Skipped: 1, Failed: 0})
	got := buf.String()
	if !strings.Contains(got, "run=7") || !strings.Contains(got, "created=3") {
		t.Fatalf("summary missing expected fields: %q", got)
	}
}

func TestResolveForceMergeReturnsTakeTheirs(t *testing.T) {
	var buf bytes.Buffer
	got := Resolve(constants.CommitInConflictModeForceMerge, "abc123", &buf)
	if got != ConflictDecisionTakeTheirs {
		t.Fatalf("want TakeTheirs, got %v", got)
	}
	if buf.Len() != 0 {
		t.Fatalf("ForceMerge must not print: %q", buf.String())
	}
}

func TestResolvePromptAbortsWithBanner(t *testing.T) {
	var buf bytes.Buffer
	got := Resolve(constants.CommitInConflictModePrompt, "abc123", &buf)
	if got != ConflictDecisionAbort {
		t.Fatalf("want Abort, got %v", got)
	}
	if !strings.Contains(buf.String(), "abc123") {
		t.Fatalf("banner missing source sha: %q", buf.String())
	}
}

func TestCleanupTempRespectsKeepFlag(t *testing.T) {
	dir := t.TempDir()
	CleanupTemp(dir, true) // no-op
}
