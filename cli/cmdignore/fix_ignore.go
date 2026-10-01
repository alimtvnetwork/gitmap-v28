package cmdignore

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunFixIgnoreAllFn is a hook for delegating fleet execution.
var RunFixIgnoresAllSSHFn func(args []string) error

// RunFixIgnoreAll fixes gitignore issues across all repositories.
func RunFixIgnoreAll(args []string) error {
	isAutoYes := hasYesFlag(args)
	records := resolveTargetRepos()
	if len(records) == 0 {
		fmt.Println("No repositories found to inspect.")
		return nil
	}

	issues := scanReposForIgnoreIssues(records)
	if len(issues) == 0 {
		fmt.Printf("%s✓ All %d repository .gitignore files are clean and synchronized.%s\n",
			constants.ColorGreen, len(records), constants.ColorReset)
		return nil
	}

	if !isAutoYes && !promptUserConsent(len(issues)) {
		fmt.Println("Aborted by user.")
		return nil
	}

	fixedCount := remediateIssuesList(issues)
	fmt.Printf("\n%s✓ Fixed and sanitized .gitignore across %d repository(ies).%s\n",
		constants.ColorGreen, fixedCount, constants.ColorReset)
	return nil
}

// RunFixIgnoresAllSSH executes fix-ignore-all across local host and remote fleet nodes.
func RunFixIgnoresAllSSH(args []string) error {
	if RunFixIgnoresAllSSHFn != nil {
		return RunFixIgnoresAllSSHFn(args)
	}
	return RunFixIgnoreAll(args)
}

func hasYesFlag(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "-y" || low == "--yes" {
			return true
		}
	}
	return false
}

func resolveTargetRepos() []model.ScanRecord {
	s, err := store.OpenDefault()
	if err != nil {
		return checkCurrentDirectoryOnly()
	}
	defer s.Close()
	records, queryErr := s.ListRepos()
	if queryErr != nil || len(records) == 0 {
		return checkCurrentDirectoryOnly()
	}
	return records
}

func checkCurrentDirectoryOnly() []model.ScanRecord {
	cwd, err := os.Getwd()
	if err != nil || !gitignoreagm.IsGitRepository(cwd) {
		return nil
	}
	return []model.ScanRecord{{AbsolutePath: cwd, RepoName: filepath.Base(cwd)}}
}

func scanReposForIgnoreIssues(records []model.ScanRecord) []IgnoreScanIssue {
	var issues []IgnoreScanIssue
	for _, r := range records {
		issue := inspectSingleRepo(r.AbsolutePath, r.RepoName)
		if issue.HasDuplicate || issue.HasResumeTask || issue.MissingGitmapDir {
			issues = append(issues, issue)
		}
	}
	return issues
}

func inspectSingleRepo(repoDir, repoName string) IgnoreScanIssue {
	issue := IgnoreScanIssue{RepoPath: repoDir, RepoName: repoName}
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, err := os.ReadFile(ignorePath)
	if err != nil {
		issue.MissingGitmapDir = true
		issue.HasResumeTask = gitignoreagm.HasUnignoredResumeTask(repoDir)
		return issue
	}

	_, isChanged := gitignoreagm.DeduplicateAndSanitizeGitignore(string(data))
	issue.HasDuplicate = isChanged
	issue.MissingGitmapDir = !strings.Contains(string(data), ".gitmap/")
	issue.HasResumeTask = gitignoreagm.HasUnignoredResumeTask(repoDir)
	return issue
}

func promptUserConsent(count int) bool {
	fmt.Printf("\n%sFound gitignore issues / duplicates in %d repository(ies).%s\n",
		constants.ColorYellow, count, constants.ColorReset)
	fmt.Print("Sanitize and fix .gitignore across these repositories? [Y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	trimmed := strings.ToLower(strings.TrimSpace(text))
	return trimmed == "" || trimmed == "y" || trimmed == "yes"
}

func remediateIssuesList(issues []IgnoreScanIssue) int {
	fixed := 0
	for _, issue := range issues {
		if remediateSingleRepo(issue.RepoPath) {
			fixed++
		}
	}
	return fixed
}

func remediateSingleRepo(repoDir string) bool {
	_, _ = gitignoreagm.RemediateRepo(repoDir, true)
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, err := os.ReadFile(ignorePath)
	if err != nil {
		return true
	}

	cleaned, isModified := gitignoreagm.DeduplicateAndSanitizeGitignore(string(data))
	if !isModified {
		return true
	}

	_ = os.WriteFile(ignorePath, []byte(cleaned), 0644)
	_ = exec.Command("git", "-C", repoDir, "add", ".gitignore").Run()
	_ = exec.Command("git", "-C", repoDir, "commit", "-m", "chore(git): sanitize .gitignore").Run()
	return true
}
