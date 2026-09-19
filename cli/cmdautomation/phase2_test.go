package cmdautomation

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

func TestTopologyDiscovery(t *testing.T) {
	opts := TopologyOptions{
		Dir:       ".",
		IsRefresh: true,
		TtlSec:    60,
	}
	monad := RunTopology(opts)
	if monad.IsFailure() {
		t.Fatalf("RunTopology failed: %v", monad.Err)
	}
	res := monad.Value
	if !res.IsValid {
		t.Errorf("expected res.IsValid to be true")
	}
	if res.TotalFiles <= 0 {
		t.Errorf("expected TotalFiles > 0, got %d", res.TotalFiles)
	}
}

func TestDbGenerate(t *testing.T) {
	dbPath := setupTestDb(t)
	defer os.Remove(dbPath)

	opts := DbGenerateOptions{
		DbPath:   dbPath,
		Lang:     "all",
		IsDryRun: true,
	}
	monad := RunDbGenerate(opts)
	if monad.IsFailure() {
		t.Fatalf("RunDbGenerate failed: %v", monad.Err)
	}
	res := monad.Value
	if !res.IsSuccess || res.TableCount != 1 {
		t.Errorf("expected success with 1 table, got %+v", res)
	}
}

func TestDbMigrate(t *testing.T) {
	dbPath := setupTestDb(t)
	defer os.Remove(dbPath)

	opts := DbMigrateOptions{
		DbPath:    dbPath,
		SqlScript: "CREATE TABLE UserSessions (SessionId TEXT PRIMARY KEY, IsActive INTEGER);",
		IsDryRun:  false,
	}
	monad := RunDbMigrate(opts)
	if monad.IsFailure() {
		t.Fatalf("RunDbMigrate failed: %v", monad.Err)
	}
	res := monad.Value
	if !res.IsSuccess || res.TotalApplied != 1 {
		t.Errorf("expected 1 applied migration, got %+v", res)
	}
}

func TestSchemaAudit(t *testing.T) {
	dbPath := setupTestDb(t)
	defer os.Remove(dbPath)

	opts := SchemaAuditOptions{
		DbPath:   dbPath,
		IsStrict: true,
	}
	monad := RunSchemaAudit(opts)
	if monad.IsFailure() {
		t.Fatalf("RunSchemaAudit failed: %v", monad.Err)
	}
	res := monad.Value
	if !res.IsClean {
		t.Errorf("expected clean schema for valid table, got violations: %+v", res.Violations)
	}
}

func setupTestDb(t *testing.T) string {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open test sqlite db: %v", err)
	}
	defer db.Close()
	createSampleTable(t, db)
	return dbPath
}

func createSampleTable(t *testing.T, db *sql.DB) {
	stmt := `CREATE TABLE UserAccounts (
		UserAccountsId INTEGER PRIMARY KEY AUTOINCREMENT,
		Username TEXT NOT NULL,
		IsActive BOOLEAN NOT NULL DEFAULT 1
	);`
	if _, err := db.Exec(stmt); err != nil {
		t.Fatalf("failed to create sample table: %v", err)
	}
}
