package enums

import (
	"testing"
)

func TestAddTargetType(t *testing.T) {
	testAddTargetAllAndString(t)
	testAddTargetParse(t)
	testAddTargetPredicates1(t)
	testAddTargetPredicates2(t)
}

func testAddTargetAllAndString(t *testing.T) {
	all := AllAddTargetTypes()
	if len(all) != 9 {
		t.Errorf("expected 9 target variants, got %d", len(all))
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

func testAddTargetParse(t *testing.T) {
	parsed, err := ParseAddTargetType("  COMMON-ATTR ")
	if err != nil || parsed != AddTargetCommonAttr {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseAddTargetType("invalid-target"); err == nil {
		t.Error("expected error parsing invalid add target")
	}

	if AddTargetType("unknown").IsValid() {
		t.Error("expected unknown target to not be valid")
	}
}

func testAddTargetPredicates1(t *testing.T) {
	if !AddTargetRepo.IsRepo() || AddTargetFolder.IsRepo() {
		t.Error("IsRepo predicate failure")
	}
	if !AddTargetFolder.IsFolder() || AddTargetRepo.IsFolder() {
		t.Error("IsFolder predicate failure")
	}
	if !AddTargetLfs.IsLfs() || AddTargetHook.IsLfs() {
		t.Error("IsLfs predicate failure")
	}
	if !AddTargetHook.IsHook() || AddTargetLfs.IsHook() {
		t.Error("IsHook predicate failure")
	}
	if !AddTargetCommonAttr.IsCommonAttr() || AddTargetCommonIgnore.IsCommonAttr() {
		t.Error("IsCommonAttr predicate failure")
	}
}

func testAddTargetPredicates2(t *testing.T) {
	if !AddTargetCommonIgnore.IsCommonIgnore() || AddTargetCommonAttr.IsCommonIgnore() {
		t.Error("IsCommonIgnore predicate failure")
	}
	if !AddTargetAttributes.IsAttributes() || AddTargetIgnore.IsAttributes() {
		t.Error("IsAttributes predicate failure")
	}
	if !AddTargetIgnore.IsIgnore() || AddTargetAttributes.IsIgnore() {
		t.Error("IsIgnore predicate failure")
	}
	if !AddTargetLfsInstall.IsLfsInstall() || AddTargetRepo.IsLfsInstall() {
		t.Error("IsLfsInstall predicate failure")
	}
}

func TestAddActionType(t *testing.T) {
	testAddActionAllAndString(t)
	testAddActionParse(t)
	testAddActionPredicates(t)
}

func testAddActionAllAndString(t *testing.T) {
	all := AllAddActionTypes()
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

func testAddActionParse(t *testing.T) {
	parsed, err := ParseAddActionType("  MERGE ")
	if err != nil || parsed != AddActionMerge {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseAddActionType("bogus"); err == nil {
		t.Error("expected error parsing invalid add action")
	}

	if AddActionType("invalid").IsValid() {
		t.Error("expected invalid action to not be valid")
	}
}

func testAddActionPredicates(t *testing.T) {
	if !AddActionCreate.IsCreate() || AddActionAppend.IsCreate() {
		t.Error("IsCreate predicate failure")
	}
	if !AddActionAppend.IsAppend() || AddActionMerge.IsAppend() {
		t.Error("IsAppend predicate failure")
	}
	if !AddActionMerge.IsMerge() || AddActionOverwrite.IsMerge() {
		t.Error("IsMerge predicate failure")
	}
	if !AddActionOverwrite.IsOverwrite() || AddActionDiff.IsOverwrite() {
		t.Error("IsOverwrite predicate failure")
	}
	if !AddActionDiff.IsDiff() || AddActionCreate.IsDiff() {
		t.Error("IsDiff predicate failure")
	}
}

func TestAddModeType(t *testing.T) {
	testAddModeAllAndString(t)
	testAddModeParse(t)
	testAddModePredicates(t)
}

func testAddModeAllAndString(t *testing.T) {
	all := AllAddModeTypes()
	if len(all) != 3 {
		t.Errorf("expected 3 mode variants, got %d", len(all))
	}

	for _, mode := range all {
		if !mode.IsValid() {
			t.Errorf("expected valid for %s", mode)
		}
		if mode.String() != string(mode) {
			t.Errorf("expected string representation to match %s", mode)
		}
		if !mode.HasMode() {
			t.Errorf("expected HasMode to be true for %s", mode)
		}
	}
}

func testAddModeParse(t *testing.T) {
	parsed, err := ParseAddModeType("  DRY-RUN ")
	if err != nil || parsed != AddModeDryRun {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseAddModeType("unknown-mode"); err == nil {
		t.Error("expected error parsing invalid add mode")
	}

	if AddModeType("bogus").IsValid() {
		t.Error("expected bogus mode to not be valid")
	}
}

func testAddModePredicates(t *testing.T) {
	if !AddModeNormal.IsNormal() || AddModeDryRun.IsNormal() {
		t.Error("IsNormal predicate failure")
	}
	if !AddModeDryRun.IsDryRun() || AddModeForce.IsDryRun() {
		t.Error("IsDryRun predicate failure")
	}
	if !AddModeForce.IsForce() || AddModeNormal.IsForce() {
		t.Error("IsForce predicate failure")
	}
}
