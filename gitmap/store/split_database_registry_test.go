package store

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestSplitDatabaseRegistry_RegisterAndGet(t *testing.T) {
	tempDir := t.TempDir()
	db, err := OpenAt(filepath.Join(tempDir, "root.db"))
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}

	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	entry := SplitDatabaseEntry{
		DatabaseType:   "custom",
		DatabaseKey:    "test-node",
		DatabasePath:   filepath.Join(tempDir, "custom.db"),
		SizeBytes:      1024,
		TableCount:     3,
		RecordCount:    42,
		SchemaVersion:  2,
		Status:         "active",
		IsActive:       true,
		Description:    "Test split database",
		Notes:          "Operational notes",
		Comments:       "User comments",
		LastAccessedAt: 1700000000,
		LastSyncedAt:   1700000100,
	}

	if err := db.RegisterSplitDB(entry); err != nil {
		t.Fatalf("RegisterSplitDB failed: %v", err)
	}

	got, err := db.GetSplitDB("custom", "test-node")
	if err != nil {
		t.Fatalf("GetSplitDB failed: %v", err)
	}

	if got.ID <= 0 || got.SplitDatabaseRegistryID != got.ID {
		t.Errorf("expected valid ID, got %d / %d", got.ID, got.SplitDatabaseRegistryID)
	}

	if got.RecordCount != 42 || got.TableCount != 3 || got.SizeBytes != 1024 {
		t.Errorf("unexpected counts: rec=%d tbl=%d sz=%d", got.RecordCount, got.TableCount, got.SizeBytes)
	}

	if !got.IsActive || got.IsAttached {
		t.Errorf("unexpected bools: active=%v attached=%v", got.IsActive, got.IsAttached)
	}
}

func TestSplitDatabaseRegistry_ConflictUpdate(t *testing.T) {
	tempDir := t.TempDir()
	db, err := OpenAt(filepath.Join(tempDir, "root.db"))
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}

	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	e1 := SplitDatabaseEntry{
		DatabaseType: "custom",
		DatabaseKey:  "conflict-node",
		DatabasePath: filepath.Join(tempDir, "conflict.db"),
		RecordCount:  10,
		TableCount:   2,
	}

	_ = db.RegisterSplitDB(e1)
	first, _ := db.GetSplitDB("custom", "conflict-node")

	e2 := SplitDatabaseEntry{
		DatabaseType: "custom",
		DatabaseKey:  "conflict-node",
		DatabasePath: filepath.Join(tempDir, "conflict.db"),
		RecordCount:  150,
		TableCount:   5,
		SizeBytes:    8192,
		Status:       "updated",
		IsActive:     true,
	}

	if err := db.RegisterSplitDB(e2); err != nil {
		t.Fatalf("second register failed: %v", err)
	}

	second, err := db.GetSplitDB("custom", "conflict-node")
	if err != nil {
		t.Fatalf("GetSplitDB second failed: %v", err)
	}

	if second.RecordCount != 150 || second.TableCount != 5 || second.SizeBytes != 8192 {
		t.Errorf("conflict update failed: rec=%d tbl=%d sz=%d", second.RecordCount, second.TableCount, second.SizeBytes)
	}

	if second.Status != "updated" || second.CreatedAt != first.CreatedAt {
		t.Errorf("expected status 'updated' and preserved CreatedAt")
	}
}

func TestSplitDatabaseRegistry_ListSplitDBs(t *testing.T) {
	tempDir := t.TempDir()
	db, err := OpenAt(filepath.Join(tempDir, "root.db"))
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}

	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	_ = db.RegisterSplitDB(SplitDatabaseEntry{DatabaseType: "typeA", DatabaseKey: "k1", DatabasePath: "p1"})
	_ = db.RegisterSplitDB(SplitDatabaseEntry{DatabaseType: "typeA", DatabaseKey: "k2", DatabasePath: "p2"})
	_ = db.RegisterSplitDB(SplitDatabaseEntry{DatabaseType: "typeB", DatabaseKey: "k3", DatabasePath: "p3"})

	typeAList, err := db.ListSplitDBs("typeA")
	if err != nil || len(typeAList) != 2 {
		t.Errorf("expected 2 typeA items, got %d, err: %v", len(typeAList), err)
	}

	allList, err := db.ListSplitDBs("")
	if err != nil || len(allList) < 3 {
		t.Errorf("expected >= 3 items in allList, got %d, err: %v", len(allList), err)
	}
}

func TestSplitDatabaseRegistry_SyncKnownSplitDatabases(t *testing.T) {
	tempDir := t.TempDir()
	db, err := OpenAt(filepath.Join(tempDir, "root.db"))
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}

	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	instEntry, err := db.GetSplitDB("installation", "default")
	if err != nil || instEntry.Status != "active" || instEntry.TableCount < 2 {
		t.Errorf("unexpected installation entry: %+v, err: %v", instEntry, err)
	}

	schedDir := filepath.Join(tempDir, "schedules")
	_ = os.MkdirAll(schedDir, 0755)
	schedDBPath := filepath.Join(schedDir, "job-alpha.db")
	conn, err := sql.Open("sqlite", schedDBPath)
	if err != nil {
		t.Fatalf("create sched db failed: %v", err)
	}

	_, _ = conn.Exec("CREATE TABLE test_job (id INTEGER PRIMARY KEY, msg TEXT);")
	_, _ = conn.Exec("INSERT INTO test_job (msg) VALUES ('hello');")
	_ = conn.Close()

	if err := db.SyncKnownSplitDatabases(); err != nil {
		t.Fatalf("SyncKnownSplitDatabases failed: %v", err)
	}

	schedEntry, err := db.GetSplitDB("schedule", "job-alpha")
	if err != nil || schedEntry.TableCount != 1 || schedEntry.RecordCount != 1 {
		t.Errorf("unexpected schedEntry: %+v, err: %v", schedEntry, err)
	}
}

func TestSplitDatabaseRegistry_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	db, err := OpenAt(filepath.Join(tempDir, "root.db"))
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}

	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	_, err = db.GetSplitDB("nonexistent", "missing")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}
