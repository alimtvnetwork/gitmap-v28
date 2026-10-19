package cmddb

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func setupTestDb(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite test db: %v", err)
	}
	defer conn.Close()

	schema := `
	CREATE TABLE projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		status TEXT NOT NULL
	);
	CREATE TABLE tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER,
		title TEXT NOT NULL,
		priority TEXT NOT NULL,
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);
	INSERT INTO projects (name, status) VALUES ('gitmap', 'ACTIVE'), ('alpha', 'PAUSED');
	INSERT INTO tasks (project_id, title, priority) VALUES (1, 'Initial task', 'HIGH'), (1, 'Second task', 'LOW');
	`
	if _, err := conn.Exec(schema); err != nil {
		t.Fatalf("setup test schema: %v", err)
	}
	return dbPath
}

func TestRunDbAdd_CsvAndKeyVal(t *testing.T) {
	dbPath := setupTestDb(t)

	// Add via CSV
	err := RunDbAdd([]string{dbPath, "projects", "csv", "beta", "ACTIVE"})
	if err != nil {
		t.Fatalf("RunDbAdd with CSV failed: %v", err)
	}

	// Add via key-value pairs
	err = RunDbAdd([]string{dbPath, "projects", "name=gamma", "status=PAUSED"})
	if err != nil {
		t.Fatalf("RunDbAdd with key-value failed: %v", err)
	}

	// Missing args should fail with validation error
	err = RunDbAdd([]string{dbPath})
	if err == nil {
		t.Fatal("expected validation error for missing table, got nil")
	}
}

func TestRunDbEdit_CsvAndKeyVal(t *testing.T) {
	dbPath := setupTestDb(t)

	// Update via CSV
	err := RunDbEdit([]string{dbPath, "-table", "projects", "-id", "1", "csv", "gitmap-pro", "IN_PROGRESS"})
	if err != nil {
		t.Fatalf("RunDbEdit with CSV failed: %v", err)
	}

	// Update via key-value
	err = RunDbEdit([]string{dbPath, "-table", "projects", "-id", "2", "status=DONE"})
	if err != nil {
		t.Fatalf("RunDbEdit with key-value failed: %v", err)
	}

	// Updating non-existent ID should return not found
	err = RunDbEdit([]string{dbPath, "-table", "projects", "-id", "9999", "status=DONE"})
	if err == nil {
		t.Fatal("expected not found error for non-existent ID, got nil")
	}
}

func TestRunDbFind_WildcardAndRegex(t *testing.T) {
	dbPath := setupTestDb(t)

	// Wildcard prefix search
	err := RunDbFind([]string{dbPath, "-table", "projects", "--find", "git*"})
	if err != nil {
		t.Fatalf("RunDbFind prefix wildcard failed: %v", err)
	}

	// Wildcard suffix search
	err = RunDbFind([]string{dbPath, "-table", "projects", "-f", "*map"})
	if err != nil {
		t.Fatalf("RunDbFind suffix wildcard failed: %v", err)
	}

	// Regex search
	err = RunDbFind([]string{dbPath, "-table", "tasks", "--grep", "^Initial", "--columns", "title"})
	if err != nil {
		t.Fatalf("RunDbFind regex search failed: %v", err)
	}

	// JSON output mode
	err = RunDbFind([]string{dbPath, "-table", "projects", "-f", "ACTIVE", "--json"})
	if err != nil {
		t.Fatalf("RunDbFind JSON mode failed: %v", err)
	}
}

func TestRunDbView_AutoJoinAndView(t *testing.T) {
	dbPath := setupTestDb(t)

	// Automatic join and filter with normalization
	err := RunDbView([]string{
		dbPath,
		"-t", "projects as p, tasks as t",
		"--find", "p.name = gitmap and t.title like 'Init*'",
	})
	if err != nil {
		t.Fatalf("RunDbView join and find failed: %v", err)
	}

	// View creation and JSON output
	err = RunDbView([]string{
		dbPath,
		"-t", "projects as p, tasks as t",
		"--viewname", "v_proj_tasks",
		"--json",
	})
	if err != nil {
		t.Fatalf("RunDbView view creation failed: %v", err)
	}
}
