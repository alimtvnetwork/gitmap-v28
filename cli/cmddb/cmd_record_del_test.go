package cmddb

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestRecordDelDb(t *testing.T) (string, *sql.DB) {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test_record_del.db")

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open sqlite test db: %v", err)
	}

	createSql := `CREATE TABLE Subtask (
		SubtaskId INTEGER PRIMARY KEY AUTOINCREMENT,
		TaskId TEXT NOT NULL,
		Title TEXT NOT NULL,
		Status TEXT NOT NULL
	);`
	if _, err := conn.Exec(createSql); err != nil {
		t.Fatalf("failed to create Subtask table: %v", err)
	}

	seedInitialSubtasks(t, conn)
	return dbPath, conn
}

func seedInitialSubtasks(t *testing.T, conn *sql.DB) {
	t.Helper()
	insertSql := `INSERT INTO Subtask (SubtaskId, TaskId, Title, Status) VALUES 
		(1, 'task-01', 'First Subtask', 'PENDING'),
		(2, 'task-01', 'Second Subtask', 'DONE'),
		(3, 'task-02', 'Third Subtask', 'IN_PROGRESS');`
	if _, err := conn.Exec(insertSql); err != nil {
		t.Fatalf("failed to insert initial subtasks: %v", err)
	}
}

func countSubtaskRows(t *testing.T, conn *sql.DB, id int) int {
	t.Helper()
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM Subtask WHERE SubtaskId = ?", id).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count subtask rows: %v", err)
	}
	return count
}

func TestCli_DbRecordDel_Success(t *testing.T) {
	dbPath, conn := setupTestRecordDelDb(t)
	defer conn.Close()

	if count := countSubtaskRows(t, conn, 1); count != 1 {
		t.Fatalf("expected 1 row before deletion, got %d", count)
	}

	appErr := RunDbRecordDel([]string{dbPath, "-table", "Subtask", "-id", "1", "-y"})
	if appErr != nil {
		t.Fatalf("expected nil error from RunDbRecordDel, got %v", appErr)
	}

	if count := countSubtaskRows(t, conn, 1); count != 0 {
		t.Fatalf("expected 0 rows after deletion of id=1, got %d", count)
	}

	if count := countSubtaskRows(t, conn, 2); count != 1 {
		t.Fatalf("expected row id=2 to remain, got %d", count)
	}
}

func TestCli_DbRecordDel_PositionalArgs(t *testing.T) {
	dbPath, conn := setupTestRecordDelDb(t)
	defer conn.Close()

	appErr := RunDbRecordDel([]string{dbPath, "Subtask", "2", "-y"})
	if appErr != nil {
		t.Fatalf("expected nil error from positional RunDbRecordDel, got %v", appErr)
	}

	if count := countSubtaskRows(t, conn, 2); count != 0 {
		t.Fatalf("expected row id=2 to be deleted, got %d", count)
	}
}

func TestCli_DbRecordDel_JsonOutput(t *testing.T) {
	dbPath, conn := setupTestRecordDelDb(t)
	defer conn.Close()

	appErr := RunDbRecordDel([]string{dbPath, "-table", "Subtask", "-id", "3", "-y", "--json"})
	if appErr != nil {
		t.Fatalf("expected nil error with --json flag, got %v", appErr)
	}

	if count := countSubtaskRows(t, conn, 3); count != 0 {
		t.Fatalf("expected row id=3 to be deleted, got %d", count)
	}
}

func TestCli_DbRecordDel_NotFound(t *testing.T) {
	dbPath, conn := setupTestRecordDelDb(t)
	defer conn.Close()

	appErr := RunDbRecordDel([]string{dbPath, "-table", "Subtask", "-id", "9999", "-y"})
	if appErr == nil {
		t.Fatalf("expected not found error for non-existent record, got nil")
	}
}

func TestCli_DbRecordDel_MissingArgs(t *testing.T) {
	appErr := RunDbRecordDel([]string{})
	if appErr != nil {
		t.Fatalf("expected nil (help displayed) when args are empty, got %v", appErr)
	}

	appErr = RunDbRecordDel([]string{"mydb"})
	if appErr == nil {
		t.Fatalf("expected validation error when table and id are omitted, got nil")
	}
}
