package vscodepm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOptimizeProjectsAt(t *testing.T) {
	tmpDir := t.TempDir()
	projectsFile := filepath.Join(tmpDir, "projects.json")

	repo1Dir := filepath.Join(tmpDir, "repo-1")
	repo2Dir := filepath.Join(tmpDir, "repo-2")
	_ = os.MkdirAll(repo1Dir, 0o755)
	_ = os.MkdirAll(repo2Dir, 0o755)
	missingDir := filepath.Join(tmpDir, "missing-repo")

	initial := []Entry{
		{Name: "repo-1", RootPath: repo1Dir, Paths: []string{repo1Dir}, Tags: []string{"tagA"}},
		{Name: "repo-1-dup", RootPath: repo1Dir, Paths: []string{filepath.Join(repo1Dir, "sub")}, Tags: []string{"tagB"}},
		{Name: "repo-2", RootPath: repo2Dir, Paths: []string{repo2Dir}, Tags: []string{"tagC"}},
		{Name: "repo-missing", RootPath: missingDir, Paths: []string{missingDir}},
	}

	if err := writeEntriesAtomic(projectsFile, initial); err != nil {
		t.Fatalf("failed writing initial: %v", err)
	}

	summary, err := OptimizeProjectsAt(projectsFile, nil, false)
	if err != nil {
		t.Fatalf("OptimizeProjectsAt failed: %v", err)
	}

	if summary.RemovedMissing != 1 {
		t.Errorf("expected 1 missing removed, got %d", summary.RemovedMissing)
	}

	if summary.RemovedDuplicates != 1 {
		t.Errorf("expected 1 duplicate removed, got %d", summary.RemovedDuplicates)
	}

	if summary.Removed != 2 {
		t.Errorf("expected 2 total removed, got %d", summary.Removed)
	}

	if summary.Remaining != 2 {
		t.Errorf("expected 2 remaining, got %d", summary.Remaining)
	}

	entries, _ := readEntries(projectsFile)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if len(entries[0].Tags) != 2 {
		t.Errorf("expected merged tags length 2, got %d", len(entries[0].Tags))
	}
}

func TestClearProjectsAtExcept(t *testing.T) {
	tmpDir := t.TempDir()
	projectsFile := filepath.Join(tmpDir, "projects.json")

	initial := []Entry{
		{Name: "keep-me", RootPath: "D:/work/keep-me"},
		{Name: "delete-me", RootPath: "D:/work/delete-me"},
	}

	if err := writeEntriesAtomic(projectsFile, initial); err != nil {
		t.Fatalf("failed writing initial: %v", err)
	}

	summary, err := ClearProjectsAt(projectsFile, []string{"keep-me"}, false, false)
	if err != nil {
		t.Fatalf("ClearProjectsAt failed: %v", err)
	}

	if summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", summary.Removed)
	}

	if summary.Remaining != 1 {
		t.Errorf("expected 1 remaining, got %d", summary.Remaining)
	}

	entries, _ := readEntries(projectsFile)
	if len(entries) != 1 || entries[0].Name != "keep-me" {
		t.Errorf("unexpected remaining entry: %+v", entries)
	}
}
