package cmd

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// printWatchDashboard clears the screen and prints one refresh cycle.
func printWatchDashboard(records []model.ScanRecord, interval int, noFetch bool) {
	snapshots := collectAllStatuses(records, noFetch)

	fmt.Print(constants.WatchClearScreen)
	printWatchBanner(interval)
	printWatchHeader()
	printWatchRows(snapshots)
	printWatchFooter(snapshots)
}

// printWatchBanner prints the header with refresh info.
func printWatchBanner(interval int) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.WatchBannerTop, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.WatchBannerTitle, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.WatchBannerBottom, constants.ColorReset)
	fmt.Println()

	now := time.Now().UTC().Format(constants.DateDisplayLayout + " " + constants.DateUTCSuffix)
	fmt.Printf("  %s"+constants.WatchRefreshFmt+"%s\n", constants.ColorDim, interval, constants.ColorReset)
	fmt.Printf("  %s"+constants.WatchLastUpdFmt+"%s\n\n", constants.ColorDim, now, constants.ColorReset)
}

// printWatchHeader prints the table column headers.
func printWatchHeader() {
	fmt.Printf(constants.WatchHeaderFmt,
		constants.ColorWhite,
		constants.WatchTableColumns[0], constants.WatchTableColumns[1],
		constants.WatchTableColumns[2], constants.WatchTableColumns[3],
		constants.WatchTableColumns[4], constants.WatchTableColumns[5],
		constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim,
		constants.TermTableRule, constants.ColorReset)
}

// printWatchRows prints each repo's status row.
func printWatchRows(snapshots []watchSnapshot) {
	for _, snap := range snapshots {
		printWatchRow(snap)
	}
}

// printWatchRow prints a single repo row.
func printWatchRow(snap watchSnapshot) {
	if snap.Status == "error" {
		fmt.Printf(constants.WatchErrorRowFmt,
			constants.ColorDim, truncate(snap.Name, 22),
			constants.ColorYellow, constants.ColorReset)

		return
	}

	statusText := formatWatchStatus(snap.Status)
	branchText := fmt.Sprintf("%s%s%s", constants.ColorCyan, truncate(snap.Branch, 16), constants.ColorReset)
	aheadText := formatWatchCount(snap.Ahead, constants.ColorCyan)
	behindText := formatWatchCount(snap.Behind, constants.ColorYellow)
	stashText := formatWatchCount(snap.Stash, constants.ColorCyan)

	fmt.Printf(constants.WatchRowFmt,
		truncate(snap.Name, 22), statusText, branchText,
		aheadText, behindText, stashText)
}

// formatWatchStatus returns a colored status string.
func formatWatchStatus(status string) string {
	if status == "dirty" {
		return constants.ColorYellow + constants.StatusIconDirty + constants.ColorReset
	}

	return constants.ColorGreen + constants.StatusIconClean + constants.ColorReset
}

// formatWatchCount returns a colored count or a dash.
func formatWatchCount(count int, color string) string {
	if count > 0 {
		return fmt.Sprintf("%s%d%s", color, count, constants.ColorReset)
	}

	return constants.ColorDim + constants.StatusDash + constants.ColorReset
}

// printWatchFooter prints the summary line.
func printWatchFooter(snapshots []watchSnapshot) {
	summary := buildWatchSummary(snapshots)

	fmt.Printf("\n  %s%s%s\n", constants.ColorDim,
		constants.TermTableRule, constants.ColorReset)
	fmt.Printf("  "+constants.WatchSummaryFmt+"\n\n",
		summary.Total, summary.Dirty, summary.Behind, summary.Stash)
}
