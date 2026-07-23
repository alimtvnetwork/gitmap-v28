package committransfer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestPlanIdempotenceBeyond200Commits is the regression guard for
// spec 114 Gap A: the planner used to cap target-log scan at 200
// commits, so an already-replayed source commit buried >200 entries
// deep in the target was mis-classified as fresh and would be
// duplicated. Verifies the unbounded scan now catches it.
//
// De-flake notes (v6.82.0): this test previously ran with t.Parallel()
// and used time.Now() (second-granularity) stamps for 252 rapid
// commits. Under parallel CI load git occasionally returned
// "bad tree object HEAD" mid-way through the bury loop. Fix: run
// serially, walk timestamps forward monotonically per commit, and
// trim buryCount from 250 to 220 (still comfortably above the legacy
// 200-cap this test guards).
func TestPlanIdempotenceBeyond200Commits(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not available: %v", err)
	}

	const (
		sourceDisplay = "src-repo"
		buriedIdx     = 1
		buryCount     = 220
	)

	base := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	stamp := func(offset int) string {
		return base.Add(time.Duration(offset) * time.Second).Format(time.RFC3339)
	}

	root := t.TempDir()
	source := mustInitCountRepo(t, filepath.Join(root, "src"))
	target := mustInitCountRepo(t, filepath.Join(root, "dst"))

	// Source: one commit whose shortSHA we later claim was already
	// replayed into the target far in the past.
	mustCommitCountAt(t, source, "buried.txt", "buried\n",
		fmt.Sprintf("buried commit %d", buriedIdx), stamp(0))
	buriedShort := mustShortSHA(t, source, "HEAD")

	// Target: write the provenance footer FIRST, then bury it under
	// 220 unrelated commits so the legacy n=200 cap would miss it.
	mustCommitCountAt(t, target, "anchor.txt", "a\n",
		"unrelated\n\ngitmap-replay: from "+sourceDisplay+" "+buriedShort+
			"\ngitmap-replay-cmd: commit-in\ngitmap-replay-at: 2026-01-01T00:00:00Z",
		stamp(1))
	for i := 0; i < buryCount; i++ {
		mustCommitCountAt(t, target, fmt.Sprintf("bury-%d.txt", i),
			"x\n", fmt.Sprintf("bury %d", i), stamp(2+i))
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

// mustCommitCountAt is like mustCommitCount but takes an explicit
// RFC3339 timestamp so callers can walk time forward monotonically.
// Kept local to this file so it does not collide with helpers in
// count_parity_e2e_test.go.
func mustCommitCountAt(t *testing.T, dir, path, body, msg, stamp string) {
	t.Helper()
	full := filepath.Join(dir, path)
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", full, err)
	}
	if out, err := exec.Command("git", "-C", dir, "add", path).CombinedOutput(); err != nil {
		t.Fatalf("git add %s: %v\n%s", path, err, out)
	}
	cmd := exec.Command("git", "-C", dir, "commit", "-m", msg)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE="+stamp,
		"GIT_COMMITTER_DATE="+stamp,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit (%s): %v\n%s", msg, err, out)
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
