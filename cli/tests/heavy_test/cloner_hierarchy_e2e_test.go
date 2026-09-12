package heavy_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func runClonerGit(t *testing.T, workdir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	if workdir != "" {
		cmd.Dir = workdir
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v (workdir=%q) failed: %v\n--- output ---\n%s",
			args, workdir, err, string(out))
	}
}

func writeClonerFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdirall %s: %v", filepath.Dir(path), err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func initClonerLocalRemote(t *testing.T, tmp, label string) string {
	t.Helper()

	safe := strings.ReplaceAll(label, "/", "_")
	dir := filepath.Join(tmp, "remotes", safe)
	runClonerGit(t, "", "init", "-b", "main", dir)
	writeClonerFile(t, filepath.Join(dir, "README.md"), "# "+label+"\n")
	runClonerGit(t, dir, "add", ".")
	runClonerGit(t, dir, "-c", "user.email=test@example.com", "-c", "user.name=test",
		"commit", "-m", "init")

	return dir
}

func makeClonerRecordsWithRemotes(t *testing.T, tmp string, relPaths []string) []model.ScanRecord {
	t.Helper()

	records := make([]model.ScanRecord, 0, len(relPaths))
	for _, rel := range relPaths {
		remote := initClonerLocalRemote(t, tmp, rel)
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

func assertClonerHierarchy(t *testing.T, target string, relPaths []string) {
	t.Helper()

	for _, rel := range relPaths {
		dest := filepath.Join(target, filepath.FromSlash(rel))
		if !cloner.IsGitRepo(dest) {
			t.Fatalf("hierarchy lost: %q is not a git repo under %q", rel, target)
		}
	}
}

func modeClonerName(workers int) string {
	if workers <= 1 {
		return "sequential"
	}

	return "parallel"
}

// TestCloneAllPreservesNestedHierarchy asserts that cloner reproduces the exact nested layout.
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
		t.Run(modeClonerName(workers), func(t *testing.T) {
			tmp := t.TempDir()
			records := makeClonerRecordsWithRemotes(t, tmp, relPaths)
			target := filepath.Join(tmp, "out")

			summary := cloner.CloneAll(records, target, cloner.CloneOptions{
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

			assertClonerHierarchy(t, target, relPaths)
		})
	}
}
