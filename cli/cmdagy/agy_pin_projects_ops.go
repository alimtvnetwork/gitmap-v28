// Package cmdagy — agy_pin_projects_ops.go handles add and remove operations for pinned projects.
package cmdagy

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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

	return insertPinnedProject(store, p)
}

func insertPinnedProject(store *PinnedProjectsStore, p *AgyProject) (*PinnedProject, *apperror.AppError) {
	if existing := findPinnedInStore(store, p.ID); existing != nil {
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
