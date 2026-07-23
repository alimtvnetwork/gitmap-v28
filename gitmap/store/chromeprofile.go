// Package store — chromeprofile.go: lazy-migrated tables for
// `gitmap chrome-profile-*`. Self-contained (CREATE IF NOT EXISTS on
// every call) so the main Migrate() pipeline stays untouched while the
// feature stabilizes. Tables:
//
//	ChromeProfile        — one row per (Name) seen by cpc/cpe
//	ChromeProfileExport  — one row per snapshot (json or csv)
//
// Spec: spec/04-generic-cli/40-chrome-profile-copy.md §6 (Persistence).
package store

import (
	"fmt"
)

const sqlCreateChromeProfile = `
CREATE TABLE IF NOT EXISTS ChromeProfile (
    ChromeProfileId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name            TEXT NOT NULL UNIQUE,
    SourcePath      TEXT NOT NULL DEFAULT '',
    IsOffline       INTEGER NOT NULL DEFAULT 1,
    CreatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const sqlCreateChromeProfileExport = `
CREATE TABLE IF NOT EXISTS ChromeProfileExport (
    ChromeProfileExportId INTEGER PRIMARY KEY AUTOINCREMENT,
    ChromeProfileId       INTEGER NOT NULL REFERENCES ChromeProfile(ChromeProfileId) ON DELETE CASCADE,
    Format                TEXT NOT NULL,
    FilePath              TEXT NOT NULL,
    ByteSize              INTEGER NOT NULL DEFAULT 0,
    CreatedAt             TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

const sqlUpsertChromeProfile = `
INSERT INTO ChromeProfile (Name, SourcePath, IsOffline, UpdatedAt)
VALUES (?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(Name) DO UPDATE SET
    SourcePath = excluded.SourcePath,
    IsOffline  = excluded.IsOffline,
    UpdatedAt  = CURRENT_TIMESTAMP`

const sqlSelectChromeProfileId = `SELECT ChromeProfileId FROM ChromeProfile WHERE Name = ?`

const sqlInsertChromeProfileExport = `
INSERT INTO ChromeProfileExport (ChromeProfileId, Format, FilePath, ByteSize)
VALUES (?, ?, ?, ?)`

const sqlListChromeProfiles = `
SELECT cp.Name,
       COALESCE((SELECT COUNT(*) FROM ChromeProfileExport e WHERE e.ChromeProfileId = cp.ChromeProfileId), 0),
       COALESCE((SELECT MAX(CreatedAt) FROM ChromeProfileExport e WHERE e.ChromeProfileId = cp.ChromeProfileId), cp.UpdatedAt)
FROM ChromeProfile cp
ORDER BY cp.Name`

// ChromeProfileRow is a list-friendly projection of one tracked profile.
type ChromeProfileRow struct {
	Name        string
	ExportCount int
	LastSeen    string
}

// EnsureChromeProfileTables creates the two tables on demand. Safe to
// call before every operation — IF NOT EXISTS makes it idempotent.
func (db *DB) EnsureChromeProfileTables() error {
	for _, stmt := range []string{sqlCreateChromeProfile, sqlCreateChromeProfileExport} {
		if _, err := db.conn.Exec(stmt); err != nil {
			return fmt.Errorf("ensure chrome-profile tables: %w", err)
		}
	}
	return nil
}

// UpsertChromeProfile inserts or refreshes the profile row by Name and
// returns its ChromeProfileId.
func (db *DB) UpsertChromeProfile(name, sourcePath string, isOffline bool) (int64, error) {
	if err := db.EnsureChromeProfileTables(); err != nil {
		return 0, err
	}
	offline := 0
	if isOffline {
		offline = 1
	}
	if _, err := db.conn.Exec(sqlUpsertChromeProfile, name, sourcePath, offline); err != nil {
		return 0, fmt.Errorf("upsert chrome-profile: %w", err)
	}
	var id int64
	if err := db.conn.QueryRow(sqlSelectChromeProfileId, name).Scan(&id); err != nil {
		return 0, fmt.Errorf("read chrome-profile id: %w", err)
	}
	return id, nil
}

// InsertChromeProfileExport records a snapshot (json or csv).
func (db *DB) InsertChromeProfileExport(profileID int64, format, filePath string, byteSize int) error {
	if err := db.EnsureChromeProfileTables(); err != nil {
		return err
	}
	if _, err := db.conn.Exec(sqlInsertChromeProfileExport, profileID, format, filePath, byteSize); err != nil {
		return fmt.Errorf("insert chrome-profile-export: %w", err)
	}
	return nil
}

// ListChromeProfilesDB returns every tracked profile with export counts.
func (db *DB) ListChromeProfilesDB() ([]ChromeProfileRow, error) {
	if err := db.EnsureChromeProfileTables(); err != nil {
		return nil, err
	}
	rows, err := db.conn.Query(sqlListChromeProfiles)
	if err != nil {
		return nil, fmt.Errorf("list chrome-profiles: %w", err)
	}
	defer rows.Close()
	var out []ChromeProfileRow
	for rows.Next() {
		var r ChromeProfileRow
		if err := rows.Scan(&r.Name, &r.ExportCount, &r.LastSeen); err != nil {
			return nil, fmt.Errorf("scan chrome-profile row: %w", err)
		}
		out = append(out, r)
	}
	return out, nil
}
