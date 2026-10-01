package cmdignore

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunFixIgnoreAllFn is a hook for delegating fleet execution.
var RunFixIgnoresAllSSHFn func(args []string) *apperror.AppError

// RunFixIgnoreAll fixes gitignore issues across all repositories.
func RunFixIgnoreAll(args []string) *apperror.AppError {
	isAutoYes := isAutoYesArg(args)
	records := resolveTargetRepos()
	isZero := len(records) == 0
	if isZero {
		fmt.Println("No repositories found to inspect.")
		return nil
	}
	return runFixIgnoreAllWithRecords(records, isAutoYes)
}

func runFixIgnoreAllWithRecords(records []model.ScanRecord, isAutoYes bool) *apperror.AppError {
	issues := scanReposForIgnoreIssues(records)
	isClean := len(issues) == 0
	if isClean {
		fmt.Printf("%s✓ All %d repository .gitignore files are clean and synchronized.%s\n",
			constants.ColorGreen, len(records), constants.ColorReset)
		return nil
	}
	return proceedWithIssues(issues, isAutoYes)
}

func proceedWithIssues(issues []IgnoreScanIssue, isAutoYes bool) *apperror.AppError {
	hasConsent := isAutoYes
	if !hasConsent {
		hasConsent = promptUserConsent(len(issues))
	}
	if !hasConsent {
		fmt.Println("Aborted by user.")
		return nil
	}
	fixedCount := remediateIssuesList(issues)
	fmt.Printf("\n%s✓ Fixed and sanitized .gitignore across %d repository(ies).%s\n",
		constants.ColorGreen, fixedCount, constants.ColorReset)
	return nil
}

func RunFixIgnoresAllSSH(args []string) *apperror.AppError {
	hasHook := RunFixIgnoresAllSSHFn != nil
	if hasHook {
		return RunFixIgnoresAllSSHFn(args)
	}
	return RunFixIgnoreAll(args)
}

func isAutoYesArg(args []string) bool {
	for _, a := range args {
		isYes := isYesFlagString(a)
		if isYes {
			return true
		}
	}
	return false
}

func isYesFlagString(a string) bool {
	low := strings.ToLower(a)
	isDashY := low == "-y"
	isDashYes := low == "--yes"
	if isDashY {
		return true
	}
	return isDashYes
}

func resolveTargetRepos() []model.ScanRecord {
	s, err := store.OpenDefault()
	hasErr := err != nil
	if hasErr {
		return checkCurrentDirectoryOnly()
	}
	defer s.Close()
	return getRecordsFromStore(s)
}

func getRecordsFromStore(s *store.DB) []model.ScanRecord {
	records, queryErr := s.ListRepos()
	hasErr := queryErr != nil
	if hasErr {
		return checkCurrentDirectoryOnly()
	}
	isZero := len(records) == 0
	if isZero {
		return checkCurrentDirectoryOnly()
	}
	return records
}

func checkCurrentDirectoryOnly() []model.ScanRecord {
	cwd, err := os.Getwd()
	hasErr := err != nil
	if hasErr {
		return nil
	}
	isGit := gitignoreagm.IsGitRepository(cwd)
	if isGit {
		return []model.ScanRecord{{AbsolutePath: cwd, RepoName: filepath.Base(cwd)}}
	}
	return nil
}

func scanReposForIgnoreIssues(records []model.ScanRecord) []IgnoreScanIssue {
	var issues []IgnoreScanIssue
	for _, r := range records {
		issue := inspectSingleRepo(r.AbsolutePath, r.RepoName)
		hasIssue := isIssueDetected(issue)
		if hasIssue {
			issues = append(issues, issue)
		}
	}
	return issues
}

func isIssueDetected(issue IgnoreScanIssue) bool {
	if issue.HasDuplicate {
		return true
	}
	if issue.HasResumeTask {
		return true
	}
	// Note: using issue.MissingGitmapDir directly
	if issue.MissingGitmapDir {
		return true
	}
	return false
}

func inspectSingleRepo(repoDir, repoName string) IgnoreScanIssue {
	issue := IgnoreScanIssue{RepoPath: repoDir, RepoName: repoName}
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, err := os.ReadFile(ignorePath)
	hasErr := err != nil
	if hasErr {
		issue.MissingGitmapDir = true
		issue.HasResumeTask = gitignoreagm.HasUnignoredResumeTask(repoDir)
		return issue
	}
	return analyzeGitignoreData(issue, repoDir, string(data))
}

func analyzeGitignoreData(issue IgnoreScanIssue, repoDir, data string) IgnoreScanIssue {
	_, isChanged := gitignoreagm.DeduplicateAndSanitizeGitignore(data)
	issue.HasDuplicate = isChanged
	hasGitmapDir := strings.Contains(data, ".gitmap/")
	isMissing := !hasGitmapDir
	issue.MissingGitmapDir = isMissing
	issue.HasResumeTask = gitignoreagm.HasUnignoredResumeTask(repoDir)
	return issue
}

func promptUserConsent(count int) bool {
	fmt.Printf("\n%sFound gitignore issues / duplicates in %d repository(ies).%s\n",
		constants.ColorYellow, count, constants.ColorReset)
	fmt.Print("Sanitize and fix .gitignore across these repositories? [Y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return isConsentGranted(text)
}

func isConsentGranted(text string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(text))
	isEmpty := trimmed == ""
	isY := trimmed == "y"
	isYes := trimmed == "yes"
	if isEmpty {
		return true
	}
	if isY {
		return true
	}
	return isYes
}

func remediateIssuesList(issues []IgnoreScanIssue) int {
	fixed := 0
	for _, issue := range issues {
		isFixed := remediateSingleRepo(issue.RepoPath)
		if isFixed {
			fixed++
		}
	}
	return fixed
}

func remediateSingleRepo(repoDir string) bool {
	_, _ = gitignoreagm.RemediateRepo(repoDir, true)
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, err := os.ReadFile(ignorePath)
	hasErr := err != nil
	if hasErr {
		return true
	}
	return sanitizeAndCommit(repoDir, ignorePath, string(data))
}

func sanitizeAndCommit(repoDir, ignorePath, data string) bool {
	cleaned, isModified := gitignoreagm.DeduplicateAndSanitizeGitignore(data)
	if isModified {
		_ = os.WriteFile(ignorePath, []byte(cleaned), 0644)
		_ = exec.Command("git", "-C", repoDir, "add", ".gitignore").Run()
		_ = exec.Command("git", "-C", repoDir, "commit", "-m", "chore(git): sanitize .gitignore").Run()
	}
	return true
}
