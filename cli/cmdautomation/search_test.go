package cmdautomation

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRunSearch_ConcurrencyNoDeadlock(t *testing.T) {
	tempDir := t.TempDir()
	createSearchTestFiles(t, tempDir, 25)
	opts := SearchOptions{Pattern: "target_pattern_needle", Dir: tempDir, Workers: 4}

	res, err := RunSearch(opts)
	if err != nil {
		t.Fatalf("unexpected search error: %v", err)
	}
	if res.TotalHits != 25 {
		t.Fatalf("expected 25 hits, got: %d", res.TotalHits)
	}
}

func createSearchTestFiles(t *testing.T, dir string, count int) {
	for i := 0; i < count; i++ {
		file := filepath.Join(dir, fmt.Sprintf("file_%d.go", i))
		writeErr := os.WriteFile(file, []byte("target_pattern_needle\n"), 0o600)
		if writeErr != nil {
			t.Fatalf("failed to write test file: %v", writeErr)
		}
	}
}

func TestRunSearch_EmptyPattern(t *testing.T) {
	opts := SearchOptions{Pattern: ""}

	_, err := RunSearch(opts)
	if err == nil {
		t.Fatalf("expected validation error for empty pattern, got nil")
	}
}
