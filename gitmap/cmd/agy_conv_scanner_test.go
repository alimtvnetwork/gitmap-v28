package cmd

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	_ "modernc.org/sqlite"
)

func TestReadSingleConvDB(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test-conv-123.db")

	conn, err := store.OpenSQLiteDB(dbPath)
	if err != nil {
		t.Fatalf("OpenSQLiteDB failed: %v", err)
	}
	defer conn.Close()

	initTestConvTables(t, conn)
	info, hasInfo := readSingleConvDB(dbPath, "test-conv-123.db")
	if !hasInfo {
		t.Fatalf("expected hasInfo to be true, got false")
	}

	assertConvInfoValues(t, info, "test-conv-123", 2, 1)
}

func initTestConvTables(t *testing.T, conn *sql.DB) {
	schema := `CREATE TABLE steps (id INTEGER PRIMARY KEY, step_type INTEGER);
	INSERT INTO steps (step_type) VALUES (1), (2);
	CREATE TABLE trajectory_metadata_blob (id TEXT PRIMARY KEY, data BLOB);
	INSERT INTO trajectory_metadata_blob (id, data) VALUES ('main', 'repo/projects/myrepo');`
	if _, err := conn.Exec(schema); err != nil {
		t.Fatalf("init test conv tables failed: %v", err)
	}
}

func assertConvInfoValues(t *testing.T, info AgyConvInfo, wantID string, wantSteps, wantUserSteps int) {
	if info.ID != wantID {
		t.Errorf("ID mismatch: got %q, want %q", info.ID, wantID)
	}
	if info.StepCount != wantSteps {
		t.Errorf("StepCount mismatch: got %d, want %d", info.StepCount, wantSteps)
	}
	if info.UserSteps != wantUserSteps {
		t.Errorf("UserSteps mismatch: got %d, want %d", info.UserSteps, wantUserSteps)
	}
}
