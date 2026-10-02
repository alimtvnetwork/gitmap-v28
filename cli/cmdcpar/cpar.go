package cmdcpar

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunCPAR executes commit-push-all-repos across all dirty repositories.
func RunCPAR(args []string) *apperror.AppError {
	opts := parseCPAROptions(args)
	dirtyList := CollectDirtyRepositories()
	if isBatchCleanOrAborted(dirtyList, opts) {
		return nil
	}
	dispatchCPARBatch(dirtyList, opts)
	return nil
}

func isBatchCleanOrAborted(dirtyList []DirtyRepoSummary, opts cparOptions) bool {
	hasDirty := len(dirtyList) > 0
	if !hasDirty {
		fmt.Printf("%s✓ All repositories are clean.%s\n", constants.ColorGreen, constants.ColorReset)
		return true
	}
	return isReviewAborted(dirtyList, opts)
}

func isReviewAborted(dirtyList []DirtyRepoSummary, opts cparOptions) bool {
	if !opts.isReview {
		return false
	}
	isProceed := handleReviewFlow(dirtyList, opts)
	if isProceed {
		return false
	}
	fmt.Println("Aborted by user.")
	return true
}

func dispatchCPARBatch(dirtyList []DirtyRepoSummary, opts cparOptions) {
	action := resolveCPARAction(opts)
	queueId, tasksDb := enqueueCPARTaskQueue(action, "local")
	updateCPARTaskQueue(tasksDb, queueId, "running")
	runState := executeCPARAcrossDirty(dirtyList, opts)
	finalizeCPARTaskQueue(tasksDb, queueId, runState.HasFailures)
	closeCPARTaskDb(tasksDb)
}

func resolveCPARAction(opts cparOptions) string {
	if opts.isCommitOnly {
		return "cpar-commit-only"
	}
	return "cpar-commit-push"
}

func CollectDirtyRepositories() []DirtyRepoSummary {
	records := resolveAllRecords()
	var dirty []DirtyRepoSummary
	for _, r := range records {
		relPath := resolveRelativePath(r)
		diag := gitutil.InspectDirtyState(relPath)
		if diag.IsDirty {
			dirty = appendDirty(dirty, r, relPath, diag)
		}
	}
	return dirty
}

func resolveRelativePath(r model.ScanRecord) string {
	hasRel := len(r.RelativePath) > 0
	if hasRel {
		return r.RelativePath
	}
	hasName := len(r.RepoName) > 0
	if hasName {
		return r.RepoName
	}
	return "."
}

func appendDirty(dirty []DirtyRepoSummary, r model.ScanRecord, relPath string, diag gitutil.DirtyDiagnosis) []DirtyRepoSummary {
	return append(dirty, DirtyRepoSummary{
		RepoName:  r.RepoName,
		RepoPath:  relPath,
		Diagnosis: diag,
	})
}

func handleReviewFlow(dirtyList []DirtyRepoSummary, opts cparOptions) bool {
	renderDirtyReviewTable(dirtyList)
	if opts.isAutoYes {
		return true
	}
	return promptReviewConsent(opts.isCommitOnly)
}

func renderDirtyReviewTable(dirty []DirtyRepoSummary) {
	fmt.Printf("\n%s  === PENDING COMMITS REVIEW (%d repositories) ===%s\n",
		constants.ColorCyan, len(dirty), constants.ColorReset)
	for _, d := range dirty {
		printDirtyRow(d)
	}
	fmt.Println()
}

func printDirtyRow(d DirtyRepoSummary) {
	fmt.Printf("    %-30s | modified: %d | untracked: %d | staged: %d\n",
		d.RepoName, d.Diagnosis.ModifiedCount, d.Diagnosis.UntrackedCount, d.Diagnosis.StagedCount)
	renderChangeFileList(d.Diagnosis.AllFiles)
}

func renderChangeFileList(files []string) {
	hasFiles := len(files) > 0
	if !hasFiles {
		return
	}
	renderCappedFiles(files, 5)
}

func renderCappedFiles(files []string, capCount int) {
	for i, f := range files {
		hasExceeded := i >= capCount
		if hasExceeded {
			fmt.Printf("        ... and %d more file(s)\n", len(files)-capCount)
			return
		}
		fmt.Printf("        -> %s\n", f)
	}
}

func promptReviewConsent(isCommitOnly bool) bool {
	action := "commit and push"
	if isCommitOnly {
		action = "commit (without push)"
	}
	fmt.Printf("Proceed to %s across these repositories? [Y/n]: ", action)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return isConsentGranted(strings.ToLower(strings.TrimSpace(text)))
}

func isConsentGranted(low string) bool {
	return low == "" || low == "y" || low == "yes"
}

func resolveAllRecords() []model.ScanRecord {
	records := queryStoreRecords()
	hasRecords := len(records) > 0
	if hasRecords {
		return records
	}
	return []model.ScanRecord{{RelativePath: ".", RepoName: "."}}
}

