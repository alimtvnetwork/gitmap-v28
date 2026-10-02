package cmdupdate

import (
	"archive/zip"
	"bytes"
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

func TestParseFleetUpdateOptions_ZipAndIncludeOthers(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantZip    bool
		wantOthers bool
		wantPkg    string
		wantAll    bool
	}{
		{
			name:       "flag --zip and --include-others",
			args:       []string{"all", "--zip", "--include-others"},
			wantZip:    true,
			wantOthers: true,
			wantPkg:    "all",
			wantAll:    true,
		},
		{
			name:       "positional zip and --include-other",
			args:       []string{"all", "zip", "--include-other"},
			wantZip:    true,
			wantOthers: true,
			wantPkg:    "all",
			wantAll:    true,
		},
		{
			name:       "agm target with zip",
			args:       []string{"agm", "--zip", "--include-others"},
			wantZip:    true,
			wantOthers: true,
			wantPkg:    "agm",
			wantAll:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := parseFleetUpdateOptions(tc.args)
			if opts.IsZip != tc.wantZip {
				t.Errorf("IsZip: got %v, want %v", opts.IsZip, tc.wantZip)
			}
			if opts.IncludeOthers != tc.wantOthers {
				t.Errorf("IncludeOthers: got %v, want %v", opts.IncludeOthers, tc.wantOthers)
			}
			if opts.Pkg != tc.wantPkg {
				t.Errorf("Pkg: got %s, want %s", opts.Pkg, tc.wantPkg)
			}
			if opts.IsAll != tc.wantAll {
				t.Errorf("IsAll: got %v, want %v", opts.IsAll, tc.wantAll)
			}
		})
	}
}

func TestIsFleetUpdateCommand_ZipAndAliases(t *testing.T) {
	cmdAliases := []string{"uaz", "update-all-zip", "updateallzip"}
	for _, cmd := range cmdAliases {
		if !IsFleetUpdateCommand(cmd, nil) {
			t.Errorf("expected IsFleetUpdateCommand(%q) to be true", cmd)
		}
	}

	if !IsFleetUpdateCommand("update", []string{"all", "zip"}) {
		t.Errorf("expected IsFleetUpdateCommand('update', ['all', 'zip']) to be true")
	}
	if !IsFleetUpdateCommand("update", []string{"--zip"}) {
		t.Errorf("expected IsFleetUpdateCommand('update', ['--zip']) to be true")
	}
	if !IsFleetUpdateCommand("update", []string{"--include-others"}) {
		t.Errorf("expected IsFleetUpdateCommand('update', ['--include-others']) to be true")
	}
}

func TestExecuteFleetUpdate_ZipDistribution(t *testing.T) {
	origLoad := LoadFleetTargetsFn
	origExec := ExecuteFleetZipUpdateFn
	origZip := CreateUpdateZipFn
	origLive := CheckConnLivenessFn
	defer func() {
		LoadFleetTargetsFn = origLoad
		ExecuteFleetZipUpdateFn = origExec
		CreateUpdateZipFn = origZip
		CheckConnLivenessFn = origLive
		clearZipCache()
	}()

	CheckConnLivenessFn = func(ctx context.Context, ip string, port int, timeout time.Duration) (bool, string) {
		return true, "online"
	}
	mockTargets := createMockFleetTargets()
	LoadFleetTargetsFn = func() ([]FleetTarget, error) {
		return mockTargets, nil
	}

	var mu sync.Mutex
	zipUpdatedNodes := make(map[string]bool)
	ExecuteFleetZipUpdateFn = func(target FleetTarget, opts FleetUpdateOptions) (string, error) {
		mu.Lock()
		zipUpdatedNodes[target.IP] = true
		mu.Unlock()
		if !opts.IsZip {
			return `{"success": false, "details": "not zip"}`, nil
		}
		return `{"success": true, "current_version": "v6.320.0", "details": "Zip installed: v6.320.0"}`, nil
	}

	CreateUpdateZipFn = func(pkg, osType string) ([]byte, error) {
		return []byte("mock-zip-bytes"), nil
	}

	testCommands := []struct {
		cmd  string
		args []string
	}{
		{cmd: "uaz", args: []string{}},
		{cmd: "update-all-zip", args: []string{}},
		{cmd: "update", args: []string{"all", "zip"}},
	}

	for _, tc := range testCommands {
		mu.Lock()
		zipUpdatedNodes = make(map[string]bool)
		mu.Unlock()

		err := RunFleetUpdateDispatch(tc.cmd, tc.args)
		if err != nil {
			t.Fatalf("command %s %v failed: %v", tc.cmd, tc.args, err)
		}
		if len(zipUpdatedNodes) != 3 {
			t.Errorf("command %s %v: expected 3 nodes updated via zip, got %d", tc.cmd, tc.args, len(zipUpdatedNodes))
		}
	}
}

