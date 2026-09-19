package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStorageCmdDriveReport(t *testing.T) {
	err := RunStorageCmd([]string{})
	if err != nil {
		t.Fatalf("RunStorageCmd without args returned error: %v", err)
	}
}

func TestStorageCmdList(t *testing.T) {
	err := RunStorageCmd([]string{"ls"})
	if err != nil {
		t.Fatalf("RunStorageCmd ls returned error: %v", err)
	}
}

func TestGetDiskSpaceMetrics(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd failed: %v", err)
	}

	metrics, err := getDiskSpaceMetrics(cwd)
	if err != nil {
		t.Fatalf("getDiskSpaceMetrics failed: %v", err)
	}

	if metrics.TotalBytes == 0 {
		t.Errorf("expected total bytes > 0, got 0")
	}

	if metrics.DrivePath == "" {
		t.Errorf("expected non-empty drive path")
	}
}

func TestIsStorageListSubcommand(t *testing.T) {
	if !isStorageListSubcommand("ls") {
		t.Errorf("expected ls to be storage list subcommand")
	}

	if !isStorageListSubcommand("list") {
		t.Errorf("expected list to be storage list subcommand")
	}

	if !isStorageListSubcommand("db") {
		t.Errorf("expected db to be storage list subcommand")
	}

	if isStorageListSubcommand("unknown") {
		t.Errorf("expected unknown NOT to be storage list subcommand")
	}
}

func TestResolveRelativeDBPath(t *testing.T) {
	dataDir := filepath.Join(string(filepath.Separator), "base", "data")
	fullPath := filepath.Join(dataDir, "installation.db")
	rel := resolveRelativeDBPath(dataDir, fullPath)
	if rel != "installation.db" {
		t.Errorf("expected installation.db, got %s", rel)
	}

	outsidePath := filepath.Join(string(filepath.Separator), "outside", "store.db")
	relOutside := resolveRelativeDBPath(dataDir, outsidePath)
	if relOutside != outsidePath {
		t.Errorf("expected %s, got %s", outsidePath, relOutside)
	}
}

func TestIsStorageRestoreSubcommand(t *testing.T) {
	if !isStorageRestoreSubcommand("restore-db") {
		t.Errorf("expected restore-db to be storage restore subcommand")
	}
	if !isStorageRestoreSubcommand("restoredb") {
		t.Errorf("expected restoredb to be storage restore subcommand")
	}
	if !isStorageRestoreSubcommand("restore") {
		t.Errorf("expected restore to be storage restore subcommand")
	}
	if isStorageRestoreSubcommand("other") {
		t.Errorf("expected other NOT to be storage restore subcommand")
	}
}

func TestIsStorageSpaceSubcommand(t *testing.T) {
	if !isStorageSpaceSubcommand("space") {
		t.Errorf("expected space to be storage space subcommand")
	}
	if isStorageSpaceSubcommand("dbs") {
		t.Errorf("expected dbs NOT to be storage space subcommand")
	}
}

func TestStorageCmdSpaceList(t *testing.T) {
	err := RunStorageCmd([]string{"space", "ls"})
	if err != nil {
		t.Fatalf("RunStorageCmd space ls returned error: %v", err)
	}
}

func TestStorageCmdRestoreDB(t *testing.T) {
	err := RunStorageCmd([]string{"restore-db"})
	if err != nil {
		t.Fatalf("RunStorageCmd restore-db returned error: %v", err)
	}
}

func TestStorageCmdRestoreDbAlias(t *testing.T) {
	err := RunStorageCmd([]string{"restoredb"})
	if err != nil {
		t.Fatalf("RunStorageCmd restoredb returned error: %v", err)
	}
}

func TestCollectStorageEntries(t *testing.T) {
	entries := collectStorageEntries()
	if len(entries) == 0 {
		t.Errorf("expected at least 1 database entry, got 0")
	}
}

