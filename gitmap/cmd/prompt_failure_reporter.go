// Package cmd — prompt_failure_reporter.go reports failed installations.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func ReportPromptFailures(results []model.PromptInstallResult) {
	for _, r := range results {
		if !r.IsSuccess {
			fmt.Fprintf(os.Stderr, "  ✖ %s: %s\n", r.RepoPath, r.Error)
		}
	}
}
