// Package cmd — pull_error_formatter.go formats pull error explanations.
package cmdpull

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func FormatPullDiagnosisMessage(repoName, repoPath string, d gitutil.DirtyDiagnosis) string {
	return fmt.Sprintf("  ⚠ %s is dirty (%s) — pull skipped to prevent merge conflicts.\n", repoName, d.SummaryReason)
}
