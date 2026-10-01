package cmdignore

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunFixIgnoresAllSSHFn is a hook for delegating fleet execution.
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
		printIssuesSummary(issues)
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

// RunFixIgnoresAllSSH executes fleet-wide fix-ignore following the GitMap PAS Formula.
func RunFixIgnoresAllSSH(args []string) *apperror.AppError {
	hasHook := RunFixIgnoresAllSSHFn != nil
	if !hasHook {
		return RunFixIgnoreAll(args)
	}
	queueId, tDB := enqueueIgnoreTaskQueue("fix-ignore-all-ssh", "fleet")
	defer closeIgnoreTaskDB(tDB)
	updateIgnoreTaskQueue(tDB, queueId, "running")
	return executeFleetSSHWithTask(tDB, queueId, args)
}

func executeFleetSSHWithTask(tDB *store.TasksSplitDB, queueId string, args []string) *apperror.AppError {
	errFleet := RunFixIgnoresAllSSHFn(args)
	hasErr := errFleet != nil
	if hasErr {
		store.LogInternalError("FIX_IGNORE_SSH", "SSH_DELEGATION_ERROR", errFleet.Error(), "", "")
		updateIgnoreTaskQueue(tDB, queueId, "failed")
		return errFleet
	}
	updateIgnoreTaskQueue(tDB, queueId, "completed")
	return nil
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
	hasMissing := issue.MissingGitmapDir
	if hasMissing {
		return true
	}
	hasTracked := len(issue.TrackedResume) > 0
	if hasTracked {
		return true
	}
	return false
}

func inspectSingleRepo(repoDir, repoName string) IgnoreScanIssue {
	issue := IgnoreScanIssue{RepoPath: repoDir, RepoName: repoName}
	trackedFiles := findTrackedIgnoredFiles(repoDir)
	issue.TrackedResume = trackedFiles
	issue.HasResumeTask = len(trackedFiles) > 0
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, err := os.ReadFile(ignorePath)
	hasErr := err != nil
	if hasErr {
		issue.MissingGitmapDir = true
		return issue
	}
	return analyzeGitignoreData(issue, repoDir, string(data))
}

func analyzeGitignoreData(issue IgnoreScanIssue, repoDir, data string) IgnoreScanIssue {
	_, dupCount := deduplicateIgnoreContent(data)
	issue.DuplicateCount = dupCount
	issue.HasDuplicate = dupCount > 0
	hasGitmapDir := strings.Contains(data, ".gitmap/")
	hasMissing := !hasGitmapDir
	issue.MissingGitmapDir = hasMissing
	return issue
}

func findTrackedIgnoredFiles(repoDir string) []string {
	isGit := gitignoreagm.IsGitRepository(repoDir)
	if !isGit {
		return nil
	}
	out, err := exec.Command("git", "-C", repoDir, "ls-files", "-c", "-i", "--exclude-standard").Output()
	hasErr := err != nil
	if hasErr {
		return nil
	}
	return parseGitLsFilesOutput(string(out))
}

func parseGitLsFilesOutput(raw string) []string {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	var result []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		hasLen := len(trimmed) > 0
		if hasLen {
			result = append(result, trimmed)
		}
	}
	return result
}

func deduplicateIgnoreContent(content string) (string, int) {
	lines := strings.Split(content, "\n")
	seen := make(map[string]bool, len(lines))
	var cleaned []string
	dupCount := 0
	for _, raw := range lines {
		dupCount = processIgnoreLine(raw, seen, &cleaned, dupCount)
	}
	finalContent := assembleCleanedGitignore(cleaned, seen)
	return finalContent, dupCount
}

func processIgnoreLine(raw string, seen map[string]bool, cleaned *[]string, dupCount int) int {
	trimmed := strings.TrimSpace(raw)
	isEmpty := trimmed == ""
	if isEmpty {
		*cleaned = append(*cleaned, raw)
		return dupCount
	}
	isComment := strings.HasPrefix(trimmed, "#")
	if isComment {
		*cleaned = append(*cleaned, raw)
		return dupCount
	}
	norm := strings.TrimPrefix(trimmed, "/")
	hasSeen := seen[norm]
	if hasSeen {
		return dupCount + 1
	}
	seen[norm] = true
	*cleaned = append(*cleaned, raw)
	return dupCount
}

func assembleCleanedGitignore(cleaned []string, seen map[string]bool) string {
	hasGitmap := seen[".gitmap/"]
	if !hasGitmap {
		cleaned = append(cleaned, "", "# GitMap & Task Persistence", ".gitmap/", ".gitmap/backup/")
		seen[".gitmap/"] = true
		seen[".gitmap/backup/"] = true
	}
	text := strings.Join(cleaned, "\n")
	return strings.TrimRight(text, "\r\n") + "\n"
}

