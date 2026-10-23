package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// UrlCacheSplitDB manages the dedicated split database for cached web responses.
type UrlCacheSplitDB struct {
	conn  *sql.DB
	group urlSingleflight
	mu    sync.RWMutex
	Path  string
}

// UrlCacheDB is a type alias for UrlCacheSplitDB.
type UrlCacheDB = UrlCacheSplitDB

type urlCall struct {
	wg  sync.WaitGroup
	val any
	err error
}

type urlSingleflight struct {
	mu    sync.Mutex
	calls map[string]*urlCall
}

func (g *urlSingleflight) Do(key string, fn func() (any, error)) (any, error) {
	c, isNew := g.registerCall(key)
	if !isNew {
		c.wg.Wait()

		return c.val, c.err
	}

	return g.executeCall(key, c, fn)
}

func (g *urlSingleflight) registerCall(key string) (*urlCall, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.ensureCallsMap()

	return g.getOrCreateCall(key)
}

func (g *urlSingleflight) ensureCallsMap() {
	if g.calls == nil {
		g.calls = make(map[string]*urlCall)
	}
}

func (g *urlSingleflight) getOrCreateCall(key string) (*urlCall, bool) {
	if c, hasExisting := g.calls[key]; hasExisting {
		return c, false
	}

	c := new(urlCall)
	c.wg.Add(1)
	g.calls[key] = c

	return c, true
}

func (g *urlSingleflight) executeCall(key string, c *urlCall, fn func() (any, error)) (any, error) {
	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.calls, key)
	g.mu.Unlock()

	return c.val, c.err
}

// HashUrl generates a deterministic SHA-256 hex string for a given URL.
func HashUrl(rawUrl string) string {
	sum := sha256.Sum256([]byte(rawUrl))

	return hex.EncodeToString(sum[:])
}

// OpenUrlCacheDB creates or opens the dedicated url_cache.db split database.
func OpenUrlCacheDB(repoRoot string) (*UrlCacheSplitDB, error) {
	return OpenUrlCacheSplitDB(repoRoot)
}

// OpenUrlCacheSplitDB creates or opens the dedicated url_cache.db split database.
func OpenUrlCacheSplitDB(repoRoot string) (*UrlCacheSplitDB, error) {
	dbPath := ResolveUrlCacheDbPath(repoRoot)

	return OpenUrlCacheSplitDBAt(dbPath)
}

// OpenUrlCacheSplitDBAt opens or creates the URL cache database at an explicit file path.
func OpenUrlCacheSplitDBAt(dbPath string) (*UrlCacheSplitDB, error) {
	conn, err := openConfiguredCacheConn(dbPath)
	if err != nil {
		return nil, err
	}

	db := &UrlCacheSplitDB{conn: conn, Path: dbPath}
	db.registerWithRegistry()

	return db, nil
}

func openConfiguredCacheConn(dbPath string) (*sql.DB, error) {
	conn, err := openUrlCacheConn(dbPath)
	if err != nil {
		return nil, err
	}

	if err := setupUrlCacheDB(conn); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return conn, nil
}

func openUrlCacheConn(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "store.OpenUrlCacheSplitDBAt.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.OpenUrlCacheSplitDBAt.open")
	}

	conn.SetMaxOpenConns(1)

	return conn, nil
}

func setupUrlCacheDB(conn *sql.DB) error {
	if err := initUrlCachePragmas(conn); err != nil {
		return err
	}

	return initUrlCacheSchema(conn)
}

var urlCachePragmas = []string{
	"PRAGMA journal_mode = WAL;",
	"PRAGMA synchronous = NORMAL;",
	"PRAGMA busy_timeout = 5000;",
	"PRAGMA mmap_size = 268435456;",
	"PRAGMA cache_size = -64000;",
	"PRAGMA foreign_keys = ON;",
}

func initUrlCachePragmas(conn *sql.DB) error {
	for _, pragma := range urlCachePragmas {
		if _, err := conn.Exec(pragma); err != nil {
			return apperror.WrapSimple(err, "store.initUrlCachePragmas")
		}
	}

	return nil
}

