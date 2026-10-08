package cmdrm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandFilePatterns(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "patch_01.py")
	file2 := filepath.Join(tempDir, "patch_02.py")
	file3 := filepath.Join(tempDir, "other.txt")
	gitFile := filepath.Join(tempDir, ".git", "config")

	_ = os.WriteFile(file1, []byte("print('patch1')"), 0644)
	_ = os.WriteFile(file2, []byte("print('patch2')"), 0644)
	_ = os.WriteFile(file3, []byte("hello"), 0644)
	_ = os.MkdirAll(filepath.Dir(gitFile), 0755)
	_ = os.WriteFile(gitFile, []byte("git config"), 0644)

	patterns := []string{"patch_*.py", ".git/*"}
	results, err := ExpandFilePatterns(patterns, tempDir)
	if err != nil {
		t.Fatalf("ExpandFilePatterns failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 matched files, got %d: %v", len(results), results)
	}

	for _, r := range results {
		if strings.Contains(r, ".git") {
			t.Errorf("Protected .git path leaked into results: %s", r)
		}
	}
}

func TestStageAndRemoveFilesAndUndo(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "temp_test_1.tmp")
	file2 := filepath.Join(tempDir, "sub", "temp_test_2.tmp")

	_ = os.MkdirAll(filepath.Dir(file2), 0755)
	_ = os.WriteFile(file1, []byte("content 1"), 0644)
	_ = os.WriteFile(file2, []byte("content 2 in subfolder"), 0644)

	taskID := "task-240-unit-test"
	opts := RmOptions{
		Patterns: []string{"temp_test_1.tmp", "sub/temp_test_2.tmp"},
		TaskId:   taskID,
		RepoRoot: tempDir,
	}

	manifest, err := StageAndRemoveFiles(opts)
	if err != nil {
		t.Fatalf("StageAndRemoveFiles failed: %v", err)
	}

	if manifest.FileCount != 2 {
		t.Errorf("Expected 2 files in manifest, got %d", manifest.FileCount)
	}

	if _, err := os.Stat(file1); !os.IsNotExist(err) {
		t.Errorf("Expected file1 to be removed from workspace")
	}
	if _, err := os.Stat(file2); !os.IsNotExist(err) {
		t.Errorf("Expected file2 to be removed from workspace")
	}

	undoOpts := RmUndoOptions{
		TaskId:   taskID,
		RepoRoot: tempDir,
	}
	if err := ExecuteUndo(undoOpts); err != nil {
		t.Fatalf("ExecuteUndo failed: %v", err)
	}

	content1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("Failed to read restored file1: %v", err)
	}
	if string(content1) != "content 1" {
		t.Errorf("Restored file1 content mismatch: got %q", string(content1))
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("Failed to read restored file2: %v", err)
	}
	if string(content2) != "content 2 in subfolder" {
		t.Errorf("Restored file2 content mismatch: got %q", string(content2))
	}
}

func TestUndoMissingTaskDirectoryError(t *testing.T) {
	missingID := "task-nonexistent-9999"
	undoOpts := RmUndoOptions{TaskId: missingID}

	err := ExecuteUndo(undoOpts)
	if err == nil {
		t.Fatalf("Expected error for missing task, got nil")
	}

	expectedPrefix := "Backup directory for task '" + missingID + "' not found in OS temporary storage."
	if !strings.Contains(err.Error(), expectedPrefix) {
		t.Errorf("Error did not contain expected message: %v", err)
	}
}

func TestPurgeRmBackups(t *testing.T) {
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "purge_me.tmp")
	_ = os.WriteFile(file, []byte("purge content"), 0644)

	taskID := "task-purge-test"
	opts := RmOptions{
		Patterns: []string{"purge_me.tmp"},
		TaskId:   taskID,
		RepoRoot: tempDir,
	}

	if _, err := StageAndRemoveFiles(opts); err != nil {
		t.Fatalf("StageAndRemoveFiles failed: %v", err)
	}

	count, err := PurgeRmBackups(RmPurgeOptions{TaskId: taskID})
	if err != nil {
		t.Fatalf("PurgeRmBackups failed: %v", err)
	}

	if count == 0 {
		t.Errorf("Expected at least 1 purged session, got 0")
	}
}

func TestDryRunExecution(t *testing.T) {
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "dry_run_file.tmp")
	_ = os.WriteFile(file, []byte("keep me"), 0644)

	taskID := "task-dry-run-test"
	opts := RmOptions{
		Patterns: []string{"dry_run_file.tmp"},
		TaskId:   taskID,
		DryRun:   true,
		RepoRoot: tempDir,
	}

	manifest, err := StageAndRemoveFiles(opts)
	if err != nil {
		t.Fatalf("StageAndRemoveFiles in dry-run failed: %v", err)
	}

	if manifest.FileCount != 1 {
		t.Errorf("Expected 1 file in dry-run manifest, got %d", manifest.FileCount)
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Errorf("File should NOT have been removed during dry-run")
	}
}
