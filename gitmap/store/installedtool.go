package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// InstalledTool represents a tracked tool installation.
type InstalledTool struct {
	ID             int64
	Tool           string
	VersionMajor   int
	VersionMinor   int
	VersionPatch   int
	VersionBuild   int
	VersionString  string
	PackageManager string
	InstallPath    string
	InstalledAt    string
	UpdatedAt      string
}

func openSplitOrFallback(db *DB) (*InstallationSplitDB, func()) {
	splitDB, err := OpenInstallationSplitDB()
	if err != nil {
		return nil, func() {}
	}

	migrateIfConnected(db, splitDB)

	return splitDB, func() { _ = splitDB.Close() }
}

func migrateIfConnected(db *DB, splitDB *InstallationSplitDB) {
	if db == nil || db.conn == nil {
		return
	}

	_ = MigrateInstalledToolsFromRoot(db.conn, splitDB)
}

// SaveInstalledTool records a tool installation with parsed version.
func (db *DB) SaveInstalledTool(tool, version, manager string) error {
	splitDB, cleanup := openSplitOrFallback(db)
	if splitDB != nil {
		defer cleanup()

		return splitDB.SaveInstalledTool(tool, version, manager)
	}

	major, minor, patch, build := parseVersionParts(version)
	versionStr := compileVersionString(major, minor, patch, build)
	if version != "" && versionStr == "0.0.0" {
		versionStr = version
	}

	_, err := ExecWrapper(db.conn, constants.SQLInsertInstalledTool,
		tool, major, minor, patch, build, versionStr, manager, "").Destruct()

	return err
}

// GetInstalledTool retrieves a single tool record by name.
func (db *DB) GetInstalledTool(name string) (InstalledTool, error) {
	splitDB, cleanup := openSplitOrFallback(db)
	if splitDB != nil {
		defer cleanup()

		return splitDB.GetInstalledTool(name)
	}

	var t InstalledTool
	err := db.conn.QueryRow(constants.SQLSelectInstalledTool, name).Scan(
		&t.ID, &t.Tool, &t.VersionMajor, &t.VersionMinor,
		&t.VersionPatch, &t.VersionBuild, &t.VersionString,
		&t.PackageManager, &t.InstallPath, &t.InstalledAt, &t.UpdatedAt,
	)

	return t, err
}

// ListInstalledTools returns all tracked installations.
func (db *DB) ListInstalledTools() ([]InstalledTool, error) {
	splitDB, cleanup := openSplitOrFallback(db)
	if splitDB != nil {
		defer cleanup()

		return splitDB.ListInstalledTools()
	}

	rows, err := QueryWrapper(db.conn, constants.SQLSelectAllInstalled).Destruct()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanLegacyInstalledTools(rows)
}

func scanLegacyInstalledTools(rows *sql.Rows) ([]InstalledTool, error) {
	var tools []InstalledTool
	for rows.Next() {
		var t InstalledTool
		err := rows.Scan(
			&t.ID, &t.Tool, &t.VersionMajor, &t.VersionMinor,
			&t.VersionPatch, &t.VersionBuild, &t.VersionString,
			&t.PackageManager, &t.InstallPath, &t.InstalledAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tools = append(tools, t)
	}

	return tools, rows.Err()
}

// RemoveInstalledTool deletes a tool record.
func (db *DB) RemoveInstalledTool(name string) error {
	splitDB, cleanup := openSplitOrFallback(db)
	if splitDB != nil {
		defer cleanup()

		return splitDB.RemoveInstalledTool(name)
	}

	_, err := ExecWrapper(db.conn, constants.SQLDeleteInstalledTool, name).Destruct()

	return err
}

// IsToolInstalled checks if a tool exists in the database.
func (db *DB) IsToolInstalled(name string) bool {
	splitDB, cleanup := openSplitOrFallback(db)
	if splitDB != nil {
		defer cleanup()

		return splitDB.IsToolInstalled(name)
	}

	var count int
	err := db.conn.QueryRow(constants.SQLExistsInstalledTool, name).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}
