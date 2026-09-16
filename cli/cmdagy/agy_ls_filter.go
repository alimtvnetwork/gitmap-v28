// Package cmdagy — agy_ls_filter.go provides filtering logic for Antigravity projects.
package cmdagy

import (
	"path/filepath"
)

func filterAgyProjects(projects []AgyProject) []AgyProject {
	pinnedMap := getPinnedMapIfRequested()
	out := make([]AgyProject, 0, len(projects))

	for _, p := range projects {
		if isProjectFilteredOut(p, pinnedMap) {
			continue
		}

		out = append(out, p)
	}

	return out
}

func isProjectFilteredOut(p AgyProject, pinnedMap map[string]bool) bool {
	if isPinnedMismatch(p, pinnedMap) {
		return true
	}

	return !matchesAgyFilter(p)
}

func isPinnedMismatch(p AgyProject, pinnedMap map[string]bool) bool {
	if !agyLsOnlyPinned {
		return false
	}

	isPinnedByID := pinnedMap[p.ID]
	isPinnedByPath := pinnedMap[filepath.Clean(p.GetPath())]

	return !isPinnedByID && !isPinnedByPath
}

func getPinnedMapIfRequested() map[string]bool {
	if !agyLsOnlyPinned {
		return make(map[string]bool)
	}

	return getPinnedProjectsMap()
}

func matchesAgyFilter(p AgyProject) bool {
	path := p.GetPath()
	isMissing := path != "" && !checkDirExists(path)

	if isStatusMismatch(isMissing) {
		return false
	}

	return isTextFilterMatch(p.Name, path)
}

func isStatusMismatch(isMissing bool) bool {
	if agyLsOnlyMissing && !isMissing {
		return true
	}

	if agyLsOnlyActive && isMissing {
		return true
	}

	return false
}

func isTextFilterMatch(name, path string) bool {
	if agyLsFilter == "" {
		return true
	}

	return matchesFilter(name, path, agyLsFilter)
}
