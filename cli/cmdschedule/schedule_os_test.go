package cmdschedule

import (
	"errors"
	"testing"
	"time"
)

type mockExecutorState struct {
	calledCount int
	lastParams  OSActionParams
	returnErr   error
}

func (m *mockExecutorState) execute(params OSActionParams) error {
	m.calledCount++
	m.lastParams = params
	return m.returnErr
}

func setupMockExecutor() (*mockExecutorState, func()) {
	mock := &mockExecutorState{}
	orig := DefaultOSActionExecutor
	DefaultOSActionExecutor = mock.execute
	return mock, func() { DefaultOSActionExecutor = orig }
}

func assertMockExecuted(t *testing.T, m *mockExecutorState, action OSActionType, secs int64) {
	if m.calledCount != 1 {
		t.Fatalf("expected 1 call, got %d", m.calledCount)
	}
	if m.lastParams.Action != action || m.lastParams.Seconds != secs {
		t.Errorf("unexpected params: %+v", m.lastParams)
	}
}

func TestRunSchedulePowerCLIShutdownMock(t *testing.T) {
	mock, cleanup := setupMockExecutor()
	defer cleanup()

	if err := RunSchedulePowerCLI(OSActionShutdown, []string{"1s"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertMockExecuted(t, mock, OSActionShutdown, 1)
}

func TestRunSchedulePowerCLIRestartMock(t *testing.T) {
	mock, cleanup := setupMockExecutor()
	defer cleanup()

	if err := RunSchedulePowerCLI(OSActionRestart, []string{"2s"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertMockExecuted(t, mock, OSActionRestart, 2)
}

func TestRunSchedulePowerCLIDurationIntegration(t *testing.T) {
	mock, cleanup := setupMockExecutor()
	defer cleanup()

	if err := RunSchedulePowerCLI(OSActionShutdown, []string{"1:45hr"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.lastParams.Seconds != 6300 || mock.lastParams.Delay != 6300*time.Second {
		t.Errorf("expected 6300s, got %d (%v)", mock.lastParams.Seconds, mock.lastParams.Delay)
	}
}

func TestRunSchedulePowerCLIHelp(t *testing.T) {
	mock, cleanup := setupMockExecutor()
	defer cleanup()

	if err := RunSchedulePowerCLI(OSActionShutdown, []string{"--help"}); err != nil {
		t.Errorf("expected nil error on help, got %v", err)
	}
	if mock.calledCount != 0 {
		t.Errorf("expected 0 executor calls on help, got %d", mock.calledCount)
	}
}

func TestRunSchedulePowerCLIInvalidDuration(t *testing.T) {
	mock, cleanup := setupMockExecutor()
	defer cleanup()

	if err := RunSchedulePowerCLI(OSActionRestart, []string{"invalid-duration"}); err == nil {
		t.Fatal("expected error for invalid duration, got nil")
	}
	if mock.calledCount != 0 {
		t.Errorf("expected 0 calls on invalid duration, got %d", mock.calledCount)
	}
}

func TestRunSchedulePowerCLIExecutorError(t *testing.T) {
	mock, cleanup := setupMockExecutor()
	defer cleanup()
	mock.returnErr = errors.New("simulated executor failure")

	if err := RunSchedulePowerCLI(OSActionShutdown, []string{"0"}); err == nil {
		t.Fatal("expected error from executor, got nil")
	}
	if mock.calledCount != 1 {
		t.Errorf("expected 1 call before failure, got %d", mock.calledCount)
	}
}

func TestBuildWindowsActionCommands(t *testing.T) {
	winExe, winArgs := buildWindowsActionCommand(OSActionShutdown, 60)
	if winExe != "shutdown" || winArgs[0] != "/s" || winArgs[2] != "60" {
		t.Errorf("unexpected windows shutdown: %s %v", winExe, winArgs)
	}
	winRstExe, winRstArgs := buildWindowsActionCommand(OSActionRestart, 120)
	if winRstExe != "shutdown" || winRstArgs[0] != "/r" || winRstArgs[2] != "120" {
		t.Errorf("unexpected windows restart: %s %v", winRstExe, winRstArgs)
	}
}

func TestBuildUnixActionCommands(t *testing.T) {
	darExe, darArgs := buildDarwinActionCommand(OSActionShutdown)
	if darExe != "sudo" || darArgs[1] != "-h" {
		t.Errorf("unexpected darwin shutdown: %s %v", darExe, darArgs)
	}
	linExe, linArgs := buildLinuxActionCommand(OSActionRestart)
	if linExe != "sudo" || linArgs[0] != "reboot" {
		t.Errorf("unexpected linux restart: %s %v", linExe, linArgs)
	}
}

func TestBuildOSActionCommandCurrentOS(t *testing.T) {
	exe, args := BuildOSActionCommand(OSActionShutdown, 10)
	if exe == "" || len(args) == 0 {
		t.Errorf("expected non-empty exe and args, got %q %v", exe, args)
	}
}
