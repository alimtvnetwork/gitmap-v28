package cmdschedule

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type mockExecutorState struct {
	calledCount int
	lastParams  OSActionParams
	returnErr   error
	onExec      func(params OSActionParams)
}

func (m *mockExecutorState) execute(params OSActionParams) error {
	m.calledCount++
	m.lastParams = params
	if m.onExec != nil {
		m.onExec(params)
	}

	return m.returnErr
}

func setupScheduleOSTest(t *testing.T) (*mockExecutorState, func()) {
	t.Helper()
	mock := &mockExecutorState{}
	resetEnv := overrideScheduleTestEnv(t, mock)

	return mock, resetEnv
}

func overrideScheduleTestEnv(t *testing.T, mock *mockExecutorState) func() {
	origExec, origOS, origSleep, origFile := DefaultOSActionExecutor, countdownOS, countdownSleepFn, stateFileOverride
	stateFileOverride = filepath.Join(t.TempDir(), "test_power.json")
	DefaultOSActionExecutor = mock.execute
	countdownSleepFn = func(time.Duration) {}

	return func() {
		DefaultOSActionExecutor, countdownOS, countdownSleepFn, stateFileOverride = origExec, origOS, origSleep, origFile
	}
}

func captureOutput(fn func()) string {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = orig
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()

	return buf.String()
}

func hasTargetArg(args []string, target string) bool {
	for _, a := range args {
		if a == target {
			return true
		}
	}

	return false
}

func assertOSCommand(t *testing.T, osName string, action OSActionType, wantExe, wantSub string) {
	t.Helper()
	exe, args := BuildOSActionCommand(action, 60, osName)
	if exe != wantExe || !hasTargetArg(args, wantSub) {
		t.Errorf("BuildOSActionCommand(%s, %s) = (%s, %v), want %s / %s", action, osName, exe, args, wantExe, wantSub)
	}
}

func TestBuildOSActionCommand_Platforms(t *testing.T) {
	_, cleanup := setupScheduleOSTest(t)
	defer cleanup()

	assertOSCommand(t, "windows", OSActionShutdown, "shutdown", "/s")
	assertOSCommand(t, "windows", OSActionRestart, "shutdown", "/r")
	assertOSCommand(t, "linux", OSActionShutdown, "sudo", "poweroff")
	assertOSCommand(t, "linux", OSActionRestart, "sudo", "reboot")
	assertOSCommand(t, "darwin", OSActionShutdown, "sudo", "-h")
	assertOSCommand(t, "darwin", OSActionRestart, "sudo", "-r")
}

func assertCancelCommand(t *testing.T, osName string, wantExe, wantFlag string) {
	t.Helper()
	exe, args := BuildCancelOSActionCommand(osName)
	if exe != wantExe || !hasTargetArg(args, wantFlag) {
		t.Errorf("BuildCancelOSActionCommand(%s) = (%s, %v), want %s / %s", osName, exe, args, wantExe, wantFlag)
	}
}

func TestBuildCancelOSActionCommand_AllPlatforms(t *testing.T) {
	_, cleanup := setupScheduleOSTest(t)
	defer cleanup()

	assertCancelCommand(t, "windows", "shutdown", "/a")
	assertCancelCommand(t, "linux", "shutdown", "-c")
	assertCancelCommand(t, "darwin", "shutdown", "-c")
}

func TestRunSchedulePowerCLI_CancelAliases(t *testing.T) {
	mock, cleanup := setupScheduleOSTest(t)
	defer cleanup()

	for _, sub := range []string{"cancel", "abort", "stop", "rm"} {
		mock.calledCount = 0
		if err := RunSchedulePowerCLI(OSActionShutdown, []string{sub}); err != nil {
			t.Fatalf("unexpected error for %s: %v", sub, err)
		}

		if mock.calledCount != 1 || mock.lastParams.Action != OSActionCancel {
			t.Errorf("expected OSActionCancel call for %s, got: %+v", sub, mock.lastParams)
		}
	}
}

func TestRunSchedulePowerCLI_StatusAndLs(t *testing.T) {
	_, cleanup := setupScheduleOSTest(t)
	defer cleanup()

	for _, sub := range []string{"status", "ls"} {
		out := captureOutput(func() {
			_ = RunSchedulePowerCLI(OSActionShutdown, []string{sub})
		})
		if !strings.Contains(out, "No scheduled shutdown") {
			t.Errorf("expected no active schedule for %s, got: %s", sub, out)
		}
	}
}

func armScheduleForTest(action OSActionType, dur time.Duration) {
	now := time.Now()
	_ = SavePowerScheduleState(PowerScheduleState{
		Action: action, Status: "ARMED", ScheduledAt: now, TriggerAt: now.Add(dur),
	})
}