func remediateIssuesList(issues []IgnoreScanIssue) int {
	fixed := 0
	for _, issue := range issues {
		isFixed := remediateSingleRepo(issue)
		if isFixed {
			fixed++
		}
	}
	return fixed
}

func remediateSingleRepo(issue IgnoreScanIssue) bool {
	hasTracked := len(issue.TrackedResume) > 0
	if hasTracked {
		untrackCachedFiles(issue.RepoPath, issue.TrackedResume)
	}
	ignorePath := filepath.Join(issue.RepoPath, ".gitignore")
	data, err := os.ReadFile(ignorePath)
	hasErr := err != nil
	if hasErr {
		return writeDefaultGitignore(issue.RepoPath, ignorePath)
	}
	return sanitizeAndCommitRepo(issue.RepoPath, ignorePath, string(data))
}

func untrackCachedFiles(repoDir string, files []string) {
	args := append([]string{"-C", repoDir, "rm", "--cached", "-f", "--ignore-unmatch", "--"}, files...)
	_ = exec.Command("git", args...).Run()
	_ = exec.Command("git", "-C", repoDir, "commit", "-m", "chore(git): untrack ignored files from repository").Run()
}

func writeDefaultGitignore(repoDir, ignorePath string) bool {
	content := "# GitMap & Task Persistence\n.gitmap/\n.gitmap/backup/\n"
	_ = os.WriteFile(ignorePath, []byte(content), 0644)
	_ = exec.Command("git", "-C", repoDir, "add", ".gitignore").Run()
	_ = exec.Command("git", "-C", repoDir, "commit", "-m", "chore(git): add default .gitignore").Run()
	return true
}

func sanitizeAndCommitRepo(repoDir, ignorePath, data string) bool {
	cleaned, dupCount := deduplicateIgnoreContent(data)
	hasGitmapDir := strings.Contains(data, ".gitmap/")
	hasDuplicates := dupCount > 0
	isMissingGitmap := !hasGitmapDir
	hasChanges := hasDuplicates || isMissingGitmap
	if hasChanges {
		_ = os.WriteFile(ignorePath, []byte(cleaned), 0644)
		_ = exec.Command("git", "-C", repoDir, "add", ".gitignore").Run()
		_ = exec.Command("git", "-C", repoDir, "commit", "-m", "chore(git): sanitize .gitignore").Run()
	}
	return true
}

func promptUserConsent(count int) bool {
	fmt.Printf("\n%sFound gitignore issues in %d repository(ies).%s\n",
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

func printIssuesSummary(issues []IgnoreScanIssue) {
	fmt.Printf("\n%s=== GitIgnore Issues By Repository ===%s\n", constants.ColorCyan, constants.ColorReset)
	for _, issue := range issues {
		printSingleRepoIssueSummary(issue)
	}
}

func printSingleRepoIssueSummary(issue IgnoreScanIssue) {
	fmt.Printf("  * %s%s%s\n", constants.ColorBold, issue.RepoName, constants.ColorReset)
	if issue.HasDuplicate {
		fmt.Printf("      - Contains %d duplicate rule(s)\n", issue.DuplicateCount)
	}
	hasMissing := issue.MissingGitmapDir
	if hasMissing {
		fmt.Printf("      - Missing .gitmap/ ignore rule\n")
	}
	printTrackedFilesSummary(issue.TrackedResume)
}

func printTrackedFilesSummary(tracked []string) {
	hasTracked := len(tracked) > 0
	if !hasTracked {
		return
	}
	fmt.Printf("      - Contains %d tracked file(s) matching .gitignore rules:\n", len(tracked))
	for _, f := range tracked {
		fmt.Printf("          * %s\n", f)
	}
}

func enqueueIgnoreTaskQueue(action, target string) (string, *store.TasksSplitDB) {
	queueId := fmt.Sprintf("%s-%d", action, time.Now().UnixNano())
	tasksDB, err := store.OpenTasksRootSplitDB()
	hasErr := err != nil
	if hasErr {
		return "", nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(tasksDB.Conn(), "INSERT INTO TaskQueue (QueueId, Section, Action, Target, Status, CreatedAt, UpdatedAt) VALUES (?, 'ignore', ?, ?, 'pending', ?, ?)", queueId, action, target, now, now)
	return queueId, tasksDB
}

func updateIgnoreTaskQueue(db *store.TasksSplitDB, queueId, status string) {
	hasDB := db != nil
	if !hasDB {
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(db.Conn(), "UPDATE TaskQueue SET Status = ?, UpdatedAt = ? WHERE QueueId = ?", status, now, queueId)
}

func closeIgnoreTaskDB(db *store.TasksSplitDB) {
	hasDB := db != nil
	if hasDB {
		_ = db.Close()
	}
}
