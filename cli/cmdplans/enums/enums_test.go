package enums

import "testing"

func TestPlanStatusType_Lifecycle(t *testing.T) {
	statuses := AllPlanStatuses()
	if len(statuses) != 6 {
		t.Fatalf("expected 6 statuses, got %d", len(statuses))
	}

	for _, s := range statuses {
		if !s.IsValid() {
			t.Errorf("status %q should be valid", s)
		}
		if s.String() != string(s) {
			t.Errorf("expected %q, got %q", string(s), s.String())
		}
	}

	invalid := PlanStatusType("nonexistent")
	if invalid.IsValid() {
		t.Errorf("expected invalid status to return false from IsValid")
	}
}

func TestPlanStatusType_Booleans(t *testing.T) {
	if !PlanStatusPending.IsPending() || PlanStatusPending.IsCompleted() {
		t.Errorf("unexpected boolean result for pending status")
	}
	if !PlanStatusActive.IsActive() || PlanStatusActive.IsFailed() {
		t.Errorf("unexpected boolean result for active status")
	}
	if !PlanStatusInProgress.IsInProgress() || PlanStatusInProgress.IsArchived() {
		t.Errorf("unexpected boolean result for in_progress status")
	}
	if !PlanStatusCompleted.IsCompleted() || PlanStatusCompleted.IsPending() {
		t.Errorf("unexpected boolean result for completed status")
	}
	if !PlanStatusArchived.IsArchived() || PlanStatusArchived.IsActive() {
		t.Errorf("unexpected boolean result for archived status")
	}
	if !PlanStatusFailed.IsFailed() || PlanStatusFailed.IsCompleted() {
		t.Errorf("unexpected boolean result for failed status")
	}
}

func TestParsePlanStatus(t *testing.T) {
	cases := []struct {
		input    string
		expected PlanStatusType
	}{
		{"pending", PlanStatusPending},
		{"PENDING", PlanStatusPending},
		{"queued", PlanStatusPending},
		{"active", PlanStatusActive},
		{"in_progress", PlanStatusInProgress},
		{"IN PROGRESS", PlanStatusInProgress},
		{"inprogress", PlanStatusInProgress},
		{"completed", PlanStatusCompleted},
		{"COMPLETE", PlanStatusCompleted},
		{"done", PlanStatusCompleted},
		{"archived", PlanStatusArchived},
		{"archive", PlanStatusArchived},
		{"failed", PlanStatusFailed},
		{"fail", PlanStatusFailed},
		{"error", PlanStatusFailed},
	}

	for _, c := range cases {
		got, err := ParsePlanStatus(c.input)
		if err != nil {
			t.Errorf("unexpected error for input %q: %v", c.input, err)
		}
		if got != c.expected {
			t.Errorf("for input %q, expected %q, got %q", c.input, c.expected, got)
		}
	}

	if _, err := ParsePlanStatus("invalid-status"); err == nil {
		t.Errorf("expected error for invalid plan status")
	}
}

func TestPlansFilterModeType_Lifecycle(t *testing.T) {
	modes := AllPlansFilterModes()
	if len(modes) != 4 {
		t.Fatalf("expected 4 filter modes, got %d", len(modes))
	}

	for _, m := range modes {
		if !m.IsValid() {
			t.Errorf("filter mode %q should be valid", m)
		}
		if m.String() != string(m) {
			t.Errorf("expected %q, got %q", string(m), m.String())
		}
	}

	invalid := PlansFilterModeType("invalid-mode")
	if invalid.IsValid() {
		t.Errorf("expected invalid filter mode to return false")
	}
}

func TestPlansFilterModeType_Booleans(t *testing.T) {
	if !PlansFilterModeAll.IsAll() || PlansFilterModeAll.IsPending() {
		t.Errorf("unexpected boolean for all filter mode")
	}
	if !PlansFilterModeActive.IsActive() || PlansFilterModeActive.IsCompleted() {
		t.Errorf("unexpected boolean for active filter mode")
	}
	if !PlansFilterModePending.IsPending() || PlansFilterModePending.IsAll() {
		t.Errorf("unexpected boolean for pending filter mode")
	}
	if !PlansFilterModeCompleted.IsCompleted() || PlansFilterModeCompleted.IsActive() {
		t.Errorf("unexpected boolean for completed filter mode")
	}
}

func TestParsePlansFilterMode(t *testing.T) {
	cases := []struct {
		input    string
		expected PlansFilterModeType
	}{
		{"all", PlansFilterModeAll},
		{"ALL", PlansFilterModeAll},
		{"active", PlansFilterModeActive},
		{"pending", PlansFilterModePending},
		{"completed", PlansFilterModeCompleted},
	}

	for _, c := range cases {
		got, err := ParsePlansFilterMode(c.input)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", c.input, err)
		}
		if got != c.expected {
			t.Errorf("for %q expected %q, got %q", c.input, c.expected, got)
		}
	}

	if _, err := ParsePlansFilterMode("unknown"); err == nil {
		t.Errorf("expected error for unknown filter mode")
	}
}

func TestPlansOutputFormatType_Lifecycle(t *testing.T) {
	formats := AllPlansOutputFormats()
	if len(formats) != 5 {
		t.Fatalf("expected 5 output formats, got %d", len(formats))
	}

	for _, f := range formats {
		if !f.IsValid() {
			t.Errorf("format %q should be valid", f)
		}
		if f.String() != string(f) {
			t.Errorf("expected %q, got %q", string(f), f.String())
		}
	}

	invalid := PlansOutputFormatType("xml")
	if invalid.IsValid() {
		t.Errorf("expected xml to be invalid")
	}
}

