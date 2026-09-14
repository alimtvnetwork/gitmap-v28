package store

import (
	"path/filepath"
	"testing"
)

func TestProfileInstallation_Lifecycle(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "installation.db")

	splitDB, err := OpenInstallationSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenInstallationSplitDBAt failed: %v", err)
	}

	defer splitDB.Close()

	testInitSchemaIdempotency(t, splitDB)
	testProfileRecordLifecycle(t, splitDB)
	testPackageInstallationRecord(t, splitDB)
}

func testInitSchemaIdempotency(t *testing.T, splitDB *InstallationSplitDB) {
	if err := splitDB.InitSchema(); err != nil {
		t.Fatalf("InitSchema second call failed: %v", err)
	}
}

func testProfileRecordLifecycle(t *testing.T, splitDB *InstallationSplitDB) {
	isInitialInstalled := splitDB.IsProfileInstalled("dev")
	if isInitialInstalled {
		t.Errorf("Expected dev to be not installed initially")
	}

	profId, err := splitDB.RecordProfileStart("dev", "Developer Tools", "Full dev stack", 5)
	if err != nil || profId <= 0 {
		t.Fatalf("RecordProfileStart failed: id=%d, err=%v", profId, err)
	}

	isDuringInstalled := splitDB.IsProfileInstalled("dev")
	if isDuringInstalled {
		t.Errorf("Expected dev to still be not installed while in progress")
	}

	completeProfileInstallation(t, splitDB, profId)
	verifyProfileDetails(t, splitDB, profId)
	verifyProfileDeletion(t, splitDB)
}

func completeProfileInstallation(t *testing.T, splitDB *InstallationSplitDB, profId int64) {
	err := splitDB.RecordProfileCompletion(profId, 4500, true, 5, 0, "", "", "all good")
	if err != nil {
		t.Fatalf("RecordProfileCompletion failed: %v", err)
	}

	isFinalInstalled := splitDB.IsProfileInstalled("dev")
	if !isFinalInstalled {
		t.Errorf("Expected dev to be installed after completion")
	}
}

func verifyProfileDetails(t *testing.T, splitDB *InstallationSplitDB, profId int64) {
	rec, err := splitDB.GetProfileInstallation("dev")
	if err != nil || rec == nil {
		t.Fatalf("GetProfileInstallation failed: %v", err)
	}

	assertProfileRecordValid(t, rec, profId)

	list, err := splitDB.ListProfileInstallations()
	if err != nil || len(list) == 0 {
		t.Fatalf("ListProfileInstallations failed: len=%d, err=%v", len(list), err)
	}
}

func assertProfileRecordValid(t *testing.T, rec *ProfileInstallationRecord, profId int64) {
	if rec.ProfileInstallationId != profId || rec.TotalTools != 5 {
		t.Errorf("Unexpected profile record: %+v", rec)
	}

	if rec.IsSuccess {
		return
	}

	t.Errorf("Unexpected profile record is not success: %+v", rec)
}

func verifyProfileDeletion(t *testing.T, splitDB *InstallationSplitDB) {
	if err := splitDB.DeleteProfileInstallation("dev"); err != nil {
		t.Fatalf("DeleteProfileInstallation failed: %v", err)
	}

	if splitDB.IsProfileInstalled("dev") {
		t.Errorf("Expected dev to be not installed after deletion")
	}
}

func testPackageInstallationRecord(t *testing.T, splitDB *InstallationSplitDB) {
	profId, err := splitDB.RecordProfileStart("ai", "AI Stack", "AI tools", 2)
	if err != nil {
		t.Fatalf("RecordProfileStart for ai failed: %v", err)
	}

	pkgRec := buildSampleFailedPackageRecord(profId)
	if err := splitDB.RecordPackageInstallation(pkgRec); err != nil {
		t.Fatalf("RecordPackageInstallation failed: %v", err)
	}

	verifyPackageInstallationRecord(t, splitDB, profId)
}

func buildSampleFailedPackageRecord(profId int64) PackageInstallationRecord {
	return PackageInstallationRecord{
		ProfileInstallationId: profId,
		PackageName:           "ollama",
		Version:               "0.1.28",
		PackageManager:        "curl",
		InstallPath:           "/usr/local/bin/ollama",
		Status:                "failed",
		IsSuccess:             false,
		ExitCode:              1,
		DurationMs:            1200,
		Stderr:                "connection timeout",
		StackTrace:            "goroutine 1 [running]:\nmain.main()",
	}
}

func verifyPackageInstallationRecord(t *testing.T, splitDB *InstallationSplitDB, profId int64) {
	pkgs, err := splitDB.ListPackagesForProfile(profId)
	if err != nil || len(pkgs) != 1 {
		t.Fatalf("ListPackagesForProfile failed: len=%d, err=%v", len(pkgs), err)
	}

	if pkgs[0].PackageName != "ollama" || pkgs[0].StackTrace == "" {
		t.Errorf("Expected stack trace preserved in package record: %+v", pkgs[0])
	}
}
