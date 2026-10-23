package cmdpull

import (
	"database/sql"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/uuidutil"
)

const (
	sqlInsertPullTask = `INSERT INTO PullTask (
    TaskUuid, CommandType, TargetNode, MachineId, Status, Pid,
    TotalRepos, PulledRepos, EtaSeconds, StartedAt, HeartbeatAt, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`

	sqlInsertPullTaskRepo = `INSERT INTO PullTaskRepo (
    PullTaskId, RepoSlug, RepoPath, Status, CommitSha, CommitTitle,
    FilesChanged, DurationMs, ErrorMessage, PulledAt, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`

	sqlDeleteExpiredRepos = `DELETE FROM PullTaskRepo WHERE PullTaskId IN (
    SELECT PullTaskId FROM PullTask WHERE datetime(CreatedAt) < datetime('now', '-' || ? || ' days')
);`

	sqlDeleteExpiredTasks = `DELETE FROM PullTask WHERE datetime(CreatedAt) < datetime('now', '-' || ? || ' days');`

	sqlDeleteExpiredEta = `DELETE FROM PullEtaHistory WHERE datetime(CreatedAt) < datetime('now', '-' || ? || ' days');`
)

// StealStaleLock marks an orphaned or crashed running task as failed.
func (s *PullTaskSplitDB) StealStaleLock(taskID int64) *apperror.AppError {
	query := `UPDATE PullTask
		SET Status = ?, ErrorMessage = ?, CompletedAt = CURRENT_TIMESTAMP
		WHERE PullTaskId = ?`
	errMsg := "process terminated unexpectedly (stale heartbeat)"

	if _, err := s.conn.Exec(query, constants.PullTaskStatusFailed, errMsg, taskID); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.steal_stale_lock")
	}

	return nil
}

// TryAcquirePullTaskLock attempts to acquire the execution lock for a node.
func (s *PullTaskSplitDB) TryAcquirePullTaskLock(cmdType, targetNode string, totalRepos, etaSec int) (*PullTaskRecord, bool, *apperror.AppError) {
	activeRec, hasLock, err := s.FindActiveLock(targetNode)
	if err != nil {
		return nil, false, err
	}

	if hasLock {
		return activeRec, false, nil
	}

	return s.acquireFreshLock(cmdType, targetNode, totalRepos, etaSec)
}

func (s *PullTaskSplitDB) acquireFreshLock(cmdType, targetNode string, totalRepos, etaSec int) (*PullTaskRecord, bool, *apperror.AppError) {
	newRec, err := s.insertNewPullTask(cmdType, targetNode, totalRepos, etaSec)
	if err != nil {
		return nil, false, err
	}

	return newRec, true, nil
}

func (s *PullTaskSplitDB) insertNewPullTask(cmdType, targetNode string, totalRepos, etaSec int) (*PullTaskRecord, *apperror.AppError) {
	node := resolveTargetNode(targetNode)
	machineId := resolveMachineId()
	uuidStr := uuidutil.GenerateUUID()

	res, err := s.conn.Exec(sqlInsertPullTask, uuidStr, cmdType, node, machineId,
		constants.PullTaskStatusRunning, os.Getpid(), totalRepos, etaSec)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_task_db.insert_task")
	}

	return s.fetchInsertedTask(res)
}

func (s *PullTaskSplitDB) fetchInsertedTask(res sql.Result) (*PullTaskRecord, *apperror.AppError) {
	id, err := res.LastInsertId()
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_task_db.last_insert_id")
	}

	return s.GetPullTask(id)
}

// UpdateTaskHeartbeat refreshes HeartbeatAt timestamp to maintain execution lease.
func (s *PullTaskSplitDB) UpdateTaskHeartbeat(taskID int64) *apperror.AppError {
	query := `UPDATE PullTask SET HeartbeatAt = CURRENT_TIMESTAMP WHERE PullTaskId = ?`
	if _, err := s.conn.Exec(query, taskID); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.heartbeat")
	}

	return nil
}

// UpdateTaskProgress updates pulled count, remaining ETA, and refreshes heartbeat.
func (s *PullTaskSplitDB) UpdateTaskProgress(taskID int64, pulledRepos, etaSec int) *apperror.AppError {
	query := `UPDATE PullTask
		SET PulledRepos = ?, EtaSeconds = ?, HeartbeatAt = CURRENT_TIMESTAMP
		WHERE PullTaskId = ?`

	if _, err := s.conn.Exec(query, pulledRepos, etaSec, taskID); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.progress")
	}

	return nil
}

