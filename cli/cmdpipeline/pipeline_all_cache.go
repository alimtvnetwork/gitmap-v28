package cmdpipeline

import (
	"database/sql"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// SQLite hash cache for `pipeline errors all` (program 265, WS7).
// Before the per-repo inspection, the repo's HEAD commit SHA is checked
// against this cache. On a hit the expensive inspection is skipped and
// the repo's previous inspection state is restored.
// On a miss the repo is inspected and the result is recorded with the SHA.

const sqlCreatePeAllCacheTable = `CREATE TABLE IF NOT EXISTS pe_all_commit_cache (
  repo_slug          TEXT PRIMARY KEY,
  commit_sha         TEXT NOT NULL,
  checked_at         TEXT NOT NULL,
  is_clean           INTEGER NOT NULL DEFAULT 0,
  has_pipeline_data  INTEGER NOT NULL DEFAULT 0,
  data_status        TEXT NOT NULL DEFAULT '',
  failed_item_json   TEXT NOT NULL DEFAULT '',
  error_records_json TEXT NOT NULL DEFAULT '',
  step_jobs_json     TEXT NOT NULL DEFAULT ''
);`

func peAllCacheDbPath() string {
	return filepath.Join(store.BinaryDataDir(), "pe_all_commit_cache.db")
}

func openPeAllCache() (*sql.DB, error) {
	dbPath := peAllCacheDbPath()
	if mkdirErr := os.MkdirAll(filepath.Dir(dbPath), 0755); mkdirErr != nil {
		return nil, mkdirErr
	}
	conn, openErr := sql.Open("sqlite", dbPath)
	if openErr != nil {
		return nil, openErr
	}
	if _, execErr := conn.Exec(sqlCreatePeAllCacheTable); execErr != nil {
		_ = conn.Close()

		return nil, execErr
	}
	migratePeAllCacheTable(conn)

	return conn, nil
}

func migratePeAllCacheTable(conn *sql.DB) {
	cols := []string{
		"ALTER TABLE pe_all_commit_cache ADD COLUMN is_clean INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE pe_all_commit_cache ADD COLUMN has_pipeline_data INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE pe_all_commit_cache ADD COLUMN data_status TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE pe_all_commit_cache ADD COLUMN failed_item_json TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE pe_all_commit_cache ADD COLUMN error_records_json TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE pe_all_commit_cache ADD COLUMN step_jobs_json TEXT NOT NULL DEFAULT ''",
	}
	for _, col := range cols {
		_, _ = conn.Exec(col) // lint-allow: ignore-db-error
	}
}

func getCachedInspection(conn *sql.DB, slug, currentHeadSha string) (singleRepoInspection, bool) {
	var (
		sha              string
		isCleanInt       int
		hasDataInt       int
		dataStatus       string
		failedItemJSON   string
		errorRecordsJSON string
		stepJobsJSON     string
	)
	row := conn.QueryRow(`
		SELECT commit_sha, is_clean, has_pipeline_data, data_status,
		       failed_item_json, error_records_json, step_jobs_json
		FROM pe_all_commit_cache
		WHERE repo_slug = ?`, slug)
	err := row.Scan(&sha, &isCleanInt, &hasDataInt, &dataStatus,
		&failedItemJSON, &errorRecordsJSON, &stepJobsJSON)
	if err != nil || sha != currentHeadSha {
		return singleRepoInspection{}, false
	}

	ins := singleRepoInspection{
		isCached:        true,
		isClean:         isCleanInt == 1,
		hasPipelineData: hasDataInt == 1,
		dataStatus:      dataStatus,
	}

	if ins.dataStatus == "" {
		ins.dataStatus = resolveCacheDataStatus(ins.isClean, ins.hasPipelineData)
	}

	ins.failedItem = parseCachedFailedItem(failedItemJSON)
	ins.errorRecords = parseCachedErrorRecords(errorRecordsJSON)
	ins.stepJobs = parseCachedStepJobs(stepJobsJSON)

	if ins.dataStatus == allRepoStatusFailed && len(ins.errorRecords) == 0 {
		return singleRepoInspection{}, false
	}

	return ins, true
}

func resolveCacheDataStatus(isClean, hasPipelineData bool) string {
	if isClean {
		return allRepoStatusClean
	}
	if hasPipelineData {
		return allRepoStatusFailed
	}

	return allRepoStatusNoData
}

func parseCachedFailedItem(raw string) *AllPipelineFailedItem {
	if len(raw) == 0 {
		return nil
	}
	var item AllPipelineFailedItem
	if err := json.Unmarshal([]byte(raw), &item); err != nil {
		return nil
	}

	return &item
}

func parseCachedErrorRecords(raw string) []pipelinedb.PipelineErrorRecord {
	if len(raw) == 0 {
		return nil
	}
	var recs []pipelinedb.PipelineErrorRecord
	if err := json.Unmarshal([]byte(raw), &recs); err != nil {
		return nil
	}

	return recs
}

func parseCachedStepJobs(raw string) map[string]string {
	if len(raw) == 0 {
		return nil
	}
	var jobs map[string]string
	if err := json.Unmarshal([]byte(raw), &jobs); err != nil {
		return nil
	}

	return jobs
}

func marshalCachedJsonSafe(v any) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(b)
}

func recordCachedInspection(conn *sql.DB, slug, sha string, ins singleRepoInspection) error {
	failedItemJSON := marshalCachedJsonSafe(ins.failedItem)
	errorRecordsJSON := marshalCachedJsonSafe(ins.errorRecords)
	stepJobsJSON := marshalCachedJsonSafe(ins.stepJobs)

	isCleanInt := 0
	if ins.isClean {
		isCleanInt = 1
	}
	hasDataInt := 0
	if ins.hasPipelineData {
		hasDataInt = 1
	}

	_, err := conn.Exec(
		`INSERT INTO pe_all_commit_cache (
			repo_slug, commit_sha, checked_at, is_clean, has_pipeline_data,
			data_status, failed_item_json, error_records_json, step_jobs_json
		 ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(repo_slug) DO UPDATE SET
			commit_sha = excluded.commit_sha,
			checked_at = excluded.checked_at,
			is_clean = excluded.is_clean,
			has_pipeline_data = excluded.has_pipeline_data,
			data_status = excluded.data_status,
			failed_item_json = excluded.failed_item_json,
			error_records_json = excluded.error_records_json,
			step_jobs_json = excluded.step_jobs_json`,
		slug, sha, time.Now().UTC().Format(time.RFC3339),
		isCleanInt, hasDataInt, ins.dataStatus,
		failedItemJSON, errorRecordsJSON, stepJobsJSON,
	)
	return err
}

// getRepoHeadSha returns the HEAD commit SHA of the repo at absPath,
// or "" when the path is not a git checkout.
func getRepoHeadSha(absPath string) string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = absPath
	out, runErr := cmd.Output()
	if runErr != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// ClearPeAllCache removes all rows from the pe_all_commit_cache SQLite table.
func ClearPeAllCache() error {
	conn, err := openPeAllCache()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec("DELETE FROM pe_all_commit_cache")
	return err
}

// ClearPeAllCacheForRepo removes the cached inspection for a specific repo slug.
func ClearPeAllCacheForRepo(slug string) error {
	conn, err := openPeAllCache()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec("DELETE FROM pe_all_commit_cache WHERE repo_slug = ?", slug)
	return err
}
