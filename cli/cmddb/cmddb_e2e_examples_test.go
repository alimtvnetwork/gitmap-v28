package cmddb

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func setupE2EExampleDb(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "e2e_example.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("setupE2EExampleDb open failed: %v", err)
	}
	defer db.Close()

	schema := `
	CREATE TABLE projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		status TEXT NOT NULL
	);
	CREATE TABLE tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		priority TEXT NOT NULL,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);
	INSERT INTO projects (name, status) VALUES ('gitmap', 'ACTIVE');
	INSERT INTO projects (name, status) VALUES ('antigravity', 'READY');
	INSERT INTO tasks (project_id, title, priority) VALUES (1, 'Initial task', 'HIGH');
	INSERT INTO tasks (project_id, title, priority) VALUES (1, 'Testing runner', 'NORMAL');
	INSERT INTO tasks (project_id, title, priority) VALUES (2, 'Deploy agent', 'LOW');
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("setupE2EExampleDb seed failed: %v", err)
	}

	return dbPath
}

func TestE2E_UserExamples_Suite(t *testing.T) {
	dbPath := setupE2EExampleDb(t)
	alias := "e2e_proj"

	// 1. gitmap db-read "path" "sql"
	t.Run("db-read path sql", func(t *testing.T) {
		err := RunDbRead([]string{dbPath, "SELECT * FROM projects"})
		if err != nil {
			t.Fatalf("db-read path sql failed: %v", err)
		}
	})

	// 2. gitmap db-read "path" ls
	t.Run("db-read path ls", func(t *testing.T) {
		err := RunDbRead([]string{dbPath, "ls"})
		if err != nil {
			t.Fatalf("db-read path ls failed: %v", err)
		}
	})

	// 3. gitmap db-read "path" -as "aliasName"
	t.Run("db-read path -as alias", func(t *testing.T) {
		err := RunDbRead([]string{dbPath, "-as", alias})
		if err != nil {
			t.Fatalf("db-read path -as alias failed: %v", err)
		}
	})

	// 4. gitmap db-read aliasName "sql ..." --limit 20
	t.Run("db-read alias sql limit", func(t *testing.T) {
		err := RunDbRead([]string{alias, "SELECT * FROM tasks", "--limit", "20"})
		if err != nil {
			t.Fatalf("db-read alias sql limit failed: %v", err)
		}
	})

	// 5. gitmap db-read aliasName table "table1" "table2" --columns/c --json
	t.Run("db-read alias table columns json", func(t *testing.T) {
		err := RunDbRead([]string{alias, "table", "projects", "tasks", "--columns", "--json"})
		if err != nil {
			t.Fatalf("db-read alias table columns json failed: %v", err)
		}
	})

	// 6. gitmap db-read aliasName table "table1" "table2" --json
	t.Run("db-read alias table overview json", func(t *testing.T) {
		err := RunDbRead([]string{alias, "table", "projects", "tasks", "--json"})
		if err != nil {
			t.Fatalf("db-read alias table overview json failed: %v", err)
		}
	})

	// 7. gitmap db-add aliasName "table" csv "col1", "col2"
	t.Run("db-add alias table csv", func(t *testing.T) {
		err := RunDbAdd([]string{alias, "projects", "csv", "superagent", "PENDING"})
		if err != nil {
			t.Fatalf("db-add alias table csv failed: %v", err)
		}
	})

	// 8. gitmap db-edit aliasName -table"table" -id "id-first" csv "col1", "col2"
	t.Run("db-edit alias unspaced table id", func(t *testing.T) {
		err := RunDbEdit([]string{alias, `-table"projects"`, `-id"1"`, "csv", "gitmap-v2", "ACTIVE"})
		if err != nil {
			t.Fatalf("db-edit alias unspaced table id failed: %v", err)
		}
	})

	// 9. gitmap db-find aliasName -table"table, table" --find "text*"
	t.Run("db-find alias unspaced table wildcard", func(t *testing.T) {
		err := RunDbFind([]string{alias, `-table"projects, tasks"`, "--find", "gitmap*"})
		if err != nil {
			t.Fatalf("db-find alias unspaced table wildcard failed: %v", err)
		}
	})

	// 10. gitmap db-find aliasName -table "table, table" --grep "regex"
	t.Run("db-find alias regex", func(t *testing.T) {
		err := RunDbFind([]string{alias, "-table", "projects, tasks", "--grep", "^Deploy"})
		if err != nil {
			t.Fatalf("db-find alias regex failed: %v", err)
		}
	})

	// 11. gitmap db-find aliasName "table, table" --columns "col1", "col2" --find "text" --json --limit 100
	t.Run("db-find alias positional tables columns json", func(t *testing.T) {
		err := RunDbFind([]string{alias, "projects, tasks", "--columns", "name, title", "--find", "agent", "--json", "--limit", "100"})
		if err != nil {
			t.Fatalf("db-find alias positional tables columns json failed: %v", err)
		}
	})

	// 12. gitmap db-listing
	t.Run("db-listing", func(t *testing.T) {
		err := RunDbListing([]string{"--json"})
		if err != nil {
			t.Fatalf("db-listing failed: %v", err)
		}
	})

	// 13. gitmap db-view aliasName -t "table1 as t1", "table2 as t2" --viewname "view-name" --find "condition"
	t.Run("db-view multi-arg tables viewname and find", func(t *testing.T) {
		err := RunDbView([]string{
			alias,
			"-t", "projects as p,", "tasks as t",
			"--viewname", "v_e2e_test",
			"--find", "p.status = ACTIVE and t.priority = HIGH",
			"--json",
		})
		if err != nil {
			t.Fatalf("db-view multi-arg tables viewname and find failed: %v", err)
		}
	})

	// 14. Query existing view with db-view
	t.Run("db-view existing view query", func(t *testing.T) {
		err := RunDbView([]string{
			alias,
			"--viewname", "v_e2e_test",
			"--find", "status = ACTIVE",
			"--json",
		})
		if err != nil {
			t.Fatalf("db-view existing view query failed: %v", err)
		}
	})

	// 15. gitmap db-view aliasName tables
	t.Run("db-view tables listing", func(t *testing.T) {
		err := RunDbView([]string{alias, "tables"})
		if err != nil {
			t.Fatalf("db-view tables listing failed: %v", err)
		}
	})
}
