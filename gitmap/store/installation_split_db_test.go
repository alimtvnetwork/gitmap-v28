package store

import (
	"path/filepath"
	"testing"
)

func TestInstallationSplitDB_Lifecycle(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "installation.db")

	splitDB, err := OpenInstallationSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("Failed to open split DB: %v", err)
	}
	defer splitDB.Close()

	// Test SaveInstalledTool
	if err := splitDB.SaveInstalledTool("git", "2.44.0", "apt"); err != nil {
		t.Fatalf("SaveInstalledTool failed: %v", err)
	}

	// Test GetInstalledTool
	tool, err := splitDB.GetInstalledTool("git")
	if err != nil {
		t.Fatalf("GetInstalledTool failed: %v", err)
	}
	if tool.Tool != "git" || tool.VersionString != "2.44.0" {
		t.Errorf("Unexpected tool record: %+v", tool)
	}

	// Test IsToolInstalled
	if !splitDB.IsToolInstalled("git") {
		t.Errorf("Expected git to be reported as installed")
	}

	// Test RecordLog & GetLogs
	logRec := InstallationLogRecord{
		Tool:           "git",
		Action:         "install",
		Version:        "2.44.0",
		PackageManager: "apt",
		DurationMs:     1250,
		IsSuccess:      true,
		ExitCode:       0,
	}
	if err := splitDB.RecordLog(logRec); err != nil {
		t.Fatalf("RecordLog failed: %v", err)
	}

	logs, err := splitDB.GetLogs(10)
	if err != nil || len(logs) == 0 {
		t.Fatalf("GetLogs failed or returned empty: %v", err)
	}
	if logs[0].Tool != "git" || !logs[0].IsSuccess {
		t.Errorf("Unexpected log record: %+v", logs[0])
	}

	// Test RemoveInstalledTool
	if err := splitDB.RemoveInstalledTool("git"); err != nil {
		t.Fatalf("RemoveInstalledTool failed: %v", err)
	}
	if splitDB.IsToolInstalled("git") {
		t.Errorf("Expected git to be removed")
	}
}

func TestMigrateInstalledToolsFromRoot(t *testing.T) {
	tempDir := t.TempDir()
	rootPath := filepath.Join(tempDir, "root.db")
	splitPath := filepath.Join(tempDir, "installation.db")

	rootDB, err := OpenAt(rootPath)
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}
	defer rootDB.Close()

	if err := rootDB.Migrate(); err != nil {
		t.Fatalf("rootDB.Migrate failed: %v", err)
	}

	// Insert into root
	_, err = rootDB.conn.Exec(`INSERT INTO InstalledTool (Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, UpdatedAt)
VALUES ('curl', 8, 5, 0, 0, '8.5.0', 'apt', '', '2026-09-07T00:00:00Z');`)
	if err != nil {
		t.Fatalf("Insert to root failed: %v", err)
	}

	splitDB, err := OpenInstallationSplitDBAt(splitPath)
	if err != nil {
		t.Fatalf("OpenInstallationSplitDBAt failed: %v", err)
	}
	defer splitDB.Close()

	if err := MigrateInstalledToolsFromRoot(rootDB.conn, splitDB); err != nil {
		t.Fatalf("MigrateInstalledToolsFromRoot failed: %v", err)
	}

	if !splitDB.IsToolInstalled("curl") {
		t.Errorf("Expected curl to be migrated to split DB")
	}

	tool, err := splitDB.GetInstalledTool("curl")
	if err != nil {
		t.Fatalf("GetInstalledTool failed: %v", err)
	}
	if tool.VersionString != "8.5.0" {
		t.Errorf("Expected version 8.5.0, got %s", tool.VersionString)
	}
}
