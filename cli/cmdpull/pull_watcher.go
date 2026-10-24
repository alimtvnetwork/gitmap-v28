package cmdpull

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// WatchPullTask blocks and monitors an active PullTask until completion, failure, or cancellation.
func WatchPullTask(ctx context.Context, db *PullTaskSplitDB, taskID int64) *apperror.AppError {
	task, err := db.GetPullTask(taskID)
	if err != nil {
		return err
	}

	renderWatcherHeader(task)
	ticker := time.NewTicker(time.Duration(constants.PullTaskWatcherPollIntervalSec) * time.Second)
	defer ticker.Stop()

	return runWatcherPollingLoop(ctx, db, taskID, ticker.C)
}

// WatchPullTaskWithDefaults runs WatchPullTask with a background context.
func WatchPullTaskWithDefaults(db *PullTaskSplitDB, taskID int64) *apperror.AppError {
	return WatchPullTask(context.Background(), db, taskID)
}

func runWatcherPollingLoop(ctx context.Context, db *PullTaskSplitDB, taskID int64, ticks <-chan time.Time) *apperror.AppError {
	for {
		select {
		case <-ctx.Done():
			return apperror.WrapSimple(ctx.Err(), "pull_watcher.context_canceled")
		case <-ticks:
			isDone, err := checkWatcherTick(db, taskID)
			if err != nil || isDone {
				return err
			}
		}
	}
}

func checkWatcherTick(db *PullTaskSplitDB, taskID int64) (bool, *apperror.AppError) {
	return pollAndRenderTask(db, taskID)
}

func pollAndRenderTask(db *PullTaskSplitDB, taskID int64) (bool, *apperror.AppError) {
	task, err := db.GetPullTask(taskID)
	if err != nil {
		return false, err
	}

	return processPolledTask(db, task)
}

func processPolledTask(db *PullTaskSplitDB, task *PullTaskRecord) (bool, *apperror.AppError) {
	isTerminal, errTerm := handleTerminalTaskStatus(task)
	if errTerm != nil || isTerminal {
		return isTerminal, errTerm
	}

	isStale := isTaskHeartbeatStale(task.HeartbeatAt)
	if isStale {
		return handleStaleWatcherHeartbeat(db, task.PullTaskId)
	}

	return renderActiveTaskProgress(db, task)
}

func handleTerminalTaskStatus(task *PullTaskRecord) (bool, *apperror.AppError) {
	st := constants.PullTaskStatusType(task.Status)
	if st.IsCompleted() {
		renderWatcherCompletion(task)

		return true, nil
	}

	if st.IsFailed() {
		return handleFailedTaskStatus(task)
	}

	return handleCanceledTaskStatus(task)
}

func handleFailedTaskStatus(task *PullTaskRecord) (bool, *apperror.AppError) {
	renderWatcherFailure(task)
	msg := fmt.Sprintf("pull task failed: %s", task.ErrorMessage)

	return true, apperror.NewExecutionError(msg)
}

func handleCanceledTaskStatus(task *PullTaskRecord) (bool, *apperror.AppError) {
	if constants.PullTaskStatusType(task.Status).IsCanceled() {
		return true, apperror.NewSimple("pull task was canceled", "E_PULL_CANCELED")
	}

	return false, nil
}

func isTaskHeartbeatStale(heartbeatStr string) bool {
	t, err := parseSQLiteTimestamp(heartbeatStr)
	if err != nil {
		return false
	}

	elapsed := time.Since(t)
	threshold := time.Duration(constants.PullTaskStaleLockThresholdSec) * time.Second

	return elapsed > threshold
}

func handleStaleWatcherHeartbeat(db *PullTaskSplitDB, taskID int64) (bool, *apperror.AppError) {
	_ = db.StealStaleLock(taskID)
	fmt.Printf("\n  ⚠ Active pull task #%d stopped sending heartbeats (>%ds). Marked as failed.\n\n",
		taskID, constants.PullTaskStaleLockThresholdSec)

	return true, apperror.NewExecutionError("active pull task heartbeat expired")
}

