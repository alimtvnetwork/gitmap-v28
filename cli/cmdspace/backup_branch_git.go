package cmdspace

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// errDirtyTree is the exact refusal message for a dirty working tree
// (spec 243.2: no override exists).
var errDirtyTree = errors.New("Working tree has uncommitted changes. Commit or stash first — backups snapshot HEAD only.")

// runGitCapture runs git and returns trimmed stdout; stderr goes to the
// returned error context.
func runGitCapture(args ...string) (string, error) {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "git "+strings.Join(args, " ")+":")
	}

	return strings.TrimSpace(string(out)), nil
}

// runGitInherit runs git with output streaming to the user (used for
// branch creation and push progress).
func runGitInherit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "git "+strings.Join(args, " ")+":")
	}

	return nil
}

// requireCleanTree refuses the backup when the working tree has any
// uncommitted changes — backups snapshot HEAD only. Untracked gitmap
// diagnostic state (`.gitmap/`) is ignored: writeLastErrorFile drops it
// into CWD on every dispatch error, so counting it would let gitmap's own
// diagnostics block a retry of the failed command.
func requireCleanTree() error {
	status, err := runGitCapture("status", "--porcelain")
	if err != nil {
		return err
	}

	if treeHasUserDirt(status) {
		return errDirtyTree
	}

	return nil
}

// treeHasUserDirt reports whether porcelain v1 output contains any change
// other than untracked gitmap diagnostic state. Positive naming: true
// means the tree has user dirt worth refusing the backup for.
func treeHasUserDirt(status string) bool {
	for _, line := range strings.Split(status, "\n") {
		if strings.TrimSpace(line) == "" || isGitmapUntrackedDirt(line) {
			continue
		}

		return true
	}

	return false
}

// isGitmapUntrackedDirt reports whether a porcelain v1 line refers only to
// untracked gitmap diagnostic state (`?? .gitmap/` or a file under it).
// Tracked modifications under `.gitmap/` still count as user dirt.
func isGitmapUntrackedDirt(line string) bool {
	if len(line) < 4 || line[:3] != "?? " {
		return false
	}

	path := line[3:]

	return path == ".gitmap" || strings.HasPrefix(path, ".gitmap/")
}

// branchExists reports whether the local branch already exists.
func branchExists(branch string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch)

	return cmd.Run() == nil
}

// ensureBranchAbsence errors naming the existing branch unless --force
// is passed to recreate it at current HEAD.
func ensureBranchAbsence(branch string, force bool) error {
	if !branchExists(branch) {
		return nil
	}

	if force {
		return nil
	}

	return fmt.Errorf("branch %s already exists — pick a different task string or pass --force to recreate it at current HEAD", branch)
}

// createBackupBranch creates (or with --force recreates at HEAD) the
// backup branch from the current HEAD.
func createBackupBranch(branch string, force bool) error {
	if force {
		return runGitInherit("branch", "-f", branch)
	}

	return runGitInherit("branch", branch)
}
