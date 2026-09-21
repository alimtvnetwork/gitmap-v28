package cmdgit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func createCompareWorktree(repoRoot, branch string) (string, error) {
	wtPath := buildWorktreePath(repoRoot, branch)
	if err := ensureWorktreeParentDir(wtPath); err != nil {
		return "", err
	}

	if err := addDetachedWorktree(repoRoot, wtPath, branch); err != nil {
		return "", err
	}

	return wtPath, nil
}

func buildWorktreePath(repoRoot, branch string) string {
	safeBranch := sanitizeBranchForPath(branch)
	dirName := fmt.Sprintf("compare-%s-%d", safeBranch, time.Now().UnixNano())

	return filepath.Join(repoRoot, ".gitmap", "tmp", dirName)
}

func sanitizeBranchForPath(branch string) string {
	r := strings.NewReplacer("/", "-", "\\", "-", ":", "-", " ", "-")

	return r.Replace(branch)
}

func ensureWorktreeParentDir(wtPath string) error {
	err := os.MkdirAll(filepath.Dir(wtPath), 0755)
	if err != nil {
		return apperror.WrapSimple(err, "create compare worktree directory")
	}

	return nil
}

func addDetachedWorktree(repoRoot, wtPath, branch string) error {
	_, err := runGitOutput(repoRoot, "worktree", "add", "--detach", wtPath, branch)
	if err != nil {
		return apperror.WrapSimple(err, "git worktree add failed")
	}

	return nil
}

func cleanupWorktrees(repoRoot string, worktrees []string) {
	for _, wt := range worktrees {
		cleanupSingleWorktree(repoRoot, wt)
	}

	pruneWorktrees(repoRoot)
}

func cleanupSingleWorktree(repoRoot, wtPath string) {
	_ = removeGitWorktree(repoRoot, wtPath)
	_ = os.RemoveAll(wtPath)
}

func removeGitWorktree(repoRoot, wtPath string) error {
	_, err := runGitOutput(repoRoot, "worktree", "remove", "--force", wtPath)

	return err
}

func pruneWorktrees(repoRoot string) {
	_, _ = runGitOutput(repoRoot, "worktree", "prune")
}

func runGitOutput(repoRoot string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if repoRoot != "" && repoRoot != "." {
		cmd.Dir = repoRoot
	}

	out, err := cmd.Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "git "+strings.Join(args, " "))
	}

	return strings.TrimSpace(string(out)), nil
}
