package cmddb

import (
	"os"
	"path/filepath"
	"testing"
)

func createDummyDbFileForRm(t *testing.T, filename string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte("SQLite format 3\x00"), 0644); err != nil {
		t.Fatalf("failed to create dummy db file: %v", err)
	}
	return path
}

func TestCli_DbRm_ProtectedSystemAliases(t *testing.T) {
	systemAliases := []string{"master", "gitmap", "pipeline", "pipelines", "heatmap", "tasks", "agents"}
	for _, alias := range systemAliases {
		appErr := RunDbRm([]string{alias})
		if appErr == nil {
			t.Fatalf("expected error when trying to remove protected system alias '%s', got nil", alias)
		}
		if appErr.Code != "E4030" {
			t.Fatalf("expected error code E4030 for alias '%s', got %s", alias, appErr.Code)
		}
	}
}

func TestCli_DbRm_CustomAliasUnregistered(t *testing.T) {
	testFile := createDummyDbFileForRm(t, "custom_unbind.db")
	alias := "unit_test_rm_custom"

	if err := RegisterAlias(alias, testFile, "Testing alias unbinding"); err != nil {
		t.Fatalf("failed to register alias: %v", err)
	}

	appErr := RunDbRm([]string{alias, "-y"})
	if appErr != nil {
		t.Fatalf("expected nil error on removing alias '%s', got %v", alias, appErr)
	}

	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatalf("database file should NOT be removed when purge flag is absent")
	}

	_, findErr := findAliasPathInRegistry(alias)
	if findErr == nil {
		t.Fatalf("expected alias '%s' to be removed from registry, but it was found", alias)
	}
}

func TestCli_DbRm_PurgeFlagDeletesFile(t *testing.T) {
	testFile := createDummyDbFileForRm(t, "custom_purge.db")
	alias := "unit_test_rm_purge"

	if err := RegisterAlias(alias, testFile, "Testing alias file purge"); err != nil {
		t.Fatalf("failed to register alias: %v", err)
	}

	appErr := RunDbRm([]string{alias, "--purge", "-y"})
	if appErr != nil {
		t.Fatalf("expected nil error on purging alias '%s', got %v", alias, appErr)
	}

	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Fatalf("database file '%s' should have been purged from disk", testFile)
	}
}

func TestCli_DbRm_MissingOrUnknownAlias(t *testing.T) {
	appErr := RunDbRm([]string{})
	if appErr != nil {
		t.Fatalf("expected nil (help displayed) when args are empty, got %v", appErr)
	}

	appErr = RunDbRm([]string{"non_existent_alias_xyz_12345", "-y"})
	if appErr == nil {
		t.Fatalf("expected error when removing non-existent alias, got nil")
	}
}
