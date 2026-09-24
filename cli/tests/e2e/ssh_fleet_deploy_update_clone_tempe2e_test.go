//go:build tempe2e

package e2e

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func TestMacroAndPeatDeploySSH_TempE2E(t *testing.T) {
	if isTempE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}
	restoreTargets := mockMacroFleetTargets()
	t.Cleanup(restoreTargets)

	errMacro := cmdmacro.ExecuteMacroDeploySSH([]string{"deploy", "ssh", "--except", "w3,192.168.1.12", "--dry-run"})
	if errMacro != nil {
		t.Fatalf("expected macro deploy ssh to succeed, got: %v", errMacro)
	}

	errPeat := cmdmacro.ExecuteMacroDeploySSH([]string{"peat", "deploy", "ssh", "--excep", "w2", "--dry-run"})
	if errPeat != nil {
		t.Fatalf("expected peat deploy ssh to succeed, got: %v", errPeat)
	}

	errPea := cmdmacro.ExecuteMacroDeploySSH([]string{"pea", "deploy", "ssh", "--except", "1", "--dry-run"})
	if errPea != nil {
		t.Fatalf("expected pea deploy ssh to succeed, got: %v", errPea)
	}
}

func TestFleetUpdateAllAndTargeted_TempE2E(t *testing.T) {
	if isTempE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}
	restoreFleet := mockUpdateFleetProviders()
	t.Cleanup(restoreFleet)

	errAll := cmdupdate.RunFleetUpdateDispatch("ua", []string{"--except", "w3"})
	if errAll != nil {
		t.Fatalf("expected ua fleet update to succeed, got: %v", errAll)
	}

	errNamed := cmdupdate.RunFleetUpdateDispatch("update", []string{"gitmap", "--excep", "w2,192.168.1.7"})
	if errNamed != nil {
		t.Fatalf("expected targeted update with --excep to succeed, got: %v", errNamed)
	}
}

func TestFleetUpdateLS_JSONTable_TempE2E(t *testing.T) {
	if isTempE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}
	restoreFleet := mockUpdateFleetProviders()
	t.Cleanup(restoreFleet)

	errLS := cmdupdate.RunFleetUpdateDispatch("update", []string{"ls", "--except", "w3"})
	if errLS != nil {
		t.Fatalf("expected update ls across fleet to succeed, got: %v", errLS)
	}
}

func isTempE2ESkipped() bool {
	return os.Getenv("RUN_TEMP_E2E") != "1"
}

func mockMacroFleetTargets() func() {
	origLoad := cmdmacro.LoadClusterTargetsFn
	origTransfer := cmdmacro.TransferMacrosToTargetFn
	cmdmacro.LoadClusterTargetsFn = func() ([]cmdmacro.MacroDeployTarget, error) {
		return []cmdmacro.MacroDeployTarget{
			{ID: "1", Alias: "w1", IP: "192.168.1.3", Username: "administrator", Port: 22, OS: "windows"},
			{ID: "2", Alias: "w2", IP: "192.168.1.7", Username: "administrator", Port: 22, OS: "windows"},
			{ID: "3", Alias: "w3", IP: "192.168.1.12", Username: "administrator", Port: 22, OS: "linux"},
		}, nil
	}
	cmdmacro.TransferMacrosToTargetFn = func(t cmdmacro.MacroDeployTarget, m []macro.Macro, o cmdmacro.MacroDeployOptions) error {
		return nil
	}
	return func() {
		cmdmacro.LoadClusterTargetsFn = origLoad
		cmdmacro.TransferMacrosToTargetFn = origTransfer
	}
}

func mockUpdateFleetProviders() func() {
	origLoad := cmdupdate.LoadFleetTargetsFn
	origExec := cmdupdate.ExecuteRemoteUpdateFn
	origInv := cmdupdate.ExecuteRemoteInventoryFn
	cmdupdate.LoadFleetTargetsFn = buildMockUpdateTargets
	cmdupdate.ExecuteRemoteUpdateFn = buildMockRemoteUpdateExec
	cmdupdate.ExecuteRemoteInventoryFn = buildMockRemoteInventoryExec
	return func() {
		cmdupdate.LoadFleetTargetsFn = origLoad
		cmdupdate.ExecuteRemoteUpdateFn = origExec
		cmdupdate.ExecuteRemoteInventoryFn = origInv
	}
}

func buildMockUpdateTargets() ([]cmdupdate.FleetTarget, error) {
	return []cmdupdate.FleetTarget{
		{ID: "1", Alias: "w1", IP: "192.168.1.3", Username: "administrator", Port: 22, OS: "windows"},
		{ID: "2", Alias: "w2", IP: "192.168.1.7", Username: "administrator", Port: 22, OS: "windows"},
		{ID: "3", Alias: "w3", IP: "192.168.1.12", Username: "administrator", Port: 22, OS: "linux"},
	}, nil
}

func buildMockRemoteUpdateExec(t cmdupdate.FleetTarget, opts cmdupdate.FleetUpdateOptions) (string, error) {
	tel := cmdupdate.FleetUpdateTelemetry{
		Alias: t.Alias, IP: t.IP, Success: true, Updated: []string{opts.Pkg}, CurrentVersion: "v6.324.0",
	}
	raw, err := json.Marshal(tel)
	return string(raw), err
}

func buildMockRemoteInventoryExec(t cmdupdate.FleetTarget) (string, error) {
	inv := cmdupdate.FleetNodeInventory{
		Alias: t.Alias, IP: t.IP, OS: t.OS,
		Items: []cmdupdate.FleetInventoryItem{{Name: "gitmap", Version: "v6.324.0", Status: "installed"}},
	}
	raw, err := json.Marshal(inv)
	return string(raw), err
}
