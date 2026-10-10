// Package cmdfixreleasetags provides diagnostic preview tables, confirmation
// prompts, and summary banners for release tag management.
package cmdfixreleasetags

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// StdinReaderOverride allows tests to supply a custom input reader.
var StdinReaderOverride io.Reader

// ForceInteractiveOverride allows tests to override interactive stdin detection.
var ForceInteractiveOverride *bool

// RenderDiagnosticTable renders the 6-column aligned diagnostic preview table.
func RenderDiagnosticTable(records []ReleaseTagAuditRecord) {
	cfg := BuildDiagnosticTableConfig(records)
	termout.PrintTable(cfg)
}

// BuildDiagnosticTableConfig constructs the termout.TableConfig for the preview table.
func BuildDiagnosticTableConfig(records []ReleaseTagAuditRecord) termout.TableConfig {
	columns := []termout.Column{
		{Title: "TAG", MinWidth: 10, MaxWidth: 18, Align: termout.AlignLeft},
		{Title: "COMMIT", MinWidth: 8, MaxWidth: 10, Align: termout.AlignLeft},
		{Title: "RELEASE STATUS", MinWidth: 16, MaxWidth: 22, Align: termout.AlignLeft},
		{Title: "ASSETS", MinWidth: 10, MaxWidth: 16, Align: termout.AlignLeft},
		{Title: "CI/CD STATUS", MinWidth: 14, MaxWidth: 18, Align: termout.AlignLeft},
		{Title: "ACTION", MinWidth: 20, MaxWidth: 28, Align: termout.AlignLeft},
	}

	tableRows := make([]termout.Row, 0, len(records))
	for _, rec := range records {
		commitDisplay := rec.CommitSha
		if len(commitDisplay) > 7 {
			commitDisplay = commitDisplay[:7]
		}
		if commitDisplay == "" {
			commitDisplay = "-"
		}

		tableRows = append(tableRows, termout.Row{
			Cells: []string{
				rec.Tag,
				commitDisplay,
				formatReleaseStatusCell(rec),
				formatAssetsCell(rec.AssetCount),
				formatCICDStatusCell(rec),
				formatActionCell(rec),
			},
		})
	}

	return termout.TableConfig{
		Columns:      columns,
		Rows:         tableRows,
		HeaderColor:  constants.ColorCyan,
		BorderColor:  constants.ColorDim,
		HasBorders:   true,
		EllipsisText: "...",
	}
}

func formatReleaseStatusCell(rec ReleaseTagAuditRecord) string {
	if !rec.HasGitHubRelease {
		return "Missing Release"
	}
	if rec.IsDraft {
		return "Draft"
	}
	if rec.AssetCount == 0 && rec.AuditReason == ReasonMissingAssets {
		return "Published (Stale)"
	}

	return "Published"
}

func formatAssetsCell(assetCount int) string {
	if assetCount == 0 {
		return "0 Assets"
	}

	return fmt.Sprintf("%d Assets", assetCount)
}

func formatCICDStatusCell(rec ReleaseTagAuditRecord) string {
	if len(rec.WorkflowRuns) == 0 {
		return "No CI"
	}

	latest := rec.WorkflowRuns[0]
	if latest.IsSuccessful || latest.Conclusion == "success" {
		return "Passed"
	}
	if latest.IsInProgress || strings.ToLower(latest.Status) == "in_progress" || strings.ToLower(latest.Status) == "queued" {
		return "In-Progress"
	}
	if latest.Conclusion == "failure" || rec.AuditReason == ReasonCICDFailed {
		return "Failed"
	}

	return "Unknown"
}

func formatActionCell(rec ReleaseTagAuditRecord) string {
	if rec.IsEligibleForDeletion {
		return "Delete Release & Tag"
	}

	if rec.IsProtected {
		switch rec.ProtectionStatus {
		case StatusProtectedActive:
			return "Skip (Active Version)"
		case StatusProtectedLatest:
			return "Skip (Latest Healthy)"
		case StatusProtectedGrace:
			return "Skip (Grace Period)"
		default:
			return "Skip (Protected)"
		}
	}

	if rec.AuditReason == ReasonHealthy {
		return "Keep (Healthy Release)"
	}

	return "Keep"
}

// IsInteractiveStdin checks whether standard input is attached to a terminal.
func IsInteractiveStdin() bool {
	if ForceInteractiveOverride != nil {
		return *ForceInteractiveOverride
	}

	if os.Getenv("CI") != "" || os.Getenv("GITMAP_NON_INTERACTIVE") == "1" {
		return false
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != 0
}

func isInteractiveStdin() bool {
	return IsInteractiveStdin()
}

// ConfirmDeletion prompts user to confirm deletion unless already confirmed or non-interactive.
func ConfirmDeletion(candidateCount int, isConfirmed bool) (bool, error) {
	if isConfirmed {
		return true, nil
	}

	if !isInteractiveStdin() && StdinReaderOverride == nil {
		return false, fmt.Errorf("E1025: stdin is non-interactive; re-run with -y/--yes to confirm deletion")
	}

	fmt.Printf("\n%sDelete %d orphan/broken release(s) and tag(s)? [y/N]: %s",
		constants.ColorYellow, candidateCount, constants.ColorReset)

	var reader *bufio.Reader
	if StdinReaderOverride != nil {
		reader = bufio.NewReader(StdinReaderOverride)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}

	input, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("failed to read confirmation input: %w", err)
	}

	trimmed := strings.ToLower(strings.TrimSpace(input))
	if trimmed == "y" || trimmed == "yes" {
		return true, nil
	}

	return false, nil
}

func confirmDeletion(candidateCount int, isConfirmed bool) (bool, error) {
	return ConfirmDeletion(candidateCount, isConfirmed)
}

// RenderHeaderBanner prints the initial audit scanning banner.
func RenderHeaderBanner(targetDir string) {
	fmt.Printf("%s[RELEASE-TAG-AUDIT]%s Inspecting release tags in '%s'...\n",
		constants.ColorCyan, constants.ColorReset, targetDir)
}

// RenderSummaryBanner prints a clean summary of deletion or audit outcome.
func RenderSummaryBanner(totalAudited, brokenFound, deletedCount int, isDryRun bool) {
	fmt.Println()
	if isDryRun {
		fmt.Printf("%s[dry-run]%s Scanned %d tag(s), found %d orphan/broken candidate(s). No actions taken.\n",
			constants.ColorYellow, constants.ColorReset, totalAudited, brokenFound)
		return
	}

	if deletedCount > 0 {
		fmt.Printf("%s✓%s Successfully deleted %d orphan/broken release tag(s).\n",
			constants.ColorGreen, constants.ColorReset, deletedCount)
	} else if brokenFound == 0 {
		fmt.Printf("%s✓%s All %d release tag(s) are healthy.\n",
			constants.ColorGreen, constants.ColorReset, totalAudited)
	}
}
