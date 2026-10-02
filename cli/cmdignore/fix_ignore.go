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
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunFixIgnoresAllSSHFn is a hook for delegating fleet execution.
var RunFixIgnoresAllSSHFn func(args []string) *apperror.AppError

// RunFixIgnoreAll fixes gitignore issues across all repositories.
func RunFixIgnoreAll(args []string) *apperror.AppError {
	queueId, tDB := enqueueIgnoreTaskQueue("fix-ignore-all", "local")
	defer closeIgnoreTaskDB(tDB)
	updateIgnoreTaskQueue(tDB, queueId, "running")

	err := executeFixIgnoreAllInternal(args)
	hasErr := err != nil
	if hasErr {
		store.LogInternalError("FIX_IGNORE", "LOCAL_FIX_ERROR", err.Error(), "", "")
		updateIgnoreTaskQueue(tDB, queueId, "failed")
		return err
	}
	updateIgnoreTaskQueue(tDB, queueId, "completed")
	return nil
}

func executeFixIgnoreAllInternal(args []string) *apperror.AppError {
	isAutoYes := isAutoYesArg(args)
	hasForce := isForceArg(args)
	records := resolveTargetRepos()
	isZero := len(records) == 0
	if isZero {
		fmt.Println("No repositories found to inspect.")
		return nil
	}
	return runFixIgnoreAllWithRecords(records, isAutoYes, hasForce)
}

func runFixIgnoreAllWithRecords(records []model.ScanRecord, isAutoYes, hasForce bool) *apperror.AppError {
	targetRecords := resolveAuditRecords(records, hasForce)
	if len(targetRecords) == 0 {
		printAllCachedCleanNotice(len(records))
		return nil
	}
	issues := scanReposForIgnoreIssues(targetRecords)
	if len(issues) == 0 {
		printAllSynchronizedNotice(len(records))
		return nil
	}
	return proceedWithIssues(issues, isAutoYes)
}

func resolveAuditRecords(records []model.ScanRecord, hasForce bool) []model.ScanRecord {
	if hasForce {
		return records
	}
	ttl := resolveIgnoreAuditTTL()
	cold, err := FilterReposNeedingCheck(records, ttl)
	if err != nil {
		return records
	}
	return cold
}

func resolveIgnoreAuditTTL() time.Duration {
	s, err := store.OpenDefault()
	if err != nil {
		return 24 * time.Hour
	}
	defer s.Close()
	return config.GetGitIgnoreTTL(s)
}

func printAllCachedCleanNotice(total int) {
	fmt.Printf("%s✓ All %d repository .gitignore files are clean and synchronized (cached).%s\n",
		constants.ColorGreen, total, constants.ColorReset)
	fmt.Printf("  %s(Audit cached within TTL. Pass '--force' / '-f' to re-verify all repositories.)%s\n",
		constants.ColorDim, constants.ColorReset)
}

func printAllSynchronizedNotice(total int) {
	fmt.Printf("%s✓ All %d repository .gitignore files are clean and synchronized.%s\n",
		constants.ColorGreen, total, constants.ColorReset)
}

func proceedWithIssues(issues []IgnoreScanIssue, isAutoYes bool) *apperror.AppError {
	if isAutoYes {
		fixedCount := remediateIssuesList(issues)
		fmt.Printf("\n%s✓ Fixed and sanitized .gitignore across %d repository(ies).%s\n",
			constants.ColorGreen, fixedCount, constants.ColorReset)
		return nil
	}

	printIssuesSummary(issues)
	choice := promptFixChoice(len(issues))
	return handleFixChoice(choice, issues)
}

func handleFixChoice(choice string, issues []IgnoreScanIssue) *apperror.AppError {
	switch choice {
	case "all":
		fixedCount := remediateIssuesList(issues)
		fmt.Printf("\n%s✓ Fixed and sanitized .gitignore across %d repository(ies).%s\n",
			constants.ColorGreen, fixedCount, constants.ColorReset)
		return nil
	case "single":
		fixedCount := remediateSingleRepoSessions(issues)
		fmt.Printf("\n%s✓ Completed single-repo sessions: fixed %d repository(ies).%s\n",
			constants.ColorGreen, fixedCount, constants.ColorReset)
		return nil
	default:
		for _, issue := range issues {
			_ = RecordRepoCheckResult(issue.RepoPath, issue.RepoName, "skipped", 0, 0)
		}
		fmt.Println("Aborted by user.")
		return nil
	}
}

