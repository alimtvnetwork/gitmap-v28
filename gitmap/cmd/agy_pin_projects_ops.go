// Package cmd — agy_pin_projects_ops.go handles add and remove operations for pinned projects.
package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func addPinnedProjectTarget(target string) (*PinnedProject, *apperror.AppError) {
	store, loadErr := loadPinnedProjectsStore()

	if loadErr != nil {
		return nil, loadErr
	}

	p, resErr := resolveTargetAgyProject(target)

	if resErr != nil {
		return nil, resErr
	}

	existing := findPinnedInStore(store, p.ID)

	if existing != nil {
		return existing, nil
	}

	pinned := PinnedProject{
		ID:       p.ID,
		Name:     p.Name,
		Path:     p.GetPath(),
		Branch:   p.GetBranch(),
		PinnedAt: time.Now().UTC().Format(time.RFC3339),
	}
	store.Projects = append(store.Projects, pinned)
	store.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if saveErr := savePinnedProjectsStore(store); saveErr != nil {
		return nil, saveErr
	}

	return &pinned, nil
}

func removePinnedProjectTarget(target string) (*PinnedProject, *apperror.AppError) {
	store, loadErr := loadPinnedProjectsStore()

	if loadErr != nil {
		return nil, loadErr
	}

	index, removed := findPinnedIndex(store, target)

	if index == -1 {
		return nil, apperror.NewSimple(fmt.Sprintf("pinned project %q not found", target), "E9000")
	}

	store.Projects = append(store.Projects[:index], store.Projects[index+1:]...)
	store.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if saveErr := savePinnedProjectsStore(store); saveErr != nil {
		return nil, saveErr
	}

	return removed, nil
}

func clearAllPinnedProjects() (int, *apperror.AppError) {
	store, loadErr := loadPinnedProjectsStore()

	if loadErr != nil {
		return 0, loadErr
	}

	count := len(store.Projects)
	store.Projects = make([]PinnedProject, 0)
	store.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if saveErr := savePinnedProjectsStore(store); saveErr != nil {
		return 0, saveErr
	}

	return count, nil
}

func findPinnedInStore(store *PinnedProjectsStore, id string) *PinnedProject {
	for i := range store.Projects {
		if store.Projects[i].ID == id {
			return &store.Projects[i]
		}
	}

	return nil
}

func findPinnedIndex(store *PinnedProjectsStore, target string) (int, *PinnedProject) {
	low := strings.ToLower(target)
	cleanPath := filepath.Clean(target)

	for i := range store.Projects {
		p := &store.Projects[i]
		if p.ID == target || strings.HasPrefix(p.ID, target) {
			return i, p
		}

		if strings.ToLower(p.Name) == low || filepath.Clean(p.Path) == cleanPath {
			return i, p
		}
	}

	return -1, nil
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
