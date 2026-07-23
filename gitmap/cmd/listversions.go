package cmd

import (
	"strconv"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
)

// versionEntry pairs a parsed version with its changelog notes and source.
type versionEntry struct {
	Version release.Version
	Notes   []string
	Source  string
}

// runListVersions handles the "list-versions" command.
func runListVersions(args []string) {
	checkHelp("list-versions", args)
	asJSON := hasListVersionsJSONFlag(args)
	limit := parseListVersionsLimit(args)
	source := parseListVersionsSource(args)
	entries := collectVersionEntries()
	entries = filterVersionsBySource(entries, source)
	entries = applyVersionLimit(entries, limit)

	if asJSON {
		printVersionEntriesJSON(entries)

		return
	}

	printVersionEntriesTerminal(entries)
}

// parseListVersionsSource extracts the --source value from args.
func parseListVersionsSource(args []string) string {
	for i, arg := range args {
		if arg == constants.FlagSource && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
}

// filterVersionsBySource keeps only entries matching the given source (empty = all).
func filterVersionsBySource(entries []versionEntry, source string) []versionEntry {
	if source == "" {
		return entries
	}

	var filtered []versionEntry
	for _, e := range entries {
		if e.Source == source {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

// hasListVersionsJSONFlag checks if --json is present in args.
func hasListVersionsJSONFlag(args []string) bool {
	for _, arg := range args {
		if arg == constants.FlagJSON {
			return true
		}
	}

	return false
}

// parseListVersionsLimit extracts the --limit N value from args.
func parseListVersionsLimit(args []string) int {
	for i, arg := range args {
		if arg == constants.FlagLimit && i+1 < len(args) {
			n, err := strconv.Atoi(args[i+1])
			if err == nil && n > 0 {
				return n
			}
		}
	}

	return 0
}

// applyVersionLimit trims entries to at most n items (0 means no limit).
func applyVersionLimit(entries []versionEntry, n int) []versionEntry {
	if n <= 0 || n >= len(entries) {
		return entries
	}

	return entries[:n]
}

// collectVersionEntries reads tags, parses, sorts, and attaches changelog + source.
func collectVersionEntries() []versionEntry {
	versions := collectVersionTags()
	changelog := loadChangelogMap()
	sources := loadVersionSourceMap()

	entries := make([]versionEntry, len(versions))
	for i, v := range versions {
		entries[i] = versionEntry{Version: v, Notes: changelog[v.String()], Source: sources[v.String()]}
	}

	return entries
}

// loadVersionSourceMap reads the Releases table to build a tag→source map.
func loadVersionSourceMap() map[string]string {
	db, err := openDB()
	if err != nil {
		return map[string]string{}
	}
	defer db.Close()

	releases, err := db.ListReleases()
	if err != nil {
		return map[string]string{}
	}

	m := make(map[string]string, len(releases))
	for _, r := range releases {
		m[r.Tag] = r.Source
	}

	return m
}
