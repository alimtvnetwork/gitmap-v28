package cmdpull

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func collectAndRemediateIgnoreIssues(handle *IgnoreScanHandle, opts pullOptions) {
	if handle == nil {
		return
	}
	issues := handle.Collect()
	if len(issues) == 0 {
		return
	}
	deduped := DeduplicateIgnoreIssues(issues)
	isAutoYes := opts.yes || opts.autoFix
	handleIgnoreRemediation(deduped, isAutoYes, opts.isJSON)
}

func handleIgnoreRemediation(issues []IgnoreRepoIssue, isAutoYes, isJSON bool) {
	if len(issues) == 0 || isJSON {
		return
	}
	printIgnoreIssuesReport(issues)
	if isAutoYes {
		remediateAllIgnoreIssues(issues)
		return
	}
	if !isInteractiveTerminal() {
		markIssuesSkipped(issues)
		printNonInteractiveIgnoreNotice(len(issues))
		return
	}
	dispatchInteractiveIgnoreRemediation(issues)
}

func printIgnoreIssuesReport(issues []IgnoreRepoIssue) {
	fmt.Printf("\n  %s⚠%s %sDetected .gitignore issues in %d repository(ies):%s\n",
		constants.ColorYellow, constants.ColorReset,
		constants.ColorBold, len(issues), constants.ColorReset)
	for _, issue := range issues {
		printSingleRepoIgnoreIssue(issue)
	}
	fmt.Println()
}

func printSingleRepoIgnoreIssue(issue IgnoreRepoIssue) {
	fmt.Printf("    %s• %s%s\n", constants.ColorCyan, issue.RepoName, constants.ColorReset)
	for _, dup := range issue.DuplicatePatterns {
		fmt.Printf("      %s-%s Duplicate pattern in .gitignore: %s%s%s\n",
			constants.ColorDim, constants.ColorReset, constants.ColorYellow, dup, constants.ColorReset)
	}
	for _, tracked := range issue.TrackedPaths {
		fmt.Printf("      %s-%s Tracked in git index (should be ignored): %s%s%s\n",
			constants.ColorDim, constants.ColorReset, constants.ColorRed, tracked, constants.ColorReset)
	}
}

func isInteractiveTerminal() bool {
	if os.Getenv("CI") != "" || os.Getenv("GITMAP_NON_INTERACTIVE") != "" {
		return false
	}
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func printNonInteractiveIgnoreNotice(count int) {
	fmt.Printf("  %sℹ Detected ignore issues in %d repository(ies). Run 'gitmap fix-ignore-all' to remediate.%s\n\n",
		constants.ColorCyan, count, constants.ColorReset)
}

func dispatchInteractiveIgnoreRemediation(issues []IgnoreRepoIssue) {
	choice := promptIgnoreRemediationChoice()
	if isChoiceAll(choice) {
		remediateAllIgnoreIssues(issues)
		return
	}
	if isChoiceSingle(choice) {
		remediateSingleRepoInteractive(issues)
		return
	}
	markIssuesSkipped(issues)
	fmt.Printf("  %s↷ Skipped ignore resolution.%s\n\n", constants.ColorDim, constants.ColorReset)
}

func markIssuesSkipped(issues []IgnoreRepoIssue) {
	for _, issue := range issues {
		_ = store.RecordIgnoreCheckResult(issue.RepoPath, issue.RepoName, "skipped", 0, 0)
	}
}

func promptIgnoreRemediationChoice() string {
	fmt.Printf("  %s?%s Resolve all at once [y/all], one-by-one [s/single], or skip [n]?: ",
		constants.ColorCyan, constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "n"
	}
	return strings.ToLower(strings.TrimSpace(line))
}

func isChoiceAll(choice string) bool {
	return choice == "y" || choice == "yes" || choice == "all" || choice == "a"
}

func isChoiceSingle(choice string) bool {
	return choice == "s" || choice == "single" || choice == "one" || choice == "1"
}

func remediateSingleRepoInteractive(issues []IgnoreRepoIssue) {
	reader := bufio.NewReader(os.Stdin)
	remediatedCount := 0
	for _, issue := range issues {
		if !promptSingleRepoRemediation(reader, issue.RepoName) {
			_ = store.RecordIgnoreCheckResult(issue.RepoPath, issue.RepoName, "skipped", 0, 0)
			continue
		}
		if remediateSingleRepoIgnore(issue) {
			remediatedCount++
		}
	}
	fmt.Printf("\n  %s✓ Completed remediation across %d repository(ies).%s\n\n",
		constants.ColorGreen, remediatedCount, constants.ColorReset)
}

func promptSingleRepoRemediation(reader *bufio.Reader, repoName string) bool {
	fmt.Printf("  %s?%s [%s] Resolve ignore issues in %s? [Y/n]: ",
		constants.ColorCyan, constants.ColorReset, repoName, repoName)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(line))
	return ans == "" || ans == "y" || ans == "yes"
}

