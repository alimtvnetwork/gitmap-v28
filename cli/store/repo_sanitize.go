package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
)

type unknownRepoRow struct {
	repoID  int64
	absPath string
}

// SanitizeUnknownRepos heals repositories with "unknown" slug or repo name by
// deriving their identity from the physical directory name on disk.
func SanitizeUnknownRepos(db *sql.DB) error {
	if db == nil {
		return nil
	}
	rows, err := db.Query("SELECT RepoId, AbsolutePath FROM Repo WHERE Slug = 'unknown' OR RepoName = 'unknown'")
	if err != nil {
		return err
	}
	defer rows.Close()

	targets := collectUnknownRepoTargets(rows)
	for _, target := range targets {
		sanitizeSingleUnknownRepo(db, target.repoID, target.absPath)
	}

	return nil
}

func collectUnknownRepoTargets(rows *sql.Rows) []unknownRepoRow {
	var targets []unknownRepoRow
	for rows.Next() {
		var row unknownRepoRow
		if err := rows.Scan(&row.repoID, &row.absPath); err == nil {
			targets = append(targets, row)
		}
	}

	return targets
}

func sanitizeSingleUnknownRepo(db *sql.DB, repoID int64, absPath string) {
	cleanPath := filepath.Clean(strings.TrimRight(absPath, "/\\"))
	if _, statErr := os.Stat(cleanPath); statErr != nil {
		return
	}
	baseName := filepath.Base(cleanPath)
	if !isValidSanitizedRepoName(baseName) {
		return
	}
	slug := strings.ToLower(baseName)
	_, _ = db.Exec("UPDATE Repo SET RepoName = ?, Slug = ? WHERE RepoId = ?", baseName, slug, repoID)
}

func isValidSanitizedRepoName(name string) bool {
	if name == "" || name == "." || name == "/" || name == "\\" || name == "unknown" {
		return false
	}

	return true
}

// SanitizeUnknownRepos invokes repository sanitization on the DB instance.
func (db *DB) SanitizeUnknownRepos() error {
	if db == nil || db.conn == nil {
		return nil
	}

	return SanitizeUnknownRepos(db.conn)
}