// RecordRepoPullResult inserts progress for a repository under an active task.
func (s *PullTaskSplitDB) RecordRepoPullResult(rec PullTaskRepoRecord) *apperror.AppError {
	title := sanitizeCommitTitle(rec.CommitTitle)
	sha := truncateSha(rec.CommitSha)

	_, err := s.conn.Exec(sqlInsertPullTaskRepo, rec.PullTaskId, rec.RepoSlug, rec.RepoPath, rec.Status,
		sha, title, rec.FilesChanged, rec.DurationMs, rec.ErrorMessage)
	if err != nil {
		return apperror.WrapSimple(err, "pull_task_db.record_repo")
	}

	return nil
}

// CompletePullTask marks the task completed and appends runtime to PullEtaHistory.
func (s *PullTaskSplitDB) CompletePullTask(taskID int64, totalRepos int, durationMs int64, cmdType, targetNode string) *apperror.AppError {
	query := `UPDATE PullTask
		SET Status = ?, PulledRepos = TotalRepos, CompletedAt = CURRENT_TIMESTAMP, HeartbeatAt = CURRENT_TIMESTAMP
		WHERE PullTaskId = ?`

	if _, err := s.conn.Exec(query, constants.PullTaskStatusCompleted, taskID); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.complete")
	}

	durationSec := float64(durationMs) / 1000.0

	return s.insertEtaHistory(cmdType, targetNode, totalRepos, durationSec)
}

func (s *PullTaskSplitDB) insertEtaHistory(cmdType, targetNode string, totalRepos int, durationSec float64) *apperror.AppError {
	node := resolveTargetNode(targetNode)
	query := `INSERT INTO PullEtaHistory (CommandType, TargetNode, TotalRepos, DurationSeconds, CreatedAt)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`

	if _, err := s.conn.Exec(query, cmdType, node, totalRepos, durationSec); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.insert_eta_history")
	}

	return nil
}

// FailPullTask marks the task failed with an explanatory error message.
func (s *PullTaskSplitDB) FailPullTask(taskID int64, errMsg string) *apperror.AppError {
	query := `UPDATE PullTask
		SET Status = ?, ErrorMessage = ?, CompletedAt = CURRENT_TIMESTAMP, HeartbeatAt = CURRENT_TIMESTAMP
		WHERE PullTaskId = ?`

	if _, err := s.conn.Exec(query, constants.PullTaskStatusFailed, errMsg, taskID); err != nil {
		return apperror.WrapSimple(err, "pull_task_db.fail")
	}

	return nil
}

// PruneExpiredPullHistory purges records older than retentionDays.
func (s *PullTaskSplitDB) PruneExpiredPullHistory(retentionDays int) (int64, *apperror.AppError) {
	daysStr := strconv.Itoa(resolveRetentionDays(retentionDays))

	deletedCount, err := s.purgeRetentionTables(daysStr)
	if err != nil {
		return 0, err
	}

	s.vacuumIfLargePurge(deletedCount)

	return deletedCount, nil
}

func (s *PullTaskSplitDB) purgeRetentionTables(daysStr string) (int64, *apperror.AppError) {
	queries := []string{sqlDeleteExpiredRepos, sqlDeleteExpiredTasks, sqlDeleteExpiredEta}
	var total int64

	for _, q := range queries {
		cnt, err := execDeleteQuery(s.conn, q, daysStr)
		if err != nil {
			return 0, err
		}

		total += cnt
	}

	return total, nil
}

func (s *PullTaskSplitDB) vacuumIfLargePurge(totalDeleted int64) {
	if totalDeleted > 500 {
		_, _ = s.conn.Exec("PRAGMA incremental_vacuum(100);")
	}
}

func execDeleteQuery(conn *sql.DB, query, arg string) (int64, *apperror.AppError) {
	res, err := conn.Exec(query, arg)
	if err != nil {
		return 0, apperror.WrapSimple(err, "pull_task_db.prune")
	}

	rows, _ := res.RowsAffected()

	return rows, nil
}

func sanitizeCommitTitle(title string) string {
	clean := strings.TrimSpace(title)
	if len(clean) > 120 {
		return clean[:120]
	}

	return clean
}

func truncateSha(sha string) string {
	clean := strings.TrimSpace(sha)
	if len(clean) > 40 {
		return clean[:40]
	}

	return clean
}
