// Package fsutil — nested_filter.go filters out submodules and nested child repos.
package fsutil

import (
	"path/filepath"
	"strings"
)

// FilterTopLevelOnly removes any repository path that is a child of another repository path.
func FilterTopLevelOnly(repoPaths []string) []string {
	var filtered []string

	for i, candidate := range repoPaths {
		isChild := false
		cleanCand := filepath.Clean(candidate)

		for j, other := range repoPaths {
			if i == j {
				continue
			}
			cleanOther := filepath.Clean(other)
			rel, err := filepath.Rel(cleanOther, cleanCand)
			if err == nil && !strings.HasPrefix(rel, "..") && rel != "." {
				isChild = true
				break
			}
		}

		if !isChild {
			filtered = append(filtered, cleanCand)
		}
	}
	return filtered
}
