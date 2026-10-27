package store

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func openTestMemoryDB(t *testing.T) *sql.DB {
	t.Helper()

	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}

	return conn
}

func hasColumn(t *testing.T, conn *sql.DB, table, col string) bool {
	t.Helper()

	rows, err := conn.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		t.Fatalf("pragma failed: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt any
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err == nil && name == col {
			return true
		}
	}

	return false
}

func TestEnsureHeatmapTables_Fresh(t *testing.T) {
	conn := openTestMemoryDB(t)
	defer conn.Close()

	if err := EnsureHeatmapTables(conn); err != nil {
		t.Fatalf("EnsureHeatmapTables failed on fresh db: %v", err)
	}

	if !hasColumn(t, conn, "HeatmapActionLog", "ModelName") {
		t.Errorf("expected ModelName column to exist")
	}

	if !hasColumn(t, conn, "HeatmapActionLog", "ToolName") {
		t.Errorf("expected ToolName column to exist")
	}
}

func TestEnsureHeatmapTables_LegacyMigration(t *testing.T) {
	conn := openTestMemoryDB(t)
	defer conn.Close()

	legacyDDL := `CREATE TABLE HeatmapActionLog (
		HeatmapActionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
		ParentLogId INTEGER NULL,
		TaskSlug TEXT NOT NULL,
		RepoName TEXT NOT NULL DEFAULT '',
		ActionType TEXT NOT NULL,
		FilePath TEXT NOT NULL DEFAULT '',
		StartLine INTEGER NOT NULL DEFAULT 0,
		EndLine INTEGER NOT NULL DEFAULT 0,
		CodeSnippet TEXT NULL,
		CommandText TEXT NULL,
		ReasoningState TEXT NOT NULL DEFAULT 'EXPLORATION',
		ReasoningNotes TEXT NULL,
		ExpectedOutcome TEXT NULL,
		ActualOutcome TEXT NULL,
		IsExpectationMet INTEGER NOT NULL DEFAULT 1,
		ExecutionDepth INTEGER NOT NULL DEFAULT 0,
		ExecutionOrder INTEGER NOT NULL DEFAULT 1,
		Notes TEXT NULL,
		Comments TEXT NULL,
		CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := conn.Exec(legacyDDL); err != nil {
		t.Fatalf("failed to create legacy table: %v", err)
	}

	if err := EnsureHeatmapTables(conn); err != nil {
		t.Fatalf("EnsureHeatmapTables failed on legacy migration: %v", err)
	}

	if !hasColumn(t, conn, "HeatmapActionLog", "ModelName") {
		t.Errorf("expected ModelName column after legacy migration")
	}

	if !hasColumn(t, conn, "HeatmapActionLog", "ToolName") {
		t.Errorf("expected ToolName column after legacy migration")
	}
}

func TestEnsureRcHeatmapRepoTables_Fresh(t *testing.T) {
	conn := openTestMemoryDB(t)
	defer conn.Close()

	if err := EnsureRcHeatmapRepoTables(conn); err != nil {
		t.Fatalf("EnsureRcHeatmapRepoTables failed: %v", err)
	}

	if !hasColumn(t, conn, "HeatmapActionLog", "ModelName") {
		t.Errorf("expected ModelName column in rc repo db")
	}

	if !hasColumn(t, conn, "HeatmapActionLog", "ToolName") {
		t.Errorf("expected ToolName column in rc repo db")
	}
}

func TestEnsureRcHeatmapRepoTables_LegacyMigration(t *testing.T) {
	conn := openTestMemoryDB(t)
	defer conn.Close()

	legacyDDL := `CREATE TABLE HeatmapActionLog (
		HeatmapActionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
		ParentLogId INTEGER NULL,
		TaskSlug TEXT NOT NULL,
		RepoSlug TEXT NOT NULL DEFAULT '',
		ActionType TEXT NOT NULL,
		FilePath TEXT NOT NULL DEFAULT '',
		StartLine INTEGER NOT NULL DEFAULT 0,
		EndLine INTEGER NOT NULL DEFAULT 0,
		CodeSnippet TEXT NULL,
		CommandText TEXT NULL,
		ReasoningState TEXT NOT NULL DEFAULT 'EXPLORATION',
		ReasoningNotes TEXT NULL,
		ExpectedOutcome TEXT NULL,
		ActualOutcome TEXT NULL,
		IsExpectationMet INTEGER NOT NULL DEFAULT 1,
		ExecutionDepth INTEGER NOT NULL DEFAULT 0,
		ExecutionOrder INTEGER NOT NULL DEFAULT 1,
		DurationMs INTEGER NOT NULL DEFAULT 0,
		Notes TEXT NULL,
		Comments TEXT NULL,
		CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := conn.Exec(legacyDDL); err != nil {
		t.Fatalf("failed to create legacy rc table: %v", err)
	}

	if err := EnsureRcHeatmapRepoTables(conn); err != nil {
		t.Fatalf("EnsureRcHeatmapRepoTables failed on legacy migration: %v", err)
	}

	if !hasColumn(t, conn, "HeatmapActionLog", "ModelName") {
		t.Errorf("expected ModelName column after rc legacy migration")
	}

	if !hasColumn(t, conn, "HeatmapActionLog", "ToolName") {
		t.Errorf("expected ToolName column after rc legacy migration")
	}
}

func TestOpenHeatmapDBAtPath_CreatesAndMigrates(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "heatmap.db")

	db, err := OpenHeatmapDBAtPath(dbPath)
	if err != nil {
		t.Fatalf("openHeatmapDBAtPath failed: %v", err)
	}

	defer db.Close()

	if !hasColumn(t, db.Conn(), "HeatmapActionLog", "ModelName") {
		t.Errorf("expected ModelName column in opened heatmap DB")
	}
}

func TestRecordParentTask_Upsert(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "heatmap.db")

	db, err := OpenHeatmapDBAtPath(dbPath)
	if err != nil {
		t.Fatalf("openHeatmapDBAtPath failed: %v", err)
	}
	defer db.Close()

	task := HeatmapParentTaskRecord{
		TaskSlug:         "test-task-slug",
		TaskName:         "Test Task Name",
		RunDirectory:     ".ai-memory/temp-agents/01-test",
		RootDbPath:       dbPath,
		RootSpecFilePath: "02-spec/21-app/01-test/spec.md",
		TaskFolder:       "02-spec/21-app/01-test/",
		Status:           "IN_PROGRESS",
	}

	if err := db.RecordParentTask(task); err != nil {
		t.Fatalf("RecordParentTask failed: %v", err)
	}

	var name, specFile, specFolder string
	query := "SELECT TaskName, RootSpecFilePath, TaskFolder FROM HeatmapParentTask WHERE TaskSlug = ?"
	if err := db.Conn().QueryRow(query, "test-task-slug").Scan(&name, &specFile, &specFolder); err != nil {
		t.Fatalf("QueryRow failed: %v", err)
	}

	if name != "Test Task Name" || specFile != "02-spec/21-app/01-test/spec.md" || specFolder != "02-spec/21-app/01-test/" {
		t.Errorf("mismatched task record: got name=%s, specFile=%s, specFolder=%s", name, specFile, specFolder)
	}
}