func remediateAllIgnoreIssues(issues []IgnoreRepoIssue) {
	remediatedCount := 0
	for _, issue := range issues {
		if remediateSingleRepoIgnore(issue) {
			remediatedCount++
		}
	}
	fmt.Printf("\n  %s✓ Completed ignore remediation across %d repository(ies).%s\n\n",
		constants.ColorGreen, remediatedCount, constants.ColorReset)
}

func remediateSingleRepoIgnore(issue IgnoreRepoIssue) bool {
	wasUntracked := untrackRepoPaths(issue.RepoPath, issue.TrackedPaths)
	wasSanitized := sanitizeRepoGitignore(issue.RepoPath)
	if wasUntracked || wasSanitized {
		fmt.Printf("  %s✓%s [%s] Resolved ignore issues and sanitized .gitignore\n",
			constants.ColorGreen, constants.ColorReset, issue.RepoName)
		_ = store.RecordIgnoreCheckResult(issue.RepoPath, issue.RepoName, "clean", 1, 0)
		return true
	}
	_ = store.RecordIgnoreCheckResult(issue.RepoPath, issue.RepoName, "skipped", 0, 0)
	return false
}

func untrackRepoPaths(repoDir string, trackedPaths []string) bool {
	if len(trackedPaths) == 0 {
		return false
	}
	args := append([]string{"-C", repoDir, "rm", "--cached", "-r", "-f", "--ignore-unmatch", "--"}, trackedPaths...)
	_, err := gitutil.ExecGitWithTimeout(10*time.Second, repoDir, args...)
	if err != nil {
		return false
	}
	commitArgs := []string{"-C", repoDir, "commit", "-m", "chore(git): untrack ignored files from index"}
	_, _ = gitutil.ExecGitWithTimeout(10*time.Second, repoDir, commitArgs...)
	return true
}

func sanitizeRepoGitignore(repoDir string) bool {
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, _ := os.ReadFile(ignorePath)
	cleaned, isModified := gitignoreagm.DeduplicateAndSanitizeGitignore(string(data))
	if !isModified {
		return false
	}
	if writeErr := os.WriteFile(ignorePath, []byte(cleaned), 0o644); writeErr != nil {
		return false
	}
	return commitSanitizedGitignore(repoDir)
}

func commitSanitizedGitignore(repoDir string) bool {
	addArgs := []string{"-C", repoDir, "add", ".gitignore"}
	_, errAdd := gitutil.ExecGitWithTimeout(10*time.Second, repoDir, addArgs...)
	if errAdd != nil {
		return false
	}
	commitArgs := []string{"-C", repoDir, "commit", "-m", "chore(git): sanitize .gitignore patterns"}
	_, errCommit := gitutil.ExecGitWithTimeout(10*time.Second, repoDir, commitArgs...)
	return errCommit == nil
}
