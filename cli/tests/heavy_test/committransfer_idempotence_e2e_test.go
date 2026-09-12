package heavy_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/committransfer"
)

func TestCommitTransfer_PlanIdempotenceBeyond200Commits(t *testing.T) {
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

	mustCommitCountAt(t, source, "buried.txt", "buried\n",
		fmt.Sprintf("buried commit %d", buriedIdx), stamp(0))
	buriedShort := mustShortSHA(t, source, "HEAD")

	mustCommitCountAt(t, target, "anchor.txt", "a\n",
		"unrelated\n\ngitmap-replay: from "+sourceDisplay+" "+buriedShort+
			"\ngitmap-replay-cmd: commit-in\ngitmap-replay-at: 2026-01-01T00:00:00Z",
		stamp(1))
	for i := 0; i < buryCount; i++ {
		mustCommitCountAt(t, target, fmt.Sprintf("bury-%d.txt", i),
			"x\n", fmt.Sprintf("bury %d", i), stamp(2+i))
	}

	opts := committransfer.Options{
		LogPrefix: "[t]",
		Message: committransfer.MessagePolicy{
			Provenance:        true,
			SourceDisplayName: sourceDisplay,
			CommandName:       "commit-in",
		},
	}

	plan, err := committransfer.BuildPlan(source, target, opts)
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}

	if len(plan.Commits) != 1 {
		t.Fatalf("plan.Commits = %d, want 1", len(plan.Commits))
	}

	if got := plan.Commits[0].SkipCause; got != "already-replayed" {
		t.Errorf("SkipCause = %q, want %q", got, "already-replayed")
	}
}

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

func mustShortSHA(t *testing.T, dir, ref string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", dir, "rev-parse", "--short", ref).Output()
	if err != nil {
		t.Fatalf("rev-parse --short %s: %v", ref, err)
	}

	return strings.TrimSpace(string(out))
}
