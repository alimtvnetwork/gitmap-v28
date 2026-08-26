package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestPullTableSuite(t *testing.T) {
	rows := []model.PullTableRow{
		{
			RepoName:   "scripts-fixer-v20",
			Branch:     "main",
			LastSHA:    "a1b2c3d",
			PRStatus:   "synced",
			PullStatus: "SUCCESS",
			Duration:   "1.0s",
			IsDirty:    false,
		},
		{
			RepoName:   "app-frontend",
			Branch:     "dev",
			LastSHA:    "e5f6g7h",
			PRStatus:   "ahead",
			PullStatus: "UP_TO_DATE",
			Duration:   "0.5s",
			IsDirty:    true,
		},
	}

	layout := NewPullTableLayout(rows)
	if layout.MaxRepo < len("scripts-fixer-v20") {
		t.Fatalf("expected MaxRepo >= 17, got %d", layout.MaxRepo)
	}

	RenderPullBatchTable(rows)
}