var urlCacheSchemaDDLs = []string{
	sqlCreateUrlCacheEntry,
	sqlIdxUrlCacheEntryHash,
	sqlIdxUrlCacheEntryExpires,
	sqlIdxUrlCacheEntryUpdated,
	sqlCreateUrlAccessLog,
	sqlIdxUrlAccessLogHash,
	sqlIdxUrlAccessLogCreated,
}

func initUrlCacheSchema(conn *sql.DB) error {
	for _, ddl := range urlCacheSchemaDDLs {
		if _, err := conn.Exec(ddl); err != nil {
			return apperror.WrapSimple(err, "store.initUrlCacheSchema")
		}
	}

	return nil
}

func (db *UrlCacheSplitDB) registerWithRegistry() {
	master, err := OpenDefault()
	if err != nil {
		return
	}

	defer master.Close()

	_ = master.RegisterSplitDB(makeRegistryEntry(db.Path))
}

func makeRegistryEntry(dbPath string) SplitDatabaseEntry {
	return SplitDatabaseEntry{
		DatabaseType:  "repodb",
		DatabaseKey:   "url_cache",
		DatabasePath:  dbPath,
		Status:        "active",
		IsActive:      true,
		Description:   "HTTP/HTTPS URL payload cache and access log",
		SchemaVersion: 1,
	}
}

// Close terminates the database connection.
func (db *UrlCacheSplitDB) Close() error {
	if db == nil || db.conn == nil {
		return nil
	}

	return db.conn.Close()
}

// GetByHash retrieves a cache entry by its SHA-256 hash.
func (db *UrlCacheSplitDB) GetByHash(urlHash string) (*UrlCacheEntry, bool, error) {
	row := db.conn.QueryRow(sqlSelectUrlCacheEntryByHash, urlHash)
	entry, err := scanUrlCacheEntryRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, apperror.WrapSimple(err, "store.GetByHash")
	}

	return entry, true, nil
}

// GetUrlCache retrieves a cache entry by raw URL, using singleflight deduplication.
func (db *UrlCacheSplitDB) GetUrlCache(rawUrl string) (*UrlCacheEntry, bool, error) {
	urlHash := HashUrl(rawUrl)
	val, err := db.group.Do(urlHash, func() (any, error) {
		return db.fetchEntryForGroup(urlHash)
	})

	if err != nil || val == nil {
		return nil, false, err
	}

	return val.(*UrlCacheEntry), true, nil
}

func (db *UrlCacheSplitDB) fetchEntryForGroup(urlHash string) (*UrlCacheEntry, error) {
	entry, hasHit, err := db.GetByHash(urlHash)
	if err != nil || !hasHit {
		return nil, err
	}

	return entry, nil
}

// UpsertEntry inserts or updates a cache entry.
func (db *UrlCacheSplitDB) UpsertEntry(entry UrlCacheEntry) error {
	entry = normalizeCacheEntry(entry)
	args := urlCacheEntryToUpsertArgs(entry)
	_, err := db.conn.Exec(sqlUpsertUrlCacheEntry, args...)
	if err != nil {
		return apperror.WrapSimple(err, "store.UpsertEntry")
	}

	return nil
}

// SaveUrlCache is a convenience wrapper for UpsertEntry.
func (db *UrlCacheSplitDB) SaveUrlCache(entry UrlCacheEntry) error {
	return db.UpsertEntry(entry)
}

func normalizeCacheEntry(entry UrlCacheEntry) UrlCacheEntry {
	now := time.Now().Unix()
	entry = normalizeUrlsAndHash(entry)
	entry = normalizeHttpMetadata(entry)
	entry = normalizeCachePayload(entry)

	return normalizeCacheTimestamps(entry, now)
}

func normalizeUrlsAndHash(entry UrlCacheEntry) UrlCacheEntry {
	entry = normalizeUrlFields(entry)
	if entry.UrlHash == "" {
		entry.UrlHash = HashUrl(entry.EffectiveURL())
	}

	return entry
}

func normalizeUrlFields(entry UrlCacheEntry) UrlCacheEntry {
	if entry.RawUrl == "" && entry.UrlText != "" {
		entry.RawUrl = entry.UrlText
	}

	if entry.UrlText == "" && entry.RawUrl != "" {
		entry.UrlText = entry.RawUrl
	}

	return entry
}

