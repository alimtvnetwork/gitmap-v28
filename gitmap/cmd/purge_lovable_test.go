package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPurgeLovable(t *testing.T) {
	tempDir := t.TempDir()
	lovableDir := filepath.Join(tempDir, ".lovable")
	os.MkdirAll(lovableDir, 0755)
	
	trackedFile := filepath.Join(lovableDir, "tracked.txt")
	os.WriteFile(trackedFile, []byte("tracked"), 0644)
	
	untrackedFile := filepath.Join(lovableDir, "untracked.txt")
	os.WriteFile(untrackedFile, []byte("untracked"), 0644)
	
	trackedMap := map[string]bool{
		".lovable/tracked.txt": true,
	}
	
	purged, err := removeUntrackedLovable(tempDir, trackedMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if purged != 1 {
		t.Errorf("expected 1 file purged, got %d", purged)
	}
	
	if _, err := os.Stat(trackedFile); os.IsNotExist(err) {
		t.Errorf("expected tracked file to be preserved")
	}
	if _, err := os.Stat(untrackedFile); !os.IsNotExist(err) {
		t.Errorf("expected untracked file to be deleted")
	}
}
