package cmddb

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestRegistry(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "gitmap-alias-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	testDbPath := filepath.Join(tmpDir, "gitmap.db")
	t.Setenv("GITMAP_DATA_DIR", tmpDir)
	t.Setenv("LOCALAPPDATA", tmpDir)

	return testDbPath
}

func cleanupTestRegistry(tmpDir string) {
	_ = os.RemoveAll(tmpDir)
}

func TestInitAliasRegistry(t *testing.T) {
	tmpDir := setupTestRegistry(t)
	defer cleanupTestRegistry(filepath.Dir(tmpDir))

	err := InitAliasRegistry()
	if err != nil {
		t.Fatalf("InitAliasRegistry failed: %v", err)
	}

	aliases, err := ListAliases()
	if err != nil {
		t.Fatalf("ListAliases failed: %v", err)
	}

	if len(aliases) < 4 {
		t.Fatalf("expected at least 4 system aliases, got %d", len(aliases))
	}
}

func TestRegisterAlias(t *testing.T) {
	tmpDir := setupTestRegistry(t)
	defer cleanupTestRegistry(filepath.Dir(tmpDir))

	dummyDb := filepath.Join(filepath.Dir(tmpDir), "custom.db")
	_ = os.WriteFile(dummyDb, []byte("SQLite format 3\x00"), 0644)

	err := RegisterAlias("mycustom", dummyDb, "Custom database description")
	if err != nil {
		t.Fatalf("RegisterAlias failed: %v", err)
	}

	resolved, err := ResolveDbPathWithRegistry("mycustom")
	if err != nil {
		t.Fatalf("ResolveDbPathWithRegistry failed: %v", err)
	}

	cleanExpected := filepath.Clean(dummyDb)
	if resolved != cleanExpected {
		t.Fatalf("expected resolved path %s, got %s", cleanExpected, resolved)
	}
}

func TestRegisterAliasValidation(t *testing.T) {
	tmpDir := setupTestRegistry(t)
	defer cleanupTestRegistry(filepath.Dir(tmpDir))

	if err := RegisterAlias("", "path/to/db", "desc"); err == nil {
		t.Fatalf("expected error on empty alias name")
	}

	if err := RegisterAlias("valid", "", "desc"); err == nil {
		t.Fatalf("expected error on empty database path")
	}
}

func TestResolveDbPathWithRegistry(t *testing.T) {
	tmpDir := setupTestRegistry(t)
	defer cleanupTestRegistry(filepath.Dir(tmpDir))

	masterPath, err := ResolveDbPathWithRegistry("master")
	if err != nil {
		t.Fatalf("failed to resolve master: %v", err)
	}
	if !strings.HasSuffix(masterPath, "gitmap.db") {
		t.Fatalf("expected master path to end with gitmap.db, got %s", masterPath)
	}

	pipelinePath, err := ResolveDbPathWithRegistry("pipeline")
	if err != nil {
		t.Fatalf("failed to resolve pipeline: %v", err)
	}
	if pipelinePath == "" {
		t.Fatalf("expected non-empty pipeline path")
	}

	rawPath, err := ResolveDbPathWithRegistry("some/raw/path.db")
	if err != nil {
		t.Fatalf("failed to resolve raw path: %v", err)
	}
	if rawPath != filepath.Clean("some/raw/path.db") {
		t.Fatalf("expected clean raw path, got %s", rawPath)
	}
}

func TestIncrementAliasSearchCount(t *testing.T) {
	tmpDir := setupTestRegistry(t)
	defer cleanupTestRegistry(filepath.Dir(tmpDir))

	dummyDb := filepath.Join(filepath.Dir(tmpDir), "search_test.db")
	_ = os.WriteFile(dummyDb, []byte("SQLite format 3\x00"), 0644)

	_ = RegisterAlias("searchtest", dummyDb, "Testing search counts")

	IncrementAliasSearchCount("searchtest")
	IncrementAliasSearchCount("searchtest")

	aliases, err := ListAliases()
	if err != nil {
		t.Fatalf("ListAliases failed: %v", err)
	}

	foundCount := findAliasSearchCount(aliases, "searchtest")
	if foundCount != 2 {
		t.Fatalf("expected search count 2, got %d", foundCount)
	}
}

func findAliasSearchCount(aliases []DatabaseAlias, targetName string) int64 {
	for _, a := range aliases {
		if strings.EqualFold(a.AliasName, targetName) {
			return a.SearchHistoryCount
		}
	}
	return -1
}
