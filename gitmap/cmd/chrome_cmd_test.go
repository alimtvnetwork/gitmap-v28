// Package cmd — chrome_cmd_test.go: unit and routing tests for the unified `gitmap chrome` CLI.
package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChromeHelpAndZeroArgs(t *testing.T) {
	if err := runChrome([]string{}); err != nil {
		t.Fatalf("expected nil on zero args, got %v", err)
	}
	if err := runChrome([]string{"--help"}); err != nil {
		t.Fatalf("expected nil on --help, got %v", err)
	}
	if err := runChrome([]string{"help"}); err != nil {
		t.Fatalf("expected nil on help, got %v", err)
	}
}

func TestChromeInstallDryRun(t *testing.T) {
	if err := runChrome([]string{"install", "--dry-run"}); err != nil {
		t.Fatalf("expected nil on install --dry-run, got %v", err)
	}
	if err := runChrome([]string{"in", "--dry-run"}); err != nil {
		t.Fatalf("expected nil on in --dry-run, got %v", err)
	}
}

func TestChromeListEmpty(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempDir)
	if err := runChrome([]string{"list"}); err != nil {
		t.Fatalf("expected nil on list with empty user data, got %v", err)
	}
	if err := runChrome([]string{"ls"}); err != nil {
		t.Fatalf("expected nil on ls alias, got %v", err)
	}
}

func TestChromeBatchOperations(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempDir)
	seedChromeProfileTree(t, tempDir)

	// Test copy-all
	outCopyDir := filepath.Join(t.TempDir(), "copied")
	if err := runChrome([]string{"copy-all", outCopyDir}); err != nil {
		t.Fatalf("copy-all failed: %v", err)
	}

	// Test export-all
	outExportDir := filepath.Join(t.TempDir(), "exported")
	if err := runChrome([]string{"export-all", outExportDir}); err != nil {
		t.Fatalf("export-all failed: %v", err)
	}

	// Test import-all
	if err := runChrome([]string{"import-all", outExportDir}); err != nil {
		t.Fatalf("import-all failed: %v", err)
	}
}

func TestChromeSubcommandDispatch(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempDir)
	seedChromeProfileTree(t, tempDir)

	// Test which
	if err := runChrome([]string{"which"}); err != nil {
		t.Fatalf("which failed: %v", err)
	}

	// Test bookmarks export
	outBM := filepath.Join(t.TempDir(), "bm.json")
	if err := runChrome([]string{"bookmarks", "Default", "--out", outBM, "--format", "json"}); err != nil {
		t.Fatalf("bookmarks export failed: %v", err)
	}
}

func TestChromeSmartExportAndImport(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempDir)
	seedChromeProfileTree(t, tempDir)

	workDir := t.TempDir()

	// 1. Export all to a.json with explicit --format=json
	outJSON := filepath.Join(workDir, "a.json")
	if err := runChromeProfileExport([]string{outJSON, "--format=json"}); err != nil {
		t.Fatalf("export to a.json failed: %v", err)
	}
	if _, err := os.Stat(outJSON); err != nil {
		t.Fatalf("expected a.json to exist: %v", err)
	}

	// 2. Export all to a.db with auto-inferred SQLite format
	outDB := filepath.Join(workDir, "a.db")
	if err := runChromeProfileExport([]string{outDB}); err != nil {
		t.Fatalf("export to a.db failed: %v", err)
	}
	if _, err := os.Stat(outDB); err != nil {
		t.Fatalf("expected a.db to exist: %v", err)
	}

	// 3. Export all to a.sqlite with explicit --format=sqlite
	outSQLite := filepath.Join(workDir, "a.sqlite")
	if err := runChromeProfileExport([]string{outSQLite, "--format=sqlite"}); err != nil {
		t.Fatalf("export to a.sqlite failed: %v", err)
	}

	// 4. Export all to a.yaml
	outYAML := filepath.Join(workDir, "a.yaml")
	if err := runChromeProfileExport([]string{outYAML}); err != nil {
		t.Fatalf("export to a.yaml failed: %v", err)
	}

	// 5. Export single profile Default to single.json
	outSingleJSON := filepath.Join(workDir, "single.json")
	if err := runChromeProfileExport([]string{"Default", outSingleJSON}); err != nil {
		t.Fatalf("export Default to single.json failed: %v", err)
	}

	// 6. Test imports
	if err := runChromeProfileImport([]string{outJSON}); err != nil {
		t.Fatalf("import from a.json failed: %v", err)
	}
	if err := runChromeProfileImport([]string{outDB}); err != nil {
		t.Fatalf("import from a.db failed: %v", err)
	}
	if err := runChromeProfileImport([]string{outYAML}); err != nil {
		t.Fatalf("import from a.yaml failed: %v", err)
	}
}
