package cmdpending

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// PendingCommitBackupRecord represents a durable backup snapshot of a repository's pending commit status.
type PendingCommitBackupRecord struct {
	RepoPath           string `json:"repoPath"`
	RepoName           string `json:"repoName"`
	CurrentBranch      string `json:"currentBranch"`
	Version            string `json:"version"`
	ShortVersionBranch string `json:"shortVersionBranch"`
	HeadSHA            string `json:"headSha"`
	UncommittedCount   int    `json:"uncommittedCount"`
	UnpushedCount      int    `json:"unpushedCount"`
	IsDirty            bool   `json:"isDirty"`
	BackedUpAtUnix     int64  `json:"backedUpAtUnix"`
	PayloadJSON        string `json:"payloadJson,omitempty"`
}

type backupScanDest struct {
	repoPath     string
	repoName     string
	branch       string
	version      string
	shortBranch  string
	headSHA      string
	uncommitted  int
	unpushed     int
	isDirtyInt   int
	backedUpAt   int64
	payloadJSON  string
}

const (
	sqlCreatePendingBackupTable = `CREATE TABLE IF NOT EXISTS pending_commits_backup (
		repo_path            TEXT PRIMARY KEY,
		repo_name            TEXT NOT NULL,
		current_branch       TEXT NOT NULL,
		version              TEXT NOT NULL,
		short_version_branch TEXT NOT NULL,
		head_sha             TEXT NOT NULL,
		uncommitted_count    INTEGER NOT NULL,
		unpushed_count       INTEGER NOT NULL,
		is_dirty             INTEGER NOT NULL,
		backed_up_at_unix    INTEGER NOT NULL,
		payload_json         TEXT
	);`

	sqlCreatePendingBackupIndex = `CREATE INDEX IF NOT EXISTS idx_pending_backup_time ON pending_commits_backup (backed_up_at_unix);`

	sqlSelectPendingBackupByPath = `SELECT repo_path, repo_name, current_branch, version, short_version_branch, head_sha, uncommitted_count, unpushed_count, is_dirty, backed_up_at_unix, coalesce(payload_json, '') FROM pending_commits_backup WHERE repo_path = ?;`

	sqlSelectAllPendingBackups = `SELECT repo_path, repo_name, current_branch, version, short_version_branch, head_sha, uncommitted_count, unpushed_count, is_dirty, backed_up_at_unix, coalesce(payload_json, '') FROM pending_commits_backup ORDER BY backed_up_at_unix DESC;`

	sqlUpsertPendingBackup = `INSERT INTO pending_commits_backup (
		repo_path, repo_name, current_branch, version, short_version_branch,
		head_sha, uncommitted_count, unpushed_count, is_dirty, backed_up_at_unix, payload_json
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(repo_path) DO UPDATE SET
		repo_name = excluded.repo_name,
		current_branch = excluded.current_branch,
		version = excluded.version,
		short_version_branch = excluded.short_version_branch,
		head_sha = excluded.head_sha,
		uncommitted_count = excluded.uncommitted_count,
		unpushed_count = excluded.unpushed_count,
		is_dirty = excluded.is_dirty,
		backed_up_at_unix = excluded.backed_up_at_unix,
		payload_json = excluded.payload_json;`

	sqlPurgeBackupOlderThan = `DELETE FROM pending_commits_backup WHERE backed_up_at_unix < ?;`
)

func resolvePendingBackupDBPath(dbPath string) string {
	if dbPath != "" {
		return dbPath
	}

	return filepath.Join(store.BinaryDataDir(), "pending_commits_backup.db")
}

func ensureBackupDir(path string) error {
	dir := filepath.Dir(path)
	if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
		return appfault.WrapSimple(mkErr, "create pending backup dir")
	}

	return nil
}

func initBackupSchema(conn *sql.DB) error {
	if _, err := conn.Exec(sqlCreatePendingBackupTable); err != nil {
		return appfault.WrapSimple(err, "create pending backup table")
	}
	if _, err := conn.Exec(sqlCreatePendingBackupIndex); err != nil {
		return appfault.WrapSimple(err, "create pending backup index")
	}

	return nil
}

func configureAndInitBackup(conn *sql.DB) error {
	if err := store.ConfigureSQLiteConn(conn); err != nil {
		return appfault.WrapSimple(err, "configure sqlite conn")
	}

	return initBackupSchema(conn)
}

func openBackupConn(dbPath string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, appfault.WrapSimple(err, "open sqlite pending backup")
	}

	return conn, nil
}

// OpenPendingCommitsBackup opens or initializes the backup SQLite database.
func OpenPendingCommitsBackup(dbPath string) (*sql.DB, error) {
	resolved := resolvePendingBackupDBPath(dbPath)
	if err := ensureBackupDir(resolved); err != nil {
		return nil, err
	}
	conn, err := openBackupConn(resolved)
	if err != nil {
		return nil, err
	}
	if err := configureAndInitBackup(conn); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return conn, nil
}

func resolveRepoBackupPath(record RepoPendingCommitRecord) string {
	if record.RelativePath != "" {
		return record.RelativePath
	}

	return record.RepoName
}

func marshalBackupPayload(record RepoPendingCommitRecord) string {
	payloadBytes, err := json.Marshal(record)
	if err != nil {
		return "{}"
	}

	return string(payloadBytes)
}

