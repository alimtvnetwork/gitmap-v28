package cmddb

import (
	"database/sql"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func createTestDb(t *testing.T, filename string) string {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, filename)

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to create sqlite db: %v", err)
	}
	defer conn.Close()

	schema := `
	CREATE TABLE Members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT,
		is_active INTEGER NOT NULL DEFAULT 1
	);
	INSERT INTO Members (username, email, is_active) VALUES
		('alice', 'alice@example.com', 1),
		('bob', 'bob@example.com', 0),
		('charlie', 'charlie@example.com', 1);

	CREATE TABLE Projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL
	);
	INSERT INTO Projects (title) VALUES ('GitMap'), ('Antigravity');
	`
	if _, err := conn.Exec(schema); err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}

	return dbPath
}

func captureStdout(fn func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = oldStdout

	outBytes, _ := io.ReadAll(r)
	return string(outBytes)
}

func TestDbRead_Help(t *testing.T) {
	out := captureStdout(func() {
		err := RunDbRead([]string{"--help"})
		if err != nil {
			t.Fatalf("expected nil error for help, got %v", err)
		}
	})

	if !strings.Contains(out, "gitmap db-read") {
		t.Fatalf("expected help text, got: %s", out)
	}
}

func TestDbRead_Sql_Readonly(t *testing.T) {
	dbPath := createTestDb(t, "readonly_test.db")

	// 1. SELECT query succeeds
	out := captureStdout(func() {
		appErr := RunDbRead([]string{dbPath, "SELECT username, email FROM Members ORDER BY id ASC", "-l", "2"})
		if appErr != nil {
			t.Fatalf("expected nil error for SELECT query, got %v", appErr)
		}
	})

	if !strings.Contains(out, "alice") || !strings.Contains(out, "Completed in") {
		t.Fatalf("expected query output with timing, got: %s", out)
	}

	// 2. Write query is rejected by query safety
	writeErr := RunDbRead([]string{dbPath, "INSERT INTO Members (username) VALUES ('hacker')"})
	if writeErr == nil {
		t.Fatal("expected error for write query, got nil")
	}
	if writeErr.Code != "ERR_READONLY_VIOLATION" {
		t.Fatalf("expected ERR_READONLY_VIOLATION, got %s", writeErr.Code)
	}
}

func TestDbRead_Sql_Json(t *testing.T) {
	dbPath := createTestDb(t, "json_test.db")

	out := captureStdout(func() {
		appErr := RunDbRead([]string{dbPath, "SELECT username FROM Members", "--json"})
		if appErr != nil {
			t.Fatalf("expected nil error, got %v", appErr)
		}
	})

	var env map[string]any
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("failed to parse JSON envelope: %v\nOutput: %s", err, out)
	}

	if env["success"] != true {
		t.Fatalf("expected success true in json, got %v", env["success"])
	}
	data, ok := env["data"].(map[string]any)
	if !ok || data["count"] == nil {
		t.Fatalf("expected data object with count in json, got %v", env["data"])
	}
}

func TestDbRead_Ls(t *testing.T) {
	dbPath := createTestDb(t, "ls_test.db")

	// 1. Terminal tabular output
	out := captureStdout(func() {
		appErr := RunDbRead([]string{dbPath, "ls"})
		if appErr != nil {
			t.Fatalf("expected nil error, got %v", appErr)
		}
	})

	if !strings.Contains(out, "Members") || !strings.Contains(out, "Projects") {
		t.Fatalf("expected tables listed in output, got: %s", out)
	}
	if !strings.Contains(out, "Completed in") {
		t.Fatalf("expected completion time in output, got: %s", out)
	}

	// 2. JSON output
	jsonOut := captureStdout(func() {
		appErr := RunDbRead([]string{dbPath, "ls", "--json"})
		if appErr != nil {
			t.Fatalf("expected nil error, got %v", appErr)
		}
	})

	var env map[string]any
	if err := json.Unmarshal([]byte(jsonOut), &env); err != nil {
		t.Fatalf("failed to parse JSON envelope: %v\nOutput: %s", err, jsonOut)
	}
	if env["command"] != "db-read:ls" {
		t.Fatalf("expected command db-read:ls, got %v", env["command"])
	}
}

func TestDbRead_AliasRegistrationAndResolution(t *testing.T) {
	dbPath := createTestDb(t, "alias_target.db")
	aliasName := "test_read_alias"

	// 1. Register alias
	out := captureStdout(func() {
		appErr := RunDbRead([]string{dbPath, "-as", aliasName})
		if appErr != nil {
			t.Fatalf("expected nil error for alias registration, got %v", appErr)
		}
	})

	if !strings.Contains(out, "Registered alias 'test_read_alias'") {
		t.Fatalf("expected confirmation message, got: %s", out)
	}

	// 2. Query using the registered alias!
	queryOut := captureStdout(func() {
		appErr := RunDbRead([]string{aliasName, "SELECT count(*) AS total FROM Members"})
		if appErr != nil {
			t.Fatalf("expected nil error querying via alias, got %v", appErr)
		}
	})

	if !strings.Contains(queryOut, "Completed in") {
		t.Fatalf("expected query completion time, got: %s", queryOut)
	}
}

func TestDbRead_Table(t *testing.T) {
	dbPath := createTestDb(t, "table_test.db")

	// 1. Table overview without --columns
	out := captureStdout(func() {
		appErr := RunDbRead([]string{dbPath, "table", "Members"})
		if appErr != nil {
			t.Fatalf("expected nil error for table overview, got %v", appErr)
		}
	})

	if !strings.Contains(out, "Table: Members") || !strings.Contains(out, "CREATE TABLE Members") {
		t.Fatalf("expected schema and overview for Members, got: %s", out)
	}
	if !strings.Contains(out, "Completed in") {
		t.Fatalf("expected completion time, got: %s", out)
	}

	// 2. Table columns with --columns / -c
	colsOut := captureStdout(func() {
		appErr := RunDbRead([]string{dbPath, "table", "Members", "-c"})
		if appErr != nil {
			t.Fatalf("expected nil error for table columns, got %v", appErr)
		}
	})

	if !strings.Contains(colsOut, "username") || !strings.Contains(colsOut, "email") {
		t.Fatalf("expected columns listed, got: %s", colsOut)
	}

	// 3. Table columns with --json
	jsonColsOut := captureStdout(func() {
		appErr := RunDbRead([]string{dbPath, "table", "Members", "-c", "--json"})
		if appErr != nil {
			t.Fatalf("expected nil error for table columns json, got %v", appErr)
		}
	})

	var env map[string]any
	if err := json.Unmarshal([]byte(jsonColsOut), &env); err != nil {
		t.Fatalf("failed to parse JSON envelope: %v\nOutput: %s", err, jsonColsOut)
	}
	if env["command"] != "db-read:table" {
		t.Fatalf("expected command db-read:table, got %v", env["command"])
	}
}
