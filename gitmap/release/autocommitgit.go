package release

import (
	"fmt"
	"os/exec"
	"strings"

	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// commitAll stages everything and commits.
func commitAll(msg string) AutoCommitResult {
	err := stageAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitFailed, err)

		return AutoCommitResult{}
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: staged all files")
	}

	err = commitStaged(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitFailed, err)

		return AutoCommitResult{}
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: committed all: %s", msg)
	}

	fmt.Printf(constants.MsgAutoCommitAll, msg)

	err = pushCurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAutoCommitPush, err)

		return AutoCommitResult{Committed: true, AllFiles: true, Message: msg}
	}

	branch, branchErr := CurrentBranchName()
	if branchErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine current branch: %v\n", branchErr)
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: pushed all to %s", branch)
	}

	fmt.Printf(constants.MsgAutoCommitPushed, branch)

	return AutoCommitResult{Committed: true, AllFiles: true, Message: msg}
}

// stageFiles runs git add on specific files.
func stageFiles(files []string) error {
	args := append([]string{constants.GitAdd}, files...)

	return runGitCmd(args...)
}

// stageAll runs git add -A.
func stageAll() error {
	return runGitCmd(constants.GitAdd, constants.GitAddAll)
}

// commitStaged runs git commit -m <msg>.
func commitStaged(msg string) error {
	return runGitCmd(constants.GitCommit, constants.GitCommitMsg, msg)
}

// pushCurrentBranch pushes the current branch to origin.
func pushCurrentBranch() error {
	branch, err := CurrentBranchName()
	if err != nil {
		return err
	}

	pushOutput, err := runGitCmdCombined(constants.GitPush, constants.GitOrigin, branch)
	if err == nil {
		return nil
	}

	if !isNonFastForwardPushError(pushOutput) {
		return formatGitCommandError(pushOutput, err)
	}

	return syncBranchAndRetryPush(branch, pushOutput)
}

func syncBranchAndRetryPush(branch, pushOutput string) error {
	if verbose.IsEnabled() {
		verbose.Get().Log(
			"autocommit: push rejected for %s, attempting rebase sync: %s",
			branch,
			singleLineGitOutput(pushOutput),
		)
	}

	fmt.Printf(constants.MsgAutoCommitSyncRetry, branch)

	pullOutput, err := runGitCmdCombined(
		constants.GitPull,
		constants.GitPullRebaseFlag,
		constants.GitOrigin,
		branch,
	)
	if err != nil {
		abortRebaseAfterFailure()

		return fmt.Errorf(
			"remote branch advanced; pull --rebase failed: %s",
			trimGitOutput(pullOutput),
		)
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: rebase sync completed for %s", branch)
	}

	retryOutput, err := runGitCmdCombined(constants.GitPush, constants.GitOrigin, branch)
	if err != nil {
		return fmt.Errorf(
			"push retry after rebase failed: %s",
			trimGitOutput(retryOutput),
		)
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: push retry succeeded for %s", branch)
	}

	return nil
}

func runGitCmdCombined(args ...string) (string, error) {
	cmd := exec.Command(constants.GitBin, args...)
	out, err := cmd.CombinedOutput()

	return string(out), err
}

func isNonFastForwardPushError(output string) bool {
	lower := strings.ToLower(output)

	return strings.Contains(lower, "fetch first") ||
		strings.Contains(lower, "non-fast-forward") ||
		strings.Contains(lower, "failed to push some refs")
}

func formatGitCommandError(output string, err error) error {
	trimmed := trimGitOutput(output)
	if len(trimmed) > 0 {
		return fmt.Errorf("%s", trimmed)
	}

	return err
}

func trimGitOutput(output string) string {
	trimmed := strings.TrimSpace(output)
	if len(trimmed) > 0 {
		return trimmed
	}

	return "unknown git error"
}

func singleLineGitOutput(output string) string {
	return strings.Join(strings.Fields(trimGitOutput(output)), " ")
}

func abortRebaseAfterFailure() {
	_, err := runGitCmdCombined(constants.GitRebase, constants.GitRebaseAbortFlag)
	if err != nil {
		return
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("autocommit: aborted failed rebase")
	}
}