func execUpsertBackup(conn *sql.DB, path string, rec RepoPendingCommitRecord, headSHA string, nowUnix int64, payload string) error {
	_, err := conn.Exec(
		sqlUpsertPendingBackup,
		path, rec.RepoName, rec.CurrentBranch, rec.Version, rec.ShortVersionBranch,
		headSHA, rec.TotalUncommitted, rec.UnpushedCommitsCount, boolToInt(rec.IsDirty),
		nowUnix, payload,
	)
	if err != nil {
		return appfault.WrapSimple(err, "upsert pending backup status")
	}

	return nil
}

// SaveBackupPendingStatus writes or updates a durable repository backup status record.
func SaveBackupPendingStatus(conn *sql.DB, record RepoPendingCommitRecord, headSHA string, nowUnix int64) error {
	if conn == nil {
		return nil
	}
	path := resolveRepoBackupPath(record)
	payload := marshalBackupPayload(record)

	return execUpsertBackup(conn, path, record, headSHA, nowUnix, payload)
}

func buildBackupRecord(d backupScanDest) *PendingCommitBackupRecord {
	return &PendingCommitBackupRecord{
		RepoPath:           d.repoPath,
		RepoName:           d.repoName,
		CurrentBranch:      d.branch,
		Version:            d.version,
		ShortVersionBranch: d.shortBranch,
		HeadSHA:            d.headSHA,
		UncommittedCount:   d.uncommitted,
		UnpushedCount:      d.unpushed,
		IsDirty:            isNonZero(d.isDirtyInt),
		BackedUpAtUnix:     d.backedUpAt,
		PayloadJSON:        d.payloadJSON,
	}
}

func scanBackupRow(row *sql.Row) (*PendingCommitBackupRecord, error) {
	var d backupScanDest
	err := row.Scan(
		&d.repoPath, &d.repoName, &d.branch, &d.version, &d.shortBranch,
		&d.headSHA, &d.uncommitted, &d.unpushed, &d.isDirtyInt, &d.backedUpAt, &d.payloadJSON,
	)
	if err != nil {
		return nil, err
	}

	return buildBackupRecord(d), nil
}

// GetBackupPendingStatus retrieves a backup status record by repository path.
func GetBackupPendingStatus(conn *sql.DB, repoPath string) (*PendingCommitBackupRecord, bool, error) {
	if conn == nil {
		return nil, false, nil
	}
	rec, err := scanBackupRow(conn.QueryRow(sqlSelectPendingBackupByPath, repoPath))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, appfault.WrapSimple(err, "query backup record")
	}

	return rec, true, nil
}

func scanBackupRows(rows *sql.Rows) ([]PendingCommitBackupRecord, error) {
	var results []PendingCommitBackupRecord
	for rows.Next() {
		var d backupScanDest
		if err := rows.Scan(
			&d.repoPath, &d.repoName, &d.branch, &d.version, &d.shortBranch,
			&d.headSHA, &d.uncommitted, &d.unpushed, &d.isDirtyInt, &d.backedUpAt, &d.payloadJSON,
		); err != nil {
			return nil, appfault.WrapSimple(err, "scan backup row")
		}
		results = append(results, *buildBackupRecord(d))
	}

	return results, rows.Err()
}

// ListAllBackupPendingStatuses returns all stored backup status records.
func ListAllBackupPendingStatuses(conn *sql.DB) ([]PendingCommitBackupRecord, error) {
	if conn == nil {
		return nil, nil
	}
	rows, err := conn.Query(sqlSelectAllPendingBackups)
	if err != nil {
		return nil, appfault.WrapSimple(err, "query all backup records")
	}
	defer rows.Close()

	return scanBackupRows(rows)
}

func convertBackupToCache(backupRec *PendingCommitBackupRecord, nowUnix int64) PendingCommitCacheRecord {
	return PendingCommitCacheRecord{
		RepoPath:         backupRec.RepoPath,
		HeadSHA:          backupRec.HeadSHA,
		UncommittedCount: backupRec.UncommittedCount,
		UnpushedCount:    backupRec.UnpushedCount,
		IsDirty:          backupRec.IsDirty,
		CachedAtUnix:     nowUnix,
	}
}

// RestoreBackupToCache copies a backed up record into the active cache database.
func RestoreBackupToCache(backupConn, cacheConn *sql.DB, repoPath string, nowUnix int64) error {
	if backupConn == nil || cacheConn == nil {
		return nil
	}
	rec, hasRecord, err := GetBackupPendingStatus(backupConn, repoPath)
	if err != nil {
		return err
	}
	if !hasRecord {
		return appfault.NewNotFound("restore backup to cache", "BACKUP_NOT_FOUND", "backup record not found for "+repoPath)
	}

	return SaveCachedPendingStatus(cacheConn, convertBackupToCache(rec, nowUnix))
}

// PurgeBackupOlderThan removes backup records older than cutoffUnix.
func PurgeBackupOlderThan(conn *sql.DB, cutoffUnix int64) (int64, error) {
	if conn == nil {
		return 0, nil
	}
	res, err := conn.Exec(sqlPurgeBackupOlderThan, cutoffUnix)
	if err != nil {
		return 0, appfault.WrapSimple(err, "purge backup older than cutoff")
	}

	return res.RowsAffected()
}
