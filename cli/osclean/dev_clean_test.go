package osclean

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsCategorySelected(t *testing.T) {
	c := categoryCleaner{
		Name:    "go-buildcache",
		Aliases: []string{"go", "gobuild"},
	}

	if !isCategorySelected(c, nil) {
		t.Errorf("expected empty only list to select all categories")
	}
	if !isCategorySelected(c, []string{"go"}) {
		t.Errorf("expected 'go' alias to select go-buildcache")
	}
	if !isCategorySelected(c, []string{"GO-BUILDCACHE"}) {
		t.Errorf("expected case-insensitive match")
	}
	if isCategorySelected(c, []string{"npm", "yarn"}) {
		t.Errorf("expected unrelated categories to not match")
	}
}

func TestSweepTarget_DryRun(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "pkg", "mod")
	_ = os.MkdirAll(subDir, 0o755)
	file1 := filepath.Join(tempDir, "file1.txt")
	_ = os.WriteFile(file1, []byte("hello"), 0o644)
	file2 := filepath.Join(subDir, "file2.txt")
	_ = os.WriteFile(file2, []byte("world123"), 0o644)

	stats := SweepTarget(tempDir, true, false)
	if stats.DirsRemoved != 1 {
		t.Errorf("expected 1 dir in dry run, got %d", stats.DirsRemoved)
	}
	if stats.ItemsRemoved != 2 {
		t.Errorf("expected 2 files counted, got %d", stats.ItemsRemoved)
	}
	if stats.BytesFreed < 13 {
		t.Errorf("expected at least 13 bytes freed, got %d", stats.BytesFreed)
	}

	// Verify files still exist in dry-run
	if _, err := os.Stat(file1); err != nil {
		t.Errorf("expected file1 to exist after dry run")
	}
}

func TestSweepTarget_LiveWithReadOnly(t *testing.T) {
	tempDir := t.TempDir()
	roFile := filepath.Join(tempDir, "readonly.txt")
	_ = os.WriteFile(roFile, []byte("data"), 0o444)

	stats := SweepTarget(tempDir, false, true)
	if stats.ItemsRemoved != 1 {
		t.Errorf("expected 1 item removed, got %d", stats.ItemsRemoved)
	}
	if _, err := os.Stat(roFile); !os.IsNotExist(err) {
		t.Errorf("expected readonly file to be deleted")
	}
}

func TestCleanDevCaches_OnlyFilter(t *testing.T) {
	opts := DevCleanOptions{
		IsDryRun:       true,
		OnlyCategories: []string{"bun"},
	}

	res := CleanDevCaches(opts)
	if res.IsFailure() {
		t.Fatalf("CleanDevCaches failed: %v", res.AppError())
	}

	summary := res.Value
	if len(summary.Categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(summary.Categories))
	}
	if summary.Categories[0].Category != "bun-cache" {
		t.Errorf("expected bun-cache, got %s", summary.Categories[0].Category)
	}
}
