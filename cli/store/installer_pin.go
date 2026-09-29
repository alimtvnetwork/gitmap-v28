// Package store — installer_pin.go provides version pinning storage and queries.
package store

import (
	"database/sql"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SQLCreateInstallerPins creates the installer_pins table.
const SQLCreateInstallerPins = `CREATE TABLE IF NOT EXISTS installer_pins (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	slug TEXT NOT NULL UNIQUE,
	version TEXT NOT NULL,
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);`

// SQLUpsertInstallerPin inserts or updates the pinned version for a slug.
const SQLUpsertInstallerPin = `INSERT INTO installer_pins (slug, version, updated_at)
VALUES (?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(slug) DO UPDATE SET version = excluded.version, updated_at = CURRENT_TIMESTAMP;`

// SQLDeleteInstallerPin deletes a pinned version record by slug.
const SQLDeleteInstallerPin = `DELETE FROM installer_pins WHERE slug = ?;`

// SQLSelectInstallerPin selects a pinned version by slug.
const SQLSelectInstallerPin = `SELECT version FROM installer_pins WHERE slug = ? LIMIT 1;`

// SQLSelectAllInstallerPins selects all pinned versions.
const SQLSelectAllInstallerPins = `SELECT slug, version FROM installer_pins ORDER BY slug ASC;`

// RegisterInstallerPinsMigration creates the installer_pins table if it does not exist.
func RegisterInstallerPinsMigration(dbConn *sql.DB, migrationVersion int, isForce bool) error {
	if dbConn == nil {
		return apperror.NewSimple("RegisterInstallerPinsMigration", "E_INSTALLER_PINS_NIL_DB")
	}

	_, errExec := dbConn.Exec(SQLCreateInstallerPins)
	if errExec != nil {
		return apperror.WrapSimple(errExec, "RegisterInstallerPinsMigration")
	}

	return nil
}

// PinInstaller pins an installer or application target to a specific version.
func (db *DB) PinInstaller(nameOrSlug, version string) error {
	if db == nil || db.conn == nil {
		return apperror.NewSimple("PinInstaller", "E_INSTALLER_PINS_NIL_DB")
	}

	slug := strings.ToLower(strings.TrimSpace(nameOrSlug))
	ver := strings.TrimSpace(version)
	if slug == "" || ver == "" {
		return apperror.NewSimple("PinInstaller", "E_INSTALLER_PINS_INVALID_ARGS")
	}

	_, errExec := ExecWrapper(db.conn, SQLUpsertInstallerPin, slug, ver).Destruct()
	if errExec != nil {
		return apperror.WrapSimple(errExec, "PinInstaller")
	}

	return nil
}

// UnpinInstaller removes any version pinning for an installer.
func (db *DB) UnpinInstaller(nameOrSlug string) error {
	if db == nil || db.conn == nil {
		return apperror.NewSimple("UnpinInstaller", "E_INSTALLER_PINS_NIL_DB")
	}

	slug := strings.ToLower(strings.TrimSpace(nameOrSlug))
	if slug == "" {
		return apperror.NewSimple("UnpinInstaller", "E_INSTALLER_PINS_INVALID_ARGS")
	}

	_, errExec := ExecWrapper(db.conn, SQLDeleteInstallerPin, slug).Destruct()
	if errExec != nil {
		return apperror.WrapSimple(errExec, "UnpinInstaller")
	}

	return nil
}

// GetPinnedVersion retrieves the pinned version for a slug, or empty string if not pinned.
func (db *DB) GetPinnedVersion(nameOrSlug string) (string, error) {
	if db == nil || db.conn == nil {
		return "", apperror.NewSimple("GetPinnedVersion", "E_INSTALLER_PINS_NIL_DB")
	}

	slug := strings.ToLower(strings.TrimSpace(nameOrSlug))
	if slug == "" {
		return "", nil
	}

	row := QueryRowWrapper(db.conn, SQLSelectInstallerPin, slug)
	var version string
	errScan := row.Scan(&version)
	if errScan == sql.ErrNoRows {
		return "", nil
	}
	if errScan != nil {
		return "", apperror.WrapSimple(errScan, "GetPinnedVersion")
	}

	return version, nil
}

// ListPinnedVersions returns a mapping of slug to pinned version.
func (db *DB) ListPinnedVersions() (map[string]string, error) {
	pins := make(map[string]string)
	if db == nil || db.conn == nil {
		return pins, apperror.NewSimple("ListPinnedVersions", "E_INSTALLER_PINS_NIL_DB")
	}

	rows, errQuery := QueryWrapper(db.conn, SQLSelectAllInstallerPins).Destruct()
	if errQuery != nil {
		return pins, apperror.WrapSimple(errQuery, "ListPinnedVersions")
	}
	defer rows.Close()

	for rows.Next() {
		var slug, ver string
		if errScan := rows.Scan(&slug, &ver); errScan == nil {
			pins[slug] = ver
		}
	}

	return pins, nil
}
