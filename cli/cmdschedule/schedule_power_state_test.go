package cmdschedule

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestPowerStateFile(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	stateFileOverride = filepath.Join(dir, "test_power_schedule.json")
	t.Cleanup(func() {
		stateFileOverride = ""
	})
}

func saveSamplePowerState(t *testing.T, action OSActionType, now time.Time, dur time.Duration) {
	t.Helper()
	err := SavePowerScheduleState(PowerScheduleState{
		Action: action, Status: "ARMED", ScheduledAt: now, TriggerAt: now.Add(dur),
	})
	if err != nil {
		t.Fatalf("saveSamplePowerState failed: %v", err)
	}
}

func assertPersistedStatus(t *testing.T, wantStatus string) {
	t.Helper()
	persisted, readErr := readPowerScheduleState(resolvePowerStateFilePath())
	if readErr != nil || persisted.Status != wantStatus {
		t.Fatalf("expected status %s, got: %+v (err: %v)", wantStatus, persisted, readErr)
	}
}

func TestPowerScheduleState_SaveAndGet(t *testing.T) {
	setupTestPowerStateFile(t)
	now := time.Now()
	saveSamplePowerState(t, OSActionShutdown, now, 2*time.Hour)

	active, err := GetActivePowerSchedule()
	if err != nil || active == nil || active.Action != OSActionShutdown {
		t.Fatalf("expected active shutdown schedule, got: %+v (err: %v)", active, err)
	}
}

func TestPowerScheduleState_ExpireToCompleted(t *testing.T) {
	setupTestPowerStateFile(t)
	now := time.Now()
	saveSamplePowerState(t, OSActionRestart, now.Add(-10*time.Minute), -5*time.Minute)

	active, err := GetActivePowerSchedule()
	if err != nil || active != nil {
		t.Fatalf("expected nil active for expired schedule, got %+v", active)
	}

	assertPersistedStatus(t, "COMPLETED")
}

func TestPowerScheduleState_Cancel(t *testing.T) {
	setupTestPowerStateFile(t)
	saveSamplePowerState(t, OSActionRestart, time.Now(), time.Hour)
	if err := CancelActivePowerSchedule(); err != nil {
		t.Fatalf("CancelActivePowerSchedule failed: %v", err)
	}

	active, err := GetActivePowerSchedule()
	if err != nil || active != nil {
		t.Fatalf("expected nil active after cancel, got %v (err: %v)", active, err)
	}

	assertPersistedStatus(t, "CANCELED")
}

func TestPowerScheduleState_CorruptFile(t *testing.T) {
	setupTestPowerStateFile(t)
	path := resolvePowerStateFilePath()
	_ = os.WriteFile(path, []byte("{corrupt-json"), 0644)

	active, err := GetActivePowerSchedule()
	if err != nil || active != nil {
		t.Fatalf("expected nil for corrupt file, got active: %v, err: %v", active, err)
	}
}

func TestFormatRemainingDuration_PastAndZero(t *testing.T) {
	now := time.Now()
	if res := FormatRemainingDuration(now.Add(-10 * time.Second)); res != "0s" {
		t.Errorf("expected 0s for past, got %s", res)
	}

	if res := FormatRemainingDuration(now); res != "0s" {
		t.Errorf("expected 0s for now, got %s", res)
	}
}

func TestFormatRemainingDuration_SecondsAndMinutes(t *testing.T) {
	now := time.Now()
	resSec := FormatRemainingDuration(now.Add(45 * time.Second))
	if resSec != "45s" && resSec != "44s" {
		t.Errorf("expected ~45s, got %s", resSec)
	}

	resMin := FormatRemainingDuration(now.Add(15*time.Minute + 30*time.Second))
	if resMin != "15m 30s" && resMin != "15m 29s" {
		t.Errorf("expected ~15m 30s, got %s", resMin)
	}
}

func TestFormatRemainingDuration_Hours(t *testing.T) {
	now := time.Now()
	resHour := FormatRemainingDuration(now.Add(2*time.Hour + 5*time.Minute + 12*time.Second))
	if resHour != "2h 5m 12s" && resHour != "2h 5m 11s" {
		t.Errorf("expected ~2h 5m 12s, got %s", resHour)
	}
}
