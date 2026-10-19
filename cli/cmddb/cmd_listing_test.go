package cmddb

import (
	"path/filepath"
	"testing"
)

func TestRunDbListingHelp(t *testing.T) {
	appErr := RunDbListing([]string{"--help"})
	if appErr != nil {
		t.Fatalf("expected nil for help, got: %v", appErr)
	}
}

func TestRunDbListingDefault(t *testing.T) {
	tmpDir := setupTestRegistry(t)
	defer cleanupTestRegistry(filepath.Dir(tmpDir))

	_ = InitAliasRegistry()

	appErr := RunDbListing(nil)
	if appErr != nil {
		t.Fatalf("RunDbListing failed: %v", appErr)
	}
}

func TestRunDbListingJson(t *testing.T) {
	tmpDir := setupTestRegistry(t)
	defer cleanupTestRegistry(filepath.Dir(tmpDir))

	_ = InitAliasRegistry()

	appErr := RunDbListing([]string{"--json"})
	if appErr != nil {
		t.Fatalf("RunDbListing --json failed: %v", appErr)
	}
}

func TestRunDbListingAllAndRefresh(t *testing.T) {
	tmpDir := setupTestRegistry(t)
	defer cleanupTestRegistry(filepath.Dir(tmpDir))

	_ = InitAliasRegistry()

	appErr := RunDbListing([]string{"--all", "--refresh"})
	if appErr != nil {
		t.Fatalf("RunDbListing --all --refresh failed: %v", appErr)
	}
}
