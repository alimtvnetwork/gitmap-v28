package cmdplans

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/enums"
)

func TestParseIDAndSlug(t *testing.T) {
	id, slug := parseIDAndSlug("120-auth-feature.md")
	if id != 120 {
		t.Errorf("expected id 120, got %d", id)
	}
	if slug != "auth-feature" {
		t.Errorf("expected slug 'auth-feature', got %q", slug)
	}

	id2, slug2 := parseIDAndSlug("custom-task.md")
	if id2 != 0 {
		t.Errorf("expected id 0, got %d", id2)
	}
	if slug2 != "custom-task" {
		t.Errorf("expected slug 'custom-task', got %q", slug2)
	}
}

func TestExtractTitleFromLine(t *testing.T) {
	title, hasTitle := extractTitleFromLine("# Architecture Design Plan")
	if !hasTitle || title != "Architecture Design Plan" {
		t.Errorf("expected title 'Architecture Design Plan', got %q", title)
	}

	_, hasTitleInvalid := extractTitleFromLine("regular text line")
	if hasTitleInvalid {
		t.Error("expected false for regular text")
	}
}

func TestParseStatusToken(t *testing.T) {
	cases := []struct {
		input    string
		expected enums.PlanStatusType
	}{
		{"Done", StatusCompleted},
		{"Completed", StatusCompleted},
		{"In Progress", StatusInProgress},
		{"Active", StatusInProgress},
		{"Archived", StatusArchived},
		{"Failed", StatusFailed},
		{"Queued", StatusPending},
	}

	for _, c := range cases {
		got := parseStatusToken(c.input)
		if got != c.expected {
			t.Errorf("for input %q, expected %q, got %q", c.input, c.expected, got)
		}
	}
}

func TestExtractStatusFromLine(t *testing.T) {
	status, hasStatus := extractStatusFromLine("## Status: Completed")
	if !hasStatus || status != StatusCompleted {
		t.Errorf("expected StatusCompleted, got %v", status)
	}

	_, hasStatus2 := extractStatusFromLine("no status here")
	if hasStatus2 {
		t.Error("expected false when no status present")
	}
}

func TestFormatFallbackTitle(t *testing.T) {
	title := formatFallbackTitle("hello-world-subtask")
	if title != "Hello World Subtask" {
		t.Errorf("expected 'Hello World Subtask', got %q", title)
	}
}

func TestIsMetaFile(t *testing.T) {
	metaFiles := []string{
		"readme.md", "README.MD", "01-index.md",
		"ledger.md", "what-to-read.md", "summary.md", "checklist.md",
	}

	for _, f := range metaFiles {
		if !isMetaFile(f) {
			t.Errorf("expected %q to be recognized as meta file", f)
		}
	}

	if isMetaFile("01-app-spec.md") {
		t.Error("expected 01-app-spec.md not to be meta file")
	}
}

func TestDiscoverPlans_Integration(t *testing.T) {
	tempDir := t.TempDir()
	plansDir := filepath.Join(tempDir, constants.DirAiMemory, constants.DirPlans)
	pendingDir := filepath.Join(plansDir, constants.DirPending)
	completedDir := filepath.Join(plansDir, constants.DirCompleted)

	if err := os.MkdirAll(pendingDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(completedDir, 0755); err != nil {
		t.Fatal(err)
	}

	pendingFile := filepath.Join(pendingDir, "101-new-task.md")
	contentPending := "# New Task\n## Status: Pending\nDetails here"
	if err := os.WriteFile(pendingFile, []byte(contentPending), 0644); err != nil {
		t.Fatal(err)
	}

	completedFile := filepath.Join(completedDir, "099-old-task.md")
	contentCompleted := "# Old Task\n## Status: Completed\nDetails here"
	if err := os.WriteFile(completedFile, []byte(contentCompleted), 0644); err != nil {
		t.Fatal(err)
	}

	comp, pend, err := DiscoverPlans(tempDir)
	if err != nil {
		t.Fatalf("unexpected error from DiscoverPlans: %v", err)
	}
	if len(comp) != 1 {
		t.Errorf("expected 1 completed plan, got %d", len(comp))
	}
	if len(pend) != 1 {
		t.Errorf("expected 1 pending plan, got %d", len(pend))
	}
	if len(comp) > 0 && comp[0].ID != 99 {
		t.Errorf("expected completed plan id 99, got %d", comp[0].ID)
	}
	if len(pend) > 0 && pend[0].ID != 101 {
		t.Errorf("expected pending plan id 101, got %d", pend[0].ID)
	}
}
