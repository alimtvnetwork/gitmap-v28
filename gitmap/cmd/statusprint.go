package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// printStatusBanner shows the dashboard header.
func printStatusBanner(count int) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.StatusBannerTop, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.StatusBannerTitle, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.StatusBannerBottom, constants.ColorReset)
	fmt.Println()
	fmt.Printf("  %s"+constants.StatusRepoCountFmt+"%s\n", constants.ColorDim, count, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermSeparator, constants.ColorReset)
	fmt.Println()
}

// printStatusTable prints each repo's status and returns a summary.
func printStatusTable(records []model.ScanRecord) statusSummary {
	s := statusSummary{Total: len(records)}
	printStatusHeader()

	for _, rec := range records {
		printOneStatus(rec, &s)
	}

	return s
}

// printStatusTableTracked prints each repo's status with batch progress tracking.
func printStatusTableTracked(records []model.ScanRecord, prog *cloner.BatchProgress) statusSummary {
	s := statusSummary{Total: len(records)}
	printStatusHeader()

	for _, rec := range records {
		prog.BeginItem(rec.RepoName)
		printOneStatus(rec, &s)
		prog.Succeed()
	}

	return s
}

// printStatusHeader prints the table column header row.
func printStatusHeader() {
	fmt.Printf(constants.StatusHeaderFmt,
		constants.ColorWhite,
		constants.StatusTableColumns[0], constants.StatusTableColumns[1],
		constants.StatusTableColumns[2], constants.StatusTableColumns[3],
		constants.StatusTableColumns[4], constants.StatusTableColumns[5],
		constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim,
		constants.TermTableRule, constants.ColorReset)
}

// printStatusSummary shows the final totals.
func printStatusSummary(s statusSummary) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	parts := buildSummaryParts(s)
	line := strings.Join(parts, constants.SummaryJoinSep)
	fmt.Printf("  %s\n\n", line)
}

// buildSummaryParts assembles summary line segments.
func buildSummaryParts(s statusSummary) []string {
	parts := []string{fmt.Sprintf(constants.SummaryReposFmt, s.Total)}
	parts = appendSummaryPart(parts, s.Clean, constants.ColorGreen, constants.SummaryCleanFmt)
	parts = appendSummaryPart(parts, s.Dirty, constants.ColorYellow, constants.SummaryDirtyFmt)
	parts = appendSummaryPart(parts, s.Ahead, constants.ColorCyan, constants.SummaryAheadFmt)
	parts = appendSummaryPart(parts, s.Behind, constants.ColorYellow, constants.SummaryBehindFmt)
	parts = appendSummaryPart(parts, s.Stashed, "", constants.SummaryStashedFmt)
	parts = appendSummaryPart(parts, s.Missing, constants.ColorYellow, constants.SummaryMissingFmt)

	return parts
}

// appendSummaryPart conditionally appends a colored summary segment.
func appendSummaryPart(parts []string, count int, color, format string) []string {
	if count == 0 {
		return parts
	}
	if len(color) > 0 {
		colored := fmt.Sprintf("%s"+format+"%s", color, count, constants.ColorReset)

		return append(parts, colored)
	}

	return append(parts, fmt.Sprintf(format, count))
}
