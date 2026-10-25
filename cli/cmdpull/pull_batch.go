package cmdpull

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"os"
	"time"
)

func runPullBatch(opts pullOptions) error {
	records, isFound := resolvePullBatchRecords(opts)
	hasRecords := isFound && len(records) > 0
	if !hasRecords {
		return handleEmptyBatchRecords(opts.isJSON)
	}
	printResolvedPullReposUnlessJSON(len(records), opts.isJSON)
	filtered := applyPullAvailableFilter(records, opts.onlyAvailable)
	if isNoAvailablePullTargets(opts.onlyAvailable, len(filtered)) {
		return handleNoAvailablePullTargets(opts)
	}

	return executePullBatchLifecycle(filtered, opts)
}

func printResolvedPullReposUnlessJSON(count int, isJSON bool) {
	if !isJSON {
		printResolvedPullRepos(count)
	}
}

func isNoAvailablePullTargets(onlyAvailable bool, count int) bool {
	return onlyAvailable && count == 0
}

func handleNoAvailablePullTargets(opts pullOptions) error {
	if opts.isJSON {
		return renderPullBatchJSONSummary(0, nil, 0)
	}
	fmt.Print(constants.MsgPullNoAvailable)

	return nil
}

func printResolvedPullRepos(count int) {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s resolved %s%d%s repo(s) to pull\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorBold, count, constants.ColorReset)
}

func applyPullAvailableFilter(records []model.ScanRecord, isOnlyAvailable bool) []model.ScanRecord {
	if isOnlyAvailable {
		return filterByAvailableUpdates(records)
	}

	return records
}

func executePullBatchLifecycle(records []model.ScanRecord, opts pullOptions) error {
	taskID, taskDB := beginPullTask(records)
	queueId, tDB := initLocalPullTaskQueue(opts.all)
	defer closePullBatchDBs(taskDB, tDB)

	applySelectiveTransport(records, opts)
	ignoreHandle := startAsyncIgnoreScan(records)
	bar, states, dur := runPullBatchWork(records, opts)
	if opts.isJSON {
		return finalizePullBatchJSON(taskDB, taskID, tDB, queueId, opts.all, len(records), states, dur)
	}

	return finalizePullBatchStandard(taskDB, taskID, tDB, queueId, opts.all, records, states, dur, opts, bar.Failed(), ignoreHandle)
}

func closePullBatchDBs(taskDB *store.DB, tDB *store.TasksSplitDB) {
	if taskDB != nil {
		_ = taskDB.Close()
	}
	closeTaskDB(tDB)
}

func initLocalPullTaskQueue(isAll bool) (string, *store.TasksSplitDB) {
	if !isAll {
		return "", nil
	}
	queueId, tDB := enqueueTaskQueue("pull-all", "local")
	updateTaskQueue(tDB, queueId, "running")

	return queueId, tDB
}

func applySelectiveTransport(records []model.ScanRecord, opts pullOptions) {
	if !opts.all {
		maybeApplyTransportToRecords(records, opts.useSSH, opts.useHTTPS)
	}
}

func runPullBatchWork(records []model.ScanRecord, opts pullOptions) (*PullProgressBar, []*PullRepoState, time.Duration) {
	bar, sortedStates, dur := runPullBatchExecution(records, opts)
	syncPullBatchTelemetry(records, sortedStates, dur, opts)

	return bar, sortedStates, dur
}

func finalizePullBatchJSON(taskDB *store.DB, taskID int64, tDB *store.TasksSplitDB, queueId string, isAll bool, total int, states []*PullRepoState, dur time.Duration) error {
	completePendingTask(taskDB, taskID)
	finalizePullTaskQueueJSON(isAll, tDB, queueId)

	return renderPullBatchJSONSummary(total, states, dur)
}

func finalizePullBatchStandard(taskDB *store.DB, taskID int64, tDB *store.TasksSplitDB, queueId string, isAll bool, records []model.ScanRecord, states []*PullRepoState, dur time.Duration, opts pullOptions, failedCount int, ignoreHandle *IgnoreScanHandle) error {
	renderPullBatchOutput(records, states, dur, opts)
	collectAndRemediateIgnoreIssues(ignoreHandle, opts)
	finalizePullTaskQueueOutput(isAll, tDB, queueId, failedCount)

	unresolved := filterValidUnresolvedRepos(ExtractUnresolvedRepos(states))
	hasUnresolved := len(unresolved) > 0
	if hasUnresolved && !opts.isJSON {
		_ = PromptPostPullRemediation(unresolved)
	}

	validFailures := filterValidPullFailures(ExtractPullFailures(states))

	return finalizePullBatchTask(taskDB, taskID, failedCount, validFailures)
}

