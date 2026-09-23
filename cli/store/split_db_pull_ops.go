package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// InsertPullRun inserts a master pull execution record and returns the new PullRunId.
func (s *PullSplitDB) InsertPullRun(run *PullRunRecord) (int64, error) {
	query := `INSERT INTO PullRun (
		CommandType, WorkingDir, TotalRepos, PulledRepos, SkippedRepos,
		SuccessCount, FailedCount, IsEfficient, DurationMs, GitMapVersion,
		Notes, Comments
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := s.conn.Exec(query,
		run.CommandType, run.WorkingDir, run.TotalRepos, run.PulledRepos, run.SkippedRepos,
		run.SuccessCount, run.FailedCount, boolToInt(run.IsEfficient), run.DurationMs, run.GitMapVersion,
		run.Notes, run.Comments,
	)
	if err != nil {
		return 0, apperror.WrapSimple(err, "pull_split.insert_run")
	}

	return res.LastInsertId()
}

// InsertPullRepoRuns batch inserts multiple repository pull records within a transaction.
func (s *PullSplitDB) InsertPullRepoRuns(runID int64, records []PullRepoRunRecord) error {
	tx, err := s.conn.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "pull_split.begin_tx")
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := prepareRepoRunStmt(tx)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err := executeRepoRunsBatch(stmt, runID, records); err != nil {
		return err
	}

	return commitRepoRunsTx(tx)
}

func prepareRepoRunStmt(tx *sql.Tx) (*sql.Stmt, error) {
	query := `INSERT INTO PullRepoRun (
		PullRunId, RepoPath, RepoName, PullStatus, FilesChanged,
		LastCommitSha, PreviousCommitSha, CommitMessage, CommitAuthor,
		IsActive, HasChanges, DurationMs, ErrorMessage, Notes, Comments
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_split.prepare_repo_run")
	}

	return stmt, nil
}

func executeRepoRunsBatch(stmt *sql.Stmt, runID int64, records []PullRepoRunRecord) error {
	for _, r := range records {
		_, err := stmt.Exec(
			runID, r.RepoPath, r.RepoName, r.PullStatus, r.FilesChanged,
			r.LastCommitSha, r.PreviousCommitSha, r.CommitMessage, r.CommitAuthor,
			boolToInt(r.IsActive), boolToInt(r.HasChanges), r.DurationMs, r.ErrorMessage, r.Notes, r.Comments,
		)
		if err != nil {
			return apperror.WrapSimple(err, "pull_split.exec_repo_run")
		}
	}

	return nil
}

func commitRepoRunsTx(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return apperror.WrapSimple(err, "pull_split.commit_tx")
	}

	return nil
}

// GetRecentRepoPullHistory returns the last N records for a given repository path.
func (s *PullSplitDB) GetRecentRepoPullHistory(repoPath string, limit int) ([]PullRepoRunRecord, error) {
	query := `SELECT PullRepoRunId, PullRunId, RepoPath, RepoName, PullStatus,
		FilesChanged, LastCommitSha, PreviousCommitSha, CommitMessage, CommitAuthor,
		IsActive, HasChanges, DurationMs, ErrorMessage, Notes, Comments, CreatedAt
	FROM PullRepoRun
	WHERE RepoPath = ?
	ORDER BY PullRepoRunId DESC
	LIMIT ?`

	rows, err := s.conn.Query(query, repoPath, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_split.query_history")
	}
	defer rows.Close()

	return scanRepoRunRecords(rows)
}

func scanRepoRunRecords(rows *sql.Rows) ([]PullRepoRunRecord, error) {
	var results []PullRepoRunRecord
	for rows.Next() {
		rec, err := scanSingleRepoRunRow(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, rec)
	}

	return results, nil
}

func scanSingleRepoRunRow(rows *sql.Rows) (PullRepoRunRecord, error) {
	var r PullRepoRunRecord
	var isActiveInt, hasChangesInt int
	var errMsg, notes, comments sql.NullString

	err := rows.Scan(
		&r.PullRepoRunId, &r.PullRunId, &r.RepoPath, &r.RepoName, &r.PullStatus,
		&r.FilesChanged, &r.LastCommitSha, &r.PreviousCommitSha, &r.CommitMessage, &r.CommitAuthor,
		&isActiveInt, &hasChangesInt, &r.DurationMs, &errMsg, &notes, &comments, &r.CreatedAt,
	)
	if err != nil {
		return r, apperror.WrapSimple(err, "pull_split.scan_row")
	}

	r.IsActive = isActiveInt == 1
	r.HasChanges = hasChangesInt == 1
	r.ErrorMessage = errMsg.String
	r.Notes = notes.String
	r.Comments = comments.String

	return r, nil
}

// EvaluateRepoInactivity determines if a repo is inactive based on minRuns and a 24h window.
func (s *PullSplitDB) EvaluateRepoInactivity(repoPath string, minRuns int, windowHours int) (RepoInactivityStatus, error) {
	status := RepoInactivityStatus{RepoPath: repoPath, IsInactive: false}
	history, err := s.GetRecentRepoPullHistory(repoPath, minRuns)
	if err != nil {
		return status, err
	}

	status.TotalRunsInspected = len(history)
	if len(history) < minRuns {
		status.Reason = fmt.Sprintf("insufficient run history (%d/%d runs)", len(history), minRuns)
		return status, nil
	}

	status.LatestPullTimestamp = history[0].CreatedAt
	status.RepoName = history[0].RepoName
	if isHistoryWindowExpired(history[0].CreatedAt, windowHours) {
		status.Reason = "latest pull is older than 24 hours"
		return status, nil
	}

	zeroCount, hasAnyChanges := countZeroChangeRuns(history)
	status.ConsecutiveZeroRunsCount = zeroCount
	if hasAnyChanges {
		status.Reason = "changes detected in recent runs"
		return status, nil
	}

	status.IsInactive = true
	status.Reason = fmt.Sprintf("0 changes across last %d pulls within 24 hours", zeroCount)
	return status, nil
}

func parseCreatedAtTime(createdAt string) (time.Time, bool) {
	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err == nil {
		return t, true
	}

	t2, err2 := time.Parse(time.RFC3339, createdAt)
	if err2 == nil {
		return t2, true
	}

	return time.Time{}, false
}

func isHistoryWindowExpired(createdAt string, windowHours int) bool {
	t, isParsed := parseCreatedAtTime(createdAt)
	if !isParsed {
		return false
	}

	return time.Since(t) > time.Duration(windowHours)*time.Hour
}

func countZeroChangeRuns(history []PullRepoRunRecord) (int, bool) {
	zeroCount := 0
	for _, rec := range history {
		if rec.HasChanges || rec.FilesChanged > 0 {
			return zeroCount, true
		}
		zeroCount++
	}

	return zeroCount, false
}

func boolToInt(val bool) int {
	if val {
		return 1
	}

	return 0
}
