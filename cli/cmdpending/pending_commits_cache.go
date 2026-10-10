package cmdpending

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// PendingCommitCacheRecord stores cached git status metrics for a repository.
type PendingCommitCacheRecord struct {
	RepoPath         string
	HeadSHA          string
	UncommittedCount int
	UnpushedCount    int
	IsDirty          bool
	CachedAtUnix     int64
}

const (
	// DefaultCacheTTLSeconds defines the standard 90-second TTL (1.5 minutes).
	DefaultCacheTTLSeconds int64 = 90

	sqlCreatePendingCacheTable = `CREATE TABLE IF NOT EXISTS pending_commits_cache (
		repo_path          TEXT PRIMARY KEY,
		head_sha           TEXT NOT NULL,
		uncommitted_count  INTEGER NOT NULL,
		unpushed_count     INTEGER NOT NULL,
		is_dirty           INTEGER NOT NULL,
		cached_at_unix     INTEGER NOT NULL
	);`

	sqlCreatePendingCacheIndex = `CREATE INDEX IF NOT EXISTS idx_pending_cache_time ON pending_commits_cache (cached_at_unix);`

	sqlSelectPendingCache = `SELECT head_sha, uncommitted_count, unpushed_count, is_dirty, cached_at_unix FROM pending_commits_cache WHERE repo_path = ?;`

	sqlUpsertPendingCache = `INSERT INTO pending_commits_cache (repo_path, head_sha, uncommitted_count, unpushed_count, is_dirty, cached_at_unix)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(repo_path) DO UPDATE SET
	head_sha = excluded.head_sha,
	uncommitted_count = excluded.uncommitted_count,
	unpushed_count = excluded.unpushed_count,
	is_dirty = excluded.is_dirty,
	cached_at_unix = excluded.cached_at_unix;`

	sqlDeletePendingCacheByPath = `DELETE FROM pending_commits_cache WHERE repo_path = ?;`

	sqlPurgeExpiredPendingCache = `DELETE FROM pending_commits_cache WHERE cached_at_unix < ?;`
)

func resolvePendingCacheDBPath(dbPath string) string {
	if dbPath != "" {
		return dbPath
	}

	return filepath.Join(store.BinaryDataDir(), "pending_commits_cache.db")
}

func ensureCacheDir(path string) error {
	dir := filepath.Dir(path)
	if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
		return appfault.WrapSimple(mkErr, "create pending cache dir")
	}

	return nil
}

func initCacheSchema(conn *sql.DB) error {
	if _, err := conn.Exec(sqlCreatePendingCacheTable); err != nil {
		return appfault.WrapSimple(err, "create pending cache table")
	}
	if _, err := conn.Exec(sqlCreatePendingCacheIndex); err != nil {
		return appfault.WrapSimple(err, "create pending cache index")
	}

	return nil
}

func configureAndInitCache(conn *sql.DB) error {
	if err := store.ConfigureSQLiteConn(conn); err != nil {
		return appfault.WrapSimple(err, "configure sqlite conn")
	}

	return initCacheSchema(conn)
}

func openCacheConn(dbPath string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, appfault.WrapSimple(err, "open sqlite pending cache")
	}

	return conn, nil
}

func purgeExpiredOnInit(conn *sql.DB) {
	nowUnix := time.Now().Unix()
	_, _ = PurgeExpiredPendingCache(conn, nowUnix, DefaultCacheTTLSeconds)
}

// OpenPendingCommitsCache opens or initializes the cache SQLite database.
func OpenPendingCommitsCache(dbPath string) (*sql.DB, error) {
	resolved := resolvePendingCacheDBPath(dbPath)
	if err := ensureCacheDir(resolved); err != nil {
		return nil, err
	}
	conn, err := openCacheConn(resolved)
	if err != nil {
		return nil, err
	}
	if err := configureAndInitCache(conn); err != nil {
		_ = conn.Close()

		return nil, err
	}
	purgeExpiredOnInit(conn)

	return conn, nil
}

func isCacheExpired(nowUnix, cachedAtUnix int64) bool {
	age := nowUnix - cachedAtUnix
	if age > DefaultCacheTTLSeconds {
		return true
	}
	if age < 0 {
		return true
	}

	return false
}

func isHeadDiverged(currentHeadSHA, cachedHeadSHA string) bool {
	if currentHeadSHA == "" {
		return false
	}
	if currentHeadSHA != cachedHeadSHA {
		return true
	}

	return false
}

func isNonZero(val int) bool {
	return val > 0
}

func boolToInt(hasFlag bool) int {
	if hasFlag {
		return 1
	}

	return 0
}

func buildCacheRecord(repoPath, headSHA string, uncommitted, unpushed, isDirtyInt int, cachedAt int64) *PendingCommitCacheRecord {
	return &PendingCommitCacheRecord{
		RepoPath:         repoPath,
		HeadSHA:          headSHA,
		UncommittedCount: uncommitted,
		UnpushedCount:    unpushed,
		IsDirty:          isNonZero(isDirtyInt),
		CachedAtUnix:     cachedAt,
	}
}