func TestShowSchedulePowerStatus_ActiveRemaining(t *testing.T) {
	_, cleanup := setupScheduleOSTest(t)
	defer cleanup()

	armScheduleForTest(OSActionShutdown, 90*time.Minute)
	out := captureOutput(func() {
		_ = ShowSchedulePowerStatus(OSActionShutdown)
	})
	if !strings.Contains(out, "Active Power Schedule") || !strings.Contains(out, "Remaining:") {
		t.Fatalf("expected active schedule with remaining time, got: %s", out)
	}
}

func assertScheduleDurationSeconds(t *testing.T, input string, wantSecs int64) {
	t.Helper()
	mock, cleanup := setupScheduleOSTest(t)
	defer cleanup()

	if err := RunSchedulePowerCLI(OSActionShutdown, []string{input}); err != nil {
		t.Fatalf("RunSchedulePowerCLI failed for %s: %v", input, err)
	}

	if mock.lastParams.Seconds != wantSecs {
		t.Errorf("for input %q: expected %d seconds, got %d", input, wantSecs, mock.lastParams.Seconds)
	}
}

func TestRunSchedulePowerCLI_DurationImmediate(t *testing.T) {
	assertScheduleDurationSeconds(t, "0", 0)
	assertScheduleDurationSeconds(t, "now", 0)
}

func TestRunSchedulePowerCLI_DurationFormats(t *testing.T) {
	assertScheduleDurationSeconds(t, "1:45hr", 6300)
	assertScheduleDurationSeconds(t, "01:30:00", 5400)
	assertScheduleDurationSeconds(t, "2h", 7200)
	assertScheduleDurationSeconds(t, "30m", 1800)
	assertScheduleDurationSeconds(t, "15s", 15)
}

func attachActiveChecker(mock *mockExecutorState) *bool {
	hasChecked := false
	mock.onExec = func(params OSActionParams) {
		active, err := GetActivePowerSchedule()
		if err == nil && active != nil && active.Status == "ARMED" {
			hasChecked = true
		}
	}

	return &hasChecked
}

func TestRunSchedulePowerCLI_SequencingArmedDuringExecution(t *testing.T) {
	mock, cleanup := setupScheduleOSTest(t)
	defer cleanup()

	hasCheckedActive := attachActiveChecker(mock)
	if err := RunSchedulePowerCLI(OSActionShutdown, []string{"30m"}); err != nil {
		t.Fatalf("RunSchedulePowerCLI failed: %v", err)
	}

	if !*hasCheckedActive {
		t.Fatal("expected mock.onExec to check active state")
	}
}

func buildCancelCountdownParams(tickCount *int) OSActionParams {
	return OSActionParams{
		Action: OSActionShutdown,
		Delay:  5 * time.Second,
		OnTick: func(rem time.Duration) {
			*tickCount++
			_ = CancelActivePowerSchedule()
		},
	}
}

func TestPerformDelayCountdown_Cancellation(t *testing.T) {
	_, cleanup := setupScheduleOSTest(t)
	defer cleanup()

	countdownOS = "linux"
	armScheduleForTest(OSActionShutdown, time.Hour)
	tickCount := 0
	performDelayCountdown(buildCancelCountdownParams(&tickCount))
	if tickCount != 1 {
		t.Errorf("expected countdown to abort after 1 tick on cancel, got %d", tickCount)
	}
}

func TestRunSchedulePowerCLI_InvalidDuration(t *testing.T) {
	mock, cleanup := setupScheduleOSTest(t)
	defer cleanup()

	if err := RunSchedulePowerCLI(OSActionRestart, []string{"not-a-valid-duration"}); err == nil {
		t.Fatal("expected error for invalid duration, got nil")
	}

	if mock.calledCount != 0 {
		t.Errorf("expected 0 calls on invalid duration, got %d", mock.calledCount)
	}
}

func TestRunSchedulePowerCLI_ExecutorErrorCleansUp(t *testing.T) {
	mock, cleanup := setupScheduleOSTest(t)
	defer cleanup()
	mock.returnErr = errors.New("simulated failure")

	if err := RunSchedulePowerCLI(OSActionShutdown, []string{"10s"}); err == nil {
		t.Fatal("expected error from executor, got nil")
	}

	active, err := GetActivePowerSchedule()
	if err != nil || active != nil {
		t.Fatalf("expected nil active after executor failure, got %+v", active)
	}
}

func TestRunSchedulePowerCLI_Help(t *testing.T) {
	mock, cleanup := setupScheduleOSTest(t)
	defer cleanup()

	if err := RunSchedulePowerCLI(OSActionShutdown, []string{"--help"}); err != nil {
		t.Errorf("expected nil error on help, got %v", err)
	}

	if mock.calledCount != 0 {
		t.Errorf("expected 0 executor calls on help, got %d", mock.calledCount)
	}
}

func TestBuildOSActionCommand_CurrentOS(t *testing.T) {
	exe, args := BuildOSActionCommand(OSActionShutdown, 10)
	if exe == "" || len(args) == 0 {
		t.Errorf("expected non-empty exe and args, got %q %v", exe, args)
	}
}
