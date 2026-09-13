package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

func TestStartupSplitDB_CRUD(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "startup_test.db")

	db, err := initStartupSplitConnDirect(dbPath)
	if err != nil {
		t.Fatalf("initStartupSplitConnDirect failed: %v", err)
	}
	defer db.Close()

	item := &StartupItemRecord{
		Name:         "weekly-backup",
		TargetType:   "script",
		TargetPath:   "C:\\scripts\\backup.ps1",
		CommandArgs:  "-full",
		RunFrequency: "once-a-week",
		IsActive:     true,
		Description:  "Runs weekly backup script",
	}

	if err := db.SaveStartupItem(item); err != nil {
		t.Fatalf("SaveStartupItem failed: %v", err)
	}

	fetched, err := db.GetStartupItem("weekly-backup")
	if err != nil {
		t.Fatalf("GetStartupItem failed: %v", err)
	}
	if fetched == nil || fetched.Name != "weekly-backup" {
		t.Fatalf("expected weekly-backup, got %v", fetched)
	}
	if fetched.RunFrequency != "once-a-week" || !fetched.IsActive {
		t.Fatalf("unexpected fields: %+v", fetched)
	}

	items, err := db.ListStartupItems()
	if err != nil || len(items) != 1 {
		t.Fatalf("ListStartupItems failed: err=%v, count=%d", err, len(items))
	}

	logRec := &StartupLogRecord{
		StartupItemId: fetched.StartupItemId,
		DurationMs:    1200,
		IsSuccess:     true,
		ExitCode:      0,
		OutputSummary: "backup completed OK",
		Notes:         "automated run",
	}
	if err := db.RecordStartupLog(logRec); err != nil {
		t.Fatalf("RecordStartupLog failed: %v", err)
	}

	logs, err := db.ListStartupLogs(fetched.StartupItemId, 10)
	if err != nil || len(logs) != 1 {
		t.Fatalf("ListStartupLogs failed: err=%v, count=%d", err, len(logs))
	}
	if logs[0].OutputSummary != "backup completed OK" {
		t.Errorf("unexpected log output: %q", logs[0].OutputSummary)
	}

	if err := db.DeleteStartupItem("weekly-backup"); err != nil {
		t.Fatalf("DeleteStartupItem failed: %v", err)
	}

	itemsAfter, _ := db.ListStartupItems()
	if len(itemsAfter) != 0 {
		t.Fatalf("expected 0 items after delete, got %d", len(itemsAfter))
	}
}

func initStartupSplitConnDirect(dbPath string) (*StartupSplitDB, error) {
	_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	if err := createStartupSchema(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &StartupSplitDB{conn: conn, Path: dbPath}, nil
}