func normalizeHttpMetadata(entry UrlCacheEntry) UrlCacheEntry {
	entry = normalizeMethodAndStatus(entry)
	if entry.ContentType == "" {
		entry.ContentType = "application/octet-stream"
	}

	return entry
}

func normalizeMethodAndStatus(entry UrlCacheEntry) UrlCacheEntry {
	if entry.HttpMethod == "" {
		entry.HttpMethod = "GET"
	}

	if entry.StatusCode == 0 {
		entry.StatusCode = 200
	}

	return entry
}

func normalizeCachePayload(entry UrlCacheEntry) UrlCacheEntry {
	if entry.ContentLength == 0 && len(entry.PayloadBlob) > 0 {
		entry.ContentLength = int64(len(entry.PayloadBlob))
	}

	entry = normalizeCompressionAndHeaders(entry)

	return entry
}

func normalizeCompressionAndHeaders(entry UrlCacheEntry) UrlCacheEntry {
	if entry.CompressionType == "" {
		entry.CompressionType = "none"
	}

	if entry.HeadersJson == "" {
		entry.HeadersJson = "{}"
	}

	return entry
}

func normalizeCacheTimestamps(entry UrlCacheEntry, now int64) UrlCacheEntry {
	entry = normalizeAccessAndCached(entry, now)
	entry = normalizeLastAccessedAndCreated(entry, now)
	entry.UpdatedAt = now

	return entry
}

func normalizeAccessAndCached(entry UrlCacheEntry, now int64) UrlCacheEntry {
	if entry.AccessCount <= 0 {
		entry.AccessCount = 1
	}

	if entry.CachedAt <= 0 {
		entry.CachedAt = now
	}

	return entry
}

func normalizeLastAccessedAndCreated(entry UrlCacheEntry, now int64) UrlCacheEntry {
	if entry.LastAccessedAt <= 0 {
		entry.LastAccessedAt = now
	}

	if entry.CreatedAt <= 0 {
		entry.CreatedAt = now
	}

	return entry
}

func urlCacheEntryToUpsertArgs(e UrlCacheEntry) []any {
	compInt := boolToInt(e.IsCompressed)
	actInt := activeToInt(e.IsActive, e.UrlCacheEntryId)

	return []any{
		e.UrlHash, e.RawUrl, e.HttpMethod, e.StatusCode, e.ETag, e.LastModified,
		e.ContentType, e.ContentLength, e.PayloadBlob, e.CompressionType, compInt,
		e.HeadersJson, e.DurationMs, e.AccessCount, e.CachedAt, e.ExpiresAt,
		e.LastAccessedAt, actInt, e.Description, e.Notes, e.Comments,
		e.CreatedAt, e.UpdatedAt,
	}
}

func activeToInt(isActive bool, entryId int64) int {
	if !isActive && entryId > 0 {
		return 0
	}

	return 1
}

// TouchAccess increments the access count and updates the last accessed timestamp.
func (db *UrlCacheSplitDB) TouchAccess(id int64) error {
	now := time.Now().Unix()
	_, err := db.conn.Exec(sqlTouchAccessUrlCacheEntry, now, now, id)
	if err != nil {
		return apperror.WrapSimple(err, "store.TouchAccess")
	}

	return nil
}

// LogUrlAccess writes an audit record to UrlAccessLog.
func (db *UrlCacheSplitDB) LogUrlAccess(log UrlAccessLog) error {
	log = normalizeAccessLog(log)
	args := accessLogToInsertArgs(log)
	_, err := db.conn.Exec(sqlInsertUrlAccessLog, args...)
	if err != nil {
		return apperror.WrapSimple(err, "store.LogUrlAccess")
	}

	return nil
}

func accessLogToInsertArgs(log UrlAccessLog) []any {
	return []any{
		log.UrlCacheEntryId, log.UrlHash, log.RawUrl, boolToInt(log.IsCacheHit),
		log.StatusCode, log.BytesServed, log.DurationMs,
		log.CallerSource, log.Notes, log.Comments, log.CreatedAt,
	}
}

