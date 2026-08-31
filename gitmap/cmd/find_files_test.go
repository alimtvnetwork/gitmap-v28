package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestFindWorkspace(t *testing.T) string {
	tempDir := t.TempDir()

	_ = os.MkdirAll(filepath.Join(tempDir, "src", "core"), 0755)
	_ = os.MkdirAll(filepath.Join(tempDir, "scripts"), 0755)

	_ = os.WriteFile(filepath.Join(tempDir, "src", "01-app.ts"), []byte("code"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "src", "core", "folder.go"), []byte("code"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "src", "core", "folder_test.go"), []byte("code"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "scripts", "06-cicd-runner.py"), []byte("code"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "README.md"), []byte("doc"), 0644)

	return tempDir
}

func TestFindFilesExact(t *testing.T) {
	ws := setupTestFindWorkspace(t)
	opts := FindFilesOptions{
		Query:     "folder.go",
		Kind:      MatchExact,
		Exts:      []string{"go"},
		TargetDir: ws,
	}

	matches, err := scanAndMatchFiles(opts)
	if err != nil {
		t.Fatalf("scanAndMatchFiles error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match for folder.go, got %d: %v", len(matches), matches)
	}
}

func TestFindFilesAnyContains(t *testing.T) {
	ws := setupTestFindWorkspace(t)
	opts := FindFilesOptions{
		Query:     "runner",
		Kind:      MatchContains,
		Exts:      []string{"py", "go"},
		TargetDir: ws,
	}

	matches, err := scanAndMatchFiles(opts)
	if err != nil {
		t.Fatalf("scanAndMatchFiles error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match for runner with py/go ext, got %d: %v", len(matches), matches)
	}
}

func TestFindFilesStartsWith(t *testing.T) {
	ws := setupTestFindWorkspace(t)
	opts := FindFilesOptions{
		Query:     "01",
		Kind:      MatchStartsWith,
		Exts:      []string{"ts", "md"},
		TargetDir: ws,
	}

	matches, err := scanAndMatchFiles(opts)
	if err != nil {
		t.Fatalf("scanAndMatchFiles error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match starting with 01, got %d: %v", len(matches), matches)
	}
}

func TestFindFilesEndsWith(t *testing.T) {
	ws := setupTestFindWorkspace(t)
	opts := FindFilesOptions{
		Query:     "_test.go",
		Kind:      MatchEndsWith,
		Exts:      []string{"go"},
		TargetDir: ws,
	}

	matches, err := scanAndMatchFiles(opts)
	if err != nil {
		t.Fatalf("scanAndMatchFiles error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match ending with _test.go, got %d: %v", len(matches), matches)
	}
}

func TestFindWildcard(t *testing.T) {
	ws := setupTestFindWorkspace(t)
	opts := FindFilesOptions{
		Query:     "*runner*",
		Kind:      MatchWildcard,
		Exts:      []string{"py"},
		TargetDir: ws,
	}

	matches, err := scanAndMatchFiles(opts)
	if err != nil {
		t.Fatalf("scanAndMatchFiles error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match for wildcard *runner*, got %d: %v", len(matches), matches)
	}
}
