// Package cmd — schedule_export_import_test.go: unit tests for multi-format schedule
// export, import, auto-inference, and except exclusion filters.
package cmd

import (
	"path/filepath"
	"testing"
)

func TestScheduleExportAndImportFormats(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Seed two schedules
	_ = runScheduleAdd([]string{"exp-task-1", "echo exp1", "--every=1h"})
	_ = runScheduleAdd([]string{"exp-task-2", "echo exp2", "--every=2h"})

	// Run task 1 to produce log records
	_ = runScheduleRun([]string{"exp-task-1"})

	// 2. Export single task to JSON
	jsonPath := filepath.Join(tempDir, "task1.json")
	if err := runScheduleExport([]string{"exp-task-1", "-f", jsonPath}); err != nil {
		t.Fatalf("export json failed: %v", err)
	}

	// 3. Export single task to YAML
	yamlPath := filepath.Join(tempDir, "task1.yaml")
	if err := runScheduleExport([]string{"exp-task-1", "-f", yamlPath}); err != nil {
		t.Fatalf("export yaml failed: %v", err)
	}

	// 4. Export all to SQLite
	sqlitePath := filepath.Join(tempDir, "all.db")
	if err := runScheduleExport([]string{"export-all", "-f", sqlitePath}); err != nil {
		t.Fatalf("export sqlite failed: %v", err)
	}

	// 5. Export all to ZIP
	zipPath := filepath.Join(tempDir, "all.zip")
	if err := runScheduleExport([]string{"export-all", "-f", zipPath}); err != nil {
		t.Fatalf("export zip failed: %v", err)
	}

	// 6. Export all with except filter
	exceptJSON := filepath.Join(tempDir, "except.json")
	if err := runScheduleExport([]string{"export-all", "-except=exp-task-1", "-f", exceptJSON}); err != nil {
		t.Fatalf("export except failed: %v", err)
	}

	bundles, err := parseImportJSON(exceptJSON)
	if err != nil || len(bundles) != 1 || bundles[0].Task.Name != "exp-task-2" {
		t.Fatalf("expected only exp-task-2 in except export, got: %+v (err: %v)", bundles, err)
	}

	// 7. Delete exp-task-1 and re-import from JSON
	_ = runScheduleDelete([]string{"exp-task-1"})
	if err := runScheduleImport([]string{"-f", jsonPath}); err != nil {
		t.Fatalf("import json failed: %v", err)
	}

	db, _ := openSchedulerDB()
	taskAfterImport, err := db.GetSchedule("exp-task-1")
	db.Close()
	if err != nil || taskAfterImport.Name != "exp-task-1" {
		t.Fatalf("expected exp-task-1 to be restored: %+v (err: %v)", taskAfterImport, err)
	}

	// 8. Test status command for single task
	if err := runScheduleStatus([]string{"exp-task-1"}); err != nil {
		t.Errorf("status command failed: %v", err)
	}

	// Clean up
	_ = runScheduleDelete([]string{"exp-task-1"})
	_ = runScheduleDelete([]string{"exp-task-2"})
}
