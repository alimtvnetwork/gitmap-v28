package store

import (
	"path/filepath"
	"testing"
)

func isFalse(b bool) bool {
	return b == false
}

func setupTestSplitDB(t *testing.T) *InstallationSplitDB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "installation.db")
	db, err := OpenInstallationSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenInstallationSplitDBAt failed: %v", err)
	}

	return db
}

func TestInstallLogRecord_Helpers(t *testing.T) {
	testRecordSuccessHelpers(t)
	testRecordFailedHelpers(t)
	testRecordSkippedHelpers(t)
	testRecordRunningHelpers(t)
}

func testRecordSuccessHelpers(t *testing.T) {
	recSuccess := InstallLogRecord{Status: InstallStatusSuccessType}
	if isFalse(recSuccess.IsSuccess()) {
		t.Errorf("expected IsSuccess to report true for success")
	}

	recCompleted := InstallLogRecord{Status: "completed"}
	if isFalse(recCompleted.IsSuccess()) {
		t.Errorf("expected IsSuccess to report true for completed")
	}
}

func testRecordFailedHelpers(t *testing.T) {
	recFailed := InstallLogRecord{Status: InstallStatusFailedType}
	if isFalse(recFailed.IsFailed()) {
		t.Errorf("expected IsFailed to report true for failed")
	}

	recFailure := InstallLogRecord{Status: "failure"}
	if isFalse(recFailure.IsFailed()) {
		t.Errorf("expected IsFailed to report true for failure")
	}
}

func testRecordSkippedHelpers(t *testing.T) {
	recSkipped := InstallLogRecord{Status: InstallStatusSkippedType}
	if isFalse(recSkipped.IsSkipped()) {
		t.Errorf("expected IsSkipped to report true for skipped")
	}
}

func testRecordRunningHelpers(t *testing.T) {
	recRunning := InstallLogRecord{Status: InstallStatusRunningType}
	if isFalse(recRunning.IsRunning()) {
		t.Errorf("expected IsRunning to report true for running")
	}

	recStarted := InstallLogRecord{Status: "started"}
	if isFalse(recStarted.IsRunning()) {
		t.Errorf("expected IsRunning to report true for started")
	}
}

func TestInstallationSplitDB_InstallLogsLifecycle(t *testing.T) {
	db := setupTestSplitDB(t)
	defer db.Close()

	testInstallStartAndSuccess(t, db)
	testInstallFailure(t, db)
	testInstallSkipped(t, db)
	testDirectFallbacks(t, db)
	testListInstallLogs(t, db)
	testGetNotFound(t, db)
}

func testInstallStartAndSuccess(t *testing.T, db *InstallationSplitDB) {
	id, err := db.RecordInstallStart("tool", "antigravity", "install")
	if err != nil {
		t.Fatalf("RecordInstallStart failed: %v", err)
	}

	verifyInstallStarted(t, db, id)

	if errSuccess := db.RecordInstallSuccess("tool", "antigravity"); errSuccess != nil {
		t.Fatalf("RecordInstallSuccess failed: %v", errSuccess)
	}

	verifyInstallSucceeded(t, db, id)
}

func verifyInstallStarted(t *testing.T, db *InstallationSplitDB, id string) {
	rec, err := db.GetInstallLog(id)
	if err != nil {
		t.Fatalf("GetInstallLog failed: %v", err)
	}
	if isFalse(rec.IsRunning()) {
		t.Errorf("expected record status running, got %s", rec.Status)
	}
}

func verifyInstallSucceeded(t *testing.T, db *InstallationSplitDB, id string) {
	rec, err := db.GetInstallLog(id)
	if err != nil {
		t.Fatalf("GetInstallLog after success failed: %v", err)
	}
	if isFalse(rec.IsSuccess()) {
		t.Errorf("expected record status success, got %s", rec.Status)
	}
	if rec.EndedAt == "" {
		t.Errorf("expected ended_at to be populated")
	}
}

func testInstallFailure(t *testing.T, db *InstallationSplitDB) {
	id, err := db.RecordInstallStart("tool", "git", "install")
	if err != nil {
		t.Fatalf("RecordInstallStart for git failed: %v", err)
	}

	if errFail := db.RecordInstallFailure("tool", "git", 1, "download error"); errFail != nil {
		t.Fatalf("RecordInstallFailure failed: %v", errFail)
	}

	rec, err := db.GetInstallLog(id)
	if err != nil || rec.Status != "failed" || rec.ExitCode != 1 {
		t.Errorf("expected failed status and exitCode 1, got %+v", rec)
	}
}

func testInstallSkipped(t *testing.T, db *InstallationSplitDB) {
	id, err := db.RecordInstallStart("package", "docker", "install")
	if err != nil {
		t.Fatalf("RecordInstallStart for package failed: %v", err)
	}

	if errSkip := db.RecordInstallSkipped("package", "docker", "already installed"); errSkip != nil {
		t.Fatalf("RecordInstallSkipped failed: %v", errSkip)
	}

	rec, err := db.GetInstallLog(id)
	if err != nil || rec.Status != "skipped" || rec.ErrorMessage != "already installed" {
		t.Errorf("expected skipped status, got %+v", rec)
	}
}

func testDirectFallbacks(t *testing.T, db *InstallationSplitDB) {
	if err := db.RecordInstallSuccess("profile", "dev"); err != nil {
		t.Fatalf("direct RecordInstallSuccess failed: %v", err)
	}
	if err := db.RecordInstallFailure("profile", "cloud", 2, "timeout"); err != nil {
		t.Fatalf("direct RecordInstallFailure failed: %v", err)
	}
	if err := db.RecordInstallSkipped("profile", "ai", "no cpu"); err != nil {
		t.Fatalf("direct RecordInstallSkipped failed: %v", err)
	}
}

func testListInstallLogs(t *testing.T, db *InstallationSplitDB) {
	testListInstallLogsByType(t, db)
	testListInstallLogsAll(t, db)
}

func testListInstallLogsByType(t *testing.T, db *InstallationSplitDB) {
	toolLogs, err := db.ListInstallLogs("tool", 10)
	if err != nil {
		t.Fatalf("ListInstallLogs by tool failed: %v", err)
	}
	if len(toolLogs) == 0 {
		t.Errorf("expected tool logs to be non-empty")
	}
}

func testListInstallLogsAll(t *testing.T, db *InstallationSplitDB) {
	allLogs, err := db.ListInstallLogs("", 50)
	if err != nil {
		t.Fatalf("ListInstallLogs all failed: %v", err)
	}
	if len(allLogs) == 0 {
		t.Errorf("expected all logs to be non-empty")
	}
}

func testGetNotFound(t *testing.T, db *InstallationSplitDB) {
	rec, err := db.GetInstallLog("non-existent-uuid")
	if err == nil {
		t.Errorf("expected error for non-existent log, got record: %+v", rec)
	}
}
