package enums

import (
	"testing"
)

func TestBranchFilterType(t *testing.T) {
	testBranchFilterAllAndString(t)
	testBranchFilterParse(t)
	testBranchFilterPredicates(t)
}

func testBranchFilterAllAndString(t *testing.T) {
	all := AllBranchFilterTypes()
	if len(all) != 5 {
		t.Errorf("expected 5 filter variants, got %d", len(all))
	}

	for _, filter := range all {
		if !filter.IsValid() {
			t.Errorf("expected valid for %s", filter)
		}
		if filter.String() != string(filter) {
			t.Errorf("expected string representation to match %s", filter)
		}
		if !filter.HasFilter() {
			t.Errorf("expected HasFilter to be true for %s", filter)
		}
	}
}

func testBranchFilterParse(t *testing.T) {
	parsed, err := ParseBranchFilterType("  REMOTE ")
	if err != nil || parsed != BranchFilterRemote {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseBranchFilterType("invalid"); err == nil {
		t.Error("expected error parsing invalid filter")
	}

	if BranchFilterType("bogus").IsValid() {
		t.Error("expected bogus filter to not be valid")
	}
}

func testBranchFilterPredicates(t *testing.T) {
	if !BranchFilterAll.IsAll() || BranchFilterLocal.IsAll() {
		t.Error("IsAll predicate failure")
	}
	if !BranchFilterLocal.IsLocal() || BranchFilterRemote.IsLocal() {
		t.Error("IsLocal predicate failure")
	}
	if !BranchFilterRemote.IsRemote() || BranchFilterMerged.IsRemote() {
		t.Error("IsRemote predicate failure")
	}
	if !BranchFilterMerged.IsMerged() || BranchFilterUnmerged.IsMerged() {
		t.Error("IsMerged predicate failure")
	}
	if !BranchFilterUnmerged.IsUnmerged() || BranchFilterAll.IsUnmerged() {
		t.Error("IsUnmerged predicate failure")
	}
}

func TestBranchActionType(t *testing.T) {
	testBranchActionAllAndString(t)
	testBranchActionParse(t)
	testBranchActionPredicates(t)
}

func testBranchActionAllAndString(t *testing.T) {
	all := AllBranchActionTypes()
	if len(all) != 7 {
		t.Errorf("expected 7 action variants, got %d", len(all))
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

func testBranchActionParse(t *testing.T) {
	parsed, err := ParseBranchActionType("  DEF ")
	if err != nil || parsed != BranchActionDefaultAlias {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseBranchActionType("unknown"); err == nil {
		t.Error("expected error parsing unknown action")
	}

	if BranchActionType("invalid").IsValid() {
		t.Error("expected invalid action to not be valid")
	}
}

func testBranchActionPredicates(t *testing.T) {
	testBranchActionDefaultPredicates(t)
	testBranchActionSubPredicates(t)
}

func testBranchActionDefaultPredicates(t *testing.T) {
	if !BranchActionDefault.IsDefault() || BranchActionDefaultAlias.IsDefault() {
		t.Error("IsDefault predicate failure")
	}
	if !BranchActionDefaultAlias.IsDefaultAlias() || BranchActionDefault.IsDefaultAlias() {
		t.Error("IsDefaultAlias predicate failure")
	}
	if !BranchActionDefault.HasDefaultAction() || !BranchActionDefaultAlias.HasDefaultAction() {
		t.Error("HasDefaultAction failure")
	}
	if BranchActionList.HasDefaultAction() {
		t.Error("HasDefaultAction false positive")
	}
}

func testBranchActionSubPredicates(t *testing.T) {
	if !BranchActionList.IsList() || BranchActionCurrent.IsList() {
		t.Error("IsList predicate failure")
	}
	if !BranchActionCurrent.IsCurrent() || BranchActionHelp.IsCurrent() {
		t.Error("IsCurrent predicate failure")
	}
	if !BranchActionLatestBranch.IsLatestBranch() || !BranchActionLatestBranchAlias.IsLatestBranchAlias() {
		t.Error("IsLatestBranch predicate failure")
	}
	if !BranchActionLatestBranch.HasLatestBranchAction() || !BranchActionLatestBranchAlias.HasLatestBranchAction() {
		t.Error("HasLatestBranchAction failure")
	}
	if !BranchActionHelp.IsHelp() || BranchActionList.IsHelp() {
		t.Error("IsHelp predicate failure")
	}
}

func TestBranchSortOrderType(t *testing.T) {
	testBranchSortOrderAllAndString(t)
	testBranchSortOrderParse(t)
	testBranchSortOrderPredicates(t)
}

func testBranchSortOrderAllAndString(t *testing.T) {
	all := AllBranchSortOrderTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 sort order variants, got %d", len(all))
	}

	for _, order := range all {
		if !order.IsValid() {
			t.Errorf("expected valid for %s", order)
		}
		if order.String() != string(order) {
			t.Errorf("expected string representation to match %s", order)
		}
		if !order.HasSortOrder() {
			t.Errorf("expected HasSortOrder to be true for %s", order)
		}
	}
}

func testBranchSortOrderParse(t *testing.T) {
	parsed, err := ParseBranchSortOrderType("  COMMITTERDATE ")
	if err != nil || parsed != BranchSortOrderCommitterDate {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseBranchSortOrderType("invalid"); err == nil {
		t.Error("expected error parsing invalid sort order")
	}

	if BranchSortOrderType("bogus").IsValid() {
		t.Error("expected bogus sort order to not be valid")
	}
}

func testBranchSortOrderPredicates(t *testing.T) {
	if !BranchSortOrderRecent.IsRecent() || BranchSortOrderAlphabetical.IsRecent() {
		t.Error("IsRecent predicate failure")
	}
	if !BranchSortOrderCommitterDate.IsCommitterDate() || BranchSortOrderAuthorDate.IsCommitterDate() {
		t.Error("IsCommitterDate predicate failure")
	}
	if !BranchSortOrderAuthorDate.IsAuthorDate() || BranchSortOrderRecent.IsAuthorDate() {
		t.Error("IsAuthorDate predicate failure")
	}
	if !BranchSortOrderAlphabetical.IsAlphabetical() || BranchSortOrderCommitterDate.IsAlphabetical() {
		t.Error("IsAlphabetical predicate failure")
	}
}

func TestBranchOutputFormatType(t *testing.T) {
	testBranchOutputFormatAllAndString(t)
	testBranchOutputFormatParse(t)
	testBranchOutputFormatPredicates(t)
}

func testBranchOutputFormatAllAndString(t *testing.T) {
	all := AllBranchOutputFormatTypes()
	if len(all) != 3 {
		t.Errorf("expected 3 output format variants, got %d", len(all))
	}

	for _, format := range all {
		if !format.IsValid() {
			t.Errorf("expected valid for %s", format)
		}
		if format.String() != string(format) {
			t.Errorf("expected string representation to match %s", format)
		}
		if !format.HasOutputFormat() {
			t.Errorf("expected HasOutputFormat to be true for %s", format)
		}
	}
}

func testBranchOutputFormatParse(t *testing.T) {
	parsed, err := ParseBranchOutputFormatType("  JSON ")
	if err != nil || parsed != BranchOutputFormatJSON {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseBranchOutputFormatType("xml"); err == nil {
		t.Error("expected error parsing invalid output format")
	}

	if BranchOutputFormatType("bogus").IsValid() {
		t.Error("expected bogus output format to not be valid")
	}
}

func testBranchOutputFormatPredicates(t *testing.T) {
	if !BranchOutputFormatTerminal.IsTerminal() || BranchOutputFormatJSON.IsTerminal() {
		t.Error("IsTerminal predicate failure")
	}
	if !BranchOutputFormatJSON.IsJSON() || BranchOutputFormatCSV.IsJSON() {
		t.Error("IsJSON predicate failure")
	}
	if !BranchOutputFormatCSV.IsCSV() || BranchOutputFormatTerminal.IsCSV() {
		t.Error("IsCSV predicate failure")
	}
}