func scanPendingCacheRow(row *sql.Row) (string, int, int, int, int64, error) {
	var headSHA string
	var uncommitted, unpushed, isDirtyInt int
	var cachedAt int64
	err := row.Scan(&headSHA, &uncommitted, &unpushed, &isDirtyInt, &cachedAt)

	return headSHA, uncommitted, unpushed, isDirtyInt, cachedAt, err
}

func queryPendingCacheRecord(conn *sql.DB, repoPath string) (*PendingCommitCacheRecord, bool, error) {
	head, uncom, unpush, isDirtyInt, at, err := scanPendingCacheRow(conn.QueryRow(sqlSelectPendingCache, repoPath))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, appfault.WrapSimple(err, "query cache record")
	}

	return buildCacheRecord(repoPath, head, uncom, unpush, isDirtyInt, at), true, nil
}

func evaluateCacheRecord(conn *sql.DB, repoPath, currentHeadSHA string, nowUnix int64, record *PendingCommitCacheRecord) (*PendingCommitCacheRecord, bool, error) {
	if isCacheExpired(nowUnix, record.CachedAtUnix) {
		_ = InvalidatePendingCache(conn, repoPath)

		return nil, false, nil
	}
	if isHeadDiverged(currentHeadSHA, record.HeadSHA) {
		_ = InvalidatePendingCache(conn, repoPath)

		return nil, false, nil
	}

	return record, true, nil
}

// GetCachedPendingStatus retrieves a valid, non-expired cache record for a repo.
// Returns (nil, false, nil) on cache miss or expired record (expired records are purged).
func GetCachedPendingStatus(conn *sql.DB, repoPath, currentHeadSHA string, nowUnix int64) (*PendingCommitCacheRecord, bool, error) {
	if conn == nil {
		return nil, false, nil
	}
	_, _ = PurgeExpiredPendingCache(conn, nowUnix, DefaultCacheTTLSeconds)
	record, hasRecord, err := queryPendingCacheRecord(conn, repoPath)
	if err != nil {
		return nil, false, err
	}
	if hasRecord {
		return evaluateCacheRecord(conn, repoPath, currentHeadSHA, nowUnix, record)
	}

	return nil, false, nil
}

func execUpsertCache(conn *sql.DB, rec PendingCommitCacheRecord, isDirtyInt int) error {
	_, err := conn.Exec(
		sqlUpsertPendingCache,
		rec.RepoPath, rec.HeadSHA, rec.UncommittedCount, rec.UnpushedCount, isDirtyInt, rec.CachedAtUnix,
	)
	if err != nil {
		return appfault.WrapSimple(err, "upsert pending cache status")
	}

	return nil
}

// SaveCachedPendingStatus writes or updates a repository status record in SQLite.
func SaveCachedPendingStatus(conn *sql.DB, record PendingCommitCacheRecord) error {
	if conn == nil {
		return nil
	}

	return execUpsertCache(conn, record, boolToInt(record.IsDirty))
}

// PurgeExpiredPendingCache removes all records older than ttlSeconds.
func PurgeExpiredPendingCache(conn *sql.DB, nowUnix, ttlSeconds int64) (int64, error) {
	if conn == nil {
		return 0, nil
	}
	cutoff := nowUnix - ttlSeconds
	res, err := conn.Exec(sqlPurgeExpiredPendingCache, cutoff)
	if err != nil {
		return 0, appfault.WrapSimple(err, "purge expired pending cache")
	}

	return res.RowsAffected()
}

// InvalidatePendingCache deletes an individual repository record by path.
func InvalidatePendingCache(conn *sql.DB, repoPath string) error {
	if conn == nil {
		return nil
	}
	if _, err := conn.Exec(sqlDeletePendingCacheByPath, repoPath); err != nil {
		return appfault.WrapSimple(err, "delete pending cache by path")
	}

	return nil
}

func prepareFullRecForBackup(record PendingCommitCacheRecord, fullRec RepoPendingCommitRecord) RepoPendingCommitRecord {
	if fullRec.RelativePath == "" && record.RepoPath != "" {
		fullRec.RelativePath = record.RepoPath
	}
	if fullRec.RepoName == "" && record.RepoPath != "" {
		fullRec.RepoName = record.RepoPath
	}

	return fullRec
}

// SaveStatusWithBackup writes the record to both the ephemeral cache DB and the durable backup DB.
func SaveStatusWithBackup(cacheConn, backupConn *sql.DB, record PendingCommitCacheRecord, fullRec RepoPendingCommitRecord, nowUnix int64) error {
	if err := SaveCachedPendingStatus(cacheConn, record); err != nil {
		return err
	}
	if backupConn == nil {
		return nil
	}
	adjustedRec := prepareFullRecForBackup(record, fullRec)

	return SaveBackupPendingStatus(backupConn, adjustedRec, record.HeadSHA, nowUnix)
}
