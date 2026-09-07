package store

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
