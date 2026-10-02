package cmdpull

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// PullFailureSummary holds concise failure information for a repository.
type PullFailureSummary struct {
	RepoName  string
	RepoPath  string
	ErrorText string
	ErrorType string
}

// PromptInteractiveBatchFix displays an interactive prompt for remediating failed repositories.
func PromptInteractiveBatchFix(failures []PullFailureSummary) error {
	hasFailures := len(failures) > 0
	if !hasFailures || !isInteractiveTerminal() {
		return nil
	}
	printRemediationPromptMenu()
	choice := readPromptChoice()

	return dispatchRemediationChoice(choice, failures)
}

func printRemediationPromptMenu() {
	fmt.Printf("\n  %s[?]%s %sRemediate failed repositories?%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, constants.ColorReset)
	fmt.Println("    [1/a] Resolve all failed repositories at once")
	fmt.Println("    [2/s] Step through one by one")
	fmt.Println("    [q/n] Skip / Exit")
	fmt.Printf("  %sChoice [1/2/q]:%s ", constants.ColorCyan, constants.ColorReset)
}

func readPromptChoice() string {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "q"
	}

	return strings.ToLower(strings.TrimSpace(line))
}

func dispatchRemediationChoice(choice string, failures []PullFailureSummary) error {
	if isChoiceResolveAll(choice) {
		return RunBatchRemediateAll(failures)
	}
	if isChoiceStepByStep(choice) {
		return RunStepByStepRemediate(failures)
	}

	return nil
}

func isChoiceResolveAll(choice string) bool {
	return choice == "1" || choice == "a" || choice == "all" || choice == "y" || choice == "yes"
}

func isChoiceStepByStep(choice string) bool {
	return choice == "2" || choice == "s" || choice == "step"
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
	fmt.Printf("  %s?%s Remediate this repository? [Y/n/q]: ", constants.ColorCyan, constants.ColorReset)
	line, _ := reader.ReadString('\n')
	ans := strings.ToLower(strings.TrimSpace(line))
	if ans == "q" || ans == "quit" {
		return false, false
	}
	if ans == "n" || ans == "no" {
		return true, false
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
	fmt.Printf("  %s→%s Remediating %s%s%s...\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, f.RepoName, constants.ColorReset)
	out, err := execPullAutoMerge(f.RepoPath)
	if err == nil {
		fmt.Printf("  %s✓%s [%s] Successfully updated.\n",
			constants.ColorGreen, constants.ColorReset, f.RepoName)
		return true
	}

	return handleRemediationFailure(f.RepoPath, f.RepoName, string(out))
}

func execPullAutoMerge(repoPath string) ([]byte, error) {
	cmd := exec.Command("git", "-C", repoPath, "pull", "--progress", "--no-rebase", "--no-edit", "--autostash")
	cmd.Env = gitutil.BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)

	return cmd.CombinedOutput()
}

func handleRemediationFailure(repoPath, repoName, out string) bool {
	if strings.Contains(strings.ToLower(out), "conflict") {
		abortCmd := exec.Command("git", "-C", repoPath, "merge", "--abort")
		abortCmd.Env = gitutil.BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)
		_ = abortCmd.Run()
		fmt.Printf("  %s✗%s [%s] Merge conflict detected (merge aborted).\n",
			constants.ColorRed, constants.ColorReset, repoName)

		return false
	}
	fmt.Printf("  %s✗%s [%s] Remediation failed.\n",
		constants.ColorRed, constants.ColorReset, repoName)

	return false
}

// RunBatchPullFix loads last recorded pull failures from SQLite and runs remediation.
func RunBatchPullFix(args []string) error {
	failures, err := loadRecentPullFailuresFromDB()
	if err != nil {
		return err
	}
	hasFailures := len(failures) > 0
	if !hasFailures {
		fmt.Printf("  %sℹ No recent pull failures found to remediate.%s\n",
			constants.ColorCyan, constants.ColorReset)
		return nil
	}
	fmt.Printf("  %sFound %d failed repository(ies) from recent runs.%s\n",
		constants.ColorCyan, len(failures), constants.ColorReset)

	return dispatchBatchFixArgs(args, failures)
}

func loadRecentPullFailuresFromDB() ([]PullFailureSummary, *apperror.AppError) {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return nil, apperror.WrapSimple(err, "open_pull_db")
	}
	defer db.Close()
	records, err := db.QueryAllLatestPullErrors(50)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query_pull_errors")
	}

	return convertRecordsToFailures(records), nil
}

func convertRecordsToFailures(records []store.PullErrorRecord) []PullFailureSummary {
	var failures []PullFailureSummary
	seen := make(map[string]bool, len(records))
	for _, r := range records {
		hasPath := r.RepoPath != ""
		if hasPath && !seen[r.RepoPath] {
			seen[r.RepoPath] = true
			failures = append(failures, PullFailureSummary{
				RepoName:  r.RepoSlug,
				RepoPath:  r.RepoPath,
				ErrorText: r.ErrorText,
				ErrorType: r.ErrorType,
			})
		}
	}

	return failures
}

func dispatchBatchFixArgs(args []string, failures []PullFailureSummary) error {
	if hasBatchFlag(args, "--all", "-a", "--yes", "-y") {
		return RunBatchRemediateAll(failures)
	}
	if hasBatchFlag(args, "--step", "-s", "--one-by-one") {
		return RunStepByStepRemediate(failures)
	}

	return PromptInteractiveBatchFix(failures)
}

func hasBatchFlag(args []string, targets ...string) bool {
	for _, arg := range args {
		low := strings.ToLower(strings.TrimSpace(arg))
		for _, target := range targets {
			if low == target {
				return true
			}
		}
	}

	return false
}

// ExtractPullFailures collects failures from slice of repo states.
func ExtractPullFailures(states []*PullRepoState) []PullFailureSummary {
	var failures []PullFailureSummary
	for _, s := range states {
		if s.ErrorMsg != "" || s.Step == PullStepTypeError {
			failures = append(failures, PullFailureSummary{
				RepoName:  s.RepoName,
				RepoPath:  s.RepoPath,
				ErrorText: s.ErrorMsg,
				ErrorType: string(s.Step),
			})
		}
	}

	return failures
}
