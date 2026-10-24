package enums

import (
	"testing"
)

func TestGitHookType(t *testing.T) {
	testGitHookAllAndString(t)
	testGitHookParse(t)
	testGitHookPredicates(t)
}

func testGitHookAllAndString(t *testing.T) {
	all := AllGitHookTypes()
	if len(all) != 8 {
		t.Errorf("expected 8 git hook types, got %d", len(all))
	}

	for _, hook := range all {
		if !hook.IsValid() {
			t.Errorf("expected valid for %s", hook)
		}
		if hook.String() != string(hook) {
			t.Errorf("expected string representation to match %s", hook)
		}
		if !hook.HasHook() {
			t.Errorf("expected HasHook to be true for %s", hook)
		}
	}
}

func testGitHookParse(t *testing.T) {
	parsed, err := ParseGitHookType("  COMMIT-MSG ")
	if err != nil || parsed != GitHookCommitMsg {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseGitHookType("invalid-hook"); err == nil {
		t.Error("expected error parsing invalid hook")
	}

	if GitHookType("unknown").IsValid() {
		t.Error("expected unknown hook to not be valid")
	}
}

func testGitHookPredicates(t *testing.T) {
	if !GitHookCommitMsg.IsCommitMsg() || GitHookPreCommit.IsCommitMsg() {
		t.Error("IsCommitMsg predicate failure")
	}
	if !GitHookPreCommit.IsPreCommit() || GitHookPrePush.IsPreCommit() {
		t.Error("IsPreCommit predicate failure")
	}
	if !GitHookPrePush.IsPrePush() || GitHookPostCheckout.IsPrePush() {
		t.Error("IsPrePush predicate failure")
	}
	if !GitHookPostCheckout.IsPostCheckout() || GitHookPostMerge.IsPostCheckout() {
		t.Error("IsPostCheckout predicate failure")
	}
	if !GitHookPostMerge.IsPostMerge() || GitHookPreRebase.IsPostMerge() {
		t.Error("IsPostMerge predicate failure")
	}
	if !GitHookPreRebase.IsPreRebase() || GitHookPrepareCommitMsg.IsPreRebase() {
		t.Error("IsPreRebase predicate failure")
	}
	if !GitHookPrepareCommitMsg.IsPrepareCommitMsg() || GitHookPostCommit.IsPrepareCommitMsg() {
		t.Error("IsPrepareCommitMsg predicate failure")
	}
	if !GitHookPostCommit.IsPostCommit() || GitHookCommitMsg.IsPostCommit() {
		t.Error("IsPostCommit predicate failure")
	}
}

func TestHookActionType(t *testing.T) {
	testHookActionAllAndString(t)
	testHookActionParse(t)
	testHookActionPredicates(t)
}

func testHookActionAllAndString(t *testing.T) {
	all := AllHookActionTypes()
	if len(all) != 9 {
		t.Errorf("expected 9 hook action types, got %d", len(all))
	}

	for _, act := range all {
		if !act.IsValid() {
			t.Errorf("expected valid for %s", act)
		}
		if act.String() != string(act) {
			t.Errorf("expected string representation to match %s", act)
		}
		if !act.HasAction() {
			t.Errorf("expected HasAction to be true for %s", act)
		}
	}
}

func testHookActionParse(t *testing.T) {
	parsed, err := ParseHookActionType("  INSTALL ")
	if err != nil || parsed != HookActionInstall {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseHookActionType("unknown-action"); err == nil {
		t.Error("expected error parsing unknown hook action")
	}

	if HookActionType("invalid").IsValid() {
		t.Error("expected invalid action to not be valid")
	}
}

func testHookActionPredicates(t *testing.T) {
	if !HookActionInstall.IsInstall() || HookActionUninstall.IsInstall() {
		t.Error("IsInstall predicate failure")
	}
	if !HookActionUninstall.IsUninstall() || HookActionRun.IsUninstall() {
		t.Error("IsUninstall predicate failure")
	}
	if !HookActionRun.IsRun() || HookActionCheck.IsRun() {
		t.Error("IsRun predicate failure")
	}
	if !HookActionCheck.IsCheck() || HookActionList.IsCheck() {
		t.Error("IsCheck predicate failure")
	}
	if !HookActionList.IsList() || HookActionStatus.IsList() {
		t.Error("IsList predicate failure")
	}
	if !HookActionStatus.IsStatus() || HookActionHelp.IsStatus() {
		t.Error("IsStatus predicate failure")
	}
	if !HookActionHelp.IsHelp() || !HookActionHelpShort.IsHelp() || !HookActionHelpLong.IsHelp() {
		t.Error("IsHelp predicate failure")
	}
}

func TestHookPlanActionType(t *testing.T) {
	testHookPlanActionAllAndString(t)
	testHookPlanActionParse(t)
	testHookPlanActionPredicates(t)
}

func testHookPlanActionAllAndString(t *testing.T) {
	all := AllHookPlanActionTypes()
	if len(all) != 3 {
		t.Errorf("expected 3 plan action types, got %d", len(all))
	}

	for _, act := range all {
		if !act.IsValid() {
			t.Errorf("expected valid for %s", act)
		}
		if act.String() != string(act) {
			t.Errorf("expected string to match for %s", act)
		}
		if !act.HasPlanAction() {
			t.Errorf("expected HasPlanAction to be true for %s", act)
		}
	}
}

func testHookPlanActionParse(t *testing.T) {
	parsed, err := ParseHookPlanActionType("  ALREADY ")
	if err != nil || parsed != HookPlanActionAlready {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseHookPlanActionType("bad-plan"); err == nil {
		t.Error("expected error parsing bad plan action")
	}

	if HookPlanActionType("invalid").IsValid() {
		t.Error("expected invalid plan action to not be valid")
	}
}

func testHookPlanActionPredicates(t *testing.T) {
	if !HookPlanActionInstall.IsInstall() || HookPlanActionAlready.IsInstall() {
		t.Error("IsInstall predicate failure")
	}
	if !HookPlanActionAlready.IsAlready() || HookPlanActionSkip.IsAlready() {
		t.Error("IsAlready predicate failure")
	}
	if !HookPlanActionSkip.IsSkip() || HookPlanActionInstall.IsSkip() {
		t.Error("IsSkip predicate failure")
	}
}

func TestHookScopeType(t *testing.T) {
	testHookScopeAllAndString(t)
	testHookScopeParse(t)
	testHookScopePredicates(t)
}

func testHookScopeAllAndString(t *testing.T) {
	all := AllHookScopeTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 scope types, got %d", len(all))
	}

	for _, scope := range all {
		if !scope.IsValid() {
			t.Errorf("expected valid for %s", scope)
		}
		if scope.String() != string(scope) {
			t.Errorf("expected string to match for %s", scope)
		}
		if !scope.HasScope() {
			t.Errorf("expected HasScope to be true for %s", scope)
		}
	}
}

func testHookScopeParse(t *testing.T) {
	parsed, err := ParseHookScopeType("  GLOBAL ")
	if err != nil || parsed != HookScopeGlobal {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseHookScopeType("invalid-scope"); err == nil {
		t.Error("expected error parsing invalid scope")
	}

	if HookScopeType("invalid").IsValid() {
		t.Error("expected invalid scope to not be valid")
	}
}

func testHookScopePredicates(t *testing.T) {
	if !HookScopeAll.IsAll() || HookScopeRepo.IsAll() {
		t.Error("IsAll predicate failure")
	}
	if !HookScopeRepo.IsRepo() || HookScopeGlobal.IsRepo() {
		t.Error("IsRepo predicate failure")
	}
	if !HookScopeGlobal.IsGlobal() || HookScopeWorkdir.IsGlobal() {
		t.Error("IsGlobal predicate failure")
	}
	if !HookScopeWorkdir.IsWorkdir() || HookScopeAll.IsWorkdir() {
		t.Error("IsWorkdir predicate failure")
	}
}
