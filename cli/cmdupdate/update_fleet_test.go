package cmdupdate

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func createMockFleetTargets() []FleetTarget {
	return []FleetTarget{
		{
			ID:       "node-1",
			Alias:    "worker-1",
			IP:       "10.0.0.1",
			Username: "root",
			Port:     22,
			OS:       "linux",
		},
		{
			ID:       "node-2",
			Alias:    "worker-2",
			IP:       "10.0.0.2",
			Username: "ubuntu",
			Port:     22,
			OS:       "linux",
		},
		{
			ID:       "node-3",
			Alias:    "worker-3",
			IP:       "10.0.0.3",
			Username: "admin",
			Port:     22,
			OS:       "windows",
		},
	}
}

func TestExecuteFleetUpdate_AllAliasesExecuteIdentically(t *testing.T) {
	origLoad := LoadFleetTargetsFn
	origExec := ExecuteRemoteUpdateFn
	origLive := CheckConnLivenessFn
	defer func() {
		LoadFleetTargetsFn = origLoad
		ExecuteRemoteUpdateFn = origExec
		CheckConnLivenessFn = origLive
	}()

	CheckConnLivenessFn = func(ctx context.Context, ip string, port int, timeout time.Duration) (bool, string) {
		return true, "online"
	}

	mockTargets := createMockFleetTargets()
	LoadFleetTargetsFn = func() ([]FleetTarget, error) {
		return mockTargets, nil
	}

	commandInvocations := [][]string{
		{"--all", "--except", "worker-2"},
		{"all", "--except", "worker-2"},
		{"ua", "--except", "worker-2"},
	}

	for _, args := range commandInvocations {
		var mu sync.Mutex
		updatedIPs := make(map[string]bool)

		ExecuteRemoteUpdateFn = func(target FleetTarget, opts FleetUpdateOptions) (string, error) {
			mu.Lock()
			updatedIPs[target.IP] = true
			mu.Unlock()
			return `{"success": true, "updated": ["gitmap", "agm"], "details": "All apps up to date"}`, nil
		}

		err := RunFleetUpdateDispatch("update", args)
		if args[0] == "ua" {
			err = RunFleetUpdateDispatch("ua", args[1:])
		}
		if err != nil {
			t.Fatalf("args %v failed: %v", args, err)
		}

		if len(updatedIPs) != 2 {
			t.Errorf("args %v: expected 2 nodes updated, got %d", args, len(updatedIPs))
		}
		if updatedIPs["10.0.0.2"] {
			t.Errorf("args %v: expected worker-2 (10.0.0.2) to be excluded, but it was updated", args)
		}
		if !updatedIPs["10.0.0.1"] || !updatedIPs["10.0.0.3"] {
			t.Errorf("args %v: expected 10.0.0.1 and 10.0.0.3 to be updated", args)
		}
	}
}