func TestPlansOutputFormatType_Booleans(t *testing.T) {
	if !PlansOutputFormatTable.IsTable() || PlansOutputFormatTable.IsJSON() {
		t.Errorf("unexpected boolean for table format")
	}
	if !PlansOutputFormatJSON.IsJSON() || PlansOutputFormatJSON.IsMarkdown() {
		t.Errorf("unexpected boolean for json format")
	}
	if !PlansOutputFormatMarkdown.IsMarkdown() || PlansOutputFormatMarkdown.IsCSV() {
		t.Errorf("unexpected boolean for markdown format")
	}
	if !PlansOutputFormatCSV.IsCSV() || PlansOutputFormatCSV.IsYAML() {
		t.Errorf("unexpected boolean for csv format")
	}
	if !PlansOutputFormatYAML.IsYAML() || PlansOutputFormatYAML.IsTable() {
		t.Errorf("unexpected boolean for yaml format")
	}
}

func TestParsePlansOutputFormat(t *testing.T) {
	cases := []struct {
		input    string
		expected PlansOutputFormatType
	}{
		{"table", PlansOutputFormatTable},
		{"json", PlansOutputFormatJSON},
		{"JSON", PlansOutputFormatJSON},
		{"markdown", PlansOutputFormatMarkdown},
		{"md", PlansOutputFormatMarkdown},
		{"csv", PlansOutputFormatCSV},
		{"yaml", PlansOutputFormatYAML},
		{"yml", PlansOutputFormatYAML},
	}

	for _, c := range cases {
		got, err := ParsePlansOutputFormat(c.input)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", c.input, err)
		}
		if got != c.expected {
			t.Errorf("for %q expected %q, got %q", c.input, c.expected, got)
		}
	}

	if _, err := ParsePlansOutputFormat("binary"); err == nil {
		t.Errorf("expected error for binary output format")
	}
}

func TestPlanMetaFileType_Lifecycle(t *testing.T) {
	files := AllPlanMetaFiles()
	if len(files) != 6 {
		t.Fatalf("expected 6 meta files, got %d", len(files))
	}

	for _, f := range files {
		if !f.IsValid() {
			t.Errorf("meta file %q should be valid", f)
		}
		if f.String() != string(f) {
			t.Errorf("expected %q, got %q", string(f), f.String())
		}
	}

	invalid := PlanMetaFileType("random.txt")
	if invalid.IsValid() {
		t.Errorf("expected random.txt to be invalid")
	}
}

func TestPlanMetaFileType_Booleans(t *testing.T) {
	if !PlanMetaFileReadme.IsReadme() || PlanMetaFileReadme.IsSummary() {
		t.Errorf("unexpected boolean for readme")
	}
	if !PlanMetaFileSummary.IsSummary() || PlanMetaFileSummary.IsChecklist() {
		t.Errorf("unexpected boolean for summary")
	}
	if !PlanMetaFileChecklist.IsChecklist() || PlanMetaFileChecklist.IsIndex() {
		t.Errorf("unexpected boolean for checklist")
	}
	if !PlanMetaFileIndex.IsIndex() || PlanMetaFileIndex.IsLedger() {
		t.Errorf("unexpected boolean for index")
	}
	if !PlanMetaFileLedger.IsLedger() || PlanMetaFileLedger.IsWhatToRead() {
		t.Errorf("unexpected boolean for ledger")
	}
	if !PlanMetaFileWhatToRead.IsWhatToRead() || PlanMetaFileWhatToRead.IsReadme() {
		t.Errorf("unexpected boolean for what-to-read")
	}
}

func TestParsePlanMetaFile(t *testing.T) {
	cases := []struct {
		input    string
		expected PlanMetaFileType
	}{
		{"readme.md", PlanMetaFileReadme},
		{"README.MD", PlanMetaFileReadme},
		{"summary.md", PlanMetaFileSummary},
		{"checklist.md", PlanMetaFileChecklist},
		{"01-index.md", PlanMetaFileIndex},
		{"ledger.md", PlanMetaFileLedger},
		{"what-to-read.md", PlanMetaFileWhatToRead},
	}

	for _, c := range cases {
		got, err := ParsePlanMetaFile(c.input)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", c.input, err)
		}
		if got != c.expected {
			t.Errorf("for %q expected %q, got %q", c.input, c.expected, got)
		}
	}

	if _, err := ParsePlanMetaFile("subtask.md"); err == nil {
		t.Errorf("expected error for subtask.md")
	}
}

func TestIsMetaFileName(t *testing.T) {
	if !IsMetaFileName("readme.md") {
		t.Errorf("expected readme.md to be meta file")
	}
	if !IsMetaFileName("01-INDEX.md") {
		t.Errorf("expected 01-INDEX.md to be meta file")
	}
	if !IsMetaFileName("ledger.md") {
		t.Errorf("expected ledger.md to be meta file")
	}
	if !IsMetaFileName("what-to-read.md") {
		t.Errorf("expected what-to-read.md to be meta file")
	}
	if !IsMetaFileName("summary.md") {
		t.Errorf("expected summary.md to be meta file")
	}
	if !IsMetaFileName("checklist.md") {
		t.Errorf("expected checklist.md to be meta file")
	}
	if IsMetaFileName("01-feature-spec.md") {
		t.Errorf("expected 01-feature-spec.md not to be meta file")
	}
}
