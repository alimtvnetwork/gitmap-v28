package cmdpull

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// RunBatchRemediateRebase resolves all failed repositories using pull --rebase.
func RunBatchRemediateRebase(failures []PullFailureSummary) error {
	fmt.Printf("\n  %sRebasing %d failed repositories...%s\n\n",
		constants.ColorBold, len(failures), constants.ColorReset)
	successCount := 0
	for _, f := range failures {
		if remediateSingleRepoRebase(f) {
			successCount++
		}
	}
	printRemediationBatchSummary(successCount, len(failures))

	return nil
}

func remediateSingleRepoRebase(f PullFailureSummary) bool {
	if f.RepoPath == "" {
		return false
	}

	if cloner.IsMissingRepo(f.RepoPath) {
		return RemediateMissingRepo(f)
	}

	fmt.Printf("  • %s%s%s\n", constants.ColorBold, f.RepoName, constants.ColorReset)
	fmt.Printf("    %s Attempting pull --rebase...\n", constants.TreeBranch)
	out, err := execPullRebase(f.RepoPath)
	if err == nil {
		fmt.Printf("    %s %s✓%s Successfully rebased and updated.\n",
			constants.TreeCorner, constants.ColorGreen, constants.ColorReset)
		return true
	}

	return handleRebaseRemediationFailure(f.RepoPath, f.RepoName, string(out))
}

func execPullRebase(repoPath string) ([]byte, error) {
	cmd := exec.Command("git", "-C", repoPath, "pull", "--progress", "--rebase", "--autostash")
	cmd.Env = gitutil.BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)

	return cmd.CombinedOutput()
}

func handleRebaseRemediationFailure(repoPath, repoName, out string) bool {
	cloner.SafeAbortRebase(repoPath)
	relPath := formatRelativeOrCleanPath(repoPath)
	if relPath == "" || relPath == "." {
		relPath = repoName
	}

	isConflict := cloner.IsMergeConflictOutput(out) || strings.Contains(strings.ToLower(out), "conflict")
	if isConflict {
		fmt.Printf("    %s %s✗%s Rebase conflict detected (rebase safely aborted).\n",
			constants.TreeCorner, constants.ColorRed, constants.ColorReset)
		fmt.Printf("        %s[hint]%s Inspect rebase status: git -C %s status\n",
			constants.ColorYellow, constants.ColorReset, relPath)

		return false
	}

	fmt.Printf("    %s %s✗%s Rebase failed (rebase safely aborted).\n",
		constants.TreeCorner, constants.ColorRed, constants.ColorReset)
	fmt.Printf("        %s[hint]%s Inspect rebase status: git -C %s status\n",
		constants.ColorYellow, constants.ColorReset, relPath)

	return false
}
