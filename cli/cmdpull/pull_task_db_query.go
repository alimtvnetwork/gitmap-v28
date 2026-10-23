package cmdpull

import (
	"database/sql"
	"errors"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	sqlSelectActiveLockWithElapsed = `SELECT PullTaskId, TaskUuid, CommandType, TargetNode, MachineId, Status,
    Pid, TotalRepos, PulledRepos, EtaSeconds,
    COALESCE(StartedAt, ''), COALESCE(CompletedAt, ''),
    COALESCE(ErrorMessage, ''), HeartbeatAt, CreatedAt,
    (strftime('%s', 'now') - strftime('%s', HeartbeatAt)) as ElapsedSec
FROM PullTask
WHERE TargetNode = ? AND Status IN ('pending', 'running')
ORDER BY PullTaskId DESC LIMIT 1;`

	sqlSelectListRecentRepos = `SELECT PullTaskRepoId, PullTaskId, RepoSlug, RepoPath, Status,
    CommitSha, CommitTitle, FilesChanged, DurationMs,
    COALESCE(ErrorMessage, ''), COALESCE(PulledAt, ''), CreatedAt
FROM PullTaskRepo WHERE PullTaskId = ?
ORDER BY PullTaskRepoId DESC LIMIT ?;`
)

// FindActiveLock queries for an active lock on targetNode with fresh heartbeat.
func (s *PullTaskSplitDB) FindActiveLock(targetNode string) (*PullTaskRecord, bool, *apperror.AppError) {
	node := resolveTargetNode(targetNode)
	row := s.conn.QueryRow(sqlSelectActiveLockWithElapsed, node)

	rec, elapsedSec, err := scanActiveLockRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, apperror.WrapSimple(err, "pull_task_db.find_lock")
	}

	return evaluateLockFreshness(s, rec, elapsedSec)
}

func evaluateLockFreshness(s *PullTaskSplitDB, rec *PullTaskRecord, elapsedSec int64) (*PullTaskRecord, bool, *apperror.AppError) {
	isFresh := elapsedSec < int64(constants.PullTaskStaleLockThresholdSec)
	if isFresh {
		return rec, true, nil
	}

	if err := s.StealStaleLock(rec.PullTaskId); err != nil {
		return nil, false, err
	}

	return nil, false, nil
}

// GetPullTask fetches a single PullTaskRecord by its primary key.
func (s *PullTaskSplitDB) GetPullTask(taskID int64) (*PullTaskRecord, *apperror.AppError) {
	query := `SELECT PullTaskId, TaskUuid, CommandType, TargetNode, MachineId, Status,
		Pid, TotalRepos, PulledRepos, EtaSeconds,
		COALESCE(StartedAt, ''), COALESCE(CompletedAt, ''),
		COALESCE(ErrorMessage, ''), HeartbeatAt, CreatedAt
		FROM PullTask WHERE PullTaskId = ?`

	row := s.conn.QueryRow(query, taskID)
	rec, err := scanPullTaskRow(row)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_task_db.get_task")
	}

	return rec, nil
}

// ListRecentPullTaskRepos returns the most recent repository progress entries.
func (s *PullTaskSplitDB) ListRecentPullTaskRepos(taskID int64, limit int) ([]PullTaskRepoRecord, *apperror.AppError) {
	rows, err := s.conn.Query(sqlSelectListRecentRepos, taskID, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_task_db.list_repos")
	}

	defer rows.Close()

	return scanPullTaskRepoRows(rows)
}

func scanActiveLockRow(row *sql.Row) (*PullTaskRecord, int64, error) {
	var t PullTaskRecord
	var elapsed sql.NullInt64

	err := row.Scan(&t.PullTaskId, &t.TaskUuid, &t.CommandType, &t.TargetNode, &t.MachineId, &t.Status,
		&t.Pid, &t.TotalRepos, &t.PulledRepos, &t.EtaSeconds,
		&t.StartedAt, &t.CompletedAt, &t.ErrorMessage, &t.HeartbeatAt, &t.CreatedAt, &elapsed)
	if err != nil {
		return nil, 0, err
	}

	return &t, elapsed.Int64, nil
}

func scanPullTaskRow(row *sql.Row) (*PullTaskRecord, error) {
	var t PullTaskRecord

	err := row.Scan(&t.PullTaskId, &t.TaskUuid, &t.CommandType, &t.TargetNode, &t.MachineId, &t.Status,
		&t.Pid, &t.TotalRepos, &t.PulledRepos, &t.EtaSeconds,
		&t.StartedAt, &t.CompletedAt, &t.ErrorMessage, &t.HeartbeatAt, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func scanPullTaskRepoRows(rows *sql.Rows) ([]PullTaskRepoRecord, *apperror.AppError) {
	var records []PullTaskRepoRecord
	for rows.Next() {
		r, err := scanSingleRepoRow(rows)
		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil
}

func scanSingleRepoRow(rows *sql.Rows) (PullTaskRepoRecord, *apperror.AppError) {
	var r PullTaskRepoRecord
	err := rows.Scan(&r.PullTaskRepoId, &r.PullTaskId, &r.RepoSlug, &r.RepoPath, &r.Status,
		&r.CommitSha, &r.CommitTitle, &r.FilesChanged, &r.DurationMs,
		&r.ErrorMessage, &r.PulledAt, &r.CreatedAt)
	if err != nil {
		return r, apperror.WrapSimple(err, "pull_task_db.scan_repo")
	}

	return r, nil
}

func resolveTargetNode(node string) string {
	clean := strings.TrimSpace(node)
	if clean == "" {
		return "local"
	}

	return clean
}

func resolveMachineId() string {
	host, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return host
}

func resolveRetentionDays(days int) int {
	if days <= 0 {
		return constants.DefaultPullRetentionDays
	}

	return days
}
