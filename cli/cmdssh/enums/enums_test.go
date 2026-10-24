package enums

import (
	"testing"
)

func TestDeploySyncModeType(t *testing.T) {
	all := AllDeploySyncModeTypes()
	if len(all) != 6 {
		t.Errorf("expected 6 deploy sync modes, got %d", len(all))
	}

	for _, mode := range all {
		if !mode.IsValid() {
			t.Errorf("expected valid for %s", mode)
		}
		if mode.String() != string(mode) {
			t.Errorf("expected string representation to match %s", mode)
		}
	}

	if DeploySyncModeType("bogus").IsValid() {
		t.Error("expected invalid sync mode to not be valid")
	}

	parsed, err := ParseDeploySyncModeType("  SYNC-RIGHT ")
	if err != nil || parsed != SyncModeSyncRight {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseDeploySyncModeType("unknown"); err == nil {
		t.Error("expected error parsing unknown sync mode")
	}

	if !SyncModeNone.IsNone() || !SyncModeOverwrite.IsOverwrite() || !SyncModeSkip.IsSkip() {
		t.Error("sync mode predicate failure for none/overwrite/skip")
	}
	if !SyncModeSync.IsSync() || !SyncModeSyncRight.IsSyncRight() || !SyncModeSyncLeft.IsSyncLeft() {
		t.Error("sync mode predicate failure for sync/sync-right/sync-left")
	}
	if SyncModeNone.HasSyncMode() || !SyncModeSync.HasSyncMode() {
		t.Error("HasSyncMode should be false for none and true for sync")
	}
}

func TestDeployActionType(t *testing.T) {
	all := AllDeployActionTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 deploy actions, got %d", len(all))
	}

	for _, action := range all {
		if !action.IsValid() {
			t.Errorf("expected valid for %s", action)
		}
		if action.String() != string(action) {
			t.Errorf("expected string representation to match %s", action)
		}
		if !action.HasAction() {
			t.Errorf("expected HasAction to be true for %s", action)
		}
	}

	if DeployActionType("bogus").IsValid() {
		t.Error("expected invalid deploy action to not be valid")
	}

	parsed, err := ParseDeployActionType("  TRANSFER-TO-LOCAL ")
	if err != nil || parsed != ActionTransferToLocal {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseDeployActionType("unknown"); err == nil {
		t.Error("expected error parsing unknown deploy action")
	}

	if !ActionTransferToRemote.IsTransferToRemote() || !ActionTransferToLocal.IsTransferToLocal() {
		t.Error("deploy action predicate failure for transfer variants")
	}
	if !ActionSkip.IsSkip() || !ActionConflictPrompt.IsConflictPrompt() {
		t.Error("deploy action predicate failure for skip/conflict-prompt")
	}
}

func TestSSHDiagnosisType(t *testing.T) {
	all := AllSSHDiagnosisTypes()
	if len(all) != 5 {
		t.Errorf("expected 5 ssh diagnoses, got %d", len(all))
	}

	for _, diag := range all {
		if !diag.IsValid() {
			t.Errorf("expected valid for %s", diag)
		}
		if diag.String() != string(diag) {
			t.Errorf("expected string representation to match %s", diag)
		}
		if !diag.HasDiagnosis() {
			t.Errorf("expected HasDiagnosis to be true for %s", diag)
		}
	}

	if SSHDiagnosisType("bogus").IsValid() {
		t.Error("expected invalid diagnosis to not be valid")
	}

	parsed, err := ParseSSHDiagnosisType("  ssh port open & accessible ")
	if err != nil || parsed != DiagSSHPortOpen {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseSSHDiagnosisType("unknown"); err == nil {
		t.Error("expected error parsing unknown diagnosis")
	}

	if !DiagSSHPortOpen.IsOpen() || !DiagFirewallInboundDrop.IsFirewallInboundDrop() {
		t.Error("diagnosis predicate failure for open/firewall")
	}
	if !DiagDaemonNotRunning.IsDaemonNotRunning() || !DiagHostOffline.IsHostOffline() {
		t.Error("diagnosis predicate failure for daemon/host-offline")
	}
	if !DiagUnknownError.IsUnknownError() {
		t.Error("diagnosis predicate failure for unknown error")
	}
}
