package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallVersionJSONSuite(t *testing.T) {
	tempDir := t.TempDir()

	// Create dummy what-to-read.md files
	lovableDir := filepath.Join(tempDir, ".lovable")
	_ = os.MkdirAll(lovableDir, 0755)
	lovableWTR := filepath.Join(lovableDir, "what-to-read.md")
	_ = os.WriteFile(lovableWTR, []byte("# What to Read\n\n## Before writing code\n- some item\n"), 0644)

	cfg := DefaultVersionInstallConfig("2.5.0")
	if err := InstallVersionJSON(tempDir, cfg, false); err != nil {
		t.Fatalf("InstallVersionJSON failed: %v", err)
	}

	// 1. Verify version.json exists and has initial version 2.5.0
	vData, errRead := os.ReadFile(filepath.Join(tempDir, "version.json"))
	if errRead != nil || !strings.Contains(string(vData), "2.5.0") {
		t.Fatalf("version.json missing or invalid: %s (err: %v)", string(vData), errRead)
	}

	// 2. Verify .lovable/versioning.md exists
	docData, errDoc := os.ReadFile(filepath.Join(tempDir, ".lovable", "versioning.md"))
	if errDoc != nil || !strings.Contains(string(docData), "Single Source of Truth") {
		t.Fatalf(".lovable/versioning.md missing or invalid: %v", errDoc)
	}

	// 3. Verify .lovable/memory/learned/01-versioning-ssot.md exists
	memData, errMem := os.ReadFile(filepath.Join(tempDir, ".lovable", "memory", "learned", "01-versioning-ssot.md"))
	if errMem != nil || !strings.Contains(string(memData), "SSOT") {
		t.Fatalf("memory file missing or invalid: %v", errMem)
	}

	// 4. Verify what-to-read.md has been enqueued
	wtrData, _ := os.ReadFile(lovableWTR)
	if !strings.Contains(string(wtrData), "versioning.md") {
		t.Fatalf("what-to-read.md not enqueued with versioning doc: %s", string(wtrData))
	}
}
