package cmddb

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestListTablesAndViews_DeadlockRegression_MaxOpenConns1(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "deadlock_test.db")

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	defer conn.Close()

	// Strictly enforce MaxOpenConns = 1 to trigger pool starvation if rows are kept open
	conn.SetMaxOpenConns(1)

	setupSql := `
	CREATE TABLE Repositories (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		path TEXT NOT NULL
	);
	INSERT INTO Repositories (name, path) VALUES ('repo1', '/path/1'), ('repo2', '/path/2');

	CREATE TABLE Commits (
		id INTEGER PRIMARY KEY,
		repo_id INTEGER,
		hash TEXT NOT NULL
	);
	INSERT INTO Commits (repo_id, hash) VALUES (1, 'abc'), (2, 'def');

	CREATE VIEW ActiveRepos AS SELECT id, name FROM Repositories;
	`
	if _, err := conn.Exec(setupSql); err != nil {
		t.Fatalf("failed to setup test tables: %v", err)
	}

	done := make(chan struct{})
	var tables []TableInfo
	var listErr error

	go func() {
		tables, listErr = listTablesAndViews(conn)
		close(done)
	}()

	select {
	case <-done:
		if listErr != nil {
			t.Fatalf("listTablesAndViews failed: %v", listErr)
		}
		if len(tables) != 3 {
			t.Fatalf("expected 3 tables/views, got %d", len(tables))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("DEADLOCK DETECTED: listTablesAndViews hung with MaxOpenConns=1")
	}
}

func TestFetchAllTableSchemas_DeadlockRegression_MaxOpenConns1(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "schema_deadlock_test.db")

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	defer conn.Close()

	conn.SetMaxOpenConns(1)

	setupSql := `
	CREATE TABLE Users (
		id INTEGER PRIMARY KEY,
		username TEXT NOT NULL
	);
	CREATE INDEX idx_users_username ON Users (username);

	CREATE TABLE Sessions (
		token TEXT PRIMARY KEY,
		user_id INTEGER
	);
	CREATE INDEX idx_sessions_user ON Sessions (user_id);
	`
	if _, err := conn.Exec(setupSql); err != nil {
		t.Fatalf("failed to setup test tables: %v", err)
	}

	done := make(chan struct{})
	var schemas []SchemaInfo
	var schemaErr error

	go func() {
		schemas, schemaErr = fetchAllTableSchemas(conn)
		close(done)
	}()

	select {
	case <-done:
		if schemaErr != nil {
			t.Fatalf("fetchAllTableSchemas failed: %v", schemaErr)
		}
		if len(schemas) != 2 {
			t.Fatalf("expected 2 schemas, got %d", len(schemas))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("DEADLOCK DETECTED: fetchAllTableSchemas hung with MaxOpenConns=1")
	}
}
