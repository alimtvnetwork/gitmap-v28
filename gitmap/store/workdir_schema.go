// Package store — workdir_schema.go defines DDL for work_directories.
package store

const (
	SQLCreateWorkDirsTable = `
CREATE TABLE IF NOT EXISTS work_directories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    absolute_path TEXT UNIQUE NOT NULL,
    label TEXT,
    is_default INTEGER DEFAULT 0,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);`

	SQLCreateWorkDirectoryView = `
CREATE VIEW IF NOT EXISTS WorkDirectory AS
SELECT
    id AS WorkDirectoryId,
    absolute_path AS AbsolutePath,
    label AS Label,
    is_default AS IsDefault,
    created_at AS CreatedAt,
    updated_at AS UpdatedAt
FROM work_directories;`

	SQLUpsertWorkDir = `
INSERT INTO work_directories (absolute_path, label, is_default, updated_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(absolute_path) DO UPDATE SET
    label = CASE WHEN ? != '' THEN ? ELSE label END,
    is_default = CASE WHEN ? != 0 THEN ? ELSE is_default END,
    updated_at = CURRENT_TIMESTAMP;`

	SQLSelectAllWorkDirs = `
SELECT id, absolute_path, label, is_default, created_at, updated_at
FROM work_directories
ORDER BY is_default DESC, updated_at DESC;`

	SQLSetDefaultWorkDir = `UPDATE work_directories SET is_default = CASE WHEN id = ? OR absolute_path = ? THEN 1 ELSE 0 END;`

	SQLDeleteWorkDir = `DELETE FROM work_directories WHERE id = ? OR absolute_path = ?;`
)
