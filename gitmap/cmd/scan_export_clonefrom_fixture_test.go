package cmd

// Fixture helpers for the scan→export→clone-from integration test.
// Split out of scan_export_clonefrom_integration_test.go so that
// file stays under the project's 200-line per-file cap. Public
// surface is only used by sibling _test.go files in this package;
// nothing here is reachable from production code.

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// seedScanTree creates two worktrees under one root, each with a
// real git remote pointing at the shared bare repo via file://.
// Two repos exercises mapper.BuildRecords' loop and clone-from's
// per-row Execute path on a non-trivial input.
func seedScanTree(t *testing.T, bare string) string {
	t.Helper()
	root := t.TempDir()
	for _, name := range []string{"alpha", "beta"} {
		seedOneWorktree(t, root, name, bare)
	}

	return root
}

// seedOneWorktree clones the bare repo into root/name and pins the
// remote URL to file://<bare> so mapper.BuildRecords sees a valid
// remote for both ModeHTTPS and ModeSSH lookups.
func seedOneWorktree(t *testing.T, root, name, bare string) {
	t.Helper()
	dest := filepath.Join(root, name)
	runIntegrationGit(t, root, "clone", "-q", "file://"+bare, name)
	runIntegrationGit(t, dest, "remote", "set-url", "origin", "file://"+bare)
}

// requireGitForIntegration mirrors the unit-test helper: skip the
// whole integration when git isn't on PATH so minimal CI containers
// stay green. Also skips when the sandbox forbids `git add` (some
// hardened environments inject a wrapper that refuses index writes);
// the integration here can't proceed without commits in the bare repo.
func requireGitForIntegration(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not on PATH: %v", err)
	}
	probe := t.TempDir()
	if err := exec.Command("git", "-C", probe, "init", "-q").Run(); err != nil {
		t.Skipf("git init blocked in sandbox: %v", err)
	}
	if err := os.WriteFile(filepath.Join(probe, "x"), []byte("x"), 0o644); err != nil {
		t.Skipf("tempdir write blocked: %v", err)
	}
	if err := exec.Command("git", "-C", probe, "add", "x").Run(); err != nil {
		t.Skipf("git add blocked in sandbox: %v", err)
	}
}

// makeIntegrationBareRepo builds a one-commit bare repo and returns
// its absolute path. Lives in cmd/ rather than reusing
// clonefrom.makeBareRepo because that helper is unexported in the
// clonefrom test package.
func makeIntegrationBareRepo(t *testing.T) string {
	t.Helper()
	work := t.TempDir()
	bare := filepath.Join(t.TempDir(), "src.git")
	runIntegrationGit(t, work, "init", "-q")
	runIntegrationGit(t, work, "config", "user.email", "t@e")
	runIntegrationGit(t, work, "config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(work, "README"), []byte("hi"), 0o644); err != nil {
		t.Fatalf("seed README: %v", err)
	}
	runIntegrationGit(t, work, "add", ".")
	runIntegrationGit(t, work, "commit", "-q", "-m", "init")
	runIntegrationGit(t, work, "clone", "--bare", "-q", work, bare)

	return bare
}

// runIntegrationGit fatals on error with the combined output
// included so a CI failure points at the offending git invocation.
func runIntegrationGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, string(out))
	}
}
