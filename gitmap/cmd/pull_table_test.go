package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestPullTableSuite(t *testing.T) {
	rows := []model.PullTableRow{
		{
			RepoName:     "scripts-fixer-v20",
			Branch:       "main",
			LatestBranch: "main",
			LastSHA:      "a1b2c3d",
			PRStatus:     "synced",
			PullStatus:   "SUCCESS",
			Duration:     "1.0s",
			IsDirty:      false,
		},
		{
			RepoName:     "app-frontend",
			Branch:       "dev",
			LatestBranch: "release/v1.35.0",
			LastSHA:      "e5f6g7h",
			PRStatus:     "ahead",
			PullStatus:   "UP_TO_DATE",
			Duration:     "0.5s",
			IsDirty:      true,
		},
	}

	layout := NewPullTableLayout(rows)
	isExpectedCap := layout.MaxRepo >= len("scripts-fixer-v20")
	if isExpectedCap == false {
		t.Fatalf("expected MaxRepo >= 17, got %d", layout.MaxRepo)
	}

	RenderPullBatchTable(rows)
}

func TestFormatBranchNamePrefixOmission(t *testing.T) {
	featureBranch := "feature/initial-migration"
	gotFeature := formatBranchName(featureBranch, 20)
	isExpectedFeature := gotFeature == "initial-migration"
	if isExpectedFeature == false {
		t.Fatalf("expected initial-migration, got %q", gotFeature)
	}

	releaseBranch := "release/v1.35.0"
	gotRelease := formatBranchName(releaseBranch, 16)
	isExpectedRelease := gotRelease == "v1.35.0"
	if isExpectedRelease == false {
		t.Fatalf("expected v1.35.0, got %q", gotRelease)
	}
}

func TestFormatBranchNameMiddleTruncate(t *testing.T) {
	longBranch := "dependabot/go_modules/go-minor-patch-0d16155629"
	formatted := formatBranchName(longBranch, 18)
	hasExpectedLen := len(formatted) == 18
	if hasExpectedLen == false {
		t.Fatalf("expected length 18, got %d (%q)", len(formatted), formatted)
	}
	hasEllipsis := strings.Contains(formatted, "...")
	if hasEllipsis == false {
		t.Fatalf("expected ellipsis in %q", formatted)
	}
	hasEndDigits := strings.HasSuffix(formatted, "55629")
	if hasEndDigits == false {
		t.Fatalf("expected ending 55629 in %q", formatted)
	}
}

func TestFormatRepoNameMiddleTruncate(t *testing.T) {
	longRepo := "ai-empathy-prompt-tuner-v1"
	formatted := formatRepoName(longRepo, 20)
	hasExpectedLen := len(formatted) == 20
	if hasExpectedLen == false {
		t.Fatalf("expected length 20, got %d (%q)", len(formatted), formatted)
	}
	hasEndSuffix := strings.HasSuffix(formatted, "er-v1")
	if hasEndSuffix == false {
		t.Fatalf("expected ending er-v1 in %q", formatted)
	}
}

func TestCalcAnsiPadding(t *testing.T) {
	plainString := "hello"
	plainPadding := calcAnsiPadding(plainString, 10)
	isPlainTen := plainPadding == 10
	if isPlainTen == false {
		t.Fatalf("expected plain padding 10, got %d", plainPadding)
	}

	coloredString := "\x1b[32mhello\x1b[0m"
	coloredPadding := calcAnsiPadding(coloredString, 10)
	hasExtraBytes := coloredPadding > 10
	if hasExtraBytes == false {
		t.Fatalf("expected colored padding > 10, got %d", coloredPadding)
	}
}

func TestPullTableUserScreenshotSimulation(t *testing.T) {
	screenshotRows := []model.PullTableRow{
		{
			RepoName:     "ai-empathy-prompt-tuner-v1",
			Branch:       "main",
			LatestBranch: "main",
			PRStatus:     "local",
			PullStatus:   "UP_TO_DATE",
			LastSHA:      "f1d94df",
			Duration:     "1.0s",
		},
		{
			RepoName:     "prompt-architect-v2",
			Branch:       "main",
			LatestBranch: "release/v1.35.0",
			PRStatus:     "local",
			PullStatus:   "UP_TO_DATE",
			LastSHA:      "03b798a",
			Duration:     "1.0s",
		},
		{
			RepoName:     "enum-v10",
			Branch:       "main",
			LatestBranch: "dependabot/go_modules/go-minor-patch-0d16155629",
			PRStatus:     "6 Open PRs",
			PullStatus:   "UP_TO_DATE",
			LastSHA:      "9ff44cf",
			Duration:     "1.0s",
		},
		{
			RepoName:     "pathhelper",
			Branch:       "develop",
			LatestBranch: "feature/initial-migration",
			PRStatus:     "local",
			PullStatus:   "UP_TO_DATE",
			LastSHA:      "5e7599b",
			Duration:     "1.0s",
		},
		{
			RepoName:     "strhelper",
			Branch:       "benchmark-examples-padding-repeat",
			LatestBranch: "training/viva-golang-v2",
			PRStatus:     "local",
			PullStatus:   "UP_TO_DATE",
			LastSHA:      "9b06402",
			Duration:     "1.0s",
		},
	}

	RenderPullBatchTable(screenshotRows)
}
