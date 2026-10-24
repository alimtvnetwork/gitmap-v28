package enums

import (
	"testing"
)

func TestSyncTargetType(t *testing.T) {
	allTargets := AllSyncTargetTypes()
	if len(allTargets) != 5 {
		t.Errorf("expected 5 targets, got %d", len(allTargets))
	}

	for _, target := range allTargets {
		if !target.IsValid() {
			t.Errorf("expected %s to be valid", target)
		}
		if target.String() != string(target) {
			t.Errorf("expected %s, got %s", string(target), target.String())
		}
		if !target.HasTargetScope() {
			t.Errorf("expected %s to have target scope", target)
		}
	}

	if SyncTargetType("invalid").IsValid() {
		t.Errorf("expected invalid target to report false")
	}

	if SyncTargetType("").HasTargetScope() {
		t.Errorf("expected empty target to not have target scope")
	}

	if !SyncTargetAll.IsAll() || SyncTargetDirty.IsAll() {
		t.Errorf("unexpected IsAll result")
	}
	if !SyncTargetDirty.IsDirty() || SyncTargetClean.IsDirty() {
		t.Errorf("unexpected IsDirty result")
	}
	if !SyncTargetClean.IsClean() || SyncTargetDirty.IsClean() {
		t.Errorf("unexpected IsClean result")
	}
	if !SyncTargetStashed.IsStashed() || SyncTargetDirty.IsStashed() {
		t.Errorf("unexpected IsStashed result")
	}
	if !SyncTargetMissing.IsMissing() || SyncTargetDirty.IsMissing() {
		t.Errorf("unexpected IsMissing result")
	}

	parsed, err := ParseSyncTargetType("  DIRTY ")
	if err != nil || parsed != SyncTargetDirty {
		t.Errorf("failed to parse dirty: %v, %s", err, parsed)
	}

	_, errInvalid := ParseSyncTargetType("unknown")
	if errInvalid == nil {
		t.Errorf("expected error parsing unknown target")
	}
}

func TestSyncRepoStatusType(t *testing.T) {
	allStatuses := AllSyncRepoStatusTypes()
	if len(allStatuses) != 5 {
		t.Errorf("expected 5 statuses, got %d", len(allStatuses))
	}

	for _, s := range allStatuses {
		if !s.IsValid() {
			t.Errorf("expected %s to be valid", s)
		}
		if s.String() != string(s) {
			t.Errorf("expected %s, got %s", string(s), s.String())
		}
	}

	if SyncRepoStatusType("invalid").IsValid() {
		t.Errorf("expected invalid status to report false")
	}

	if !SyncRepoStatusSynced.IsSynced() || SyncRepoStatusModified.IsSynced() {
		t.Errorf("unexpected IsSynced result")
	}
	if !SyncRepoStatusModified.IsModified() || SyncRepoStatusSynced.IsModified() {
		t.Errorf("unexpected IsModified result")
	}
	if !SyncRepoStatusConflict.IsConflict() || !SyncRepoStatusConflict.HasConflict() {
		t.Errorf("unexpected conflict result")
	}
	if SyncRepoStatusSynced.HasConflict() {
		t.Errorf("unexpected HasConflict result for synced")
	}
	if !SyncRepoStatusSkipped.IsSkipped() || SyncRepoStatusSynced.IsSkipped() {
		t.Errorf("unexpected IsSkipped result")
	}
	if !SyncRepoStatusError.IsError() || !SyncRepoStatusError.HasError() {
		t.Errorf("unexpected error status result")
	}
	if SyncRepoStatusSynced.HasError() {
		t.Errorf("unexpected HasError result for synced")
	}

	parsed, err := ParseSyncRepoStatusType("  CONFLICT ")
	if err != nil || parsed != SyncRepoStatusConflict {
		t.Errorf("failed to parse conflict: %v, %s", err, parsed)
	}

	_, errInvalid := ParseSyncRepoStatusType("unknown")
	if errInvalid == nil {
		t.Errorf("expected error parsing unknown status")
	}
}

func TestSyncActionType(t *testing.T) {
	allActions := AllSyncActionTypes()
	if len(allActions) != 5 {
		t.Errorf("expected 5 actions, got %d", len(allActions))
	}

	for _, a := range allActions {
		if !a.IsValid() {
			t.Errorf("expected %s to be valid", a)
		}
		if a.String() != string(a) {
			t.Errorf("expected %s, got %s", string(a), a.String())
		}
	}

	if SyncActionType("invalid").IsValid() {
		t.Errorf("expected invalid action to report false")
	}

	if !SyncActionFetch.IsFetch() || SyncActionPull.IsFetch() {
		t.Errorf("unexpected IsFetch result")
	}
	if !SyncActionPull.IsPull() || SyncActionFetch.IsPull() {
		t.Errorf("unexpected IsPull result")
	}
	if !SyncActionMerge.IsMerge() || SyncActionFetch.IsMerge() {
		t.Errorf("unexpected IsMerge result")
	}
	if !SyncActionRebase.IsRebase() || !SyncActionRebase.HasRebase() {
		t.Errorf("unexpected rebase result")
	}
	if SyncActionFetch.HasRebase() {
		t.Errorf("unexpected HasRebase result for fetch")
	}
	if !SyncActionCheckout.IsCheckout() || SyncActionFetch.IsCheckout() {
		t.Errorf("unexpected IsCheckout result")
	}

	parsed, err := ParseSyncActionType("  REBASE ")
	if err != nil || parsed != SyncActionRebase {
		t.Errorf("failed to parse rebase: %v, %s", err, parsed)
	}

	_, errInvalid := ParseSyncActionType("unknown")
	if errInvalid == nil {
		t.Errorf("expected error parsing unknown action")
	}
}

func TestSyncModeType(t *testing.T) {
	allModes := AllSyncModeTypes()
	if len(allModes) != 3 {
		t.Errorf("expected 3 modes, got %d", len(allModes))
	}

	for _, m := range allModes {
		if !m.IsValid() {
			t.Errorf("expected %s to be valid", m)
		}
		if m.String() != string(m) {
			t.Errorf("expected %s, got %s", string(m), m.String())
		}
	}

	if SyncModeType("invalid").IsValid() {
		t.Errorf("expected invalid mode to report false")
	}

	if !SyncModeFastForward.IsFastForward() || SyncModeRebase.IsFastForward() {
		t.Errorf("unexpected IsFastForward result")
	}
	if !SyncModeRebase.IsRebase() || SyncModeFastForward.IsRebase() {
		t.Errorf("unexpected IsRebase result")
	}
	if !SyncModeDryRun.IsDryRun() || SyncModeFastForward.IsDryRun() {
		t.Errorf("unexpected IsDryRun result")
	}
	if !SyncModeFastForward.HasMutation() || !SyncModeRebase.HasMutation() {
		t.Errorf("expected mutation for fast-forward and rebase")
	}
	if SyncModeDryRun.HasMutation() {
		t.Errorf("expected no mutation for dry-run")
	}

	parsed, err := ParseSyncModeType("  DRY-RUN ")
	if err != nil || parsed != SyncModeDryRun {
		t.Errorf("failed to parse dry-run: %v, %s", err, parsed)
	}

	_, errInvalid := ParseSyncModeType("unknown")
	if errInvalid == nil {
		t.Errorf("expected error parsing unknown mode")
	}
}
