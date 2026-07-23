package committransfer

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestPlanIdempotenceBeyond200Commits is the regression guard for
// spec 114 Gap A: the planner used to cap target-log scan at 200
// commits, so an already-replayed source commit buried >200 entries
// deep in the target was mis-classified as fresh and would be
// duplicated. Verifies the unbounded scan now catches it.
func TestPlanIdempotenceBeyond200Commits(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not available: %v", err)
	}

	const (
		sourceDisplay = "src-repo"
		burySource    = "src-old"
		buriedIdx     = 1
		buryCount     = 250
	)

	root := t.TempDir()
	source := mustInitCountRepo(t, filepath.Join(root, "src"))
	target := mustInitCountRepo(t, filepath.Join(root, "dst"))

	// Source: one commit whose shortSHA we will later "claim" was
	// already replayed into the target far in the past.
	mustCommitCount(t, source, "buried.txt", "buried\n",
		fmt.Sprintf("buried commit %d", buriedIdx))
	buriedShort := mustShortSHA(t, source, "HEAD")

	// Target: write the provenance footer FIRST, then bury it under
	// 250 unrelated commits so the legacy n=200 cap would miss it.
	mustCommitCount(t, target, "anchor.txt", "a\n",
		"unrelated\n\ngitmap-replay: from "+sourceDisplay+" "+buriedShort+
			"\ngitmap-replay-cmd: commit-in\ngitmap-replay-at: 2026-01-01T00:00:00Z")
	for i := 0; i < buryCount; i++ {
		mustCommitCount(t, target, fmt.Sprintf("bury-%d.txt", i),
			"x\n", fmt.Sprintf("bury %d", i))
	}

	opts := Options{
		LogPrefix: "[t]",
		Message: MessagePolicy{
			Provenance:        true,
			SourceDisplayName: sourceDisplay,
			CommandName:       "commit-in",
		},
	}

	plan, err := BuildPlan(source, target, opts)
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	if len(plan.Commits) != 1 {
		t.Fatalf("plan.Commits = %d, want 1", len(plan.Commits))
	}
	if got := plan.Commits[0].SkipCause; got != "already-replayed" {
		t.Errorf("SkipCause = %q, want %q (idempotence broke; legacy 200-cap regression)",
			got, "already-replayed")
	}
}

// mustShortSHA returns the abbreviated SHA git computes for ref in dir.
func mustShortSHA(t *testing.T, dir, ref string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", dir, "rev-parse", "--short", ref).Output()
	if err != nil {
		t.Fatalf("rev-parse --short %s: %v", ref, err)
	}

	return strings.TrimSpace(string(out))
}
