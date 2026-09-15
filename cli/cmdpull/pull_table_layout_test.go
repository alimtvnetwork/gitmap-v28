package cmdpull

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestPullTableLayoutColumns(t *testing.T) {
	rows := []model.PullTableRow{
		{
			RepoName:     "gitmap",
			Branch:       "main",
			LatestBranch: "main",
			CommitRange:  "a1b2c3d..e5f6a7b",
			Changes:      "+10/-2 (3)",
			PRStatus:     "synced",
			PullStatus:   "FAST_FORWARD",
			Duration:     "0.8s",
		},
	}

	layout := NewPullTableLayoutWithWidth(rows, 120)
	isWide := layout.IsWide
	if !isWide {
		t.Fatalf("expected wide layout for width 120")
	}

	isRangeValid := layout.MaxRange > 0
	if !isRangeValid {
		t.Fatalf("expected positive MaxRange in wide layout")
	}

	isChangesValid := layout.MaxChanges > 0
	if !isChangesValid {
		t.Fatalf("expected positive MaxChanges in wide layout")
	}

	layout.PrintHeader()
	layout.PrintRow(rows[0])
}

func TestPullTableLayoutCompactColumns(t *testing.T) {
	rows := []model.PullTableRow{
		{
			RepoName:     "gitmap",
			Branch:       "main",
			LatestBranch: "main",
			CommitRange:  "a1b2c3d",
			Changes:      "up-to-date",
			PullStatus:   "UP_TO_DATE",
			Duration:     "0.5s",
		},
	}

	layout := NewPullTableLayoutWithWidth(rows, 80)
	isCompact := !layout.IsWide
	if !isCompact {
		t.Fatalf("expected compact layout for width 80")
	}

	isDividerWithinLimit := layout.DividerLen+2 <= 80
	if !isDividerWithinLimit {
		t.Fatalf("expected divider <= 80, got %d", layout.DividerLen+2)
	}

	layout.PrintHeader()
	layout.PrintRow(rows[0])
}

func TestResolveRowCommitRangeFallback(t *testing.T) {
	rowWithRange := model.PullTableRow{
		CommitRange: "a1b2c3d..def5678",
		LastSHA:     "def5678",
	}
	resRange := resolveRowCommitRange(rowWithRange, 20)
	if resRange != "a1b2c3d..def5678" {
		t.Fatalf("expected full range, got %q", resRange)
	}

	rowWithSHAOnly := model.PullTableRow{
		LastSHA: "def5678",
	}
	resSHA := resolveRowCommitRange(rowWithSHAOnly, 15)
	if resSHA != "def5678" {
		t.Fatalf("expected fallback to LastSHA, got %q", resSHA)
	}
}

func TestResolveRowChangesFallback(t *testing.T) {
	rowWithChanges := model.PullTableRow{
		Changes: "+5/-1",
	}
	resChanges := resolveRowChanges(rowWithChanges, 10)
	if resChanges != "+5/-1" {
		t.Fatalf("expected +5/-1, got %q", resChanges)
	}

	rowEmpty := model.PullTableRow{}
	resEmpty := resolveRowChanges(rowEmpty, 10)
	if resEmpty != "-" {
		t.Fatalf("expected '-', got %q", resEmpty)
	}
}
