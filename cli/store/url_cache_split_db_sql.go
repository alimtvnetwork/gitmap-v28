package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var urlCachePragmas = []string{
	"PRAGMA journal_mode = WAL;",
	"PRAGMA synchronous = NORMAL;",
	"PRAGMA busy_timeout = 5000;",
	"PRAGMA mmap_size = 268435456;",
	"PRAGMA cache_size = -64000;",
	"PRAGMA foreign_keys = ON;",
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

func accessLogToInsertArgs(log UrlAccessLog) []any {
	return []any{
		log.UrlCacheEntryId, log.UrlHash, log.RawUrl, boolToInt(log.IsCacheHit),
		log.StatusCode, log.BytesServed, log.DurationMs,
		log.CallerSource, log.Notes, log.Comments, log.CreatedAt,
	}
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
