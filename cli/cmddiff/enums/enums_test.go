package enums

import (
	"testing"
)

func TestDiffViewModeType(t *testing.T) {
	testDiffViewModeAllAndString(t)
	testDiffViewModeParse(t)
	testDiffViewModePredicates(t)
}

func testDiffViewModeAllAndString(t *testing.T) {
	all := AllDiffViewModeTypes()
	if len(all) != 5 {
		t.Errorf("expected 5 view mode variants, got %d", len(all))
	}

	for _, mode := range all {
		if !mode.IsValid() {
			t.Errorf("expected valid for %s", mode)
		}
		if mode.String() != string(mode) {
			t.Errorf("expected string representation to match %s", mode)
		}
		if !mode.HasViewMode() {
			t.Errorf("expected HasViewMode to be true for %s", mode)
		}
	}
}

func testDiffViewModeParse(t *testing.T) {
	parsed, err := ParseDiffViewModeType("  SIDE-BY-SIDE ")
	if err != nil || parsed != DiffViewModeSideBySide {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseDiffViewModeType("unknown"); err == nil {
		t.Error("expected error parsing unknown view mode")
	}

	if DiffViewModeType("invalid").IsValid() {
		t.Error("expected invalid view mode to not be valid")
	}
}

func testDiffViewModePredicates(t *testing.T) {
	if !DiffViewModeUnified.IsUnified() || DiffViewModeStat.IsUnified() {
		t.Error("IsUnified predicate failure")
	}
	if !DiffViewModeSideBySide.IsSideBySide() || DiffViewModeSummary.IsSideBySide() {
		t.Error("IsSideBySide predicate failure")
	}
	if !DiffViewModeStat.IsStat() || DiffViewModeCompact.IsStat() {
		t.Error("IsStat predicate failure")
	}
	if !DiffViewModeSummary.IsSummary() || DiffViewModeUnified.IsSummary() {
		t.Error("IsSummary predicate failure")
	}
	if !DiffViewModeCompact.IsCompact() || DiffViewModeSideBySide.IsCompact() {
		t.Error("IsCompact predicate failure")
	}
}

func TestDiffOutputFormatType(t *testing.T) {
	testDiffOutputFormatAllAndString(t)
	testDiffOutputFormatParse(t)
	testDiffOutputFormatPredicates(t)
}

func testDiffOutputFormatAllAndString(t *testing.T) {
	all := AllDiffOutputFormatTypes()
	if len(all) != 3 {
		t.Errorf("expected 3 output format variants, got %d", len(all))
	}

	for _, fmtType := range all {
		if !fmtType.IsValid() {
			t.Errorf("expected valid for %s", fmtType)
		}
		if fmtType.String() != string(fmtType) {
			t.Errorf("expected string representation to match %s", fmtType)
		}
		if !fmtType.HasOutputFormat() {
			t.Errorf("expected HasOutputFormat to be true for %s", fmtType)
		}
	}
}

func testDiffOutputFormatParse(t *testing.T) {
	parsed, err := ParseDiffOutputFormatType("  JSON ")
	if err != nil || parsed != DiffOutputFormatJSON {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseDiffOutputFormatType("yaml"); err == nil {
		t.Error("expected error parsing unknown output format")
	}

	if DiffOutputFormatType("bogus").IsValid() {
		t.Error("expected invalid output format to not be valid")
	}
}

func testDiffOutputFormatPredicates(t *testing.T) {
	if !DiffOutputFormatText.IsText() || DiffOutputFormatJSON.IsText() {
		t.Error("IsText predicate failure")
	}
	if !DiffOutputFormatJSON.IsJSON() || DiffOutputFormatRaw.IsJSON() {
		t.Error("IsJSON predicate failure")
	}
	if !DiffOutputFormatRaw.IsRaw() || DiffOutputFormatText.IsRaw() {
		t.Error("IsRaw predicate failure")
	}
}

func TestDiffScopeType(t *testing.T) {
	testDiffScopeAllAndString(t)
	testDiffScopeParse(t)
	testDiffScopePredicates(t)
}

func testDiffScopeAllAndString(t *testing.T) {
	all := AllDiffScopeTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 scope variants, got %d", len(all))
	}

	for _, sc := range all {
		if !sc.IsValid() {
			t.Errorf("expected valid for %s", sc)
		}
		if sc.String() != string(sc) {
			t.Errorf("expected string representation to match %s", sc)
		}
		if !sc.HasScope() {
			t.Errorf("expected HasScope to be true for %s", sc)
		}
	}
}

func testDiffScopeParse(t *testing.T) {
	parsed, err := ParseDiffScopeType("  WORKING-TREE ")
	if err != nil || parsed != DiffScopeWorkingTree {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseDiffScopeType("remote"); err == nil {
		t.Error("expected error parsing unknown scope")
	}

	if DiffScopeType("unknown").IsValid() {
		t.Error("expected invalid scope to not be valid")
	}
}

func testDiffScopePredicates(t *testing.T) {
	if !DiffScopeWorkingTree.IsWorkingTree() || DiffScopeStaged.IsWorkingTree() {
		t.Error("IsWorkingTree predicate failure")
	}
	if !DiffScopeStaged.IsStaged() || DiffScopeBranch.IsStaged() {
		t.Error("IsStaged predicate failure")
	}
	if !DiffScopeBranch.IsBranch() || DiffScopeCommit.IsBranch() {
		t.Error("IsBranch predicate failure")
	}
	if !DiffScopeCommit.IsCommit() || DiffScopeWorkingTree.IsCommit() {
		t.Error("IsCommit predicate failure")
	}
}

func TestDiffFilterType(t *testing.T) {
	testDiffFilterAllAndString(t)
	testDiffFilterParse(t)
	testDiffFilterPredicates(t)
}

func testDiffFilterAllAndString(t *testing.T) {
	all := AllDiffFilterTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 filter variants, got %d", len(all))
	}

	for _, flt := range all {
		if !flt.IsValid() {
			t.Errorf("expected valid for %s", flt)
		}
		if flt.String() != string(flt) {
			t.Errorf("expected string representation to match %s", flt)
		}
		if !flt.HasFilter() {
			t.Errorf("expected HasFilter to be true for %s", flt)
		}
	}
}

func testDiffFilterParse(t *testing.T) {
	parsed, err := ParseDiffFilterType("  ONLY-CONFLICTS ")
	if err != nil || parsed != DiffFilterOnlyConflicts {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseDiffFilterType("excluded"); err == nil {
		t.Error("expected error parsing unknown filter")
	}

	if DiffFilterType("invalid").IsValid() {
		t.Error("expected invalid filter to not be valid")
	}
}

func testDiffFilterPredicates(t *testing.T) {
	if !DiffFilterAll.IsAll() || DiffFilterOnlyConflicts.IsAll() {
		t.Error("IsAll predicate failure")
	}
	if !DiffFilterOnlyConflicts.IsOnlyConflicts() || DiffFilterOnlyMissing.IsOnlyConflicts() {
		t.Error("IsOnlyConflicts predicate failure")
	}
	if !DiffFilterOnlyMissing.IsOnlyMissing() || DiffFilterIncludeIdentical.IsOnlyMissing() {
		t.Error("IsOnlyMissing predicate failure")
	}
	if !DiffFilterIncludeIdentical.IsIncludeIdentical() || DiffFilterAll.IsIncludeIdentical() {
		t.Error("IsIncludeIdentical predicate failure")
	}
}
