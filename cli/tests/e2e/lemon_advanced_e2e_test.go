//go:build e2e

package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlemon"
)

// captureAdvancedOutput captures stdout output emitted by the given function.
func captureAdvancedOutput(t *testing.T, action func()) string {
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

// backupAdvancedStateFile copies state file to backupPath if it exists.
func backupAdvancedStateFile(statePath, backupPath string) bool {
	if _, statErr := os.Stat(statePath); statErr != nil {
		return false
	}
	data, err := os.ReadFile(statePath)
	if err != nil {
		return false
	}
	return os.WriteFile(backupPath, data, 0o644) == nil
}

// TestLemonExecutor_TimeoutAndBufferSafety verifies timeout clamping and buffer truncation boundaries.
func TestLemonExecutor_TimeoutAndBufferSafety(t *testing.T) {
	if timeoutZero := cmdlemon.ResolveTimeout(0); timeoutZero != 60*time.Second {
		t.Errorf("expected 0 to resolve to 60s, got: %v", timeoutZero)
	}
	if timeoutNeg := cmdlemon.ResolveTimeout(-15 * time.Second); timeoutNeg != 60*time.Second {
		t.Errorf("expected negative duration to resolve to 60s, got: %v", timeoutNeg)
	}
	if timeoutMax := cmdlemon.ResolveTimeout(450 * time.Second); timeoutMax != 300*time.Second {
		t.Errorf("expected >300s to clamp to 300s, got: %v", timeoutMax)
	}
	if timeoutNormal := cmdlemon.ResolveTimeout(120 * time.Second); timeoutNormal != 120*time.Second {
		t.Errorf("expected 120s to be preserved, got: %v", timeoutNormal)
	}

	exactData := []byte("ABCDE")
	outExact, isTruncExact := cmdlemon.TruncateBuffer(exactData, 5)
	if isTruncExact || outExact != "ABCDE" {
		t.Errorf("exact buffer truncation mismatch: out=%q, isTrunc=%v", outExact, isTruncExact)
	}

	plus1Data := []byte("ABCDEF")
	outPlus1, isTruncPlus1 := cmdlemon.TruncateBuffer(plus1Data, 5)
	if !isTruncPlus1 || !strings.Contains(outPlus1, "...[output truncated]") || !strings.HasPrefix(outPlus1, "ABCDE") {
		t.Errorf("plus1 buffer truncation mismatch: out=%q, isTrunc=%v", outPlus1, isTruncPlus1)
	}
}

// TestLemonDaemon_CorruptedStateRecovery verifies resilient recovery across invalid state files.
func TestLemonDaemon_CorruptedStateRecovery(t *testing.T) {
	stateDir := filepath.Dir(cmdlemon.StateFilePath)
	_ = os.MkdirAll(stateDir, 0o755)
	statePath := cmdlemon.StateFilePath
	backupPath := statePath + ".corruptbackup"
	hasBackup := backupAdvancedStateFile(statePath, backupPath)
	defer func() {
		if hasBackup {
			_ = os.Rename(backupPath, statePath)
		} else {
			_ = os.Remove(statePath)
		}
	}()

	corruptedPayloads := []struct {
		name    string
		content []byte
	}{
		{"truncated JSON", []byte(`{"isRunning": true, "pollCount": 42, "lastExec`)},
		{"zero bytes", []byte("")},
		{"invalid syntax", []byte(`{INVALID_DATA_NOT_JSON: ???}`)},
	}

	for _, cp := range corruptedPayloads {
		if writeErr := os.WriteFile(statePath, cp.content, 0o644); writeErr != nil {
			t.Fatalf("failed to write test corrupted state (%s): %v", cp.name, writeErr)
		}

		outText := captureAdvancedOutput(t, func() {
			if err := cmdlemon.RunLemonCLI([]string{"status"}); err != nil {
				t.Errorf("RunLemonCLI('status') failed on %s: %v", cp.name, err)
			}
		})
		if !strings.Contains(outText, "GitMap Lemon Daemon Status") {
			t.Errorf("expected status header on %s, got: %s", cp.name, outText)
		}

		outJSON := captureAdvancedOutput(t, func() {
			if err := cmdlemon.RunLemonCLI([]string{"status", "--json"}); err != nil {
				t.Errorf("RunLemonCLI('status', '--json') failed on %s: %v", cp.name, err)
			}
		})
		var recoveredState cmdlemon.DaemonState
		if err := json.Unmarshal([]byte(outJSON), &recoveredState); err != nil {
			t.Errorf("failed to unmarshal JSON status on %s: %v", cp.name, err)
		}

		daemon := cmdlemon.NewLemonDaemon(cmdlemon.DaemonConfig{IsOnce: true})
		if daemon == nil {
			t.Errorf("expected non-nil daemon instance on %s", cp.name)
		}
	}
}

// TestLemonDaemon_ConcurrentExecutionSafety verifies thread safety and atomic tracking across parallel executions.
func TestLemonDaemon_ConcurrentExecutionSafety(t *testing.T) {
	stateDir := filepath.Dir(cmdlemon.StateFilePath)
	_ = os.MkdirAll(stateDir, 0o755)
	statePath := cmdlemon.StateFilePath
	backupPath := statePath + ".concurrentbackup"
	hasBackup := backupAdvancedStateFile(statePath, backupPath)
	defer func() {
		if hasBackup {
			_ = os.Rename(backupPath, statePath)
		} else {
			_ = os.Remove(statePath)
		}
	}()

	var activePollCount atomic.Int64
	daemon := cmdlemon.NewLemonDaemon(cmdlemon.DaemonConfig{
		IsOnce: true,
		PollFunc: func(ctx context.Context) error {
			activePollCount.Add(1)
			time.Sleep(3 * time.Millisecond)
			return nil
		},
	})

	const concurrentWorkers = 12
	var wg sync.WaitGroup
	wg.Add(concurrentWorkers)

	for idx := 0; idx < concurrentWorkers; idx++ {
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = daemon.Run(ctx)
		}()
	}

	wg.Wait()

	state := daemon.GetState()
	if activePollCount.Load() != concurrentWorkers {
		t.Errorf("expected %d active polls, got: %d", concurrentWorkers, activePollCount.Load())
	}
	if state.PollCount != concurrentWorkers {
		t.Errorf("expected PollCount %d, got: %d", concurrentWorkers, state.PollCount)
	}
}

