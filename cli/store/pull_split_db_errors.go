package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	sqlCreatePullErrors = `CREATE TABLE IF NOT EXISTS pull_errors (
    error_id        TEXT PRIMARY KEY,
    repo_slug       TEXT NOT NULL,
    repo_path       TEXT NOT NULL,
    node_id         TEXT NOT NULL,
    node_version    TEXT NOT NULL,
    error_type      TEXT NOT NULL,
    error_text      TEXT NOT NULL,
    stack_trace     TEXT NULL,
    remediation_cmd TEXT NULL,
    created_at      DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_pull_errors_repo ON pull_errors(repo_slug);
CREATE INDEX IF NOT EXISTS idx_pull_errors_created ON pull_errors(created_at DESC);`
)

// PullErrorRecord represents a diagnosed pull failure stored in the split pull database.
type PullErrorRecord struct {
	ErrorID        string    `json:"error_id"`
	RepoSlug       string    `json:"repo_slug"`
	RepoPath       string    `json:"repo_path"`
	NodeID         string    `json:"node_id"`
	NodeVersion    string    `json:"node_version"`
	ErrorType      string    `json:"error_type"`
	ErrorText      string    `json:"error_text"`
	StackTrace     string    `json:"stack_trace"`
	RemediationCmd string    `json:"remediation_cmd"`
	CreatedAt      time.Time `json:"created_at"`
}

// EnsurePullErrorsTable initializes the pull_errors table and its indexes.
func (s *PullSplitDB) EnsurePullErrorsTable() error {
	if s.conn == nil {
		return apperror.NewValidationError("database connection is nil")
	}

	if _, err := s.conn.Exec(sqlCreatePullErrors); err != nil {
		return apperror.WrapSimple(err, "pull_errors.init_schema")
	}

	return nil
}

// InsertPullError saves a diagnosed pull failure record into pull_errors.
func (s *PullSplitDB) InsertPullError(rec PullErrorRecord) error {
	if s.conn == nil {
		return apperror.NewValidationError("database connection is nil")
	}

	rec = preparePullErrorRecord(rec)
	return s.execInsertPullError(rec)
}

func preparePullErrorRecord(rec PullErrorRecord) PullErrorRecord {
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = time.Now().UTC()
	}
	if rec.ErrorID == "" {
		rec.ErrorID = fmt.Sprintf("err-%d", time.Now().UnixNano())
	}
	return rec
}

func (s *PullSplitDB) execInsertPullError(rec PullErrorRecord) error {
	query := `INSERT OR REPLACE INTO pull_errors
		(error_id, repo_slug, repo_path, node_id, node_version, error_type, error_text, stack_trace, remediation_cmd, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.conn.Exec(query,
		rec.ErrorID, rec.RepoSlug, rec.RepoPath, rec.NodeID, rec.NodeVersion,
		rec.ErrorType, rec.ErrorText, rec.StackTrace, rec.RemediationCmd,
		rec.CreatedAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return apperror.WrapSimple(err, "pull_errors.insert")
	}
	return nil
}

// QueryAllLatestPullErrors retrieves the most recent pull errors across all repositories.
func (s *PullSplitDB) QueryAllLatestPullErrors(limit int) ([]PullErrorRecord, error) {
	return s.QueryLatestPullErrors("all", limit)
}

// QueryLatestPullErrors retrieves the most recent pull errors, optionally filtered by repoSlug.
func (s *PullSplitDB) QueryLatestPullErrors(repoSlug string, limit int) ([]PullErrorRecord, error) {
	if s.conn == nil {
		return nil, apperror.NewValidationError("database connection is nil")
	}

	limit = normalizePullErrorLimit(limit)
	query, args := buildPullErrorQuery(repoSlug, limit)
	rows, err := s.conn.Query(query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "pull_errors.query")
	}
	defer rows.Close()

	return scanPullErrorRows(rows)
}

func normalizePullErrorLimit(limit int) int {
	if limit <= 0 {
		return 50
	}
	return limit
}

func buildPullErrorQuery(repoSlug string, limit int) (string, []any) {
	if isAllReposQuery(repoSlug) {
		query := `SELECT error_id, repo_slug, repo_path, node_id, node_version, error_type, error_text, stack_trace, remediation_cmd, created_at
			FROM pull_errors ORDER BY created_at DESC LIMIT ?`
		return query, []any{limit}
	}

	query := `SELECT error_id, repo_slug, repo_path, node_id, node_version, error_type, error_text, stack_trace, remediation_cmd, created_at
		FROM pull_errors WHERE repo_slug = ? OR repo_path LIKE ? ORDER BY created_at DESC LIMIT ?`
	return query, []any{repoSlug, "%" + repoSlug + "%", limit}
}

func isAllReposQuery(repoSlug string) bool {
	clean := strings.TrimSpace(strings.ToLower(repoSlug))
	return clean == "" || clean == "all"
}

func scanPullErrorRows(rows *sql.Rows) ([]PullErrorRecord, error) {
	var records []PullErrorRecord
	for rows.Next() {
		rec, isOk := scanSinglePullError(rows)
		if isOk {
			records = append(records, rec)
		}
	}
	return records, rows.Err()
}

func scanSinglePullError(rows *sql.Rows) (PullErrorRecord, bool) {
	var rec PullErrorRecord
	var createdAtStr string
	var stackTrace, remediationCmd *string

	err := rows.Scan(
		&rec.ErrorID, &rec.RepoSlug, &rec.RepoPath, &rec.NodeID, &rec.NodeVersion,
		&rec.ErrorType, &rec.ErrorText, &stackTrace, &remediationCmd, &createdAtStr,
	)
	if err != nil {
		return rec, false
	}

	rec = assignNullableFields(rec, stackTrace, remediationCmd, createdAtStr)
	return rec, true
}

func assignNullableFields(rec PullErrorRecord, stack, remed *string, createdStr string) PullErrorRecord {
	if stack != nil {
		rec.StackTrace = *stack
	}
	if remed != nil {
		rec.RemediationCmd = *remed
	}
	rec.CreatedAt = parseFlexibleDBTimestamp(createdStr)
	return rec
}

// ParseFlexibleDBTimestamp parses timestamps using fallback layouts including RFC3339, RFC3339Nano, and SQLite.
func ParseFlexibleDBTimestamp(s string) time.Time {
	return parseFlexibleDBTimestamp(s)
}

func parseFlexibleDBTimestamp(s string) time.Time {
	clean := strings.TrimSpace(s)
	if clean == "" {
		return time.Now().UTC()
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, clean); err == nil {
			return t.UTC()
		}
	}
	return time.Now().UTC()
}
