package store

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	DefaultFreshnessCooldownMinutes = 5
	DefaultActivityWindowHours      = 24
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
// It inspects the local git trace to see if there are commits within the 24h window.
func (s *PullSplitDB) EvaluateRepoActivityStatus(repoPath string, windowHours int, cooldownMinutes int) (RepoInactivityStatus, error) {
	status := RepoInactivityStatus{RepoPath: repoPath, IsInactive: false}

	// ALWAYS check if there's a cooldown from a very recent pull.
	// This prevents spamming `pae` repeatedly.
	history, err := s.GetRecentActualRepoPullHistory(repoPath, 1)
	if err == nil && len(history) > 0 {
		latestTime, isParsed := parseCreatedAtTime(history[0].CreatedAt)
		if isParsed {
			elapsed := time.Since(latestTime)
			cooldownDur := time.Duration(cooldownMinutes) * time.Minute
			if cooldownMinutes > 0 && elapsed <= cooldownDur {
				status.IsInactive = true
				status.Reason = fmt.Sprintf("already checked recently (<%dm ago)", cooldownMinutes)
				return status, nil
			}
		}
	}

	// 1. Run git log -1 --format=%ct to get the latest commit timestamp
	cmd := exec.Command("git", "-C", repoPath, "log", "-1", "--format=%ct")
	out, err := cmd.Output()
	if err != nil {
		// If git log fails, we assume it's active so we pull it (must establish baseline or error later)
		status.Reason = "git log failed (must establish baseline)"
		return status, nil
	}

	strOut := strings.TrimSpace(string(out))
	timestamp, err := strconv.ParseInt(strOut, 10, 64)
	if err != nil {
		status.Reason = "git log timestamp parse failed"
		return status, nil
	}

	commitTime := time.Unix(timestamp, 0)
	elapsed := time.Since(commitTime)
	windowDur := time.Duration(windowHours) * time.Hour

	// 2. If the last commit is older than the window, mark as inactive
	if elapsed > windowDur {
		status.IsInactive = true
		status.Reason = fmt.Sprintf("no commits in local git trace within %dh window", windowHours)
		return status, nil
	}

	status.Reason = "active repository (recent commits in git trace)"
	return status, nil
}

