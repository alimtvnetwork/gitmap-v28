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
		_ = os.Remove(stateFileOverride)
	})
}

func TestPowerScheduleState_SaveAndGet(t *testing.T) {
	setupTestPowerStateFile(t)
	now := time.Now()
	state := PowerScheduleState{
		Action:          OSActionShutdown,
		Status:          "ARMED",
		ScheduledAt:     now,
		TriggerAt:       now.Add(2 * time.Hour),
		DurationSeconds: 7200,
		RawDuration:     "2h",
	}
	if err := SavePowerScheduleState(state); err != nil {
		t.Fatalf("SavePowerScheduleState failed: %v", err)
	}
	active, err := GetActivePowerSchedule()
	if err != nil {
		t.Fatalf("GetActivePowerSchedule failed: %v", err)
	}
	if active == nil {
		t.Fatal("expected active power schedule, got nil")
	}
	if active.Action != OSActionShutdown {
		t.Fatalf("expected action shutdown, got %v", active.Action)
	}
}

func TestPowerScheduleState_Cancel(t *testing.T) {
	setupTestPowerStateFile(t)
	now := time.Now()
	state := PowerScheduleState{
		Action:      OSActionRestart,
		Status:      "ARMED",
		ScheduledAt: now,
		TriggerAt:   now.Add(1 * time.Hour),
	}
	_ = SavePowerScheduleState(state)
	if err := CancelActivePowerSchedule(); err != nil {
		t.Fatalf("CancelActivePowerSchedule failed: %v", err)
	}
	active, err := GetActivePowerSchedule()
	if err != nil || active != nil {
		t.Fatalf("expected nil active after cancel, got %v (err: %v)", active, err)
	}
}

func TestFormatRemainingDuration(t *testing.T) {
	now := time.Now()
	res := FormatRemainingDuration(now.Add(1*time.Hour + 45*time.Minute + 10*time.Second))
	if res == "0s" {
		t.Fatalf("expected non-zero remaining, got %s", res)
	}
	zero := FormatRemainingDuration(now.Add(-10 * time.Second))
	if zero != "0s" {
		t.Fatalf("expected 0s for past time, got %s", zero)
	}
}
