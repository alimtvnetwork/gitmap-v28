package heavy_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
)

func makeRepo(t *testing.T, dir string, uniqueBody bool) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	body := "shared\n"
	if uniqueBody {
		body = "unique-" + filepath.Base(dir) + "\n"
	}

	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	run := func(args ...string) {
		c := exec.Command("git", append([]string{"-C", dir}, args...)...)
		c.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t",
			"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t",
		)
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	run("init", "-q", "-b", "main")
	run("add", ".")
	run("commit", "-q", "-m", "init")
}

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not on PATH")
	}
}

func TestHygieneIntegrationScansAndProbes(t *testing.T) {
	requireGit(t)
	root := t.TempDir()
	a := filepath.Join(root, "a")
	b := filepath.Join(root, "b")
	makeRepo(t, a, false)
	makeRepo(t, b, false)
	if err := os.MkdirAll(filepath.Join(root, "not-a-repo"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	repos := cmd.ScanForReposParallel(root)
	if len(repos) != 2 {
		t.Fatalf("scanForReposParallel got %d repos, want 2: %v", len(repos), repos)
	}

	for _, r := range repos {
		if _, ok := cmd.LastCommitTime(r); !ok {
			t.Fatalf("lastCommitTime(%s) failed", r)
		}

		if sz := cmd.DirSize(filepath.Join(r, ".git")); sz <= 0 {
			t.Fatalf("dirSize(%s) = %d, want > 0", r, sz)
		}
	}

	groups := map[string][]string{}
	for _, r := range repos {
		sha, ok := cmd.HeadTreeSHA(r)
		if !ok {
			t.Fatalf("headTreeSHA(%s) failed", r)
		}

		groups[sha] = append(groups[sha], r)
	}

	dupes := cmd.FilterDuplicateGroups(groups)
	if len(dupes) != 1 {
		t.Fatalf("expected 1 duplicate group, got %d", len(dupes))
	}
}

func TestHygieneIntegrationOrphanProbe(t *testing.T) {
	requireGit(t)
	root := t.TempDir()
	a := filepath.Join(root, "a")
	makeRepo(t, a, true)
	c := exec.Command("git", "-C", a, "remote", "add", "origin", "git@github.com:owner/repo.git")
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("remote add: %v\n%s", err, out)
	}

	u, ok := cmd.OriginURL(a)
	if !ok || u != "git@github.com:owner/repo.git" {
		t.Fatalf("originURL = %q ok=%v", u, ok)
	}

	if got := cmd.GitURLToHTTPS(u); got != "https://github.com/owner/repo" {
		t.Fatalf("gitURLToHTTPS = %q", got)
	}
}
