package cloner

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// TestNormalizeWorkers locks in the clamping rules so that the
// --max-concurrency flag never spins up more goroutines than there is
// work to do, and never silently degrades a 0/negative request to
// something other than sequential.
func TestNormalizeWorkers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		requested int
		jobs      int
		want      int
	}{
		{"zero clamps to one", 0, 5, 1},
		{"negative clamps to one", -7, 5, 1},
		{"one stays one", 1, 5, 1},
		{"under job count is preserved", 3, 5, 3},
		{"equal to job count is preserved", 5, 5, 5},
		{"over job count clamps down", 99, 5, 5},
		{"empty job list keeps requested", 4, 0, 4},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := normalizeWorkers(tc.requested, tc.jobs)
			if got != tc.want {
				t.Fatalf("normalizeWorkers(%d,%d) = %d, want %d",
					tc.requested, tc.jobs, got, tc.want)
			}
		})
	}
}

// TestCloneAllPreservesNestedHierarchy is the canonical assertion that
// `gitmap clone` reproduces the exact nested folder layout captured by
// `gitmap scan`. We point each ScanRecord at a tiny local "remote" repo
// (init + commit) and verify each clone lands at its recorded
// RelativePath under targetDir, including multi-segment paths.
//
// Sequential and parallel runners are exercised through the same entry
// point so any future refactor that breaks hierarchy preservation in
// either path will fail this test.
func TestCloneAllPreservesNestedHierarchy(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not on PATH")
	}

	relPaths := []string{
		"flat-repo",
		"group-a/repo-a1",
		"group-a/repo-a2",
		"group-b/sub/repo-b1",
		"deep/very/deep/leaf",
	}

	for _, workers := range []int{1, 3} {
		t.Run(modeName(workers), func(t *testing.T) {
			tmp := t.TempDir()
			records := makeRecordsWithRemotes(t, tmp, relPaths)
			target := filepath.Join(tmp, "out")

			summary := cloneAll(records, target, CloneOptions{
				Quiet:          true,
				MaxConcurrency: workers,
			})

			if summary.Failed != 0 {
				t.Fatalf("workers=%d: unexpected failures: %+v", workers, summary.Errors)
			}
			if summary.Succeeded != len(relPaths) {
				t.Fatalf("workers=%d: succeeded=%d want=%d",
					workers, summary.Succeeded, len(relPaths))
			}
			assertHierarchy(t, target, relPaths)
		})
	}
}

// modeName maps a worker count to a stable subtest label.
func modeName(workers int) string {
	if workers <= 1 {
		return "sequential"
	}

	return "parallel"
}

// makeRecordsWithRemotes creates a bare-ish source repo for each
// rel-path and returns ScanRecords pointing at them via file:// URLs.
func makeRecordsWithRemotes(t *testing.T, tmp string, relPaths []string) []model.ScanRecord {
	t.Helper()

	records := make([]model.ScanRecord, 0, len(relPaths))
	for _, rel := range relPaths {
		remote := initLocalRemote(t, tmp, rel)
		records = append(records, model.ScanRecord{
			RepoName:     filepath.Base(rel),
			RelativePath: rel,
			HTTPSUrl:     "file://" + filepath.ToSlash(remote),
			BranchSource: "default",
			Branch:       "main",
		})
	}

	return records
}

// initLocalRemote builds a repo with a single commit on `main` and
// returns its absolute path. Used as the clone source so the test runs
// fully offline.
func initLocalRemote(t *testing.T, tmp, label string) string {
	t.Helper()

	safe := strings.ReplaceAll(label, "/", "_")
	dir := filepath.Join(tmp, "remotes", safe)
	runGit(t, "", "init", "-b", "main", dir)
	writeFile(t, filepath.Join(dir, "README.md"), "# "+label+"\n")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "-c", "user.email=test@example.com", "-c", "user.name=test",
		"commit", "-m", "init")

	return dir
}

// assertHierarchy verifies that every relPath landed under target as a
// real git repo. Fails fast on the first missing path so the error
// message identifies the specific layout violation.
func assertHierarchy(t *testing.T, target string, relPaths []string) {
	t.Helper()

	for _, rel := range relPaths {
		dest := filepath.Join(target, filepath.FromSlash(rel))
		if !IsGitRepo(dest) {
			t.Fatalf("hierarchy lost: %q is not a git repo under %q", rel, target)
		}
	}
}
