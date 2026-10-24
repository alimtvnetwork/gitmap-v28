package enums

import (
	"testing"
)

func TestStaleOutputFormatType(t *testing.T) {
	testStaleOutputFormatAllAndString(t)
	testStaleOutputFormatParse(t)
	testStaleOutputFormatPredicates(t)
}

func testStaleOutputFormatAllAndString(t *testing.T) {
	all := AllStaleOutputFormatTypes()
	if len(all) != 3 {
		t.Errorf("expected 3 format variants, got %d", len(all))
	}

	for _, format := range all {
		if !format.IsValid() {
			t.Errorf("expected valid for %s", format)
		}
		if format.String() != string(format) {
			t.Errorf("expected string representation to match %s", format)
		}
		if !format.HasFormat() {
			t.Errorf("expected HasFormat to be true for %s", format)
		}
	}
}

func testStaleOutputFormatParse(t *testing.T) {
	parsed, err := ParseStaleOutputFormatType("  JSON ")
	if err != nil || parsed != StaleOutputFormatJSON {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	empty, err := ParseStaleOutputFormatType("")
	if err != nil || empty != StaleOutputFormatTable {
		t.Errorf("expected empty string to default to table, got: %v, err: %v", empty, err)
	}

	if _, err := ParseStaleOutputFormatType("invalid"); err == nil {
		t.Error("expected error parsing invalid stale output format")
	}

	if StaleOutputFormatType("bogus").IsValid() {
		t.Error("expected bogus format to not be valid")
	}
}

func testStaleOutputFormatPredicates(t *testing.T) {
	if !StaleOutputFormatTable.IsTable() || StaleOutputFormatJSON.IsTable() {
		t.Error("IsTable predicate failure")
	}
	if !StaleOutputFormatJSON.IsJSON() || StaleOutputFormatCSV.IsJSON() {
		t.Error("IsJSON predicate failure")
	}
	if !StaleOutputFormatCSV.IsCSV() || StaleOutputFormatTable.IsCSV() {
		t.Error("IsCSV predicate failure")
	}
}

func TestStaleActionType(t *testing.T) {
	testStaleActionAllAndString(t)
	testStaleActionParse(t)
	testStaleActionPredicates(t)
}

func testStaleActionAllAndString(t *testing.T) {
	all := AllStaleActionTypes()
	if len(all) != 3 {
		t.Errorf("expected 3 action variants, got %d", len(all))
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

func testStaleActionParse(t *testing.T) {
	parsed, err := ParseStaleActionType("  ARCHIVE ")
	if err != nil || parsed != StaleActionArchive {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseStaleActionType("invalid"); err == nil {
		t.Error("expected error parsing invalid stale action")
	}

	if StaleActionType("bogus").IsValid() {
		t.Error("expected bogus action to not be valid")
	}
}

func testStaleActionPredicates(t *testing.T) {
	if !StaleActionInspect.IsInspect() || StaleActionArchive.IsInspect() {
		t.Error("IsInspect predicate failure")
	}
	if !StaleActionArchive.IsArchive() || StaleActionPreview.IsArchive() {
		t.Error("IsArchive predicate failure")
	}
	if !StaleActionPreview.IsPreview() || StaleActionInspect.IsPreview() {
		t.Error("IsPreview predicate failure")
	}
}

func TestStaleFilterType(t *testing.T) {
	testStaleFilterAllAndString(t)
	testStaleFilterParse(t)
	testStaleFilterPredicates(t)
}

func testStaleFilterAllAndString(t *testing.T) {
	all := AllStaleFilterTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 filter variants, got %d", len(all))
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

func testStaleFilterParse(t *testing.T) {
	parsed, err := ParseStaleFilterType("  STALE ")
	if err != nil || parsed != StaleFilterStale {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseStaleFilterType("invalid"); err == nil {
		t.Error("expected error parsing invalid stale filter")
	}

	if StaleFilterType("bogus").IsValid() {
		t.Error("expected bogus filter to not be valid")
	}
}

func testStaleFilterPredicates(t *testing.T) {
	if !StaleFilterAll.IsAll() || StaleFilterStale.IsAll() {
		t.Error("IsAll predicate failure")
	}
	if !StaleFilterStale.IsStale() || StaleFilterActive.IsStale() {
		t.Error("IsStale predicate failure")
	}
	if !StaleFilterActive.IsActive() || StaleFilterArchived.IsActive() {
		t.Error("IsActive predicate failure")
	}
	if !StaleFilterArchived.IsArchived() || StaleFilterAll.IsArchived() {
		t.Error("IsArchived predicate failure")
	}
}
