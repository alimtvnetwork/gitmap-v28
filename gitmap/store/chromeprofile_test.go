// Package store — chromeprofile_test.go: verifies the lazy-migrated
// ChromeProfile / ChromeProfileExport tables. Covers upsert idempotency,
// export insert + read-back, and cascading delete with artifact-path
// extraction (the CLI relies on the returned slice to rm files on disk).
package store

import (
	"path/filepath"
	"testing"
)

func freshChromeProfileDB(t *testing.T) *DB {
	t.Helper()
	db, err := OpenAt(filepath.Join(t.TempDir(), "cp.sqlite"))
	if err != nil {
		t.Fatalf("OpenAt: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.EnsureChromeProfileTables(); err != nil {
		t.Fatalf("EnsureChromeProfileTables: %v", err)
	}
	return db
}

func TestChromeProfileUpsertIdempotent(t *testing.T) {
	db := freshChromeProfileDB(t)
	id1, err := db.UpsertChromeProfile("Default", "/tmp/src", true)
	if err != nil {
		t.Fatalf("upsert #1: %v", err)
	}
	id2, err := db.UpsertChromeProfile("Default", "/tmp/src2", true)
	if err != nil {
		t.Fatalf("upsert #2: %v", err)
	}
	if id1 != id2 {
		t.Fatalf("expected same id on re-upsert, got %d vs %d", id1, id2)
	}
}

func TestChromeProfileExportInsertAndList(t *testing.T) {
	db := freshChromeProfileDB(t)
	id, err := db.UpsertChromeProfile("Work", "/tmp/work", true)
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if err := db.InsertChromeProfileExport(id, "json", "/snap/work.json", 1234); err != nil {
		t.Fatalf("insert json: %v", err)
	}
	if err := db.InsertChromeProfileExport(id, "csv", "/snap/work.csv", 234); err != nil {
		t.Fatalf("insert csv: %v", err)
	}
	rows, err := db.ListChromeProfilesDB()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("want 1 row, got %d", len(rows))
	}
	if rows[0].Name != "Work" || rows[0].ExportCount != 2 {
		t.Fatalf("unexpected row: %+v", rows[0])
	}
}

func TestChromeProfileDeleteReturnsArtifacts(t *testing.T) {
	db := freshChromeProfileDB(t)
	id, _ := db.UpsertChromeProfile("Gone", "/tmp/x", true)
	_ = db.InsertChromeProfileExport(id, "json", "/snap/gone.json", 1)
	_ = db.InsertChromeProfileExport(id, "csv", "/snap/gone.csv", 1)

	paths, err := db.DeleteChromeProfile("Gone")
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("want 2 artifact paths, got %d (%v)", len(paths), paths)
	}
	if db.ChromeProfileExists("Gone") {
		t.Fatalf("row should be gone")
	}
}
