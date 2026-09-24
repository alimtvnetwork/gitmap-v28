package cmdssh

import (
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestRunParallelFleetExecution(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "node-1", IPAddress: "192.168.1.10", OS: "linux"},
		{Alias: "node-2", IPAddress: "192.168.1.11", OS: "linux"},
		{Alias: "node-3", IPAddress: "192.168.1.12", OS: "linux"},
	}

	opts := FleetParallelOptions{
		Target:   "all-nodes",
		Except:   "node-2",
		TaskName: "Test Fleet Task",
	}

	var execCount int32
	results := RunParallelFleetExecution(conns, opts, func(c db.SSHConnection) (string, error) {
		atomic.AddInt32(&execCount, 1)
		time.Sleep(20 * time.Millisecond)
		if c.Alias == "node-3" {
			return "", errors.New("simulated node failure")
		}
		return "node ok", nil
	})

	if len(results) != 2 {
		t.Fatalf("expected 2 executed nodes, got %d", len(results))
	}

	if atomic.LoadInt32(&execCount) != 2 {
		t.Fatalf("expected 2 executions, got %d", execCount)
	}

	// Verify node-1 succeeded
	if !results[0].Success || results[0].Alias != "node-1" {
		t.Errorf("node-1 expected success, got %+v", results[0])
	}

	// Verify node-3 failed
	if results[1].Success || results[1].Alias != "node-3" {
		t.Errorf("node-3 expected failure, got %+v", results[1])
	}
	if !strings.Contains(results[1].Error.Error(), "simulated node failure") {
		t.Errorf("expected failure error, got %v", results[1].Error)
	}
}

func TestParseAgyFleetArgs(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantTarget string
		wantExcept string
		wantCmd    string
	}{
		{
			name:       "default all",
			args:       []string{"status"},
			wantTarget: "all",
			wantExcept: "",
			wantCmd:    "status",
		},
		{
			name:       "with except shorthand",
			args:       []string{"-e", "node-2", "prompt", "ls"},
			wantTarget: "all",
			wantExcept: "node-2",
			wantCmd:    "prompt ls",
		},
		{
			name:       "explicit target and except",
			args:       []string{"-t", "workers", "--except=node-3", "status"},
			wantTarget: "workers",
			wantExcept: "node-3",
			wantCmd:    "status",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			target, except, agyArgs := parseAgyFleetArgs(tc.args)
			if target != tc.wantTarget {
				t.Errorf("target got %s, want %s", target, tc.wantTarget)
			}
			if except != tc.wantExcept {
				t.Errorf("except got %s, want %s", except, tc.wantExcept)
			}
			cmdStr := strings.Join(agyArgs, " ")
			if cmdStr != tc.wantCmd {
				t.Errorf("cmd got %s, want %s", cmdStr, tc.wantCmd)
			}
		})
	}
}
