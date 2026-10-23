package cmdpull

import (
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestPullTableNarrowWidth72(t *testing.T) {
	rows := []model.PullTableRow{
		{
			RepoName:     "ai-empathy-prompt-tuner",
			Branch:       "main",
			LatestBranch: "main",
			PRStatus:     "0 PRs",
			PullStatus:   "UP_TO_DATE",
			LastSHA:      "f1d94df",
			Duration:     "1.0s",
		},
	}

	layout := NewPullTableLayoutWithWidth(rows, 72)
	if layout.IsWide {
		t.Errorf("expected compact layout for width 72, got wide")
	}

	if layout.DividerLen+2 > 72 {
		t.Errorf("expected total row width <= 72, got %d", layout.DividerLen+2)
	}
}

func TestPullTableStandardWidth80(t *testing.T) {
	rows := []model.PullTableRow{
		{
			RepoName:     "prompt-architect",
			Branch:       "main",
			LatestBranch: "release/v1.35.0",
			PRStatus:     "0 PRs",
			PullStatus:   "UP_TO_DATE",
			LastSHA:      "03be798",
			Duration:     "1.0s",
		},
	}

	layout := NewPullTableLayoutWithWidth(rows, 80)
	if layout.IsWide {
		t.Errorf("expected compact layout for width 80, got wide")
	}

	if layout.DividerLen+2 > 80 {
		t.Errorf("expected total row width <= 80, got %d", layout.DividerLen+2)
	}
}

func TestPullTableWideWidth120(t *testing.T) {
	rows := []model.PullTableRow{
		{
			RepoName:     "core",
			Branch:       "main",
			LatestBranch: "release/v1.5.7",
			PRStatus:     "0 PRs",
			PullStatus:   "DIRTY",
			LastSHA:      "77c6edf",
			Duration:     "1.0s",
			IsDirty:      true,
		},
	}

	layout := NewPullTableLayoutWithWidth(rows, 120)
	if !layout.IsWide {
		t.Errorf("expected wide layout for width 120, got compact")
	}

	if layout.DividerLen+2 > 120 {
		t.Errorf("expected total row width <= 120, got %d", layout.DividerLen+2)
	}

	if layout.MaxLatestBr <= 0 {
		t.Errorf("expected positive MaxLatestBr in wide layout, got %d", layout.MaxLatestBr)
	}
}

func TestFormatCombinedBranch(t *testing.T) {
	same := formatCombinedBranch("main", "main", 20)
	if same != "main" {
		t.Errorf("expected 'main', got %q", same)
	}

	diff := formatCombinedBranch("main", "release/v1.35.0", 20)
	isExpectedDiff := diff == "main→v1.35.0" || diff == "main->v1.35.0"
	if !isExpectedDiff {
		t.Errorf("expected 'main→v1.35.0' or 'main->v1.35.0', got %q", diff)
	}
}

func TestPadVisualWithArrow(t *testing.T) {
	padded := PadVisual("main→v1.35.0", 17)
	visWidth := runewidth.StringWidth(stripANSI(padded))
	if visWidth != 17 {
		t.Errorf("expected visual width 17, got %d for %q", visWidth, padded)
	}
}

func TestFormatPullStatusCheckmarks(t *testing.T) {
	activeStatus := formatPullStatus("active", false)
	if !strings.Contains(activeStatus, "✔ active") {
		t.Errorf("expected '✔ active', got %q", activeStatus)
	}

	upToDateStatus := formatPullStatus("UP_TO_DATE", false)
	if !strings.Contains(upToDateStatus, "✔ active") {
		t.Errorf("expected '✔ active', got %q", upToDateStatus)
	}

	dirtyStatus := formatPullStatus("DIRTY", true)
	if !strings.Contains(dirtyStatus, "● dirty") {
		t.Errorf("expected '● dirty', got %q", dirtyStatus)
	}
}

func TestLeadTruncate(t *testing.T) {
	short := "main"
	if leadTruncate(short, 10) != "main" {
		t.Errorf("expected 'main', got %q", leadTruncate(short, 10))
	}

	long := "backup/plan-29-auto-aliasing-92927"
	truncated := leadTruncate(long, 12)
	if len(truncated) != 12 {
		t.Errorf("expected len 12, got len %d (%q)", len(truncated), truncated)
	}
	if !strings.HasPrefix(truncated, "...") {
		t.Errorf("expected prefix '...', got %q", truncated)
	}
	if !strings.HasSuffix(truncated, "92927") {
		t.Errorf("expected suffix '92927', got %q", truncated)
	}
}

func TestFormatLatestBranchName(t *testing.T) {
	backupBr := "backup/plan-29-auto-aliasing-92927"
	got := formatLatestBranchName(backupBr, 13)
	if !strings.HasPrefix(got, "...") || !strings.HasSuffix(got, "92927") {
		t.Errorf("expected '...92927', got %q", got)
	}

	depBr := "dependabot/go_modules_2026-09-19_175629"
	gotDep := formatLatestBranchName(depBr, 14)
	if !strings.HasPrefix(gotDep, "...") || !strings.HasSuffix(gotDep, "175629") {
		t.Errorf("expected '...175629', got %q", gotDep)
	}
}

func TestFormatPRCell(t *testing.T) {
	if formatPRCell("30 Open PRs") != "30" {
		t.Errorf("expected '30', got %q", formatPRCell("30 Open PRs"))
	}
	if formatPRCell("6 Open PRs") != "06" {
		t.Errorf("expected '06', got %q", formatPRCell("6 Open PRs"))
	}
	if formatPRCell("1 Open PR") != "01" {
		t.Errorf("expected '01', got %q", formatPRCell("1 Open PR"))
	}
	if formatPRCell("local") != "-" {
		t.Errorf("expected '-', got %q", formatPRCell("local"))
	}
	if formatPRCell("0 PRs") != "-" {
		t.Errorf("expected '-', got %q", formatPRCell("0 PRs"))
	}
}

func TestPullTableWideUserScreenshotSimulation(t *testing.T) {
	rows := []model.PullTableRow{
		{
			RepoName:     "Antigravity-Manager",
			Branch:       "main",
			LatestBranch: "v4.29.0",
			CommitRange:  "06479c6..7e0",
			Changes:      "+112/-1 (1)",
			PRStatus:     "30 Open PRs",
			PullStatus:   "FAST_FORWARD",
		},
		{
			RepoName:     "enum-v10",
			Branch:       "main",
			LatestBranch: "dependabot/go_modules_2026-09-19_175629",
			CommitRange:  "9ff44cf",
			Changes:      "up-to-date",
			PRStatus:     "6 Open PRs",
			PullStatus:   "UP_TO_DATE",
		},
		{
			RepoName:     "coding-guidelines-v24",
			Branch:       "main",
			LatestBranch: "backup/plan-29-auto-aliasing-92927",
			CommitRange:  "7972239..8b2",
			Changes:      "+304/-2 (2)",
			PRStatus:     "local",
			PullStatus:   "FAST_FORWARD",
		},
	}

	layout := NewPullTableLayoutWithWidth(rows, 100)
	layout.PrintHeader()
	for _, r := range rows {
		layout.PrintRow(r)
	}
}
