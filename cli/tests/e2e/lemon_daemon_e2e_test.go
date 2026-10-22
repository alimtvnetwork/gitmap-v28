//go:build e2e

package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlemon"
)

// captureStdout redirects os.Stdout and returns captured output string.
func captureStdout(t *testing.T, action func()) string {
	t.Helper()
	pipeReader, pipeWriter, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("failed to create pipe: %v", pipeErr)
	}
	savedStdout := os.Stdout
	os.Stdout = pipeWriter
	outChan := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, pipeReader)
		outChan <- buf.String()
	}()
	action()
	_ = pipeWriter.Close()
	os.Stdout = savedStdout
	captured := <-outChan
	_ = pipeReader.Close()
	return captured
}

// TestLemonCLI_RoutingAndHelp tests help flags and unknown subcommand rejection.
func TestLemonCLI_RoutingAndHelp(t *testing.T) {
	for _, args := range [][]string{{}, {"help"}, {"-h"}, {"--help"}} {
		if err := cmdlemon.RunLemonCLI(args); err != nil {
			t.Errorf("RunLemonCLI(%v) returned unexpected error: %v", args, err)
		}
	}
	err := cmdlemon.RunLemonCLI([]string{"unknown-cmd"})
	if err == nil {
		t.Fatal("expected error for unknown subcommand, got nil")
	}
	hasUnknownMsg := strings.Contains(err.Error(), "unknown lemon subcommand")
	if !hasUnknownMsg {
		t.Errorf("expected error message to contain 'unknown lemon subcommand', got: %v", err)
	}
}

func backupTestState(statePath, backupPath string) bool {
	if _, statErr := os.Stat(statePath); statErr != nil {
		return false
	}
	return os.Rename(statePath, backupPath) == nil
}

// TestLemonCLI_StatusOutput tests status commands with and without state file.
func TestLemonCLI_StatusOutput(t *testing.T) {
	stateDir := filepath.Dir(cmdlemon.StateFilePath)
	_ = os.MkdirAll(stateDir, 0o755)
	t.Cleanup(func() { _ = os.RemoveAll(stateDir) })
	statePath := cmdlemon.StateFilePath
	backupPath := statePath + ".testbackup"
	hasBackup := backupTestState(statePath, backupPath)
	defer func() {
		if hasBackup {
			_ = os.Rename(backupPath, statePath)
		}
	}()
	_ = os.Remove(statePath)

	missingTextOut := captureStdout(t, func() {
		if statusErr := cmdlemon.RunLemonCLI([]string{"status"}); statusErr != nil {
			t.Errorf("status failed when state file missing: %v", statusErr)
		}
	})
	hasStatusHeader := strings.Contains(missingTextOut, "GitMap Lemon Daemon Status")
	if !hasStatusHeader {
		t.Errorf("expected status header in output, got: %s", missingTextOut)
	}

	missingJSONOut := captureStdout(t, func() {
		if statusJSONErr := cmdlemon.RunLemonCLI([]string{"status", "--json"}); statusJSONErr != nil {
			t.Errorf("status --json failed when state file missing: %v", statusJSONErr)
		}
	})
	var parsedMissingState cmdlemon.DaemonState
	if err := json.Unmarshal([]byte(missingJSONOut), &parsedMissingState); err != nil {
		t.Fatalf("expected valid JSON when state file missing: %v", err)
	}

	dummyState := cmdlemon.DaemonState{Pid: 12345, PollCount: 42, LastPollTime: time.Now().UTC(), LastExecutedCommand: "status"}
	encodedData, encodeErr := json.Marshal(dummyState)
	if encodeErr != nil {
		t.Fatalf("failed to encode dummy state: %v", encodeErr)
	}
	if err := os.WriteFile(statePath, encodedData, 0o644); err != nil {
		t.Fatalf("failed to write test state file: %v", err)
	}

	validJSONOut := captureStdout(t, func() {
		if err := cmdlemon.RunLemonCLI([]string{"status", "--json"}); err != nil {
			t.Errorf("status --json failed with state file: %v", err)
		}
	})
	var parsedValidState cmdlemon.DaemonState
	if err := json.Unmarshal([]byte(validJSONOut), &parsedValidState); err != nil {
		t.Fatalf("expected valid JSON output with state file: %v", err)
	}
	hasExpectedPid := parsedValidState.Pid == 12345
	if !hasExpectedPid {
		t.Errorf("expected PID 12345, got %d", parsedValidState.Pid)
	}
}

// TestLemonCLI_PollOnceExecution tests single-pass execution and state file update.
func TestLemonCLI_PollOnceExecution(t *testing.T) {
	stateDir := filepath.Dir(cmdlemon.StateFilePath)
	_ = os.MkdirAll(stateDir, 0o755)
	t.Cleanup(func() { _ = os.RemoveAll(stateDir) })
	beforePollTime := time.Now().Add(-2 * time.Second)
	if err := cmdlemon.RunLemonCLI([]string{"poll"}); err != nil {
		t.Fatalf("unexpected error running poll: %v", err)
	}

	stateData, readErr := os.ReadFile(cmdlemon.StateFilePath)
	if readErr != nil {
		t.Fatalf("failed to read state file after poll: %v", readErr)
	}
	var state cmdlemon.DaemonState
	if err := json.Unmarshal(stateData, &state); err != nil {
		t.Fatalf("failed to unmarshal state file: %v", err)
	}
	hasPositivePid := state.Pid > 0
	if !hasPositivePid {
		t.Errorf("expected positive PID, got %d", state.Pid)
	}
	hasUpdatedTime := state.LastPollTime.After(beforePollTime)
	if !hasUpdatedTime {
		t.Errorf("expected LastPollTime (%v) to be after %v", state.LastPollTime, beforePollTime)
	}
	if state.IsRunning {
		t.Errorf("expected IsRunning to be false after poll execution")
	}
	if err := cmdlemon.RunLemonCLI([]string{"run", "--once"}); err != nil {
		t.Fatalf("unexpected error running run --once: %v", err)
	}
}

