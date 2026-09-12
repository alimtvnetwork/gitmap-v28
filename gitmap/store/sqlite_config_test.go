package store

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func queryPragmaInt(t *testing.T, conn *sql.DB, query string) int {
	t.Helper()
	var val int
	if err := conn.QueryRow(query).Scan(&val); err != nil {
		t.Fatalf("query %s failed: %v", query, err)
	}

	return val
}

func queryPragmaString(t *testing.T, conn *sql.DB, query string) string {
	t.Helper()
	var val string
	if err := conn.QueryRow(query).Scan(&val); err != nil {
		t.Fatalf("query %s failed: %v", query, err)
	}

	return val
}

func verifyPragmaValues(t *testing.T, conn *sql.DB) {
	t.Helper()
	if timeout := queryPragmaInt(t, conn, "PRAGMA busy_timeout;"); timeout != 5000 {
		t.Errorf("expected busy_timeout 5000, got %d", timeout)
	}

	if fk := queryPragmaInt(t, conn, "PRAGMA foreign_keys;"); fk != 1 {
		t.Errorf("expected foreign_keys 1, got %d", fk)
	}

	if syncMode := queryPragmaInt(t, conn, "PRAGMA synchronous;"); syncMode != 1 {
		t.Errorf("expected synchronous 1 (NORMAL), got %d", syncMode)
	}

	if journal := queryPragmaString(t, conn, "PRAGMA journal_mode;"); journal != "wal" {
		t.Errorf("expected journal_mode wal, got %s", journal)
	}
}

func TestConfigureSQLiteConn_PragmasAndPooling(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_config.db")

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open sqlite test db: %v", err)
	}
	defer conn.Close()

	if err := ConfigureSQLiteConn(conn); err != nil {
		t.Fatalf("ConfigureSQLiteConn failed: %v", err)
	}

	if maxOpen := conn.Stats().MaxOpenConnections; maxOpen != 1 {
		t.Errorf("expected MaxOpenConnections 1, got %d", maxOpen)
	}

	verifyPragmaValues(t, conn)
}

func TestConfigureSQLiteConn_ClosedDbError(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	_ = conn.Close()

	if err := ConfigureSQLiteConn(conn); err == nil {
		t.Errorf("expected error when configuring closed connection, got nil")
	}
}

func TestOpenSQLiteDB_Success(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_open_factory.db")

	conn, appErr := OpenSQLiteDB(dbPath)
	if appErr != nil {
		t.Fatalf("OpenSQLiteDB failed: %v", appErr)
	}
	defer conn.Close()

	if maxOpen := conn.Stats().MaxOpenConnections; maxOpen != 1 {
		t.Errorf("expected MaxOpenConnections 1, got %d", maxOpen)
	}

	verifyPragmaValues(t, conn)
}

func TestOpenSQLiteDB_ConfigError(t *testing.T) {
	invalidPath := filepath.Join(t.TempDir(), "nonexistent_dir", "sub", "test.db")

	conn, appErr := OpenSQLiteDB(invalidPath)
	if appErr == nil {
		conn.Close()
		t.Fatalf("expected error for invalid directory path, got nil")
	}

	if conn != nil {
		t.Errorf("expected nil connection on error, got %v", conn)
	}
}