func queryStoreRecords() []model.ScanRecord {
	s, err := store.OpenDefault()
	hasErr := err != nil
	if hasErr {
		return nil
	}
	defer s.Close()
	records, _ := s.ListRepos()
	return records
}

func executeCPARAcrossDirty(dirty []DirtyRepoSummary, opts cparOptions) CPARRunState {
	fmt.Printf("\n%s  ▶ Executing CPAR across %d dirty repo(s)...%s\n",
		constants.ColorCyan, len(dirty), constants.ColorReset)
	state := CPARRunState{}
	for _, d := range dirty {
		isSuccess := commitAndMaybePush(d, opts)
		state = updateRunState(state, isSuccess)
	}
	fmt.Println()
	return state
}

func updateRunState(state CPARRunState, isSuccess bool) CPARRunState {
	if isSuccess {
		state.SuccessCount++
		return state
	}
	state.HasFailures = true
	state.FailureCount++
	return state
}

func commitAndMaybePush(d DirtyRepoSummary, opts cparOptions) bool {
	_ = exec.Command("git", "-C", d.RepoPath, "add", "-A").Run()
	cmd := exec.Command("git", "-C", d.RepoPath, "commit", "-m", opts.commitMsg)
	cmd.Env = gitutil.BuildSafeGitEnv()
	out, err := cmd.CombinedOutput()
	hasErr := err != nil
	if hasErr {
		logCommitFailure(d, string(out), err)
		return false
	}
	return handlePush(d, opts)
}

func logCommitFailure(d DirtyRepoSummary, output string, err error) {
	fmt.Printf("      %-30s %scommit failed%s\n", d.RepoName, constants.ColorRed, constants.ColorReset)
	rec := store.InternalErrorRecord{
		ErrorCode:  "E_CPAR_COMMIT",
		ErrorType:  "COMMIT_ERROR",
		Command:    "cpar",
		Message:    fmt.Sprintf("git commit failed for %s: %v", d.RepoName, err),
		Details:    output,
		SourceFile: "cli/cmdcpar/cpar.go",
	}
	store.LogInternalErrorRecord(rec)
}

func handlePush(d DirtyRepoSummary, opts cparOptions) bool {
	if opts.isCommitOnly {
		fmt.Printf("      %-30s %scommitted (commit-only)%s\n", d.RepoName, constants.ColorGreen, constants.ColorReset)
		return true
	}
	return pushToRemote(d)
}

func pushToRemote(d DirtyRepoSummary) bool {
	pushCmd := exec.Command("git", "-C", d.RepoPath, "push")
	pushCmd.Env = gitutil.BuildSafeGitEnv()
	out, pushErr := pushCmd.CombinedOutput()
	hasPushErr := pushErr != nil
	if hasPushErr {
		logPushFailure(d, string(out), pushErr)
		return false
	}
	fmt.Printf("      %-30s %scommitted & pushed%s\n", d.RepoName, constants.ColorGreen, constants.ColorReset)
	return true
}

func logPushFailure(d DirtyRepoSummary, output string, err error) {
	fmt.Printf("      %-30s %scommitted (push failed)%s\n", d.RepoName, constants.ColorYellow, constants.ColorReset)
	rec := store.InternalErrorRecord{
		ErrorCode:  "E_CPAR_PUSH",
		ErrorType:  "PUSH_ERROR",
		Command:    "cpar",
		Message:    fmt.Sprintf("git push failed for %s: %v", d.RepoName, err),
		Details:    output,
		SourceFile: "cli/cmdcpar/cpar.go",
	}
	store.LogInternalErrorRecord(rec)
}

func enqueueCPARTaskQueue(action, target string) (string, *store.TasksSplitDB) {
	queueId := fmt.Sprintf("%s-%d", action, time.Now().UnixNano())
	tasksDb, err := store.OpenTasksRootSplitDB()
	hasErr := err != nil
	if hasErr {
		return "", nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(tasksDb.Conn(), "INSERT INTO TaskQueue (QueueId, Section, Action, Target, Status, CreatedAt, UpdatedAt) VALUES (?, 'cpar', ?, ?, 'pending', ?, ?)", queueId, action, target, now, now)
	return queueId, tasksDb
}

func updateCPARTaskQueue(db *store.TasksSplitDB, queueId, status string) {
	hasDb := db != nil
	if !hasDb {
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(db.Conn(), "UPDATE TaskQueue SET Status = ?, UpdatedAt = ? WHERE QueueId = ?", status, now, queueId)
}

func finalizeCPARTaskQueue(db *store.TasksSplitDB, queueId string, hasFailures bool) {
	if hasFailures {
		updateCPARTaskQueue(db, queueId, "failed")
		return
	}
	updateCPARTaskQueue(db, queueId, "completed")
}

func closeCPARTaskDb(db *store.TasksSplitDB) {
	hasDb := db != nil
	if hasDb {
		_ = db.Close()
	}
}
