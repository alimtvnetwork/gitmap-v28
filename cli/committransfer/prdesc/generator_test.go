package prdesc

import (
	"strings"
	"testing"
	"time"
)

func TestGeneratePRDescription_Sections(t *testing.T) {
	meta := sampleTestMetadata()
	out := GeneratePRDescription(meta)
	assertContainsSection(t, out, "## Pull Request: Add user auth")
	assertContainsSection(t, out, "### Executive Summary")
	assertContainsSection(t, out, "### Component Impact")
	assertContainsSection(t, out, "### Commits Included")
	assertContainsSection(t, out, "### File Diff Summary")
	assertContainsSection(t, out, "### Safety Checklist")
}

func TestGeneratePRDescription_EmptySummary(t *testing.T) {
	meta := sampleTestMetadata()
	meta.ExecutiveSummary = ""
	out := GeneratePRDescription(meta)
	assertContainsSection(t, out, "Automated feature branch replay and merge.")
}

func TestGeneratePRDescription_DefaultSafetyChecks(t *testing.T) {
	meta := sampleTestMetadata()
	meta.SafetyChecks = nil
	out := GeneratePRDescription(meta)
	assertContainsSection(t, out, "- [x] Clean replay without merge conflicts")
	assertContainsSection(t, out, "- [x] Author identity and timestamps preserved")
}

func sampleTestMetadata() PRMetadata {
	return PRMetadata{
		PRNumber: 42, Title: "Add user auth", SourceBranch: "feature/auth",
		TargetBranch: "main", Author: "Alice <alice@example.com>",
		Timestamp:        time.Date(2026, 9, 20, 10, 0, 0, 0, time.UTC),
		ExecutiveSummary: "Implements JWT auth tokens.",
		Commits:          []PRCommitInfo{{ShortSHA: "a1b2c3d", Author: "Alice", Subject: "feat(auth): add login"}},
		Impacts:          []PRComponentImpact{{Component: "auth", Impact: "New feature", Files: 3}},
		DiffSummary:      PRFileDiffSummary{FilesChanged: 3, Additions: 120, Deletions: 5},
	}
}

func assertContainsSection(t *testing.T, output, expected string) {
	t.Helper()
	if !strings.Contains(output, expected) {
		t.Errorf("expected output to contain %q, but was not found in:\n%s", expected, output)
	}
}
