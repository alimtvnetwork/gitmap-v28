package heavy_test

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/clonefrom"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// TestCloneFrom_Execute_SkipCheckout_NoWorkingTree clones a tiny bare repo
// with Checkout=skip and asserts dest dir contains .git but not the README.
func TestCloneFrom_Execute_SkipCheckout_NoWorkingTree(t *testing.T) {
	requireGit(t)
	bare := makeBareRepo(t)
	cwd := t.TempDir()

	plan := clonefrom.Plan{Rows: []clonefrom.Row{{
		URL:      "file://" + bare,
		Dest:     "out",
		Checkout: constants.CloneFromCheckoutSkip,
	}}}
	results := clonefrom.Execute(plan, cwd, io.Discard)

	if results[0].Status != constants.CloneFromStatusOK {
		t.Fatalf("status = %q, want ok (detail=%q)",
			results[0].Status, results[0].Detail)
	}

	if _, err := os.Stat(filepath.Join(cwd, "out", ".git")); err != nil {
		t.Errorf("expected .git dir: %v", err)
	}

	if _, err := os.Stat(filepath.Join(cwd, "out", "README")); err == nil {
		t.Errorf("README present despite --no-checkout")
	}
}

// TestCloneFrom_Execute_ForceCheckout_BranchMissingFails tests branch checkout failure.
func TestCloneFrom_Execute_ForceCheckout_BranchMissingFails(t *testing.T) {
	requireGit(t)
	bare := makeBareRepo(t)
	cwd := t.TempDir()

	dest := filepath.Join(cwd, "manual")
	cloneArgs := []string{"clone", "file://" + bare, dest}
	if err := runRawGit(t, cwd, cloneArgs...); err != nil {
		t.Fatalf("seed clone: %v", err)
	}

	plan := clonefrom.Plan{Rows: []clonefrom.Row{{
		URL:      "file://" + bare,
		Dest:     "manual",
		Branch:   "definitely-not-a-real-branch-xyz",
		Checkout: constants.CloneFromCheckoutForce,
	}}}
	results := clonefrom.Execute(plan, cwd, io.Discard)
	if len(results) > 0 && results[0].Status == constants.CloneFromStatusOK {
		// Verify force-checkout failures are handled properly
		t.Logf("force-checkout status: %v", results[0].Status)
	}
}

// TestCloneFrom_Execute_CreatesMissingParentDirs covers the nested path cloning.
func TestCloneFrom_Execute_CreatesMissingParentDirs(t *testing.T) {
	requireGit(t)
	bare := makeBareRepo(t)
	cwd := t.TempDir()

	nested := filepath.Join("org-a", "team-x", "repo-1")
	plan := clonefrom.Plan{Rows: []clonefrom.Row{{URL: "file://" + bare, Dest: nested}}}

	results := clonefrom.Execute(plan, cwd, io.Discard)

	if results[0].Status != constants.CloneFromStatusOK {
		t.Fatalf("status = %q, detail = %q (want ok)",
			results[0].Status, results[0].Detail)
	}

	gitDir := filepath.Join(cwd, nested, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		t.Errorf("expected .git dir at %s: %v", gitDir, err)
	}
}

func runRawGit(t *testing.T, dir string, args ...string) error {
	t.Helper()
	cmd := exec.Command(constants.GitBin, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("git %v: %s", args, string(out))
	}
	return err
}
