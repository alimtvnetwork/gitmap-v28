package e2e

import (
	"strings"
	"testing"
	"time"
)

// TestHappyPathSingleInputReplaysAllCommits is the canonical commit-in
// happy path: one source repo (auto-init in place) + one input repo
// with three commits → orchestrator should produce three commits in
// the destination, in the same order, with both author + committer
// dates byte-equal to the source.
func TestHappyPathSingleInputReplaysAllCommits(t *testing.T) {
	src, input := buildHappyPathFixture(t)

	res := Run(t, NewRawArgs(src.Path, input.Path))
	if res.ExitCode != 0 {
		t.Fatalf("exit=%d, want 0\nstderr=%s", res.ExitCode, res.Stderr)
	}

	dst := src.LogFirstParent(t)
	if len(dst) != 3 {
		t.Fatalf("dst commits=%d, want 3\nstderr=%s", len(dst), res.Stderr)
	}
	wantSubjects := []string{"first", "second", "third"}
	for i, want := range wantSubjects {
		if dst[i].Subject != want {
			t.Errorf("dst[%d].Subject=%q, want %q", i, dst[i].Subject, want)
		}
	}
	// Date replication: spec §3 hard rule — both AuthorDate AND
	// CommitterDate must equal the source commit's. The fixture pinned
	// the third commit at 2024-03-10T09:00:00Z.
	if !strings.HasPrefix(dst[2].AuthorDate, "2024-03-10T09:00:00") {
		t.Errorf("third AuthorDate=%q, want 2024-03-10T09:00:00 prefix", dst[2].AuthorDate)
	}
	if !strings.HasPrefix(dst[2].CommitDate, "2024-03-10T09:00:00") {
		t.Errorf("third CommitDate=%q, want 2024-03-10T09:00:00 prefix", dst[2].CommitDate)
	}
}

// TestSecondRunDedupesAllInputCommits replays the same input twice
// against the same source. The first run produces 3 commits; the
// second run must produce zero new commits (every SourceSha hits
// ShaMap → DuplicateSourceSha skip per spec §3.1 stage 10).
func TestSecondRunDedupesAllInputCommits(t *testing.T) {
	src, input := buildHappyPathFixture(t)

	first := Run(t, NewRawArgs(src.Path, input.Path))
	if first.ExitCode != 0 {
		t.Fatalf("first run exit=%d\nstderr=%s", first.ExitCode, first.Stderr)
	}
	src.AssertCommitCount(t, 3)

	second := Run(t, NewRawArgs(src.Path, input.Path))
	if second.ExitCode != 0 {
		t.Fatalf("second run exit=%d\nstderr=%s", second.ExitCode, second.Stderr)
	}
	// Dedupe means zero new commits — destination is unchanged.
	src.AssertCommitCount(t, 3)
	// Summary line must report the dedupes as skipped, not created.
	if !strings.Contains(second.Stderr, "skipped=3") {
		t.Errorf("second run summary missing skipped=3\nstderr=%s", second.Stderr)
	}
	if !strings.Contains(second.Stderr, "created=0") {
		t.Errorf("second run summary missing created=0\nstderr=%s", second.Stderr)
	}
}

// buildHappyPathFixture is the shared two-repo setup used by both
// happy-path tests. Returned `src` is an empty repo (auto-init by
// orchestrator); `input` has three commits at fixed timestamps.
func buildHappyPathFixture(t *testing.T) (src, input *Repo) {
	t.Helper()
	src = NewRepo(t, "src")
	input = NewRepo(t, "input")
	input.Commit("a.txt", "1\n", "first", time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC))
	input.Commit("b.txt", "2\n", "second", time.Date(2024, 2, 10, 9, 0, 0, 0, time.UTC))
	input.Commit("c.txt", "3\n", "third", time.Date(2024, 3, 10, 9, 0, 0, 0, time.UTC))
	return src, input
}
