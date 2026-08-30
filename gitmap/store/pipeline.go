package store

import (
	"database/sql"
	"time"
)

const (
	SQLCreatePipelineRunsTable = `
CREATE TABLE IF NOT EXISTS pipeline_runs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id INTEGER UNIQUE NOT NULL,
    repo TEXT NOT NULL,
    workflow_name TEXT NOT NULL,
    status TEXT NOT NULL,
    conclusion TEXT DEFAULT '',
    branch TEXT DEFAULT '',
    sha TEXT DEFAULT '',
    eta_seconds INTEGER DEFAULT 0,
    error_log TEXT DEFAULT '',
    url TEXT DEFAULT '',
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);`
)

// PipelineRun represents a CI/CD pipeline workflow execution record.
type PipelineRun struct {
	ID           int64  `json:"id"`
	RunID        int64  `json:"runId"`
	Repo         string `json:"repo"`
	WorkflowName string `json:"workflowName"`
	Status       string `json:"status"`
	Conclusion   string `json:"conclusion"`
	Branch       string `json:"branch"`
	Sha          string `json:"sha"`
	EtaSeconds   int    `json:"etaSeconds"`
	ErrorLog     string `json:"errorLog"`
	URL          string `json:"url"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// InitPipelineTable creates the pipeline_runs table if it doesn't exist.
func (db *DB) InitPipelineTable() error {
	_, err := db.conn.Exec(SQLCreatePipelineRunsTable)

	return err
}

// InsertOrUpdatePipelineRun persists a pipeline run record.
func (db *DB) InsertOrUpdatePipelineRun(run PipelineRun) error {
	_ = db.InitPipelineTable()

	query := `
INSERT INTO pipeline_runs (
    run_id, repo, workflow_name, status, conclusion, branch, sha, eta_seconds, error_log, url, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(run_id) DO UPDATE SET
    status = excluded.status,
    conclusion = excluded.conclusion,
    eta_seconds = excluded.eta_seconds,
    error_log = CASE WHEN excluded.error_log != '' THEN excluded.error_log ELSE pipeline_runs.error_log END,
    url = excluded.url,
    updated_at = excluded.updated_at;`

	nowStr := time.Now().UTC().Format(time.RFC3339)

	_, err := db.conn.Exec(
		query,
		run.RunID,
		run.Repo,
		run.WorkflowName,
		run.Status,
		run.Conclusion,
		run.Branch,
		run.Sha,
		run.EtaSeconds,
		run.ErrorLog,
		run.URL,
		nowStr,
	)

	return err
}

// GetLatestPipelineRun retrieves the most recent pipeline run for a repo.
func (db *DB) GetLatestPipelineRun(repo string) (*PipelineRun, error) {
	_ = db.InitPipelineTable()

	query := `
SELECT id, run_id, repo, workflow_name, status, conclusion, branch, sha, eta_seconds, error_log, url, created_at, updated_at
FROM pipeline_runs
WHERE repo = ? OR ? = ''
ORDER BY id DESC
LIMIT 1;`

	row := db.conn.QueryRow(query, repo, repo)

	return scanPipelineRun(row)
}

// GetLatestPipelineError retrieves the most recent failed pipeline run containing error logs.
func (db *DB) GetLatestPipelineError(repo string) (*PipelineRun, error) {
	_ = db.InitPipelineTable()

	query := `
SELECT id, run_id, repo, workflow_name, status, conclusion, branch, sha, eta_seconds, error_log, url, created_at, updated_at
FROM pipeline_runs
WHERE (repo = ? OR ? = '') AND (conclusion = 'failure' OR error_log != '')
ORDER BY id DESC
LIMIT 1;`

	row := db.conn.QueryRow(query, repo, repo)

	return scanPipelineRun(row)
}

// ListRecentPipelineRuns returns recent runs up to limit.
func (db *DB) ListRecentPipelineRuns(repo string, limit int) ([]PipelineRun, error) {
	_ = db.InitPipelineTable()

	query := `
SELECT id, run_id, repo, workflow_name, status, conclusion, branch, sha, eta_seconds, error_log, url, created_at, updated_at
FROM pipeline_runs
WHERE repo = ? OR ? = ''
ORDER BY id DESC
LIMIT ?;`

	if limit <= 0 {
		limit = 10
	}

	rows, err := db.conn.Query(query, repo, repo, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanPipelineRows(rows)
}

func scanPipelineRun(row *sql.Row) (*PipelineRun, error) {
	var run PipelineRun

	err := row.Scan(
		&run.ID,
		&run.RunID,
		&run.Repo,
		&run.WorkflowName,
		&run.Status,
		&run.Conclusion,
		&run.Branch,
		&run.Sha,
		&run.EtaSeconds,
		&run.ErrorLog,
		&run.URL,
		&run.CreatedAt,
		&run.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &run, nil
}

func scanPipelineRows(rows *sql.Rows) ([]PipelineRun, error) {
	var runs []PipelineRun

	for rows.Next() {
		var run PipelineRun

		err := rows.Scan(
			&run.ID,
			&run.RunID,
			&run.Repo,
			&run.WorkflowName,
			&run.Status,
			&run.Conclusion,
			&run.Branch,
			&run.Sha,
			&run.EtaSeconds,
			&run.ErrorLog,
			&run.URL,
			&run.CreatedAt,
			&run.UpdatedAt,
		)

		if err == nil {
			runs = append(runs, run)
		}
	}

	return runs, rows.Err()
}
