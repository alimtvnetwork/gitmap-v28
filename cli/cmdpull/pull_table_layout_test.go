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
	if isWide == false {
		t.Fatalf("expected wide layout for width 120")
	}

	isRepoValid := layout.MaxRepo > 0
	if isRepoValid == false {
		t.Fatalf("expected positive MaxRepo in wide layout")
	}

	isBranchValid := layout.MaxBranch > 0
	if isBranchValid == false {
		t.Fatalf("expected positive MaxBranch in wide layout")
	}

	isLatestValid := layout.MaxLatestBr > 0
	if isLatestValid == false {
		t.Fatalf("expected positive MaxLatestBr in wide layout")
	}

	if layout.MaxRelease <= 0 {
		t.Fatalf("expected positive MaxRelease in wide layout")
	}

	if layout.MaxSHA <= 0 {
		t.Fatalf("expected positive MaxSHA in wide layout")
	}

	isDividerClean := layout.DividerLen <= 120
	if isDividerClean == false {
		t.Fatalf("expected divider <= 120 for neat layout, got %d", layout.DividerLen)
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
	isCompact := layout.IsWide == false
	if isCompact == false {
		t.Fatalf("expected compact layout for width 80")
	}

	isDividerWithinLimit := layout.DividerLen+2 <= 80
	if isDividerWithinLimit == false {
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

func TestPullTableUltraWideWidth198(t *testing.T) {
	rows := []model.PullTableRow{{RepoName: "gitmap", Branch: "main"}}
	layout := NewPullTableLayoutWithWidth(rows, 198)
	if !layout.IsWide {
		t.Fatalf("expected wide layout")
	}

	if layout.MaxRepo != 71 {
		t.Fatalf("expected MaxRepo=71, got %d", layout.MaxRepo)
	}

	if layout.MaxLatestBr != 53 {
		t.Fatalf("expected MaxLatestBr=53, got %d", layout.MaxLatestBr)
	}

	if layout.MaxSHA != 8 {
		t.Fatalf("expected MaxSHA=8, got %d", layout.MaxSHA)
	}
}