// TestLemonCLI_FlagParsingVariations tests interval flag variations and fallbacks.
func TestLemonCLI_FlagParsingVariations(t *testing.T) {
	stateDir := filepath.Dir(cmdlemon.StateFilePath)
	_ = os.MkdirAll(stateDir, 0o755)
	t.Cleanup(func() { _ = os.RemoveAll(stateDir) })
	testCases := []struct {
		name     string
		args     []string
		rawVal   string
		expected time.Duration
	}{
		{"--interval 12s", []string{"--interval", "12s"}, "12s", 12 * time.Second},
		{"-i 15s", []string{"-i", "15s"}, "15s", 15 * time.Second},
		{"--interval=20s", []string{"--interval=20s"}, "20s", 20 * time.Second},
		{"--interval missing arg", []string{"--interval"}, "", cmdlemon.DefaultPollingInterval},
		{"--interval invalid", []string{"--interval", "invalid"}, "invalid", cmdlemon.DefaultPollingInterval},
	}
	guard := cmdlemon.NewIntervalGuard(io.Discard)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dur, _ := guard.ParseAndResolveInterval(tc.rawVal)
			hasExpectedGuardDur := dur == tc.expected
			if !hasExpectedGuardDur {
				t.Errorf("guard resolution expected %v, got %v", tc.expected, dur)
			}
			cliArgs := append([]string{"run", "--once"}, tc.args...)
			if err := cmdlemon.RunLemonCLI(cliArgs); err != nil {
				t.Errorf("RunLemonCLI(%v) returned unexpected error: %v", cliArgs, err)
			}
		})
	}
}

// TestLemonDaemon_ContextTermination tests clean shutdown without leaking goroutines.
func TestLemonDaemon_ContextTermination(t *testing.T) {
	initialGoroutines := runtime.NumGoroutine()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pollCount := int64(0)
	cfg := cmdlemon.DaemonConfig{
		Interval: 5 * time.Second,
		IsOnce:   false,
		PollFunc: func(pCtx context.Context) error {
			atomic.AddInt64(&pollCount, 1)
			return nil
		},
	}
	daemon := cmdlemon.NewLemonDaemon(cfg)
	errChan := make(chan error, 1)
	go func() { errChan <- daemon.Run(ctx) }()

	isStarted := false
	for i := 0; i < 100; i++ {
		if daemon.GetState().IsRunning {
			isStarted = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !isStarted {
		t.Fatal("expected daemon.GetState().IsRunning to become true")
	}

	cancel()
	select {
	case runErr := <-errChan:
		if runErr != nil {
			t.Errorf("expected clean shutdown with nil error, got: %v", runErr)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("daemon did not terminate within timeout")
	}

	isRunningAfterShutdown := daemon.GetState().IsRunning
	if isRunningAfterShutdown {
		t.Errorf("expected daemon.GetState().IsRunning to be false after shutdown")
	}
	time.Sleep(50 * time.Millisecond)
	finalGoroutines := runtime.NumGoroutine()
	hasLeakedGoroutines := finalGoroutines > initialGoroutines+2
	if hasLeakedGoroutines {
		t.Errorf("detected possible goroutine leak: initial=%d, final=%d", initialGoroutines, finalGoroutines)
	}
}

// TestLemonDaemon_IntervalBoundaryConditions tests exact 5s, sub-5s, and 0/negative clamping.
func TestLemonDaemon_IntervalBoundaryConditions(t *testing.T) {
	testCases := []struct {
		requested  time.Duration
		expected   time.Duration
		wantNotice bool
	}{
		{5 * time.Second, 5 * time.Second, false},
		{4999 * time.Millisecond, 5 * time.Second, true},
		{0, 8 * time.Second, false},
		{-1 * time.Second, 8 * time.Second, false},
		{-5 * time.Second, 8 * time.Second, false},
	}
	for _, tc := range testCases {
		var buf bytes.Buffer
		guard := cmdlemon.NewIntervalGuard(&buf)
		resolved := guard.ResolveInterval(tc.requested)
		hasExpectedInterval := resolved == tc.expected
		if !hasExpectedInterval {
			t.Errorf("for %v: expected %v, got %v", tc.requested, tc.expected, resolved)
		}
		hasNotice := strings.Contains(buf.String(), "Notice: interval clamped to minimum 5s safety threshold")
		if tc.wantNotice && !hasNotice {
			t.Errorf("for %v: expected clamp notice, got: %s", tc.requested, buf.String())
		}
		if !tc.wantNotice && hasNotice {
			t.Errorf("for %v: expected no clamp notice, got: %s", tc.requested, buf.String())
		}
	}
}