func normalizeAccessLog(log UrlAccessLog) UrlAccessLog {
	log = normalizeLogTimeAndStatus(log)
	if log.CallerSource == "" {
		log.CallerSource = "cli"
	}

	return log
}

func normalizeLogTimeAndStatus(log UrlAccessLog) UrlAccessLog {
	if log.CreatedAt <= 0 {
		log.CreatedAt = time.Now().Unix()
	}

	if log.StatusCode == 0 {
		log.StatusCode = 200
	}

	return log
}

// PurgeExpired evicts records with ExpiresAt > 0 and ExpiresAt <= now.
func (db *UrlCacheSplitDB) PurgeExpired(now int64) (int64, error) {
	res, err := db.conn.Exec(sqlPurgeExpiredEntries, now)
	if err != nil {
		return 0, apperror.WrapSimple(err, "store.PurgeExpired")
	}

	rowsAffected, _ := res.RowsAffected()

	return rowsAffected, nil
}

// PurgeExpiredUrlCache purges records expired as of the current unix epoch.
func (db *UrlCacheSplitDB) PurgeExpiredUrlCache() (int64, error) {
	return db.PurgeExpired(time.Now().Unix())
}

// DeleteByHash removes a cache entry by its hash.
func (db *UrlCacheSplitDB) DeleteByHash(urlHash string) error {
	_, err := db.conn.Exec(sqlDeleteUrlCacheEntryByHash, urlHash)
	if err != nil {
		return apperror.WrapSimple(err, "store.DeleteByHash")
	}

	return nil
}

// DeleteUrlCache removes a cache entry by its raw URL.
func (db *UrlCacheSplitDB) DeleteUrlCache(rawUrl string) error {
	return db.DeleteByHash(HashUrl(rawUrl))
}

// ListEntries returns a list of cached entries ordered by update timestamp descending.
func (db *UrlCacheSplitDB) ListEntries(limit int) ([]UrlCacheEntry, error) {
	return db.ListUrlCache(limit, 0)
}

// ListUrlCache queries cached entries with pagination limits.
func (db *UrlCacheSplitDB) ListUrlCache(limit, offset int) ([]UrlCacheEntry, error) {
	limit, offset = normalizePagination(limit, offset)
	rows, err := db.conn.Query(sqlListUrlCacheEntries, limit, offset)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.ListUrlCache")
	}

	defer rows.Close()

	return scanUrlCacheEntryRows(rows)
}

func normalizePagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 50
	}

	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

type rawCacheFields struct {
	e             UrlCacheEntry
	compressedInt int
	activeInt     int
	etag          sql.NullString
	lastMod       sql.NullString
}

func (f *rawCacheFields) destPointers() []any {
	return []any{
		&f.e.UrlCacheEntryId, &f.e.UrlHash, &f.e.RawUrl, &f.e.HttpMethod, &f.e.StatusCode,
		&f.etag, &f.lastMod, &f.e.ContentType, &f.e.ContentLength, &f.e.PayloadBlob,
		&f.e.CompressionType, &f.compressedInt, &f.e.HeadersJson, &f.e.DurationMs,
		&f.e.AccessCount, &f.e.CachedAt, &f.e.ExpiresAt, &f.e.LastAccessedAt,
		&f.activeInt, &f.e.Description, &f.e.Notes, &f.e.Comments, &f.e.CreatedAt, &f.e.UpdatedAt,
	}
}

func (f *rawCacheFields) finalize() *UrlCacheEntry {
	f.e.ETag = f.etag.String
	f.e.LastModified = f.lastMod.String
	f.e.IsCompressed = f.compressedInt == 1
	f.e.IsActive = f.activeInt == 1
	f.e.UrlText = f.e.RawUrl

	return &f.e
}

func scanUrlCacheEntryRow(row *sql.Row) (*UrlCacheEntry, error) {
	var f rawCacheFields
	if err := row.Scan(f.destPointers()...); err != nil {
		return nil, err
	}

	return f.finalize(), nil
}

type rawListFields struct {
	e             UrlCacheEntry
	compressedInt int
	activeInt     int
	etag          sql.NullString
	lastMod       sql.NullString
}

