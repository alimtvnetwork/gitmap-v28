package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// computeOneStatus returns a single repo's status row or missing indicator.
func computeOneStatus(rec model.ScanRecord, s *statusSummary) statusRow {
	_, err := os.Stat(rec.AbsolutePath)
	if err == nil {
		return computeRepoStatus(rec, s)
	}

	s.Missing++
	return statusRow{
		Missing:  true,
		RepoName: rec.RepoName,
	}
}

// computeRepoStatus returns the status row for a repo that exists on disk.
func computeRepoStatus(rec model.ScanRecord, s *statusSummary) statusRow {
	rs := gitutil.Status(rec.AbsolutePath)
	
	return statusRow{
		Missing:   false,
		RepoName:  rec.RepoName,
		Branch:    rs.Branch,
		StateIcon: formatStateIcon(rs.Dirty, s),
		SyncText:  formatSyncText(rs.Ahead, rs.Behind, s),
		StashText: formatStashText(rs.StashCount, s),
		FilesText: formatFileCounts(rs),
	}
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
