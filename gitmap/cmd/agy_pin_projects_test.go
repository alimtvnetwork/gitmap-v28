// Package cmd — agy_pin_projects_test.go tests pinned Antigravity projects functionality.
package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAgyPinProjectsStoreLoadSave(t *testing.T) {
	tmp := t.TempDir()
	storeFile := filepath.Join(tmp, "pinned_projects.json")

	store := &PinnedProjectsStore{
		Version:   "1.0.0",
		UpdatedAt: "2026-09-04T00:00:00Z",
		Projects: []PinnedProject{
			{
				ID:       "proj-1",
				Name:     "Test Project",
				Path:     tmp,
				Branch:   "main",
				PinnedAt: "2026-09-04T00:00:00Z",
			},
		},
	}

	t.Setenv("USERPROFILE", tmp)
	t.Setenv("HOME", tmp)

	if err := savePinnedProjectsStore(store); err != nil {
		t.Fatalf("savePinnedProjectsStore failed: %v", err)
	}

	loaded, loadErr := loadPinnedProjectsStore()
	if loadErr != nil {
		t.Fatalf("loadPinnedProjectsStore failed: %v", loadErr)
	}

	if len(loaded.Projects) != 1 || loaded.Projects[0].ID != "proj-1" {
		t.Fatalf("unexpected loaded projects: %+v", loaded.Projects)
	}

	_ = os.Remove(storeFile)
}

func TestAgyPinProjectsAddAndRemove(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("USERPROFILE", tmp)
	t.Setenv("HOME", tmp)

	pinned, err := addPinnedProjectTarget(tmp)
	if err != nil {
		t.Fatalf("addPinnedProjectTarget failed: %v", err)
	}

	if pinned.Name != filepath.Base(tmp) {
		t.Errorf("expected name %s, got %s", filepath.Base(tmp), pinned.Name)
	}

	// Idempotency: re-adding should return existing without error
	dup, dupErr := addPinnedProjectTarget(tmp)
	if dupErr != nil {
		t.Fatalf("re-adding pinned project returned error: %v", dupErr)
	}
	if dup.ID != pinned.ID {
		t.Errorf("expected duplicate to have same ID")
	}

	store, _ := loadPinnedProjectsStore()
	if len(store.Projects) != 1 {
		t.Fatalf("expected 1 project in store, got %d", len(store.Projects))
	}

	// Remove project
	removed, rmErr := removePinnedProjectTarget(pinned.ID)
	if rmErr != nil {
		t.Fatalf("removePinnedProjectTarget failed: %v", rmErr)
	}
	if removed.ID != pinned.ID {
		t.Errorf("expected removed ID %s, got %s", pinned.ID, removed.ID)
	}

	storeAfter, _ := loadPinnedProjectsStore()
	if len(storeAfter.Projects) != 0 {
		t.Errorf("expected 0 projects after remove, got %d", len(storeAfter.Projects))
	}
}

func TestAgyPinProjectsClearAll(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("USERPROFILE", tmp)
	t.Setenv("HOME", tmp)

	_, _ = addPinnedProjectTarget(tmp)
	store, _ := loadPinnedProjectsStore()
	if len(store.Projects) != 1 {
		t.Fatalf("failed to add project for clear test")
	}

	cleared, clearErr := clearAllPinnedProjects()
	if clearErr != nil {
		t.Fatalf("clearAllPinnedProjects failed: %v", clearErr)
	}
	if cleared != 1 {
		t.Errorf("expected 1 cleared, got %d", cleared)
	}

	storeAfter, _ := loadPinnedProjectsStore()
	if len(storeAfter.Projects) != 0 {
		t.Errorf("expected 0 projects after clear, got %d", len(storeAfter.Projects))
	}
}
