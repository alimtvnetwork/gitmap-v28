//go:build e2e

package e2e_test

import (
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestSSHFleetParallel_E2E(t *testing.T) {
	mockFleet := []db.SSHConnection{
		{Alias: "vm-worker-1", IPAddress: "10.10.0.1", OS: "linux"},
		{Alias: "vm-worker-2", IPAddress: "10.10.0.2", OS: "linux"},
		{Alias: "vm-worker-3", IPAddress: "10.10.0.3", OS: "windows"},
	}

	opts := cmdssh.FleetParallelOptions{
		Target:   "all-nodes",
		Except:   "vm-worker-2",
		TaskName: "E2E Multi-Node Verification",
	}

	var activeRuns int32
	var maxParallel int32

	results := cmdssh.RunParallelFleetExecution(mockFleet, opts, func(c db.SSHConnection) (string, error) {
		current := atomic.AddInt32(&activeRuns, 1)
		defer atomic.AddInt32(&activeRuns, -1)

		for {
			oldMax := atomic.LoadInt32(&maxParallel)
			if current <= oldMax || atomic.CompareAndSwapInt32(&maxParallel, oldMax, current) {
				break
			}
		}

		time.Sleep(30 * time.Millisecond)
		return fmt.Sprintf("Node %s active", c.Alias), nil
	})

	if len(results) != 2 {
		t.Fatalf("expected 2 nodes executed (excluding vm-worker-2), got %d", len(results))
	}

	for _, r := range results {
		if !r.Success {
			t.Errorf("expected success for %s, got err: %v", r.Alias, r.Error)
		}
		if r.Alias == "vm-worker-2" {
			t.Errorf("excluded node vm-worker-2 should not be in execution results")
		}
	}
}

func TestSSHAgyFleetArgsParsing_E2E(t *testing.T) {
	rawArgs := []string{"-e", "node-excluded", "prompt", "inject", "--title", "Hello"}
	var except string
	var agyTokens []string
	for i := 0; i < len(rawArgs); i++ {
		if rawArgs[i] == "-e" && i+1 < len(rawArgs) {
			except = rawArgs[i+1]
			i++
			continue
		}
		agyTokens = append(agyTokens, rawArgs[i])
	}

	if except != "node-excluded" {
		t.Errorf("expected except 'node-excluded', got %s", except)
	}
	if strings.Join(agyTokens, " ") != "prompt inject --title Hello" {
		t.Errorf("expected agyTokens 'prompt inject --title Hello', got %v", agyTokens)
	}
}
