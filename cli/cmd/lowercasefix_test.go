package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRenameCandidates(t *testing.T) {
	tempDir := t.TempDir()
	readmePath := filepath.Join(tempDir, "README.md")
	_ = os.WriteFile(readmePath, []byte("# Hello"), 0644)
	subDir := filepath.Join(tempDir, "docs")
	_ = os.MkdirAll(subDir, 0755)
	subPath := filepath.Join(subDir, "GUIDE.MD")
	_ = os.WriteFile(subPath, []byte("# Guide"), 0644)
	lowerPath := filepath.Join(tempDir, "lower.txt")
	_ = os.WriteFile(lowerPath, []byte("ok"), 0644)

	opts := LowerCaseFixOptions{Patterns: []string{"*.md"}}
	pairs, scanned, err := findRenameCandidates(tempDir, opts)
	if err != nil {
		t.Fatalf("findRenameCandidates error: %v", err)
	}
	if scanned < 3 {
		t.Errorf("expected at least 3 files scanned, got %d", scanned)
	}
	if len(pairs) != 2 {
		t.Errorf("expected 2 rename candidates, got %d", len(pairs))
	}
}

func TestPerformTwoStepFSRename(t *testing.T) {
	tempDir := t.TempDir()
	src := filepath.Join(tempDir, "FILE.txt")
	tmp := filepath.Join(tempDir, "FILE.txt.tmp-lcf")
	dst := filepath.Join(tempDir, "file.txt")
	_ = os.WriteFile(src, []byte("content"), 0644)

	err := performTwoStepFSRename(src, tmp, dst)
	if err != nil {
		t.Fatalf("performTwoStepFSRename error: %v", err)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Errorf("expected %s to exist: %v", dst, err)
	}
}

func TestExecuteLowerCaseFix_Workflow(t *testing.T) {
	tempDir := t.TempDir()
	origWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(origWd) }()
	_ = os.Chdir(tempDir)

	readmePath := filepath.Join(tempDir, "README.md")
	_ = os.WriteFile(readmePath, []byte("# Readme"), 0644)
	llmPath := filepath.Join(tempDir, "LLM.md")
	_ = os.WriteFile(llmPath, []byte("# LLM"), 0644)
	lowerPath := filepath.Join(tempDir, "already_lower.md")
	_ = os.WriteFile(lowerPath, []byte("# Lower"), 0644)

	// Dry run test
	optsDry := LowerCaseFixOptions{
		Patterns:   []string{"*.md"},
		IsDryRun:   true,
		IsNoCommit: true,
	}
	if err := ExecuteLowerCaseFix(optsDry); err != nil {
		t.Fatalf("ExecuteLowerCaseFix dry-run failed: %v", err)
	}
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Errorf("expected README.md to still exist after dry-run")
	}

	// Real execution test (no commit since tempDir isn't a git repo)
	optsReal := LowerCaseFixOptions{
		Patterns:   []string{"*.md"},
		IsDryRun:   false,
		IsNoCommit: true,
	}
	if err := ExecuteLowerCaseFix(optsReal); err != nil {
		t.Fatalf("ExecuteLowerCaseFix real failed: %v", err)
	}

	newReadme := filepath.Join(tempDir, "readme.md")
	if _, err := os.Stat(newReadme); err != nil {
		t.Errorf("expected readme.md to exist: %v", err)
	}
	newLLM := filepath.Join(tempDir, "llm.md")
	if _, err := os.Stat(newLLM); err != nil {
		t.Errorf("expected llm.md to exist: %v", err)
	}
}