func TestExecuteFleetUpdate_SingleAppWithExcept(t *testing.T) {
	origLoad := LoadFleetTargetsFn
	origExec := ExecuteRemoteUpdateFn
	origLive := CheckConnLivenessFn
	defer func() {
		LoadFleetTargetsFn = origLoad
		ExecuteRemoteUpdateFn = origExec
		CheckConnLivenessFn = origLive
	}()

	CheckConnLivenessFn = func(ctx context.Context, ip string, port int, timeout time.Duration) (bool, string) {
		return true, "online"
	}

	mockTargets := createMockFleetTargets()
	LoadFleetTargetsFn = func() ([]FleetTarget, error) {
		return mockTargets, nil
	}

	var mu sync.Mutex
	updatedTargets := make(map[string]string)

	ExecuteRemoteUpdateFn = func(target FleetTarget, opts FleetUpdateOptions) (string, error) {
		mu.Lock()
		updatedTargets[target.Alias] = opts.Pkg
		mu.Unlock()
		return `{"success": true, "current_version": "v20.11.0"}`, nil
	}

	err := ExecuteFleetUpdate([]string{"node", "--except", "node-1,10.0.0.3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(updatedTargets) != 1 {
		t.Fatalf("expected exactly 1 node updated, got %d", len(updatedTargets))
	}
	if pkg, ok := updatedTargets["worker-2"]; !ok || pkg != "node" {
		t.Errorf("expected worker-2 updated with package 'node', got pkg=%q, ok=%v", pkg, ok)
	}
}

func TestExecuteFleetUpdateLS_ParsesVariousJSONFormats(t *testing.T) {
	origLoad := LoadFleetTargetsFn
	origInv := ExecuteRemoteInventoryFn
	origLive := CheckConnLivenessFn
	defer func() {
		LoadFleetTargetsFn = origLoad
		ExecuteRemoteInventoryFn = origInv
		CheckConnLivenessFn = origLive
	}()

	CheckConnLivenessFn = func(ctx context.Context, ip string, port int, timeout time.Duration) (bool, string) {
		return true, "online"
	}

	mockTargets := createMockFleetTargets()
	LoadFleetTargetsFn = func() ([]FleetTarget, error) {
		return mockTargets, nil
	}

	ExecuteRemoteInventoryFn = func(target FleetTarget) (string, error) {
		switch target.ID {
		case "node-1":
			// Array format
			return `[
				{"name": "gitmap", "version": "v6.319.0", "status": "active"},
				{"name": "node", "version": "v20.11.0", "status": "active"}
			]`, nil
		case "node-2":
			// Key-value map format
			return `{"gitmap": "v6.319.0", "agm": "v2.1.0", "python": "3.11.8"}`, nil
		case "node-3":
			// Full object format with items
			return `{
				"node_id": "node-3",
				"alias": "worker-3",
				"items": [
					{"name": "gitmap", "version": "v6.319.0", "status": "active", "path": "C:\\gitmap.exe"}
				]
			}`, nil
		default:
			return "", errors.New("unknown target")
		}
	}

	err := ExecuteFleetUpdateLS([]string{"--except", "nonexistent"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFleetUpdateLS_HandlesErrorsAndMalformedJSON(t *testing.T) {
	origLoad := LoadFleetTargetsFn
	origInv := ExecuteRemoteInventoryFn
	origLive := CheckConnLivenessFn
	defer func() {
		LoadFleetTargetsFn = origLoad
		ExecuteRemoteInventoryFn = origInv
		CheckConnLivenessFn = origLive
	}()

	CheckConnLivenessFn = func(ctx context.Context, ip string, port int, timeout time.Duration) (bool, string) {
		return true, "online"
	}

	mockTargets := createMockFleetTargets()
	LoadFleetTargetsFn = func() ([]FleetTarget, error) {
		return mockTargets, nil
	}

	ExecuteRemoteInventoryFn = func(target FleetTarget) (string, error) {
		if target.ID == "node-2" {
			return "500 Internal Server Error: Daemon crashed", nil
		}
		if target.ID == "node-3" {
			return "", errors.New("ssh dial timeout")
		}
		return `[{"name": "gitmap", "version": "v6.319.0", "status": "active"}]`, nil
	}

	// Must handle errors without crashing or unhandled panic
	err := ExecuteFleetUpdateLS([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFleetUpdate_OfflineNodeReporting(t *testing.T) {
	origLoad := LoadFleetTargetsFn
	origExec := ExecuteRemoteUpdateFn
	origLive := CheckConnLivenessFn
	defer func() {
		LoadFleetTargetsFn = origLoad
		ExecuteRemoteUpdateFn = origExec
		CheckConnLivenessFn = origLive
	}()

	mockTargets := createMockFleetTargets()
	LoadFleetTargetsFn = func() ([]FleetTarget, error) {
		return mockTargets, nil
	}

	// Mock worker-2 (10.0.0.2) as offline
	CheckConnLivenessFn = func(ctx context.Context, ip string, port int, timeout time.Duration) (bool, string) {
		if ip == "10.0.0.2" {
			return false, "connection timed out"
		}
		return true, "online"
	}

	var mu sync.Mutex
	updatedIPs := make(map[string]bool)
	ExecuteRemoteUpdateFn = func(target FleetTarget, opts FleetUpdateOptions) (string, error) {
		mu.Lock()
		updatedIPs[target.IP] = true
		mu.Unlock()
		return `{"success": true, "details": "ok"}`, nil
	}

	err := RunFleetUpdateDispatch("update", []string{"all"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updatedIPs["10.0.0.2"] {
		t.Errorf("expected offline worker-2 (10.0.0.2) NOT to be updated via SSH")
	}
	if !updatedIPs["10.0.0.1"] || !updatedIPs["10.0.0.3"] {
		t.Errorf("expected online nodes 10.0.0.1 and 10.0.0.3 to be updated")
	}
}

func TestParseFleetUpdateTelemetry(t *testing.T) {
	target := FleetTarget{ID: "node-1", Alias: "worker-1", IP: "192.168.1.5"}

	t.Run("valid json object", func(t *testing.T) {
		raw := `{"success": true, "updated": ["gitmap"], "current_version": "v6.319.0", "details": "Success"}`
		res := ParseFleetUpdateTelemetry(raw, target, nil)
		if !res.Success {
			t.Errorf("expected success=true, got false")
		}
		if res.CurrentVersion != "v6.319.0" {
			t.Errorf("expected v6.319.0, got %s", res.CurrentVersion)
		}
	})

	t.Run("array format", func(t *testing.T) {
		raw := `[{"name": "gitmap", "status": "updated"}, {"name": "node", "status": "ok"}]`
		res := ParseFleetUpdateTelemetry(raw, target, nil)
		if !res.Success {
			t.Errorf("expected success=true for all ok/updated")
		}
		if len(res.Updated) != 2 {
			t.Errorf("expected 2 updated, got %d", len(res.Updated))
		}
	})

	t.Run("raw fallback string", func(t *testing.T) {
		raw := "All components successfully upgraded"
		res := ParseFleetUpdateTelemetry(raw, target, nil)
		if !res.Success {
			t.Errorf("expected success=true for nil error")
		}
		if res.Details != raw {
			t.Errorf("expected details %q, got %q", raw, res.Details)
		}
	})
}
