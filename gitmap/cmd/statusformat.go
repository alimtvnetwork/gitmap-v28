package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// printOneStatus prints a single repo's status row or missing indicator.
func printOneStatus(rec model.ScanRecord, s *statusSummary) {
	_, err := os.Stat(rec.AbsolutePath)
	if err == nil {
		printRepoStatus(rec, s)

		return
	}

	printMissingRepo(rec.RepoName, s)
}

// printRepoStatus prints the status row for a repo that exists on disk.
func printRepoStatus(rec model.ScanRecord, s *statusSummary) {
	rs := gitutil.Status(rec.AbsolutePath)
	stateIcon := formatStateIcon(rs.Dirty, s)
	syncText := formatSyncText(rs.Ahead, rs.Behind, s)
	stashText := formatStashText(rs.StashCount, s)
	filesText := formatFileCounts(rs)
	branchText := fmt.Sprintf("%s%s%s", constants.ColorCyan, truncate(rs.Branch, 12), constants.ColorReset)

	fmt.Printf(constants.StatusRowFmt,
		truncate(rec.RepoName, 22),
		branchText, stateIcon, syncText, stashText, filesText)
}

// printMissingRepo prints a row for a repo not found on disk.
func printMissingRepo(name string, s *statusSummary) {
	truncated := truncate(name, 22)
	fmt.Printf(constants.StatusMissingFmt,
		constants.ColorDim, truncated,
		constants.ColorYellow, constants.ColorReset)
	s.Missing++
}

// formatStateIcon returns the clean/dirty indicator and updates summary.
func formatStateIcon(dirty bool, s *statusSummary) string {
	if dirty {
		s.Dirty++

		return constants.ColorYellow + constants.StatusIconDirty + constants.ColorReset
	}
	s.Clean++

	return constants.ColorGreen + constants.StatusIconClean + constants.ColorReset
}

// formatSyncText returns the ahead/behind indicator and updates summary.
func formatSyncText(ahead, behind int, s *statusSummary) string {
	if ahead > 0 && behind > 0 {
		s.Ahead++
		s.Behind++

		return fmt.Sprintf("%s"+constants.StatusSyncBothFmt+"%s", constants.ColorYellow, ahead, behind, constants.ColorReset)
	}

	return formatSyncSingle(ahead, behind, s)
}

// formatSyncSingle handles one-directional or no sync difference.
func formatSyncSingle(ahead, behind int, s *statusSummary) string {
	if ahead > 0 {
		s.Ahead++

		return fmt.Sprintf("%s"+constants.StatusSyncUpFmt+"%s", constants.ColorCyan, ahead, constants.ColorReset)
	}
	if behind > 0 {
		s.Behind++

		return fmt.Sprintf("%s"+constants.StatusSyncDownFmt+"%s", constants.ColorYellow, behind, constants.ColorReset)
	}

	return constants.ColorDim + constants.StatusSyncDash + constants.ColorReset
}

// formatStashText returns the stash indicator and updates summary.
func formatStashText(stashCount int, s *statusSummary) string {
	if stashCount > 0 {
		s.Stashed++

		return fmt.Sprintf("%s"+constants.StatusStashFmt+"%s", constants.ColorCyan, stashCount, constants.ColorReset)
	}

	return constants.ColorDim + constants.StatusDash + constants.ColorReset
}

// formatFileCounts returns staged/modified/untracked counts.
func formatFileCounts(rs gitutil.RepoStatus) string {
	if rs.Dirty {
		return buildFileCountParts(rs)
	}

	dash := constants.ColorDim + constants.StatusDash + constants.ColorReset

	return dash
}

// buildFileCountParts assembles the file count display parts.
func buildFileCountParts(rs gitutil.RepoStatus) string {
	parts := make([]string, 0, 3)
	if rs.Staged > 0 {
		parts = append(parts, fmt.Sprintf("%s"+constants.StatusStagedFmt+"%s", constants.ColorGreen, rs.Staged, constants.ColorReset))
	}
	if rs.Modified > 0 {
		parts = append(parts, fmt.Sprintf("%s"+constants.StatusModifiedFmt+"%s", constants.ColorYellow, rs.Modified, constants.ColorReset))
	}
	if rs.Untracked > 0 {
		parts = append(parts, fmt.Sprintf("%s"+constants.StatusUntrackedFmt+"%s", constants.ColorDim, rs.Untracked, constants.ColorReset))
	}

	return strings.Join(parts, constants.StatusFileCountSep)
}
