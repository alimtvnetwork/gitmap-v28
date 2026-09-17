package cmdos

import (
	"testing"
)

func setupMockFileRemover() (*[]string, func()) {
	var captured []string
	origRemover := defaultFileRemover
	origPrompt := promptAICleanConfirmFn

	defaultFileRemover = func(path string) (int, int64) {
		captured = append(captured, path)
		return 1, 512
	}
	promptAICleanConfirmFn = func(msg string) (bool, error) { return true, nil }

	cleanup := func() {
		defaultFileRemover = origRemover
		promptAICleanConfirmFn = origPrompt
	}
	return &captured, cleanup
}

func TestRunOSAICleanCLI_DryRun(t *testing.T) {
	_, cleanup := setupMockFileRemover()
	defer cleanup()

	if err := RunOSAICleanCLI([]string{"--dry-run"}); err != nil {
		t.Fatalf("unexpected error on dry-run: %v", err)
	}
}

func TestRunOSAICleanCLI_JSON(t *testing.T) {
	_, cleanup := setupMockFileRemover()
	defer cleanup()

	if err := RunOSAICleanCLI([]string{"--json"}); err != nil {
		t.Fatalf("unexpected error on json: %v", err)
	}
}

func TestRunOSAICleanCLI_Help(t *testing.T) {
	if err := RunOSAICleanCLI([]string{"--help"}); err != nil {
		t.Fatalf("unexpected error on help: %v", err)
	}
}

func TestPurgeCategoryFiles_Mocked(t *testing.T) {
	captured, cleanup := setupMockFileRemover()
	defer cleanup()

	testPaths := []string{"/mock/brain/file1.json", "/mock/brain/file2.json"}
	count, bytes := purgeCategoryFiles(testPaths)
	if count != 2 || bytes != 1024 {
		t.Errorf("expected 2 files and 1024 bytes, got %d files %d bytes", count, bytes)
	}
	if len(*captured) != 2 {
		t.Errorf("expected 2 captured files, got %d", len(*captured))
	}
}