func renderActiveTaskProgress(db *PullTaskSplitDB, task *PullTaskRecord) (bool, *apperror.AppError) {
	repos, err := db.ListRecentPullTaskRepos(task.PullTaskId, 3)
	if err != nil {
		return false, err
	}

	renderRepoTreeSnapshot(repos)
	renderOverallSummaryLine(task)

	return false, nil
}

func renderWatcherHeader(task *PullTaskRecord) {
	fmt.Printf("\n  → gitmap %s: Active task detected (Task #%d, PID %d)\n",
		task.CommandType, task.PullTaskId, task.Pid)
	fmt.Printf("  Waiting for pull completion... (polling every %ds)\n\n",
		constants.PullTaskWatcherPollIntervalSec)
}

func renderWatcherCompletion(task *PullTaskRecord) {
	fmt.Printf("\n  ✔ Pull task #%d completed successfully (%d/%d repos synchronized).\n\n",
		task.PullTaskId, task.PulledRepos, task.TotalRepos)
}

func renderWatcherFailure(task *PullTaskRecord) {
	fmt.Printf("\n  ✘ Pull task #%d failed: %s\n\n", task.PullTaskId, task.ErrorMessage)
}

func renderRepoTreeSnapshot(repos []PullTaskRepoRecord) {
	reversed := reverseRepoRecords(repos)
	for _, repo := range reversed {
		line := formatRepoTreeLine(repo)
		fmt.Printf("  ├── %s\n", line)
	}
}

func reverseRepoRecords(repos []PullTaskRepoRecord) []PullTaskRepoRecord {
	n := len(repos)
	out := make([]PullTaskRepoRecord, n)
	for i, r := range repos {
		out[n-1-i] = r
	}

	return out
}

func formatRepoTreeLine(repo PullTaskRepoRecord) string {
	durSec := float64(repo.DurationMs) / 1000.0
	slugPadded := fmt.Sprintf("%-25s", repo.RepoSlug)

	if constants.PullTaskRepoStatusType(repo.Status).IsPulling() {
		return fmt.Sprintf("%s ➜ pulling...", slugPadded)
	}

	return formatCompletedRepoLine(repo, slugPadded, durSec)
}

func formatCompletedRepoLine(repo PullTaskRepoRecord, slugPadded string, durSec float64) string {
	repoSt := constants.PullTaskRepoStatusType(repo.Status)
	if repoSt.IsSuccess() {
		return formatSuccessRepoLine(slugPadded, repo.CommitSha, repo.CommitTitle, durSec)
	}

	if repoSt.IsFailed() {
		return fmt.Sprintf("%s ✘ failed (%s)", slugPadded, repo.ErrorMessage)
	}

	return fmt.Sprintf("%s ↷ skipped", slugPadded)
}

func formatSuccessRepoLine(slug, sha, title string, durSec float64) string {
	hasCommit := len(sha) > 0
	if !hasCommit {
		return fmt.Sprintf("%s ✔ up-to-date (%.1fs)", slug, durSec)
	}

	shortSha := truncateShortSha(sha)

	return fmt.Sprintf("%s ✔ pulled (hash: %s - %q) (%.1fs)", slug, shortSha, title, durSec)
}

func truncateShortSha(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}

	return sha
}

func renderOverallSummaryLine(task *PullTaskRecord) {
	elapsed := calculateElapsedSeconds(task.StartedAt)
	etaFormatted := FormatEtaDuration(task.EtaSeconds)

	fmt.Printf("  └── Overall: %d/%d completed · Elapsed: %ds · ETA: ~%s remaining (8-run rolling avg)\n",
		task.PulledRepos, task.TotalRepos, elapsed, etaFormatted)
}

func calculateElapsedSeconds(startedAtStr string) int {
	t, err := parseSQLiteTimestamp(startedAtStr)
	if err != nil {
		return 0
	}

	return int(time.Since(t).Seconds())
}

func parseSQLiteTimestamp(ts string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05",
	}

	for _, layout := range formats {
		if t, err := time.Parse(layout, ts); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("unsupported timestamp format")
}