func TestExecuteFleetUpdate_IncludeOthers(t *testing.T) {
	origLoad := LoadFleetTargetsFn
	origCluster := LoadClusterTargetsFn
	origExec := ExecuteRemoteUpdateFn
	origLive := CheckConnLivenessFn
	defer func() {
		LoadFleetTargetsFn = origLoad
		LoadClusterTargetsFn = origCluster
		ExecuteRemoteUpdateFn = origExec
		CheckConnLivenessFn = origLive
	}()

	CheckConnLivenessFn = func(ctx context.Context, ip string, port int, timeout time.Duration) (bool, string) {
		return true, "online"
	}
	LoadFleetTargetsFn = func() ([]FleetTarget, error) {
		return []FleetTarget{
			{ID: "node-1", Alias: "worker-1", IP: "10.0.0.1", OS: "linux"},
		}, nil
	}
	LoadClusterTargetsFn = func() ([]FleetTarget, error) {
		return []FleetTarget{
			{ID: "cluster-1", Alias: "cluster-worker", IP: "10.0.0.99", OS: "windows"},
			{ID: "node-1-dup", Alias: "worker-1-dup", IP: "10.0.0.1", OS: "linux"},
		}, nil
	}

	var mu sync.Mutex
	updatedIPs := make(map[string]bool)
	ExecuteRemoteUpdateFn = func(target FleetTarget, opts FleetUpdateOptions) (string, error) {
		mu.Lock()
		updatedIPs[target.IP] = true
		mu.Unlock()
		return `{"success": true, "details": "ok"}`, nil
	}

	err := ExecuteFleetUpdate([]string{"all"})
	if err != nil {
		t.Fatalf("ExecuteFleetUpdate failed: %v", err)
	}
	if len(updatedIPs) != 1 || !updatedIPs["10.0.0.1"] {
		t.Errorf("expected only 10.0.0.1 without --include-others, got %v", updatedIPs)
	}

	mu.Lock()
	updatedIPs = make(map[string]bool)
	mu.Unlock()

	err = ExecuteFleetUpdate([]string{"all", "--include-others"})
	if err != nil {
		t.Fatalf("ExecuteFleetUpdate with --include-others failed: %v", err)
	}
	if len(updatedIPs) != 2 {
		t.Errorf("expected 2 unique nodes with --include-others, got %d (%v)", len(updatedIPs), updatedIPs)
	}
	if !updatedIPs["10.0.0.1"] || !updatedIPs["10.0.0.99"] {
		t.Errorf("expected 10.0.0.1 and 10.0.0.99 updated, got %v", updatedIPs)
	}
}

func TestCreateUpdatePackageZip(t *testing.T) {
	data, err := createUpdatePackageZip("gitmap", "windows")
	if err != nil {
		t.Fatalf("createUpdatePackageZip failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected non-empty zip bytes")
	}

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("failed to read zip: %v", err)
	}
	foundBin := false
	foundScript := false
	for _, f := range r.File {
		if f.Name == "gitmap.exe" {
			foundBin = true
		}
		if f.Name == "install_remote.ps1" {
			foundScript = true
		}
	}
	if !foundBin {
		t.Errorf("expected gitmap.exe in zip")
	}
	if !foundScript {
		t.Errorf("expected install_remote.ps1 in zip")
	}
}
