package cmdpull

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// PullFailureSummary holds concise failure information for a repository.
type PullFailureSummary struct {
	RepoID    int64
	RepoName  string
	RepoPath  string
	RemoteURL string
	ErrorText string
	ErrorType string
}

// PromptInteractiveBatchFix displays an interactive prompt for remediating failed repositories.
func PromptInteractiveBatchFix(failures []PullFailureSummary) error {
	hasFailures := len(failures) > 0
	if !hasFailures || !isInteractiveTerminal() {
		return nil
	}
	for {
		printRemediationPromptMenu()
		choice := readPromptChoice()
		if isChoiceHelp(choice) {
			renderRemediationHelpGuide()
			continue
		}

		return dispatchRemediationChoice(choice, failures)
	}
}

// RunBatchRemediateAll resolves all failed repositories.
func RunBatchRemediateAll(failures []PullFailureSummary) error {
	fmt.Printf("\n  %sResolving %d failed repositories...%s\n\n",
		constants.ColorBold, len(failures), constants.ColorReset)
	successCount := 0
	for _, f := range failures {
		if remediateSingleRepo(f) {
			successCount++
		}
	}
	printRemediationBatchSummary(successCount, len(failures))

	return nil
}

func printRemediationBatchSummary(successCount, total int) {
	fmt.Printf("\n  %s✓ Batch remediation complete: %d/%d resolved.%s\n\n",
		constants.ColorGreen, successCount, total, constants.ColorReset)
}

// RunStepByStepRemediate presents each failed repository to the user one by one.
func RunStepByStepRemediate(failures []PullFailureSummary) error {
	reader := bufio.NewReader(os.Stdin)
	successCount := 0
	for i, f := range failures {
		shouldContinue, isFixed := promptAndRemediateSingle(reader, i, len(failures), f)
		if isFixed {
			successCount++
		}
		if !shouldContinue {
			break
		}
	}
	printStepRemediationSummary(successCount, len(failures))

	return nil
}

func promptAndRemediateSingle(reader *bufio.Reader, idx, total int, f PullFailureSummary) (bool, bool) {
	printSingleRepoFailureHeader(idx+1, total, f)
	if cloner.IsMissingRepo(f.RepoPath) {
		return PromptMissingRepoAction(reader, f)
	}

	fmt.Printf("  %s?%s Remediate this repository? [Y/r/n/q] (Y=auto-merge, r=rebase, n=skip, q=quit): ", constants.ColorCyan, constants.ColorReset)
	line, _ := reader.ReadString('\n')
	ans := strings.ToLower(strings.TrimSpace(line))
	if ans == "q" || ans == "quit" {
		return false, false
	}
	if ans == "n" || ans == "no" {
		return true, false
	}
	if ans == "r" || ans == "rebase" {
		return true, remediateSingleRepoRebase(f)
	}

	return true, remediateSingleRepo(f)
}

func printSingleRepoFailureHeader(idx, total int, f PullFailureSummary) {
	fmt.Printf("\n  [%d/%d] Repository: %s%s%s\n", idx, total, constants.ColorBold, f.RepoName, constants.ColorReset)
	fmt.Printf("    Path:  %s\n", f.RepoPath)
	fmt.Printf("    Error: %s\n", f.ErrorText)
}

func printStepRemediationSummary(successCount, total int) {
	fmt.Printf("\n  %s✓ Step-by-step remediation complete: %d/%d resolved.%s\n\n",
		constants.ColorGreen, successCount, total, constants.ColorReset)
}

func remediateSingleRepo(f PullFailureSummary) bool {
	if f.RepoPath == "" {
		return false
	}

	if cloner.IsMissingRepo(f.RepoPath) {
		return RemediateMissingRepo(f)
	}

	fmt.Printf("  • %s%s%s\n", constants.ColorBold, f.RepoName, constants.ColorReset)
	applyPreemptiveCredentialFix(f)
	fmt.Printf("    %s Attempting auto-merge...\n", constants.TreeBranch)
	out, err := execPullAutoMerge(f.RepoPath)
	if err == nil {
		handleSuccessfulRemediation(f)
		return true
	}

	if !isCredentialOrHTTPSFailure(string(out)) {
		return handleRemediationFailure(f.RepoPath, f.RepoName, string(out))
	}

	isSuccess, outRetry := tryRetryAfterCredentialRemediation(f)
	if isSuccess {
		return true
	}

	return handleRemediationFailure(f.RepoPath, f.RepoName, string(outRetry))
}

