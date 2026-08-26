package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestPromptDashboardSuite(t *testing.T) {
	layout := NewPromptStatusTableLayout()
	layout.PrintHeader()
	layout.PrintRow("/tmp/my-repo", model.PromptArchitectMetadata{
		Version:     "v2.0.0",
		InstalledAt: "2026-08-26T17:30:00Z",
		Status:      "active",
	})

	results := []model.PromptInstallResult{
		{
			RepoPath:  "/tmp/my-repo",
			IsSuccess: true,
			Version:   "v2.0.0",
		},
	}
	RenderPromptInstallSummary(results)
	ReportPromptFailures(results)
}
