package store

import (
	"database/sql"
	"path/filepath"
	"time"
)

// DevtoolsCacheRecord represents a persisted cache directory entry.
type DevtoolsCacheRecord struct {
	ID             int64  `json:"id"`
	Path           string `json:"path"`
	Ecosystem      string `json:"ecosystem"`
	SizeBytes      int64  `json:"sizeBytes"`
	FilesCount     int    `json:"filesCount"`
	DirsCount      int    `json:"dirsCount"`
	LastVerifiedAt int64  `json:"lastVerifiedAt"`
	IsCustom       bool   `json:"isCustom"`
	IsActive       bool   `json:"isActive"`
}

// DevtoolsCacheSplitDB wraps an isolated SQLite database connection for devtools cache.
type DevtoolsCacheSplitDB struct {
	conn *sql.DB
	Path string
}

// OpenDevtoolsCacheSplitDB opens the split SQLite database at .gitmap/data/devtools/cache/sql.db.
func OpenDevtoolsCacheSplitDB() (*DevtoolsCacheSplitDB, error) {
	dbPath := ResolveSplitDbPath("devtools", "cache", "")
	conn, err := OpenSQLiteDB(dbPath)
	if err != nil {
		return nil, err
	}
	s := &DevtoolsCacheSplitDB{conn: conn, Path: dbPath}
	if err := s.EnsureDevtoolsCacheTable(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return s, nil
}

// EnsureDevtoolsCacheTable creates the devtools_cache_paths table and indices if missing.
func (s *DevtoolsCacheSplitDB) EnsureDevtoolsCacheTable() error {
	schema := `CREATE TABLE IF NOT EXISTS devtools_cache_paths (
id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT NOT NULL UNIQUE, ecosystem TEXT NOT NULL,
size_bytes INTEGER NOT NULL DEFAULT 0, files_count INTEGER NOT NULL DEFAULT 0, dirs_count INTEGER NOT NULL DEFAULT 0,
last_verified_at INTEGER NOT NULL DEFAULT 0, is_custom INTEGER NOT NULL DEFAULT 0, is_active INTEGER NOT NULL DEFAULT 1);
CREATE INDEX IF NOT EXISTS idx_devtools_cache_ecosystem ON devtools_cache_paths(ecosystem);
CREATE INDEX IF NOT EXISTS idx_devtools_cache_active ON devtools_cache_paths(is_active);`
	_, err := s.conn.Exec(schema)
	return err
}

// SaveDiscoveredPaths saves or updates discovered cache records in the split DB.
func (s *DevtoolsCacheSplitDB) SaveDiscoveredPaths(records []DevtoolsCacheRecord) error {
	now := time.Now().Unix()
	for _, r := range records {
		if err := s.saveRecordRow(r, now); err != nil {
			return err
		}
	}
	return nil
}

func (s *DevtoolsCacheSplitDB) saveRecordRow(r DevtoolsCacheRecord, now int64) error {
	q := `INSERT INTO devtools_cache_paths (path, ecosystem, size_bytes, files_count, dirs_count, last_verified_at, is_custom, is_active)
VALUES (?, ?, ?, ?, ?, ?, ?, 1) ON CONFLICT(path) DO UPDATE SET ecosystem=excluded.ecosystem, size_bytes=excluded.size_bytes, files_count=excluded.files_count, dirs_count=excluded.dirs_count, last_verified_at=excluded.last_verified_at, is_custom=excluded.is_custom, is_active=1`
	custom := 0
	if r.IsCustom {
		custom = 1
	}
	_, err := s.conn.Exec(q, filepath.ToSlash(r.Path), r.Ecosystem, r.SizeBytes, r.FilesCount, r.DirsCount, now, custom)
	return err
}

// GetCachedPaths retrieves active records optionally filtered by ecosystem.
func (s *DevtoolsCacheSplitDB) GetCachedPaths(ecosystem string) ([]DevtoolsCacheRecord, error) {
	q := "SELECT id, path, ecosystem, size_bytes, files_count, dirs_count, last_verified_at, is_custom, is_active FROM devtools_cache_paths WHERE is_active = 1"
	if ecosystem != "" {
		return s.queryCachedRows(q+" AND ecosystem = ? ORDER BY path ASC", ecosystem)
	}
	return s.queryCachedRows(q + " ORDER BY path ASC")
}

func (s *DevtoolsCacheSplitDB) queryCachedRows(q string, args ...any) ([]DevtoolsCacheRecord, error) {
	rows, err := s.conn.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDevtoolsCacheRows(rows)
}

func scanDevtoolsCacheRows(rows *sql.Rows) ([]DevtoolsCacheRecord, error) {
	var records []DevtoolsCacheRecord
	for rows.Next() {
		var r DevtoolsCacheRecord
		var custom, active int
		if err := rows.Scan(&r.ID, &r.Path, &r.Ecosystem, &r.SizeBytes, &r.FilesCount, &r.DirsCount, &r.LastVerifiedAt, &custom, &active); err != nil {
			return nil, err
		}
		r.IsCustom = custom == 1
		r.IsActive = active == 1
		records = append(records, r)
	}
	return records, rows.Err()
}

// InvalidateCache marks records inactive, optionally filtered by ecosystem.
func (s *DevtoolsCacheSplitDB) InvalidateCache(ecosystem string) error {
	if ecosystem != "" {
		_, err := s.conn.Exec("UPDATE devtools_cache_paths SET is_active = 0 WHERE ecosystem = ?", ecosystem)
		return err
	}
	_, err := s.conn.Exec("UPDATE devtools_cache_paths SET is_active = 0")
	return err
}

// Close closes the underlying SQLite connection.
func (s *DevtoolsCacheSplitDB) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
