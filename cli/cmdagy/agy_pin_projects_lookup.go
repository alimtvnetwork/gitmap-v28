// Package cmdagy — agy_pin_projects_lookup.go handles matching and querying pinned projects.
package cmdagy

import (
	"path/filepath"
	"strconv"
	"strings"
)

func findPinnedInStore(store *PinnedProjectsStore, id string) *PinnedProject {
	for i := range store.Projects {
		if store.Projects[i].ID == id {
			return &store.Projects[i]
		}
	}

	return nil
}

func findPinnedIndex(store *PinnedProjectsStore, target string) (int, *PinnedProject) {
	if idx, p := findPinnedBySeq(store, target); idx != -1 {
		return idx, p
	}
	return findPinnedByNameOrPath(store, target)
}

func findPinnedBySeq(store *PinnedProjectsStore, target string) (int, *PinnedProject) {
	trimmed := strings.TrimLeft(target, "0")
	seq, err := strconv.Atoi(trimmed)
	if trimmed != "" && err == nil && seq >= 1 && seq <= len(store.Projects) {
		return seq - 1, &store.Projects[seq-1]
	}
	return -1, nil
}

func findPinnedByNameOrPath(store *PinnedProjectsStore, target string) (int, *PinnedProject) {
	low := strings.ToLower(target)
	cleanPath := filepath.Clean(target)
	for i := range store.Projects {
		p := &store.Projects[i]
		if isPinnedProjectMatch(p, target, low, cleanPath) {
			return i, p
		}
	}
	return -1, nil
}

func isPinnedProjectMatch(p *PinnedProject, target, low, cleanPath string) bool {
	if p.ID == target || strings.HasPrefix(strings.ToLower(p.ID), low) {
		return true
	}
	return strings.ToLower(p.Name) == low || strings.HasPrefix(strings.ToLower(p.Name), low) || filepath.Clean(p.Path) == cleanPath
}

func getPinnedProjectsMap() map[string]bool {
	pinnedMap := make(map[string]bool)
	store, loadErr := loadPinnedProjectsStore()

	if loadErr != nil {
		return pinnedMap
	}

	for _, p := range store.Projects {
		pinnedMap[p.ID] = true
		pinnedMap[filepath.Clean(p.Path)] = true
	}

	return pinnedMap
}
