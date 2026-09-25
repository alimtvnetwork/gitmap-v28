package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SearchPullTraces queries historical commit traces, messages, commit hashes, or repo names from PullRepoRun.
func (s *PullSplitDB) SearchPullTraces(term string, limit int) ([]PullRepoRunRecord, error) {
	boundedLimit := resolvePullSearchLimit(limit)
	pattern := "%" + term + "%"
	query := `SELECT PullRepoRunId, PullRunId, RepoPath, RepoName, PullStatus,
		FilesChanged, LastCommitSha, PreviousCommitSha, CommitMessage, CommitAuthor,
		IsActive, HasChanges, DurationMs, ErrorMessage, Notes, Comments, CreatedAt
	FROM PullRepoRun
	WHERE Notes LIKE ? OR CommitMessage LIKE ? OR LastCommitSha LIKE ? OR RepoName LIKE ?
	ORDER BY PullRepoRunId DESC
	LIMIT ?`

	rows, err := s.conn.Query(query, pattern, pattern, pattern, pattern, boundedLimit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_split.search_traces")
	}
	defer rows.Close()

	return scanRepoRunRecords(rows)
}

func resolvePullSearchLimit(limit int) int {
	if limit <= 0 {
		return 50
	}
	if limit > 500 {
		return 500
	}

	return limit
}
