// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ChangelogEntry represents one version section in CHANGELOG.md.
//
// Notes is preserved for backward compatibility (flat top-level bullets).
// Title and Bullets are populated by the new structured parser used by the
// pretty-printed `changelog --latest` console output.
type ChangelogEntry struct {
	Version string
	Title   string
	Notes   []string
	Bullets []ChangelogBullet
}

// ChangelogBullet represents a single bullet line with its indent depth and
// whether it was an ordered-list item. Depth 0 = top-level, 1 = nested, etc.
type ChangelogBullet struct {
	Depth   int
	Ordered bool
	Marker  string // "-", "*", or "1." — preserved for ordered numbering
	Text    string
}

// ReadChangelog reads concise changelog entries from CHANGELOG.md.
func ReadChangelog() ([]ChangelogEntry, error) {
	file, err := os.Open(constants.ChangelogFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries, err := parseChangelogStream(file)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no version sections found in %s", constants.ChangelogFile)
	}

	return entries, nil
}

// FindChangelogEntry returns a changelog entry by version.
func FindChangelogEntry(entries []ChangelogEntry, version string) (ChangelogEntry, bool) {
	target := NormalizeVersion(version)
	for _, entry := range entries {
		if NormalizeVersion(entry.Version) == target {
			return entry, true
		}
	}

	return ChangelogEntry{}, false
}

// NormalizeVersion normalizes a changelog version string to v-prefixed form.
func NormalizeVersion(version string) string {
	v := strings.TrimSpace(version)
	v = strings.TrimPrefix(v, "gitmap")
	v = strings.TrimSpace(v)
	if len(v) == 0 {
		return ""
	}
	if strings.HasPrefix(v, "v") {
		return v
	}

	return "v" + v
}

// parseVersionHeader extracts the version token from a markdown heading.
func parseVersionHeader(header string) string {
	raw := strings.TrimSpace(strings.TrimPrefix(header, "## "))
	if len(raw) == 0 {
		return ""
	}

	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return ""
	}

	version := strings.Trim(parts[0], "[]")
	if len(version) == 0 {
		return ""
	}
	if strings.HasPrefix(version, "v") {
		return version
	}

	return "v" + version
}