func (f *rawListFields) destPointers() []any {
	return []any{
		&f.e.UrlCacheEntryId, &f.e.UrlHash, &f.e.RawUrl, &f.e.HttpMethod, &f.e.StatusCode,
		&f.etag, &f.lastMod, &f.e.ContentType, &f.e.ContentLength,
		&f.e.CompressionType, &f.compressedInt, &f.e.HeadersJson, &f.e.DurationMs,
		&f.e.AccessCount, &f.e.CachedAt, &f.e.ExpiresAt, &f.e.LastAccessedAt,
		&f.activeInt, &f.e.Description, &f.e.Notes, &f.e.Comments, &f.e.CreatedAt, &f.e.UpdatedAt,
	}
}

func (f *rawListFields) finalize() UrlCacheEntry {
	f.e.ETag = f.etag.String
	f.e.LastModified = f.lastMod.String
	f.e.IsCompressed = f.compressedInt == 1
	f.e.IsActive = f.activeInt == 1
	f.e.UrlText = f.e.RawUrl

	return f.e
}

func scanUrlCacheEntryRows(rows *sql.Rows) ([]UrlCacheEntry, error) {
	entries := make([]UrlCacheEntry, 0, 16)
	for rows.Next() {
		var f rawListFields
		if err := rows.Scan(f.destPointers()...); err != nil {
			return nil, apperror.WrapSimple(err, "store.scanUrlCacheEntryRows")
		}

		entries = append(entries, f.finalize())
	}

	return entries, rows.Err()
}

