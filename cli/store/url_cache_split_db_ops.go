package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

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

	entry, ok := val.(*UrlCacheEntry)
	if !ok || entry == nil {
		return nil, false, nil
	}

	return entry, true, nil
}

func (db *UrlCacheSplitDB) fetchEntryForGroup(urlHash string) (any, error) {
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
