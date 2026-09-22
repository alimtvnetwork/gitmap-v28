package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectDataTreeDatabases(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-tree-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tasksDir := filepath.Join(tempDir, "tasks")
	_ = os.MkdirAll(tasksDir, 0755)
	dummyDb := filepath.Join(tasksDir, "sql.db")
	_ = os.WriteFile(dummyDb, []byte(""), 0644)

	var found []SplitDatabaseEntry
	CollectDataTreeDatabases(tempDir, func(e SplitDatabaseEntry) {
		found = append(found, e)
	})

	if len(found) != 1 {
		t.Fatalf("expected 1 database found, got %d", len(found))
	}

	if found[0].DatabaseType != "tasks" {
		t.Errorf("expected database type tasks, got %s", found[0].DatabaseType)
	}

	if found[0].DatabaseKey != "tasks" {
		t.Errorf("expected database key tasks, got %s", found[0].DatabaseKey)
	}
}

func TestClassifyDiscoveredDb(t *testing.T) {
	dbType, key, _ := classifyDiscoveredDb("/data/search/sql.db", "sql.db")
	if dbType != "search" || key != "search" {
		t.Errorf("expected search/search, got %s/%s", dbType, key)
	}

	dbType, key, _ = classifyDiscoveredDb("/data/ssh/ssh-tasks.db", "ssh-tasks.db")
	if dbType != "tasks" || key != "ssh-tasks" {
		t.Errorf("expected tasks/ssh-tasks, got %s/%s", dbType, key)
	}

	dbType, key, _ = classifyDiscoveredDb("/data/pipeline/my-repo/sql.db", "sql.db")
	if dbType != "pipeline" || key != "my-repo" {
		t.Errorf("expected pipeline/my-repo, got %s/%s", dbType, key)
	}
}
