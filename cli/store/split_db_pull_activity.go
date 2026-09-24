package store

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	DefaultFreshnessCooldownMinutes = 5
	DefaultActivityWindowHours     = 24
)

// GetRecentActualRepoPullHistory returns the last N genuine pull execution records (excluding skipped-inactive).
func (s *PullSplitDB) GetRecentActualRepoPullHistory(repoPath string, limit int) ([]PullRepoRunRecord, error) {
	query := `SELECT PullRepoRunId, PullRunId, RepoPath, RepoName, PullStatus,
		FilesChanged, LastCommitSha, PreviousCommitSha, CommitMessage, CommitAuthor,
		IsActive, HasChanges, DurationMs, ErrorMessage, Notes, Comments, CreatedAt
	FROM PullRepoRun
	WHERE RepoPath = ? AND PullStatus != 'skipped-inactive'
	ORDER BY PullRepoRunId DESC
	LIMIT ?`

	rows, err := s.conn.Query(query, repoPath, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_split.query_actual_history")
	}
	defer rows.Close()

	return scanRepoRunRecords(rows)
}

// EvaluateRepoActivityStatus evaluates whether a repository should be pulled or skipped in efficient mode.
func (s *PullSplitDB) EvaluateRepoActivityStatus(repoPath string, windowHours int, cooldownMinutes int) (RepoInactivityStatus, error) {
	status := RepoInactivityStatus{RepoPath: repoPath, IsInactive: false}
	history, err := s.GetRecentActualRepoPullHistory(repoPath, 5)
	if err != nil {
		return status, err
	}

	status.TotalRunsInspected = len(history)
	if len(history) == 0 {
		status.Reason = "no prior pull history (must establish baseline)"
		return status, nil
	}

	latest := history[0]
	status.LatestPullTimestamp = latest.CreatedAt
	status.RepoName = latest.RepoName

	latestTime, isParsed := parseCreatedAtTime(latest.CreatedAt)
	if !isParsed {
		status.Reason = "unparseable latest pull timestamp"
		return status, nil
	}

	elapsed := time.Since(latestTime)
	return evaluateElapsedActivity(status, latest, elapsed, windowHours, cooldownMinutes)
}

func evaluateElapsedActivity(status RepoInactivityStatus, latest PullRepoRunRecord, elapsed time.Duration, windowHours, cooldownMinutes int) (RepoInactivityStatus, error) {
	cooldownDur := time.Duration(cooldownMinutes) * time.Minute
	if cooldownMinutes > 0 && elapsed <= cooldownDur {
		status.IsInactive = true
		status.Reason = fmt.Sprintf("already checked recently (<%dm ago)", cooldownMinutes)
		return status, nil
	}

	windowDur := time.Duration(windowHours) * time.Hour
	if elapsed > windowDur {
		status.Reason = fmt.Sprintf("latest pull is older than %d hours (daily refresh)", windowHours)
		return status, nil
	}

	if !latest.HasChanges && latest.FilesChanged == 0 {
		status.IsInactive = true
		status.Reason = fmt.Sprintf("0 changes on latest pull within %dh window", windowHours)
		return status, nil
	}

	status.Reason = "active repository (recent changes within 24h)"
	return status, nil
}
