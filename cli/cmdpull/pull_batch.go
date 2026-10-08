package cmdpull

import (
	"fmt"
	"time"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"os"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
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

	return finalizePullBatchTask(taskDB, taskID, failedCount, ExtractPullFailures(states))
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
	var remItems []RemediationItem
	for _, rec := range records {
		diag := gitutil.InspectDirtyState(rec.AbsolutePath)
		if diag.IsDirty {
			remItems = append(remItems, buildRemediationItem(rec, diag))
		}
	}
	handlePullRemediation(remItems, opts)
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
		errMsg := fmt.Sprintf("pull batch finished with %d failure(s)", failCount)
		failPendingTask(taskDB, taskID, errMsg)
		fmt.Fprintf(os.Stderr, "\n  %s%s%s\n\n", constants.ColorRed, errMsg, constants.ColorReset)
		handleBatchFailuresRemediation(failures)
		cliexit.Exit(1)

		return nil
	}
	completePendingTask(taskDB, taskID)

	return nil
}

func handleBatchFailuresRemediation(failures []PullFailureSummary) {
	hasFailures := len(failures) > 0
	if hasFailures {
		_ = PromptInteractiveBatchFix(failures)
	}
}
