package osclean

import (
	"os"
	"path/filepath"
	"testing"
)

func createTestTempFixture(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	filePath := filepath.Join(dir, "dummy_temp.log")
	if err := os.WriteFile(filePath, []byte("temp content"), 0600); err != nil {
		t.Fatalf("failed to create dummy temp file: %v", err)
	}
	return dir, filePath
}

func setupMockResolver(dir string) func() {
	orig := defaultTempDirResolver
	defaultTempDirResolver = func() []string {
		return []string{dir}
	}
	return func() { defaultTempDirResolver = orig }
}

func TestCleanTempDirectories_DryRun(t *testing.T) {
	dir, filePath := createTestTempFixture(t)
	cleanup := setupMockResolver(dir)
	defer cleanup()

	res := CleanTempDirectories(CleanOptions{IsDryRun: true})
	if res.IsFailure() {
		t.Fatalf("unexpected failure: %v", res.AppError())
	}
	if res.Value.RemovedFilesCount != 1 {
		t.Errorf("expected 1 file counted in dry run, got %d", res.Value.RemovedFilesCount)
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("file should not be removed in dry-run mode")
	}
}

func TestCleanTempDirectories_Execution(t *testing.T) {
	dir, filePath := createTestTempFixture(t)
	cleanup := setupMockResolver(dir)
	defer cleanup()

	res := CleanTempDirectories(CleanOptions{IsDryRun: false})
	if res.IsFailure() {
		t.Fatalf("unexpected failure: %v", res.AppError())
	}
	if res.Value.RemovedFilesCount != 1 {
		t.Errorf("expected 1 file removed, got %d", res.Value.RemovedFilesCount)
	}
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Errorf("file should be removed in execution mode")
	}
}
