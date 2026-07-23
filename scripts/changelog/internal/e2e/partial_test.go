package e2e

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/scripts/changelog/internal/gitlog"
	"github.com/alimtvnetwork/gitmap-v28/scripts/changelog/internal/group"
	"github.com/alimtvnetwork/gitmap-v28/scripts/changelog/internal/render"
)

// TestPartialRangeOnlyIncludesSlicedCommits proves the --since / --release-tag
// path yields exactly the commits between the two bounds — and that the
// resulting Markdown / TypeScript fragments still come out byte-identical
// to a pinned fixture across reruns. This is the regression gate for
// `make changelog SINCE=... RELEASE_TAG=...`.
func TestPartialRangeOnlyIncludesSlicedCommits(t *testing.T) {
	repoRoot := t.TempDir()

	r := newGitRepo(t, repoRoot)
	seedFiles(t, repoRoot)
	seedHistory(r)
	r.finalize()

	commits, err := gitlog.CommitsInRange(repoRoot, "v0.9.0", "")
	if err != nil {
		t.Fatalf("CommitsInRange: %v", err)
	}

	for _, c := range commits {
		if c.Subject == "chore: bootstrap" {
			t.Fatalf("partial range leaked a commit from before --since")
		}
	}

	gotMD, gotTS := renderEntry(t, repoRoot, commits, "v1.0.0", "2026-04-24")

	assertFixture(t, "changelog.md.golden", gotMD)
	assertFixture(t, "changelog.ts.golden", gotTS)
}

// renderEntry mirrors what the runner does without re-running writer
// (we already cover writer in TestFixtureMarkdownAndTypeScript). It
// turns commits into the rendered fragments so we can diff in-memory.
func renderEntry(t *testing.T, _ string, commits []gitlog.Commit, version, date string) (string, string) {
	t.Helper()

	groups, _ := group.ByPrefix(commits)

	entry := render.Entry{Version: version, Date: date, Groups: groups}

	return render.Markdown(entry), "  " + render.TypeScript(entry)[2:]
}
