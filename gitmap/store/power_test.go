package store

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/power"
)

func TestPowerProfileAndHistory_CRUD(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "gitmap_test.db")

	db, err := OpenAt(dbPath)
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}

	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	setting := power.Settings{
		Platform:              "windows",
		DisplayTimeoutMinutes: 0,
		SleepTimeoutMinutes:   0,
		IsNeverSleep:          true,
		IsLockDisabled:        true,
	}

	// Save profile
	if err := db.SavePowerProfile("never-sleep", setting, true); err != nil {
		t.Fatalf("SavePowerProfile failed: %v", err)
	}

	// Get active setting
	active, err := db.GetActivePowerSetting()
	if err != nil {
		t.Fatalf("GetActivePowerSetting failed: %v", err)
	}

	if !active.IsNeverSleep || active.DisplayTimeoutMinutes != 0 {
		t.Errorf("Unexpected active power setting: %+v", active)
	}

	// Record history
	if err := db.RecordPowerHistory("never-sleep", setting, "CLI activation"); err != nil {
		t.Fatalf("RecordPowerHistory failed: %v", err)
	}

	history, err := db.ListPowerHistory(10)
	if err != nil || len(history) == 0 {
		t.Fatalf("ListPowerHistory failed or empty: %v", err)
	}

	if history[0].Action != "never-sleep" || !history[0].IsNeverSleep {
		t.Errorf("Unexpected history entry: %+v", history[0])
	}
}
