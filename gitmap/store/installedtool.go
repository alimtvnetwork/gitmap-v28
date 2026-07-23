package store

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

// SaveInstalledTool records a tool installation with parsed version.
func (db *DB) SaveInstalledTool(tool, version, manager string) error {
	major, minor, patch, build := parseVersionParts(version)
	versionStr := compileVersionString(major, minor, patch, build)

	if version != "" && versionStr == "0.0.0" {
		versionStr = version
	}

	_, err := db.conn.Exec(constants.SQLInsertInstalledTool,
		tool, major, minor, patch, build, versionStr, manager, "")

	return err
}

// GetInstalledTool retrieves a single tool record by name.
func (db *DB) GetInstalledTool(name string) (InstalledTool, error) {
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
	rows, err := db.conn.Query(constants.SQLSelectAllInstalled)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	_, err := db.conn.Exec(constants.SQLDeleteInstalledTool, name)

	return err
}

// IsToolInstalled checks if a tool exists in the database.
func (db *DB) IsToolInstalled(name string) bool {
	var count int

	err := db.conn.QueryRow(constants.SQLExistsInstalledTool, name).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// parseVersionParts splits a version string into major, minor, patch, build.
func parseVersionParts(version string) (int, int, int, int) {
	s := strings.TrimPrefix(version, "v")
	if s == "" {
		return 0, 0, 0, 0
	}

	parts := strings.Split(s, ".")
	major := atoiSafe(safeIndex(parts, 0))
	minor := atoiSafe(safeIndex(parts, 1))
	patch := atoiSafe(safeIndex(parts, 2))
	build := atoiSafe(safeIndex(parts, 3))

	return major, minor, patch, build
}

// compileVersionString builds a version string from parts.
func compileVersionString(major, minor, patch, build int) string {
	if build > 0 {
		return fmt.Sprintf("%d.%d.%d.%d", major, minor, patch, build)
	}

	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

// CompareVersions compares two installed tools by version.
// Returns -1 if a < b, 0 if equal, 1 if a > b.
func CompareVersions(a, b InstalledTool) int {
	if a.VersionMajor != b.VersionMajor {
		return intCmp(a.VersionMajor, b.VersionMajor)
	}
	if a.VersionMinor != b.VersionMinor {
		return intCmp(a.VersionMinor, b.VersionMinor)
	}
	if a.VersionPatch != b.VersionPatch {
		return intCmp(a.VersionPatch, b.VersionPatch)
	}

	return intCmp(a.VersionBuild, b.VersionBuild)
}

// intCmp returns -1, 0, or 1.
func intCmp(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}

	return 0
}

// atoiSafe converts string to int, returning 0 on error.
func atoiSafe(s string) int {
	// Strip pre-release suffix (e.g. "3-rc1" → "3").
	if idx := strings.IndexAny(s, "-+"); idx >= 0 {
		s = s[:idx]
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return n
}

// safeIndex returns the element at index or empty string.
func safeIndex(parts []string, idx int) string {
	if idx < len(parts) {
		return parts[idx]
	}

	return ""
}

// FormatInstalledAt formats the InstalledAt field for display.
func (t InstalledTool) FormatInstalledAt() string {
	parsed, err := time.Parse("2006-01-02 15:04:05", t.InstalledAt)
	if err != nil {
		return t.InstalledAt
	}

	return parsed.Format("02-Jan-2006")
}