// TestLemonDaemon_SubMillisecondAndFractionalIntervals verifies interval clamping across fractional and microsecond values.
func TestLemonDaemon_SubMillisecondAndFractionalIntervals(t *testing.T) {
	buf := &bytes.Buffer{}
	guard := cmdlemon.NewIntervalGuard(buf)

	// Fractional 4.9999s clamped with notice
	durClamped, errClamped := guard.ParseAndResolveInterval("4.9999s")
	if errClamped != nil || durClamped != cmdlemon.MinPollingInterval {
		t.Errorf("expected 4.9999s to resolve to 5s, got: %v (err: %v)", durClamped, errClamped)
	}
	if !strings.Contains(buf.String(), "Notice: interval clamped to minimum 5s safety threshold") {
		t.Errorf("expected clamp notice for 4.9999s, got: %s", buf.String())
	}

	// Fractional 5.0001s preserved without notice
	buf.Reset()
	durPreserved, errPreserved := guard.ParseAndResolveInterval("5.0001s")
	expectedPreserved := 5000100 * time.Microsecond
	if errPreserved != nil || durPreserved != expectedPreserved {
		t.Errorf("expected 5.0001s to be preserved (%v), got: %v (err: %v)", expectedPreserved, durPreserved, errPreserved)
	}
	if buf.Len() > 0 {
		t.Errorf("expected no notice for 5.0001s, got: %s", buf.String())
	}

	// Microsecond 500us clamped with notice
	buf.Reset()
	durMicro, errMicro := guard.ParseAndResolveInterval("500us")
	if errMicro != nil || durMicro != cmdlemon.MinPollingInterval {
		t.Errorf("expected 500us to resolve to 5s, got: %v (err: %v)", durMicro, errMicro)
	}
	if !strings.Contains(buf.String(), "Notice: interval clamped") {
		t.Errorf("expected clamp notice for 500us, got: %s", buf.String())
	}

	// Nanosecond 250ns clamped with notice
	buf.Reset()
	durNano, errNano := guard.ParseAndResolveInterval("250ns")
	if errNano != nil || durNano != cmdlemon.MinPollingInterval {
		t.Errorf("expected 250ns to resolve to 5s, got: %v (err: %v)", durNano, errNano)
	}
	if !strings.Contains(buf.String(), "Notice: interval clamped") {
		t.Errorf("expected clamp notice for 250ns, got: %s", buf.String())
	}
}
