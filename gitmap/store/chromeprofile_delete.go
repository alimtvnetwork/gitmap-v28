// Package store — chromeprofile_delete.go: removes a ChromeProfile +
// its cascaded ChromeProfileExport rows. Returns the artifact file
// paths that were tracked so the CLI can rm() them on disk.
package store

import "fmt"

const sqlSelectChromeProfileExports = `
SELECT e.FilePath FROM ChromeProfileExport e
JOIN ChromeProfile p ON p.ChromeProfileId = e.ChromeProfileId
WHERE p.Name = ?`

const sqlDeleteChromeProfile = `DELETE FROM ChromeProfile WHERE Name = ?`

const sqlDeleteChromeProfileExports = `
DELETE FROM ChromeProfileExport
WHERE ChromeProfileId IN (SELECT ChromeProfileId FROM ChromeProfile WHERE Name = ?)`

// DeleteChromeProfile removes the named profile and its artifact rows.
// Returns the list of artifact file paths so the caller can clean disk.
func (db *DB) DeleteChromeProfile(name string) ([]string, error) {
	if err := db.EnsureChromeProfileTables(); err != nil {
		return nil, err
	}
	paths, err := db.collectChromeArtifactPaths(name)
	if err != nil {
		return nil, err
	}
	if _, err := db.conn.Exec(sqlDeleteChromeProfileExports, name); err != nil {
		return nil, fmt.Errorf("delete chrome-profile exports: %w", err)
	}
	if _, err := db.conn.Exec(sqlDeleteChromeProfile, name); err != nil {
		return nil, fmt.Errorf("delete chrome-profile: %w", err)
	}
	return paths, nil
}

// ChromeProfileExists reports whether a row with the given name exists.
func (db *DB) ChromeProfileExists(name string) bool {
	if err := db.EnsureChromeProfileTables(); err != nil {
		return false
	}
	var id int64
	return db.conn.QueryRow(sqlSelectChromeProfileId, name).Scan(&id) == nil
}

// collectChromeArtifactPaths reads the FilePath column for every export
// row tied to the named profile.
func (db *DB) collectChromeArtifactPaths(name string) ([]string, error) {
	rows, err := db.conn.Query(sqlSelectChromeProfileExports, name)
	if err != nil {
		return nil, fmt.Errorf("query chrome-profile artifacts: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, fmt.Errorf("scan artifact: %w", err)
		}
		out = append(out, p)
	}
	return out, nil
}
