package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	// SQLCreateReleaseCache defines the DDL for the ReleaseCache table in installation.db or gitmap.db.
	SQLCreateReleaseCache = `CREATE TABLE IF NOT EXISTS ReleaseCache (
    ReleaseCacheId   INTEGER PRIMARY KEY AUTOINCREMENT,
    Tag              TEXT NOT NULL UNIQUE,
    Version          TEXT NOT NULL,
    HasExecutables   INTEGER NOT NULL DEFAULT 0,
    AssetUrl         TEXT NOT NULL DEFAULT '',
    AssetSize        INTEGER NOT NULL DEFAULT 0,
    Platform         TEXT NOT NULL DEFAULT '',
    Arch             TEXT NOT NULL DEFAULT '',
    ChecksumUrl      TEXT NOT NULL DEFAULT '',
    CheckedAt        TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ExpiresAt        TEXT NOT NULL,
    CreatedAt        TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt        TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS IdxReleaseCache_Tag ON ReleaseCache(Tag);
CREATE INDEX IF NOT EXISTS IdxReleaseCache_Platform_Arch ON ReleaseCache(Platform, Arch);
CREATE INDEX IF NOT EXISTS IdxReleaseCache_ExpiresAt ON ReleaseCache(ExpiresAt);`

	sqlSelectReleaseCacheByTag = `SELECT ReleaseCacheId, Tag, Version, HasExecutables, AssetUrl, AssetSize, Platform, Arch, ChecksumUrl, CheckedAt, ExpiresAt, CreatedAt, UpdatedAt
FROM ReleaseCache
WHERE Tag = ? AND (Platform = ? OR Platform = '' OR ? = '') AND (Arch = ? OR Arch = '' OR ? = '') AND datetime(ExpiresAt) > datetime('now')
ORDER BY ReleaseCacheId DESC LIMIT 1;`

	sqlUpsertReleaseCache = `INSERT INTO ReleaseCache (
    Tag, Version, HasExecutables, AssetUrl, AssetSize, Platform, Arch, ChecksumUrl, CheckedAt, ExpiresAt, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT(Tag) DO UPDATE SET
    Version = excluded.Version,
    HasExecutables = excluded.HasExecutables,
    AssetUrl = excluded.AssetUrl,
    AssetSize = excluded.AssetSize,
    Platform = excluded.Platform,
    Arch = excluded.Arch,
    ChecksumUrl = excluded.ChecksumUrl,
    CheckedAt = excluded.CheckedAt,
    ExpiresAt = excluded.ExpiresAt,
    UpdatedAt = CURRENT_TIMESTAMP;`

	sqlListCachedReleases = `SELECT ReleaseCacheId, Tag, Version, HasExecutables, AssetUrl, AssetSize, Platform, Arch, ChecksumUrl, CheckedAt, ExpiresAt, CreatedAt, UpdatedAt
FROM ReleaseCache
ORDER BY ReleaseCacheId DESC LIMIT ?;`
)