func promptFixChoice(count int) string {
	fmt.Printf("\n%sFound gitignore issues in %d repository(ies).%s\n",
		constants.ColorYellow, count, constants.ColorReset)
	fmt.Println("How would you like to proceed?")
	fmt.Println("  [1/a] Resolve all at once (recommended)")
	fmt.Println("  [2/s] Step through single-repo sessions")
	fmt.Println("  [q/n] Cancel / Abort")
	fmt.Print("Choose option [1/2/q]: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return parseFixChoice(text)
}

func parseFixChoice(input string) string {
	trimmed := strings.ToLower(strings.TrimSpace(input))
	isAll := trimmed == "" || trimmed == "1" || trimmed == "a" || trimmed == "all" || trimmed == "y" || trimmed == "yes"
	if isAll {
		return "all"
	}
	isSingle := trimmed == "2" || trimmed == "s" || trimmed == "single"
	if isSingle {
		return "single"
	}
	return "cancel"
}

func remediateSingleRepoSessions(issues []IgnoreScanIssue) int {
	reader := bufio.NewReader(os.Stdin)
	fixed := 0
	for idx, issue := range issues {
		fmt.Printf("\n[%d/%d] %s%s%s (%s)\n",
			idx+1, len(issues), constants.ColorBold, issue.RepoName, constants.ColorReset, issue.RepoPath)
		printSingleRepoIssueSummary(issue)
		fmt.Print("Apply fix to this repository? [Y/n/q]: ")
		ans, _ := reader.ReadString('\n')
		trimmed := strings.ToLower(strings.TrimSpace(ans))
		isQuit := trimmed == "q" || trimmed == "quit"
		if isQuit {
			fmt.Println("Exiting single-repo sessions.")
			break
		}
		isSkip := trimmed == "n" || trimmed == "no"
		if isSkip {
			_ = RecordRepoCheckResult(issue.RepoPath, issue.RepoName, "skipped", 0, 0)
			fmt.Printf("Skipped %s.\n", issue.RepoName)
			continue
		}
		isFixed := remediateSingleRepo(issue)
		if isFixed {
			fixed++
			fmt.Printf("%s✓ Remediated %s%s\n", constants.ColorGreen, issue.RepoName, constants.ColorReset)
		}
	}
	return fixed
}

// RunFixIgnoresAllSSH executes fleet-wide fix-ignore following the GitMap PAS Formula.
func RunFixIgnoresAllSSH(args []string) *apperror.AppError {
	queueId, tDB := enqueueIgnoreTaskQueue("fix-ignore-all-ssh", "fleet")
	defer closeIgnoreTaskDB(tDB)
	updateIgnoreTaskQueue(tDB, queueId, "running")

	hasHook := RunFixIgnoresAllSSHFn != nil
	if hasHook {
		return executeFleetSSHWithTask(tDB, queueId, args)
	}
	return executeLocalFallbackWithTask(tDB, queueId, args)
}

func executeLocalFallbackWithTask(tDB *store.TasksSplitDB, queueId string, args []string) *apperror.AppError {
	errLocal := RunFixIgnoreAll(args)
	hasErr := errLocal != nil
	if hasErr {
		store.LogInternalError("FIX_IGNORE_SSH", "LOCAL_FALLBACK_ERROR", errLocal.Error(), "", "")
		updateIgnoreTaskQueue(tDB, queueId, "failed")
		return errLocal
	}
	updateIgnoreTaskQueue(tDB, queueId, "completed")
	return nil
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

func isForceArg(args []string) bool {
	for _, a := range args {
		isForce := isForceFlagString(a)
		if isForce {
			return true
		}
	}
	return false
}

func isForceFlagString(a string) bool {
	low := strings.ToLower(a)
	isDashF := low == "-f"
	isDashForce := low == "--force"
	if isDashF {
		return true
	}
	return isDashForce
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
		} else {
			_ = RecordRepoCheckResult(r.AbsolutePath, r.RepoName, "clean", 0, 0)
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
	start := time.Now()
	issue := IgnoreScanIssue{RepoPath: repoDir, RepoName: repoName}
	trackedFiles := findTrackedIgnoredFiles(repoDir)
	issue.TrackedResume = trackedFiles
	issue.HasResumeTask = len(trackedFiles) > 0
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, err := os.ReadFile(ignorePath)
	if err != nil {
		issue.MissingGitmapDir = true
		recordInspectedCache(repoDir, repoName, true, time.Since(start))
		return issue
	}
	analyzed := analyzeGitignoreData(issue, repoDir, string(data))
	recordInspectedCache(repoDir, repoName, isIssueDetected(analyzed), time.Since(start))
	return analyzed
}

func recordInspectedCache(repoDir, repoName string, hasIssue bool, dur time.Duration) {
	status := "clean"
	if hasIssue {
		status = "has_issues"
	}
	_ = RecordRepoCheckResult(repoDir, repoName, status, 0, dur)
}

func analyzeGitignoreData(issue IgnoreScanIssue, repoDir, data string) IgnoreScanIssue {
	_, dupCount := deduplicateIgnoreContent(data)
	issue.DuplicateCount = dupCount
	issue.HasDuplicate = dupCount > 0
	isMissingBackup := !strings.Contains(data, ".gitmap/backup")
	issue.MissingGitmapDir = isMissingBackup
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
	line := strings.TrimRight(raw, "\r")
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		*cleaned = append(*cleaned, line)
		return dupCount
	}
	norm := strings.Trim(trimmed, "/ \t\r\n")
	if norm == "" {
		*cleaned = append(*cleaned, line)
		return dupCount
	}
	if seen[norm] {
		return dupCount + 1
	}
	seen[norm] = true
	*cleaned = append(*cleaned, line)
	return dupCount
}

func assembleCleanedGitignore(cleaned []string, seen map[string]bool) string {
	hasBackup := seen[".gitmap/backup"] || seen[".gitmap/backup/"]
	if !hasBackup {
		cleaned = append(cleaned, "", "# GitMap & Task Persistence", ".gitmap/backup/")
		seen[".gitmap/backup"] = true
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
	if err != nil {
		return handleRemediateDefault(issue.RepoPath, issue.RepoName, ignorePath)
	}
	return handleRemediateSanitize(issue.RepoPath, issue.RepoName, ignorePath, string(data))
}

func handleRemediateDefault(repoPath, repoName, ignorePath string) bool {
	isFixed := writeDefaultGitignore(repoPath, ignorePath)
	if isFixed {
		_ = RecordRepoCheckResult(repoPath, repoName, "clean", 1, 0)
	}
	return isFixed
}

func handleRemediateSanitize(repoPath, repoName, ignorePath, data string) bool {
	isFixed := sanitizeAndCommitRepo(repoPath, ignorePath, data)
	if isFixed {
		_ = RecordRepoCheckResult(repoPath, repoName, "clean", 1, 0)
	}
	return isFixed
}

func untrackCachedFiles(repoDir string, files []string) {
	tracked := filterTrackedFilesInIndex(repoDir, files)
	isZero := len(tracked) == 0
	if isZero {
		return
	}
	args := append([]string{"-C", repoDir, "rm", "--cached", "-f", "--"}, tracked...)
	_ = exec.Command("git", args...).Run()
	_ = exec.Command("git", "-C", repoDir, "commit", "-m", "chore(git): untrack ignored files from repository").Run()
}

func filterTrackedFilesInIndex(repoDir string, files []string) []string {
	var tracked []string
	for _, f := range files {
		out, err := exec.Command("git", "-C", repoDir, "ls-files", "--", f).Output()
		hasOut := err == nil && len(strings.TrimSpace(string(out))) > 0
		if hasOut {
			tracked = append(tracked, f)
		}
	}
	return tracked
}

func writeDefaultGitignore(repoDir, ignorePath string) bool {
	content := "# GitMap & Task Persistence\n.gitmap/backup/\n"
	_ = os.WriteFile(ignorePath, []byte(content), 0644)
	_ = exec.Command("git", "-C", repoDir, "add", ".gitignore").Run()
	_ = exec.Command("git", "-C", repoDir, "commit", "-m", "chore(git): add default .gitignore").Run()
	return true
}

func sanitizeAndCommitRepo(repoDir, ignorePath, data string) bool {
	cleaned, dupCount := deduplicateIgnoreContent(data)
	hasBackupDir := strings.Contains(data, ".gitmap/backup")
	hasDuplicates := dupCount > 0
	isMissingBackup := !hasBackupDir
	hasChanges := hasDuplicates || isMissingBackup || (cleaned != data)
	if hasChanges {
		_ = os.WriteFile(ignorePath, []byte(cleaned), 0644)
		_ = exec.Command("git", "-C", repoDir, "add", ".gitignore").Run()
		_ = exec.Command("git", "-C", repoDir, "commit", "-m", "chore(git): sanitize .gitignore").Run()
	}
	return true
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
		fmt.Printf("      - Missing .gitmap/backup/ ignore rule\n")
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
