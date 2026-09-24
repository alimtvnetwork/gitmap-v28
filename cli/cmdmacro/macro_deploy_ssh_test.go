package cmdmacro

import (
	"errors"
	"sync"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func createMockTargets() []MacroDeployTarget {
	return []MacroDeployTarget{
		{
			ID:       "node-1",
			Alias:    "worker-alpha",
			IP:       "192.168.1.10",
			Username: "root",
			Port:     22,
			OS:       "linux",
		},
		{
			ID:       "node-2",
			Alias:    "worker-beta",
			IP:       "192.168.1.20",
			Username: "admin",
			Port:     22,
			OS:       "linux",
		},
		{
			ID:       "node-3",
			Alias:    "worker-gamma",
			IP:       "192.168.1.30",
			Username: "root",
			Port:     22,
			OS:       "windows",
		},
	}
}

func TestExecuteMacroDeploySSH_AllReachableNodes(t *testing.T) {
	origLoad := LoadClusterTargetsFn
	origTransfer := TransferMacrosToTargetFn
	defer func() {
		LoadClusterTargetsFn = origLoad
		TransferMacrosToTargetFn = origTransfer
	}()

	mockTargets := createMockTargets()
	LoadClusterTargetsFn = func() ([]MacroDeployTarget, error) {
		return mockTargets, nil
	}

	var mu sync.Mutex
	deployedIPs := make(map[string]bool)

	TransferMacrosToTargetFn = func(target MacroDeployTarget, macros []macro.Macro, opts MacroDeployOptions) error {
		mu.Lock()
		deployedIPs[target.IP] = true
		mu.Unlock()
		return nil
	}

	err := ExecuteMacroDeploySSH([]string{"deploy", "ssh"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(deployedIPs) != 3 {
		t.Errorf("expected 3 nodes deployed, got %d", len(deployedIPs))
	}
	if !deployedIPs["192.168.1.10"] || !deployedIPs["192.168.1.20"] || !deployedIPs["192.168.1.30"] {
		t.Errorf("expected all 3 mock IPs deployed, got %+v", deployedIPs)
	}
}

func TestExecuteMacroDeploySSH_ExceptExclusionFilters(t *testing.T) {
	origLoad := LoadClusterTargetsFn
	origTransfer := TransferMacrosToTargetFn
	defer func() {
		LoadClusterTargetsFn = origLoad
		TransferMacrosToTargetFn = origTransfer
	}()

	mockTargets := createMockTargets()
	LoadClusterTargetsFn = func() ([]MacroDeployTarget, error) {
		return mockTargets, nil
	}

	testCases := []struct {
		name         string
		args         []string
		expectedLeft []string
	}{
		{
			name:         "exclude by alias",
			args:         []string{"deploy", "ssh", "--except", "worker-beta"},
			expectedLeft: []string{"192.168.1.10", "192.168.1.30"},
		},
		{
			name:         "exclude by ip",
			args:         []string{"deploy", "ssh", "--except", "192.168.1.10"},
			expectedLeft: []string{"192.168.1.20", "192.168.1.30"},
		},
		{
			name:         "exclude by id",
			args:         []string{"deploy", "ssh", "--except", "node-3"},
			expectedLeft: []string{"192.168.1.10", "192.168.1.20"},
		},
		{
			name:         "exclude multiple comma separated",
			args:         []string{"deploy", "ssh", "--except", "worker-alpha,192.168.1.30"},
			expectedLeft: []string{"192.168.1.20"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			deployed := make(map[string]bool)

			TransferMacrosToTargetFn = func(target MacroDeployTarget, macros []macro.Macro, opts MacroDeployOptions) error {
				mu.Lock()
				deployed[target.IP] = true
				mu.Unlock()
				return nil
			}

			err := ExecuteMacroDeploySSH(tc.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(deployed) != len(tc.expectedLeft) {
				t.Errorf("expected %d nodes deployed, got %d (%+v)", len(tc.expectedLeft), len(deployed), deployed)
			}
			for _, expectedIP := range tc.expectedLeft {
				if !deployed[expectedIP] {
					t.Errorf("expected IP %s deployed, but was not in %+v", expectedIP, deployed)
				}
			}
		})
	}
}

func TestExecuteMacroDeploySSH_AliasesExecuteIdentically(t *testing.T) {
	origLoad := LoadClusterTargetsFn
	origTransfer := TransferMacrosToTargetFn
	defer func() {
		LoadClusterTargetsFn = origLoad
		TransferMacrosToTargetFn = origTransfer
	}()

	mockTargets := createMockTargets()
	LoadClusterTargetsFn = func() ([]MacroDeployTarget, error) {
		return mockTargets, nil
	}

	aliasInvocations := [][]string{
		{"macro", "deploy", "ssh", "--except", "worker-alpha"},
		{"peat", "deploy", "ssh", "--except", "worker-alpha"},
		{"pea", "deploy", "ssh", "--except", "worker-alpha"},
	}

	for _, args := range aliasInvocations {
		var mu sync.Mutex
		deployed := make(map[string]bool)

		TransferMacrosToTargetFn = func(target MacroDeployTarget, macros []macro.Macro, opts MacroDeployOptions) error {
			mu.Lock()
			deployed[target.IP] = true
			mu.Unlock()
			return nil
		}

		err := ExecuteMacroDeploySSH(args)
		if err != nil {
			t.Fatalf("invocation with %v returned unexpected error: %v", args, err)
		}

		if len(deployed) != 2 {
			t.Errorf("invocation %v: expected 2 nodes deployed, got %d", args, len(deployed))
		}
		if deployed["192.168.1.10"] {
			t.Errorf("invocation %v: expected worker-alpha (192.168.1.10) to be excluded, but it was deployed", args)
		}
		if !deployed["192.168.1.20"] || !deployed["192.168.1.30"] {
			t.Errorf("invocation %v: expected 192.168.1.20 and 192.168.1.30 deployed", args)
		}
	}
}

func TestExecuteMacroDeploySSH_PartialFailure(t *testing.T) {
	origLoad := LoadClusterTargetsFn
	origTransfer := TransferMacrosToTargetFn
	defer func() {
		LoadClusterTargetsFn = origLoad
		TransferMacrosToTargetFn = origTransfer
	}()

	mockTargets := createMockTargets()
	LoadClusterTargetsFn = func() ([]MacroDeployTarget, error) {
		return mockTargets, nil
	}

	TransferMacrosToTargetFn = func(target MacroDeployTarget, macros []macro.Macro, opts MacroDeployOptions) error {
		if target.IP == "192.168.1.20" {
			return errors.New("connection timed out")
		}
		return nil
	}

	// Should record failure in summary table without crashing
	err := ExecuteMacroDeploySSH([]string{"deploy", "ssh"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
