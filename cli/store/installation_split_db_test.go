package store

import (
	"path/filepath"
	"strings"
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

	testSaveAndGetTool(t, splitDB)
	testRecordLogAndGet(t, splitDB)
	testRemoveTool(t, splitDB)
}

func testSaveAndGetTool(t *testing.T, splitDB *InstallationSplitDB) {
	if err := splitDB.SaveInstalledTool("git", "2.44.0", "apt"); err != nil {
		t.Fatalf("SaveInstalledTool failed: %v", err)
	}

	tool, err := splitDB.GetInstalledTool("git")
	if err != nil {
		t.Fatalf("GetInstalledTool failed: %v", err)
	}

	if tool.Tool != "git" || tool.VersionString != "2.44.0" {
		t.Errorf("Unexpected tool record: %+v", tool)
	}

	isInstalled := splitDB.IsToolInstalled("git")
	if !isInstalled {
		t.Errorf("Expected git to be reported as installed")
	}
}

func testRecordLogAndGet(t *testing.T, splitDB *InstallationSplitDB) {
	logRec := InstallationLogRecord{
		Tool: "git", Action: "install", Version: "2.44.0",
		PackageManager: "apt", DurationMs: 1250, IsSuccess: true,
	}

	if err := splitDB.RecordLog(logRec); err != nil {
		t.Fatalf("RecordLog failed: %v", err)
	}

	logs, err := splitDB.GetLogs(10)
	if err != nil || len(logs) == 0 {
		t.Fatalf("GetLogs failed: %v", err)
	}

	if logs[0].Tool != "git" || !logs[0].IsSuccess {
		t.Errorf("Unexpected log record: %+v", logs[0])
	}
}

func testRemoveTool(t *testing.T, splitDB *InstallationSplitDB) {
	if err := splitDB.RemoveInstalledTool("git"); err != nil {
		t.Fatalf("RemoveInstalledTool failed: %v", err)
	}

	isInstalled := splitDB.IsToolInstalled("git")
	if isInstalled {
		t.Errorf("Expected git to be removed")
	}
}

func TestSanitizeLogOutput(t *testing.T) {
	shortStr := "normal command output"
	if SanitizeLogOutput(shortStr) != shortStr {
		t.Errorf("SanitizeLogOutput changed short string")
	}

	largeStr := strings.Repeat("A", 70000)
	sanitized := SanitizeLogOutput(largeStr)

	if len(sanitized) != maxLogOutputBytes {
		t.Errorf("Expected capped length %d, got %d", maxLogOutputBytes, len(sanitized))
	}
}

func TestRecordExecutionAndGetFailedLogs(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "installation.db")

	splitDB, err := OpenInstallationSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenInstallationSplitDBAt failed: %v", err)
	}

	defer splitDB.Close()

	testRecordExecutionFlow(t, splitDB)
}

func testRecordExecutionFlow(t *testing.T, splitDB *InstallationSplitDB) {
	err1 := splitDB.RecordExecution("rust", "install", "1.77.0", "cargo", 500, true, 0, "ok", "", "cargo build", "notes", "none")
	if err1 != nil {
		t.Fatalf("RecordExecution success failed: %v", err1)
	}

	err2 := splitDB.RecordExecution("ollama", "install", "latest", "curl", 1200, false, 1, "", "download timeout", "curl -fsSL https://ollama.com/install.sh", "err", "none")
	if err2 != nil {
		t.Fatalf("RecordExecution failure failed: %v", err2)
	}

	failedLogs, err := splitDB.GetFailedLogs(10)
	if err != nil || len(failedLogs) != 1 {
		t.Fatalf("GetFailedLogs unexpected count: len=%d err=%v", len(failedLogs), err)
	}

	if failedLogs[0].Tool != "ollama" || failedLogs[0].ExitCode != 1 {
		t.Errorf("Unexpected failed log: %+v", failedLogs[0])
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

	seedRootDB(t, rootDB)
	verifyMigratedRootTools(t, rootDB, splitPath)
}

func seedRootDB(t *testing.T, rootDB *DB) {
	_, err := rootDB.conn.Exec(`INSERT INTO InstalledTool (Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, UpdatedAt)
VALUES ('curl', 8, 5, 0, 0, '8.5.0', 'apt', '', '2026-09-07T00:00:00Z');`)
	if err != nil {
		t.Fatalf("Insert to root failed: %v", err)
	}
}

func verifyMigratedRootTools(t *testing.T, rootDB *DB, splitPath string) {
	splitDB, err := OpenInstallationSplitDBAt(splitPath)
	if err != nil {
		t.Fatalf("OpenInstallationSplitDBAt failed: %v", err)
	}

	defer splitDB.Close()

	if err := MigrateInstalledToolsFromRoot(rootDB.conn, splitDB); err != nil {
		t.Fatalf("MigrateInstalledToolsFromRoot failed: %v", err)
	}

	isCurlInstalled := splitDB.IsToolInstalled("curl")
	if !isCurlInstalled {
		t.Errorf("Expected curl to be migrated to split DB")
	}

	tool, err := splitDB.GetInstalledTool("curl")
	if err != nil || tool.VersionString != "8.5.0" {
		t.Errorf("Unexpected tool record: %+v", tool)
	}
}