func applyPreemptiveCredentialFix(f PullFailureSummary) {
	if !isCredentialOrHTTPSFailure(f.ErrorText) && !isSameOwnerHTTPSRepo(f.RepoPath) {
		return
	}

	errRem := remediateCredentialAndConvertToSSH(f.RepoPath, f.RepoName)
	if errRem != nil {
		fmt.Printf("    %s %s⚠%s Credential remediation notice: %v\n",
			constants.TreeBranch, constants.ColorYellow, constants.ColorReset, errRem)
	}
}

func handleSuccessfulRemediation(f PullFailureSummary) {
	errClear := clearResolvedPullError(f.RepoName, f.RepoPath)
	if errClear != nil {
		fmt.Printf("    %s %s⚠%s Warning clearing pull error: %v\n",
			constants.TreeBranch, constants.ColorYellow, constants.ColorReset, errClear)
	}

	fmt.Printf("    %s %s✓%s Successfully updated.\n",
		constants.TreeCorner, constants.ColorGreen, constants.ColorReset)
}

func tryRetryAfterCredentialRemediation(f PullFailureSummary) (bool, []byte) {
	errRem := remediateCredentialAndConvertToSSH(f.RepoPath, f.RepoName)
	if errRem != nil {
		fmt.Printf("    %s %s⚠%s SSH credential retry notice: %v\n",
			constants.TreeBranch, constants.ColorYellow, constants.ColorReset, errRem)
	}

	outRetry, errRetry := execPullAutoMerge(f.RepoPath)
	if errRetry == nil {
		handleSuccessfulRemediationSSH(f)
		return true, nil
	}

	return false, outRetry
}

func handleSuccessfulRemediationSSH(f PullFailureSummary) {
	errClear := clearResolvedPullError(f.RepoName, f.RepoPath)
	if errClear != nil {
		fmt.Printf("    %s %s⚠%s Warning clearing pull error: %v\n",
			constants.TreeBranch, constants.ColorYellow, constants.ColorReset, errClear)
	}

	fmt.Printf("    %s %s✓%s Successfully updated via SSH.\n",
		constants.TreeCorner, constants.ColorGreen, constants.ColorReset)
}

func isCredentialOrHTTPSFailure(msg string) bool {
	return isWincredmanFailure(msg) || isCacheFailure(msg) || isAuthFailure(msg) ||
		strings.Contains(strings.ToLower(msg), "credential")
}

func execPullAutoMerge(repoPath string) ([]byte, error) {
	return gitutil.RunGitCombinedInDir(repoPath, "pull", "--progress", "--no-rebase", "--no-edit", "--autostash")
}

func handleRemediationFailure(repoPath, repoName, out string) bool {
	relPath := formatRelativeOrCleanPath(repoPath)
	if relPath == "" || relPath == "." {
		relPath = repoName
	}

	isConflict := cloner.IsMergeConflictOutput(out) || strings.Contains(strings.ToLower(out), "conflict")
	if isConflict {
		cloner.SafeAbortMerge(repoPath)
		fmt.Printf("    %s %s✗%s Merge conflict detected (merge safely aborted to protect working tree).\n",
			constants.TreeCorner, constants.ColorRed, constants.ColorReset)
		fmt.Printf("        %s[hint]%s Resolve manually: git -C %s status\n",
			constants.ColorYellow, constants.ColorReset, relPath)

		return false
	}

	isDiverged := cloner.IsDivergedOutput(out)
	if isDiverged {
		cloner.SafeAbortMerge(repoPath)
		fmt.Printf("    %s %s✗%s Diverged branch detected (merge safely aborted to protect working tree).\n",
			constants.TreeCorner, constants.ColorRed, constants.ColorReset)
		fmt.Printf("        %s[hint]%s Resolve manually: git -C %s pull --rebase\n",
			constants.ColorYellow, constants.ColorReset, relPath)

		return false
	}

	cloner.SafeAbortMerge(repoPath)
	fmt.Printf("    %s %s✗%s Remediation failed.\n",
		constants.TreeCorner, constants.ColorRed, constants.ColorReset)
	fmt.Printf("        %s[hint]%s Resolve manually: git -C %s status\n",
		constants.ColorYellow, constants.ColorReset, relPath)

	return false
}