func filterValidPullFailures(failures []PullFailureSummary) []PullFailureSummary {
	var valid []PullFailureSummary
	for _, f := range failures {
		if isValidRepoNameCandidate(f.RepoName) {
			valid = append(valid, f)
		}
	}

	return valid
}

func filterValidUnresolvedRepos(repos []UnresolvedPullRepo) []UnresolvedPullRepo {
	var valid []UnresolvedPullRepo
	for _, r := range repos {
		if isValidRepoNameCandidate(r.RepoName) {
			valid = append(valid, r)
		}
	}

	return valid
}

// AppendPullErrorIfValid records a pull failure to split DB only if repoSlug is valid.
func AppendPullErrorIfValid(db *store.PullSplitDB, rec store.PullErrorRecord) error {
	if !isValidRepoNameCandidate(rec.RepoSlug) {
		return nil
	}
	if db == nil {
		return nil
	}

	return db.InsertPullError(rec)
}

// AppendFailureIfValid appends a failure record only if the repo slug is a valid candidate.
func AppendFailureIfValid(failures []PullFailureSummary, f PullFailureSummary) []PullFailureSummary {
	if !isValidRepoNameCandidate(f.RepoName) {
		return failures
	}

	return append(failures, f)
}

func finalizePullTaskQueueJSON(isAll bool, tDB *store.TasksSplitDB, queueId string) {
	if isAll {
		updateTaskQueue(tDB, queueId, "completed")
	}
}

func finalizePullTaskQueueOutput(isAll bool, tDB *store.TasksSplitDB, queueId string, failedCount int) {
	if !isAll {
		return
	}
	status := "completed"
	if failedCount > 0 {
		status = "failed"
	}
	updateTaskQueue(tDB, queueId, status)
}

func handlePullRemediationForRecords(records []model.ScanRecord, opts pullOptions) {
	remItems := collectRemediationItems(records)
	handlePullRemediation(remItems, opts)
}

func collectRemediationItems(records []model.ScanRecord) []RemediationItem {
	var items []RemediationItem
	for _, rec := range records {
		item, hasItem := evaluateRemediationRecord(rec)
		if hasItem {
			items = append(items, item)
		}
	}

	return items
}

func evaluateRemediationRecord(rec model.ScanRecord) (RemediationItem, bool) {
	if !isValidRepoNameCandidate(rec.RepoName) {
		return RemediationItem{}, false
	}
	diag := gitutil.InspectDirtyState(rec.AbsolutePath)
	if !diag.IsDirty {
		return RemediationItem{}, false
	}

	return buildRemediationItem(rec, diag), true
}

func buildRemediationItem(rec model.ScanRecord, diag gitutil.DirtyDiagnosis) RemediationItem {
	recipes := gitutil.GenerateRemediationRecipes(rec.AbsolutePath, diag)

	return RemediationItem{
		RepoName:      rec.RepoName,
		RepoPath:      rec.AbsolutePath,
		SummaryReason: diag.SummaryReason,
		Recipes:       recipes,
		Files:         diag.AllFiles,
	}
}

func finalizePullBatchTask(taskDB *store.DB, taskID int64, failCount int, failures []PullFailureSummary) error {
	if failCount > 0 {
		return handleFailedBatchTask(taskDB, taskID, failCount, failures)
	}
	completePendingTask(taskDB, taskID)

	return nil
}

func handleFailedBatchTask(taskDB *store.DB, taskID int64, failCount int, failures []PullFailureSummary) error {
	errMsg := fmt.Sprintf("pull batch finished with %d failure(s)", failCount)
	failPendingTask(taskDB, taskID, errMsg)
	fmt.Fprintf(os.Stderr, "\n  %s%s%s\n\n", constants.ColorRed, errMsg, constants.ColorReset)
	handleBatchFailuresRemediation(failures)
	cliexit.Exit(1)

	return nil
}

func handleBatchFailuresRemediation(failures []PullFailureSummary) {
	valid := filterValidPullFailures(failures)
	hasFailures := len(valid) > 0
	if hasFailures {
		_ = PromptInteractiveBatchFix(valid)
	}
}
