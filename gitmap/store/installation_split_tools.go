package store

import (
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	sqlInsertSplitTool = `INSERT OR REPLACE INTO InstalledTool
(Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, UpdatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	sqlSelectSplitTool = `SELECT InstalledToolId, Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, InstalledAt, UpdatedAt
FROM InstalledTool WHERE Tool = ?;`

	sqlSelectAllSplitTools = `SELECT InstalledToolId, Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, InstalledAt, UpdatedAt
FROM InstalledTool ORDER BY Tool;`

	sqlDeleteSplitTool = `DELETE FROM InstalledTool WHERE Tool = ?;`
	sqlExistsSplitTool = `SELECT COUNT(*) FROM InstalledTool WHERE Tool = ?;`
)

// SaveInstalledTool records a tool installation into installation.db.
func (s *InstallationSplitDB) SaveInstalledTool(tool, version, manager string) error {
	major, minor, patch, build := parseVersionParts(version)
	versionStr := compileVersionString(major, minor, patch, build)
	if version != "" && versionStr == "0.0.0" {
		versionStr = version
	}

	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.conn.Exec(sqlInsertSplitTool, tool, major, minor, patch, build, versionStr, manager, "", now)
	if err != nil {
		return apperror.WrapSimple(err, "installation_split.saveInstalledTool")
	}

	return nil
}

// GetInstalledTool retrieves a single tool record by name from installation.db.
func (s *InstallationSplitDB) GetInstalledTool(name string) (InstalledTool, error) {
	var t InstalledTool
	row := s.conn.QueryRow(sqlSelectSplitTool, name)
	err := row.Scan(&t.ID, &t.Tool, &t.VersionMajor, &t.VersionMinor,
		&t.VersionPatch, &t.VersionBuild, &t.VersionString, &t.PackageManager,
		&t.InstallPath, &t.InstalledAt, &t.UpdatedAt)
	if err != nil {
		return t, apperror.WrapSimple(err, "installation_split.getInstalledTool")
	}

	return t, nil
}

// ListInstalledTools returns all tracked installations from installation.db.
func (s *InstallationSplitDB) ListInstalledTools() ([]InstalledTool, error) {
	rows, err := s.conn.Query(sqlSelectAllSplitTools)
	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.listInstalledTools")
	}

	defer rows.Close()

	return scanSplitTools(rows)
}

func scanSplitTools(rows *sql.Rows) ([]InstalledTool, error) {
	var tools []InstalledTool
	for rows.Next() {
		var t InstalledTool
		err := rows.Scan(&t.ID, &t.Tool, &t.VersionMajor, &t.VersionMinor,
			&t.VersionPatch, &t.VersionBuild, &t.VersionString, &t.PackageManager,
			&t.InstallPath, &t.InstalledAt, &t.UpdatedAt)
		if err != nil {
			return nil, apperror.WrapSimple(err, "installation_split.scanSplitTools")
		}

		tools = append(tools, t)
	}

	return tools, nil
}

// RemoveInstalledTool deletes a tool record from installation.db.
func (s *InstallationSplitDB) RemoveInstalledTool(name string) error {
	_, err := s.conn.Exec(sqlDeleteSplitTool, name)
	if err != nil {
		return apperror.WrapSimple(err, "installation_split.removeInstalledTool")
	}

	_ = s.RecordLog(InstallationLogRecord{
		Tool:      name,
		Action:    "uninstall",
		IsSuccess: true,
	})

	return nil
}

// IsToolInstalled reports whether a tool exists in installation.db.
func (s *InstallationSplitDB) IsToolInstalled(name string) bool {
	var count int
	err := s.conn.QueryRow(sqlExistsSplitTool, name).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}
