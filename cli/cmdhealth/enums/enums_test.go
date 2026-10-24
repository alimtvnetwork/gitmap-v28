package enums

import (
	"testing"
)

func TestHealthCategoryType(t *testing.T) {
	testHealthCategoryAllAndString(t)
	testHealthCategoryParse(t)
	testHealthCategoryPredicates(t)
}

func testHealthCategoryAllAndString(t *testing.T) {
	all := AllHealthCategories()
	if len(all) != 6 {
		t.Errorf("expected 6 categories, got %d", len(all))
	}
	if len(AllHealthCategoryTypes()) != 6 {
		t.Error("AllHealthCategoryTypes length mismatch")
	}

	for _, cat := range all {
		if !cat.IsValid() {
			t.Errorf("expected valid for %s", cat)
		}
		if cat.String() != string(cat) {
			t.Errorf("expected string representation to match %s", cat)
		}
		if !cat.HasCategory() {
			t.Errorf("expected HasCategory to be true for %s", cat)
		}
	}
}

func testHealthCategoryParse(t *testing.T) {
	parsed, err := ParseHealthCategory("  PIPELINE ")
	if err != nil || parsed != HealthCategoryPipeline {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	parsedType, errType := ParseHealthCategoryType("git")
	if errType != nil || parsedType != HealthCategoryGit {
		t.Errorf("unexpected parse type result: %v, err: %v", parsedType, errType)
	}

	if _, err := ParseHealthCategory("unknown"); err == nil {
		t.Error("expected error parsing unknown category")
	}

	if HealthCategoryType("invalid").IsValid() {
		t.Error("expected invalid category to not be valid")
	}
}

func testHealthCategoryPredicates(t *testing.T) {
	if !HealthCategoryGit.IsGit() || HealthCategoryPipeline.IsGit() {
		t.Error("IsGit predicate failure")
	}
	if !HealthCategoryPipeline.IsPipeline() || HealthCategoryFleet.IsPipeline() {
		t.Error("IsPipeline predicate failure")
	}
	if !HealthCategoryFleet.IsFleet() || HealthCategoryDb.IsFleet() {
		t.Error("IsFleet predicate failure")
	}
	if !HealthCategoryDb.IsDb() || HealthCategorySystem.IsDb() {
		t.Error("IsDb predicate failure")
	}
	if !HealthCategorySystem.IsSystem() || HealthCategoryNetwork.IsSystem() {
		t.Error("IsSystem predicate failure")
	}
	if !HealthCategoryNetwork.IsNetwork() || HealthCategoryGit.IsNetwork() {
		t.Error("IsNetwork predicate failure")
	}
}

func TestHealthSeverityType(t *testing.T) {
	testHealthSeverityAllAndString(t)
	testHealthSeverityParse(t)
	testHealthSeverityPredicates(t)
}

func testHealthSeverityAllAndString(t *testing.T) {
	all := AllHealthSeverities()
	if len(all) != 4 {
		t.Errorf("expected 4 severities, got %d", len(all))
	}
	if len(AllHealthSeverityTypes()) != 4 {
		t.Error("AllHealthSeverityTypes length mismatch")
	}

	for _, sev := range all {
		if !sev.IsValid() {
			t.Errorf("expected valid for %s", sev)
		}
		if sev.String() != string(sev) {
			t.Errorf("expected string representation to match %s", sev)
		}
		if !sev.HasSeverity() {
			t.Errorf("expected HasSeverity to be true for %s", sev)
		}
	}
}

func testHealthSeverityParse(t *testing.T) {
	parsed, err := ParseHealthSeverity("  WARN ")
	if err != nil || parsed != HealthSeverityWarn {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	parsedType, errType := ParseHealthSeverityType("fatal")
	if errType != nil || parsedType != HealthSeverityFatal {
		t.Errorf("unexpected parse type result: %v, err: %v", parsedType, errType)
	}

	if _, err := ParseHealthSeverity("unknown"); err == nil {
		t.Error("expected error parsing unknown severity")
	}

	if HealthSeverityType("invalid").IsValid() {
		t.Error("expected invalid severity to not be valid")
	}
}

func testHealthSeverityPredicates(t *testing.T) {
	if !HealthSeverityOk.IsOk() || HealthSeverityWarn.IsOk() {
		t.Error("IsOk predicate failure")
	}
	if !HealthSeverityWarn.IsWarn() || HealthSeverityError.IsWarn() {
		t.Error("IsWarn predicate failure")
	}
	if !HealthSeverityError.IsError() || HealthSeverityFatal.IsError() {
		t.Error("IsError predicate failure")
	}
	if !HealthSeverityFatal.IsFatal() || HealthSeverityOk.IsFatal() {
		t.Error("IsFatal predicate failure")
	}
	if !HealthSeverityError.IsUrgent() || !HealthSeverityFatal.IsUrgent() || HealthSeverityWarn.IsUrgent() {
		t.Error("IsUrgent predicate failure")
	}
	if !HealthSeverityWarn.HasDegradation() || HealthSeverityOk.HasDegradation() {
		t.Error("HasDegradation predicate failure")
	}
}

func TestHealthOutputFormatType(t *testing.T) {
	testHealthOutputFormatAllAndString(t)
	testHealthOutputFormatParse(t)
	testHealthOutputFormatPredicates(t)
}

func testHealthOutputFormatAllAndString(t *testing.T) {
	all := AllHealthOutputFormats()
	if len(all) != 4 {
		t.Errorf("expected 4 output formats, got %d", len(all))
	}
	if len(AllHealthOutputFormatTypes()) != 4 {
		t.Error("AllHealthOutputFormatTypes length mismatch")
	}

	for _, fmtVal := range all {
		if !fmtVal.IsValid() {
			t.Errorf("expected valid for %s", fmtVal)
		}
		if fmtVal.String() != string(fmtVal) {
			t.Errorf("expected string representation to match %s", fmtVal)
		}
		if !fmtVal.HasFormat() {
			t.Errorf("expected HasFormat to be true for %s", fmtVal)
		}
	}
}

func testHealthOutputFormatParse(t *testing.T) {
	parsed, err := ParseHealthOutputFormat("  JSON ")
	if err != nil || parsed != HealthOutputFormatJSON {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	parsedType, errType := ParseHealthOutputFormatType("markdown")
	if errType != nil || parsedType != HealthOutputFormatMarkdown {
		t.Errorf("unexpected parse type result: %v, err: %v", parsedType, errType)
	}

	if _, err := ParseHealthOutputFormat("unknown"); err == nil {
		t.Error("expected error parsing unknown output format")
	}

	if HealthOutputFormatType("invalid").IsValid() {
		t.Error("expected invalid output format to not be valid")
	}
}

func testHealthOutputFormatPredicates(t *testing.T) {
	if !HealthOutputFormatTerminal.IsTerminal() || HealthOutputFormatJSON.IsTerminal() {
		t.Error("IsTerminal predicate failure")
	}
	if !HealthOutputFormatJSON.IsJSON() || HealthOutputFormatSummary.IsJSON() {
		t.Error("IsJSON predicate failure")
	}
	if !HealthOutputFormatSummary.IsSummary() || HealthOutputFormatMarkdown.IsSummary() {
		t.Error("IsSummary predicate failure")
	}
	if !HealthOutputFormatMarkdown.IsMarkdown() || HealthOutputFormatTerminal.IsMarkdown() {
		t.Error("IsMarkdown predicate failure")
	}
}

func TestHealthCheckStatusType(t *testing.T) {
	testHealthCheckStatusAllAndString(t)
	testHealthCheckStatusParse(t)
	testHealthCheckStatusPredicates(t)
}

func testHealthCheckStatusAllAndString(t *testing.T) {
	all := AllHealthCheckStatuses()
	if len(all) != 4 {
		t.Errorf("expected 4 check statuses, got %d", len(all))
	}
	if len(AllHealthCheckStatusTypes()) != 4 {
		t.Error("AllHealthCheckStatusTypes length mismatch")
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

func testHealthCheckStatusParse(t *testing.T) {
	parsed, err := ParseHealthCheckStatus("  YELLOW ")
	if err != nil || parsed != HealthCheckStatusYellow {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	parsedType, errType := ParseHealthCheckStatusType("skipped")
	if errType != nil || parsedType != HealthCheckStatusSkipped {
		t.Errorf("unexpected parse type result: %v, err: %v", parsedType, errType)
	}

	if _, err := ParseHealthCheckStatus("unknown"); err == nil {
		t.Error("expected error parsing unknown check status")
	}

	if HealthCheckStatusType("invalid").IsValid() {
		t.Error("expected invalid check status to not be valid")
	}
}

func testHealthCheckStatusPredicates(t *testing.T) {
	if !HealthCheckStatusGreen.IsGreen() || HealthCheckStatusYellow.IsGreen() {
		t.Error("IsGreen predicate failure")
	}
	if !HealthCheckStatusYellow.IsYellow() || HealthCheckStatusRed.IsYellow() {
		t.Error("IsYellow predicate failure")
	}
	if !HealthCheckStatusRed.IsRed() || HealthCheckStatusSkipped.IsRed() {
		t.Error("IsRed predicate failure")
	}
	if !HealthCheckStatusSkipped.IsSkipped() || HealthCheckStatusGreen.IsSkipped() {
		t.Error("IsSkipped predicate failure")
	}
	if !HealthCheckStatusGreen.HasPassed() || HealthCheckStatusRed.HasPassed() {
		t.Error("HasPassed predicate failure")
	}
	if !HealthCheckStatusRed.HasFailed() || HealthCheckStatusGreen.HasFailed() {
		t.Error("HasFailed predicate failure")
	}
}