const (
	sqlCreateUrlCacheEntry = `CREATE TABLE IF NOT EXISTS UrlCacheEntry (
    UrlCacheEntryId INTEGER PRIMARY KEY AUTOINCREMENT,
    UrlHash         TEXT NOT NULL UNIQUE,
    RawUrl          TEXT NOT NULL,
    HttpMethod      TEXT NOT NULL DEFAULT 'GET',
    StatusCode      INTEGER NOT NULL DEFAULT 200,
    ETag            TEXT NULL,
    LastModified    TEXT NULL,
    ContentType     TEXT NOT NULL DEFAULT 'application/octet-stream',
    ContentLength   INTEGER NOT NULL DEFAULT 0,
    PayloadBlob     BLOB NOT NULL,
    CompressionType TEXT NOT NULL DEFAULT 'none',
    IsCompressed    INTEGER NOT NULL DEFAULT 0,
    HeadersJson     TEXT NOT NULL DEFAULT '{}',
    DurationMs      INTEGER NOT NULL DEFAULT 0,
    AccessCount     INTEGER NOT NULL DEFAULT 1,
    CachedAt        INTEGER NOT NULL DEFAULT 0,
    ExpiresAt       INTEGER NOT NULL DEFAULT 0,
    LastAccessedAt  INTEGER NOT NULL DEFAULT 0,
    IsActive        INTEGER NOT NULL DEFAULT 1,
    Description     TEXT NOT NULL DEFAULT '',
    Notes           TEXT NOT NULL DEFAULT '',
    Comments        TEXT NOT NULL DEFAULT '',
    CreatedAt       INTEGER NOT NULL DEFAULT 0,
    UpdatedAt       INTEGER NOT NULL DEFAULT 0
);`

	sqlIdxUrlCacheEntryHash    = `CREATE INDEX IF NOT EXISTS IdxUrlCacheEntry_UrlHash ON UrlCacheEntry(UrlHash);`
	sqlIdxUrlCacheEntryExpires = `CREATE INDEX IF NOT EXISTS IdxUrlCacheEntry_ExpiresAt ON UrlCacheEntry(ExpiresAt);`
	sqlIdxUrlCacheEntryUpdated = `CREATE INDEX IF NOT EXISTS IdxUrlCacheEntry_UpdatedAt ON UrlCacheEntry(UpdatedAt);`

	sqlCreateUrlAccessLog = `CREATE TABLE IF NOT EXISTS UrlAccessLog (
    UrlAccessLogId  INTEGER PRIMARY KEY AUTOINCREMENT,
    UrlCacheEntryId INTEGER NULL,
    UrlHash         TEXT NOT NULL,
    RawUrl          TEXT NOT NULL DEFAULT '',
    IsCacheHit      INTEGER NOT NULL DEFAULT 0,
    StatusCode      INTEGER NOT NULL DEFAULT 200,
    BytesServed     INTEGER NOT NULL DEFAULT 0,
    DurationMs      INTEGER NOT NULL DEFAULT 0,
    CallerSource    TEXT NOT NULL DEFAULT 'cli',
    Notes           TEXT NOT NULL DEFAULT '',
    Comments        TEXT NOT NULL DEFAULT '',
    CreatedAt       INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY(UrlHash) REFERENCES UrlCacheEntry(UrlHash) ON DELETE CASCADE
);`

	sqlIdxUrlAccessLogHash    = `CREATE INDEX IF NOT EXISTS IdxUrlAccessLog_UrlHash ON UrlAccessLog(UrlHash);`
	sqlIdxUrlAccessLogCreated = `CREATE INDEX IF NOT EXISTS IdxUrlAccessLog_CreatedAt ON UrlAccessLog(CreatedAt);`

	sqlSelectUrlCacheEntryByHash = `SELECT
    UrlCacheEntryId, UrlHash, RawUrl, HttpMethod, StatusCode, ETag,
    LastModified, ContentType, ContentLength, PayloadBlob, CompressionType,
    IsCompressed, HeadersJson, DurationMs, AccessCount, CachedAt, ExpiresAt,
    LastAccessedAt, IsActive, Description, Notes, Comments, CreatedAt, UpdatedAt
FROM UrlCacheEntry WHERE UrlHash = ? LIMIT 1;`

	sqlUpsertUrlCacheEntry = `INSERT INTO UrlCacheEntry (
    UrlHash, RawUrl, HttpMethod, StatusCode, ETag, LastModified, ContentType,
    ContentLength, PayloadBlob, CompressionType, IsCompressed, HeadersJson,
    DurationMs, AccessCount, CachedAt, ExpiresAt, LastAccessedAt, IsActive,
    Description, Notes, Comments, CreatedAt, UpdatedAt
) VALUES (
    ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?
) ON CONFLICT(UrlHash) DO UPDATE SET
    RawUrl = excluded.RawUrl,
    HttpMethod = excluded.HttpMethod,
    StatusCode = excluded.StatusCode,
    ETag = excluded.ETag,
    LastModified = excluded.LastModified,
    ContentType = excluded.ContentType,
    ContentLength = excluded.ContentLength,
    PayloadBlob = excluded.PayloadBlob,
    CompressionType = excluded.CompressionType,
    IsCompressed = excluded.IsCompressed,
    HeadersJson = excluded.HeadersJson,
    DurationMs = excluded.DurationMs,
    AccessCount = UrlCacheEntry.AccessCount + 1,
    ExpiresAt = excluded.ExpiresAt,
    LastAccessedAt = excluded.LastAccessedAt,
    IsActive = excluded.IsActive,
    UpdatedAt = excluded.UpdatedAt;`

	sqlTouchAccessUrlCacheEntry = `UPDATE UrlCacheEntry SET
    AccessCount = AccessCount + 1,
    LastAccessedAt = ?,
    UpdatedAt = ?
WHERE UrlCacheEntryId = ?;`

	sqlInsertUrlAccessLog = `INSERT INTO UrlAccessLog (
    UrlCacheEntryId, UrlHash, RawUrl, IsCacheHit, StatusCode, BytesServed,
    DurationMs, CallerSource, Notes, Comments, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	sqlPurgeExpiredEntries = `DELETE FROM UrlCacheEntry WHERE ExpiresAt > 0 AND ExpiresAt <= ?;`

	sqlDeleteUrlCacheEntryByHash = `DELETE FROM UrlCacheEntry WHERE UrlHash = ?;`

	sqlListUrlCacheEntries = `SELECT
    UrlCacheEntryId, UrlHash, RawUrl, HttpMethod, StatusCode, ETag,
    LastModified, ContentType, ContentLength,
    CompressionType, IsCompressed, HeadersJson, DurationMs, AccessCount,
    CachedAt, ExpiresAt, LastAccessedAt, IsActive, Description, Notes,
    Comments, CreatedAt, UpdatedAt
FROM UrlCacheEntry ORDER BY UpdatedAt DESC LIMIT ? OFFSET ?;`
)