// ReleaseCacheRecord mirrors the Split SQLite ReleaseCache table.
type ReleaseCacheRecord struct {
	ReleaseCacheId int64     `json:"releaseCacheId"`
	Tag            string    `json:"tag"`
	Version        string    `json:"version"`
	HasExecutables bool      `json:"hasExecutables"`
	AssetURL       string    `json:"assetUrl"`
	AssetSize      int64     `json:"assetSize"`
	Platform       string    `json:"platform"`
	Arch           string    `json:"arch"`
	ChecksumURL    string    `json:"checksumUrl"`
	CheckedAt      time.Time `json:"checkedAt"`
	ExpiresAt      time.Time `json:"expiresAt"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// InitReleaseCacheSchema initializes the ReleaseCache table and indexes on the given SQLite connection.
func InitReleaseCacheSchema(db *sql.DB) error {
	if _, err := db.Exec(SQLCreateReleaseCache); err != nil {
		return apperror.WrapSimple(err, "store.init_release_cache")
	}

	return nil
}

// OpenReleaseCacheDB opens or creates the installation split database configured with ReleaseCache.
func OpenReleaseCacheDB() (*sql.DB, error) {
	path := InstallationDBPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "release_cache.mkdir")
	}

	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, apperror.WrapSimple(err, "release_cache.open")
	}

	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()

		return nil, apperror.WrapSimple(err, "release_cache.config")
	}

	if err := InitReleaseCacheSchema(conn); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return conn, nil
}

// GetCachedRelease checks if a valid unexpired cache entry exists for the given tag, platform, and arch.
func GetCachedRelease(db *sql.DB, tag, platform, arch string) (*ReleaseCacheRecord, error) {
	row := db.QueryRow(sqlSelectReleaseCacheByTag, tag, platform, platform, arch, arch)

	var (
		rec        ReleaseCacheRecord
		hasExec    int
		checkedStr string
		expireStr  string
		createStr  string
		updateStr  string
	)

	err := row.Scan(
		&rec.ReleaseCacheId,
		&rec.Tag,
		&rec.Version,
		&hasExec,
		&rec.AssetURL,
		&rec.AssetSize,
		&rec.Platform,
		&rec.Arch,
		&rec.ChecksumURL,
		&checkedStr,
		&expireStr,
		&createStr,
		&updateStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, apperror.WrapSimple(err, "store.get_cached_release")
	}

	rec.HasExecutables = (hasExec == 1)
	rec.CheckedAt = parseSQLiteTime(checkedStr)
	rec.ExpiresAt = parseSQLiteTime(expireStr)
	rec.CreatedAt = parseSQLiteTime(createStr)
	rec.UpdatedAt = parseSQLiteTime(updateStr)

	return &rec, nil
}

// UpsertReleaseCache inserts or updates a release cache entry with 24-hour expiration.
func UpsertReleaseCache(db *sql.DB, record ReleaseCacheRecord) error {
	hasExec := 0
	if record.HasExecutables {
		hasExec = 1
	}

	if record.CheckedAt.IsZero() {
		record.CheckedAt = time.Now()
	}

	if record.ExpiresAt.IsZero() {
		record.ExpiresAt = record.CheckedAt.Add(24 * time.Hour)
	}

	checkedStr := record.CheckedAt.UTC().Format("2006-01-02 15:04:05")
	expiresStr := record.ExpiresAt.UTC().Format("2006-01-02 15:04:05")

	_, err := db.Exec(
		sqlUpsertReleaseCache,
		record.Tag,
		record.Version,
		hasExec,
		record.AssetURL,
		record.AssetSize,
		record.Platform,
		record.Arch,
		record.ChecksumURL,
		checkedStr,
		expiresStr,
	)
	if err != nil {
		return apperror.WrapSimple(err, "store.upsert_release_cache")
	}

	return nil
}

// ListCachedReleases returns recently checked release entries.
func ListCachedReleases(db *sql.DB, limit int) ([]ReleaseCacheRecord, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := db.Query(sqlListCachedReleases, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.list_cached_releases")
	}

	defer rows.Close()

	var records []ReleaseCacheRecord
	for rows.Next() {
		var (
			rec        ReleaseCacheRecord
			hasExec    int
			checkedStr string
			expireStr  string
			createStr  string
			updateStr  string
		)

		errScan := rows.Scan(
			&rec.ReleaseCacheId,
			&rec.Tag,
			&rec.Version,
			&hasExec,
			&rec.AssetURL,
			&rec.AssetSize,
			&rec.Platform,
			&rec.Arch,
			&rec.ChecksumURL,
			&checkedStr,
			&expireStr,
			&createStr,
			&updateStr,
		)
		if errScan != nil {
			return nil, apperror.WrapSimple(errScan, "store.list_cached_releases.scan")
		}

		rec.HasExecutables = (hasExec == 1)
		rec.CheckedAt = parseSQLiteTime(checkedStr)
		rec.ExpiresAt = parseSQLiteTime(expireStr)
		rec.CreatedAt = parseSQLiteTime(createStr)
		rec.UpdatedAt = parseSQLiteTime(updateStr)

		records = append(records, rec)
	}

	return records, nil
}

func parseSQLiteTime(s string) time.Time {
	formats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05-07:00",
	}

	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}

	return time.Time{}
}
