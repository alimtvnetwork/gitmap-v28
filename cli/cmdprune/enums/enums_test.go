package enums

import (
	"testing"
)

func TestPruneScopeType(t *testing.T) {
	testPruneScopeAllAndString(t)
	testPruneScopeParse(t)
	testPruneScopePredicates(t)
}

func testPruneScopeAllAndString(t *testing.T) {
	all := AllPruneScopeTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 scope variants, got %d", len(all))
	}

	for _, scope := range all {
		if !scope.IsValid() {
			t.Errorf("expected valid for %s", scope)
		}
		if scope.String() != string(scope) {
			t.Errorf("expected string representation to match %s", scope)
		}
		if !scope.HasScope() {
			t.Errorf("expected HasScope to be true for %s", scope)
		}
	}
}

func testPruneScopeParse(t *testing.T) {
	parsed, err := ParsePruneScopeType("  REMOTE ")
	if err != nil || parsed != PruneScopeRemote {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParsePruneScopeType("invalid"); err == nil {
		t.Error("expected error parsing invalid prune scope")
	}

	if PruneScopeType("bogus").IsValid() {
		t.Error("expected bogus scope to not be valid")
	}
}

func testPruneScopePredicates(t *testing.T) {
	if !PruneScopeLocal.IsLocal() || PruneScopeRemote.IsLocal() {
		t.Error("IsLocal predicate failure")
	}
	if !PruneScopeRemote.IsRemote() || PruneScopeLocal.IsRemote() {
		t.Error("IsRemote predicate failure")
	}
	if !PruneScopeAll.IsAll() || PruneScopeStale.IsAll() {
		t.Error("IsAll predicate failure")
	}
	if !PruneScopeStale.IsStale() || PruneScopeAll.IsStale() {
		t.Error("IsStale predicate failure")
	}
}

func TestPruneTargetType(t *testing.T) {
	testPruneTargetAllAndString(t)
	testPruneTargetParse(t)
	testPruneTargetPredicates(t)
}

func testPruneTargetAllAndString(t *testing.T) {
	all := AllPruneTargetTypes()
	if len(all) != 5 {
		t.Errorf("expected 5 target variants, got %d", len(all))
	}

	for _, target := range all {
		if !target.IsValid() {
			t.Errorf("expected valid for %s", target)
		}
		if target.String() != string(target) {
			t.Errorf("expected string representation to match %s", target)
		}
		if !target.HasTarget() {
			t.Errorf("expected HasTarget to be true for %s", target)
		}
	}
}

func testPruneTargetParse(t *testing.T) {
	parsed, err := ParsePruneTargetType("  WORKTREES ")
	if err != nil || parsed != PruneTargetWorktrees {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParsePruneTargetType("invalid"); err == nil {
		t.Error("expected error parsing invalid prune target")
	}

	if PruneTargetType("bogus").IsValid() {
		t.Error("expected bogus target to not be valid")
	}
}

func testPruneTargetPredicates(t *testing.T) {
	if !PruneTargetBranches.IsBranches() || PruneTargetTags.IsBranches() {
		t.Error("IsBranches predicate failure")
	}
	if !PruneTargetTags.IsTags() || PruneTargetBranches.IsTags() {
		t.Error("IsTags predicate failure")
	}
	if !PruneTargetRemotes.IsRemotes() || PruneTargetObjects.IsRemotes() {
		t.Error("IsRemotes predicate failure")
	}
	if !PruneTargetObjects.IsObjects() || PruneTargetRemotes.IsObjects() {
		t.Error("IsObjects predicate failure")
	}
	if !PruneTargetWorktrees.IsWorktrees() || PruneTargetBranches.IsWorktrees() {
		t.Error("IsWorktrees predicate failure")
	}
}

func TestPruneActionType(t *testing.T) {
	testPruneActionAllAndString(t)
	testPruneActionParse(t)
	testPruneActionPredicates(t)
}

func testPruneActionAllAndString(t *testing.T) {
	all := AllPruneActionTypes()
	if len(all) != 5 {
		t.Errorf("expected 5 action variants, got %d", len(all))
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
}

func testPruneActionParse(t *testing.T) {
	parsed, err := ParsePruneActionType("  DRY-RUN ")
	if err != nil || parsed != PruneActionDryRun {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParsePruneActionType("invalid"); err == nil {
		t.Error("expected error parsing invalid prune action")
	}

	if PruneActionType("bogus").IsValid() {
		t.Error("expected bogus action to not be valid")
	}
}

func testPruneActionPredicates(t *testing.T) {
	if !PruneActionDelete.IsDelete() || PruneActionDryRun.IsDelete() {
		t.Error("IsDelete predicate failure")
	}
	if !PruneActionDryRun.IsDryRun() || PruneActionPrompt.IsDryRun() {
		t.Error("IsDryRun predicate failure")
	}
	if !PruneActionPrompt.IsPrompt() || PruneActionKeep.IsPrompt() {
		t.Error("IsPrompt predicate failure")
	}
	if !PruneActionKeep.IsKeep() || PruneActionSkip.IsKeep() {
		t.Error("IsKeep predicate failure")
	}
	if !PruneActionSkip.IsSkip() || PruneActionDelete.IsSkip() {
		t.Error("IsSkip predicate failure")
	}
}

func TestPruneStatusType(t *testing.T) {
	testPruneStatusAllAndString(t)
	testPruneStatusParse(t)
	testPruneStatusPredicates(t)
}

func testPruneStatusAllAndString(t *testing.T) {
	all := AllPruneStatusTypes()
	if len(all) != 5 {
		t.Errorf("expected 5 status variants, got %d", len(all))
	}

	for _, status := range all {
		if !status.IsValid() {
			t.Errorf("expected valid for %s", status)
		}
		if status.String() != string(status) {
			t.Errorf("expected string representation to match %s", status)
		}
		if !status.HasStatus() {
			t.Errorf("expected HasStatus to be true for %s", status)
		}
	}
}

func testPruneStatusParse(t *testing.T) {
	parsed, err := ParsePruneStatusType("  ABORTED ")
	if err != nil || parsed != PruneStatusAborted {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParsePruneStatusType("invalid"); err == nil {
		t.Error("expected error parsing invalid prune status")
	}

	if PruneStatusType("bogus").IsValid() {
		t.Error("expected bogus status to not be valid")
	}
}

func testPruneStatusPredicates(t *testing.T) {
	if !PruneStatusDeleted.IsDeleted() || PruneStatusKept.IsDeleted() {
		t.Error("IsDeleted predicate failure")
	}
	if !PruneStatusKept.IsKept() || PruneStatusSkipped.IsKept() {
		t.Error("IsKept predicate failure")
	}
	if !PruneStatusSkipped.IsSkipped() || PruneStatusFailed.IsSkipped() {
		t.Error("IsSkipped predicate failure")
	}
	if !PruneStatusFailed.IsFailed() || PruneStatusAborted.IsFailed() {
		t.Error("IsFailed predicate failure")
	}
	if !PruneStatusAborted.IsAborted() || PruneStatusDeleted.IsAborted() {
		t.Error("IsAborted predicate failure")
	}
}
