package constants

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestDriverSQLiteConstant(t *testing.T) {
	if DriverSQLite != "sqlite" {
		t.Errorf("expected DriverSQLite to be 'sqlite', got %s", DriverSQLite)
	}
}

func TestApplyDefaultPragmas_Success(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_pragma.db")

	conn, err := sql.Open(DriverSQLite, dbPath)
	if err != nil {
		t.Fatalf("failed to open sqlite test db: %v", err)
	}

	defer conn.Close()

	if err := ApplyDefaultPragmas(conn); err != nil {
		t.Fatalf("ApplyDefaultPragmas failed: %v", err)
	}
}

func TestApplyDefaultPragmas_ClosedDbError(t *testing.T) {
	conn, err := sql.Open(DriverSQLite, ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	_ = conn.Close()

	if err := ApplyDefaultPragmas(conn); err == nil {
		t.Errorf("expected error when applying pragmas to closed db, got nil")
	}
}
