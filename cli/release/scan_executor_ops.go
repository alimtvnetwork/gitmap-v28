package release

import (
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func createBranch(repoDir, branchName, hash string, action *ScanCommitAction) error {
	cmd := exec.Command("git", "branch", branchName, hash)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, "createBranch", map[string]any{"branch": branchName})
	}

	action.IsBranchCreated = true

	return nil
}

func createTag(repoDir, tagName, hash string, action *ScanCommitAction) error {
	cmd := exec.Command("git", "tag", tagName, hash)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, "createTag", map[string]any{"tag": tagName})
	}

	action.IsTagCreated = true

	return nil
}

func isRefPresent(repoDir, refPath string) (bool, error) {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", refPath)
	cmd.Dir = repoDir
	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	if _, isExit := err.(*exec.ExitError); isExit {
		return false, nil
	}

	return false, apperror.Wrap(err, "isRefPresent", map[string]any{"ref": refPath})
}
