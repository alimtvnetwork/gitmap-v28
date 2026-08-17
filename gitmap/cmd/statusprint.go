package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
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

type statusRow struct {
	Missing   bool
	RepoName  string
	Branch    string
	StateIcon string
	SyncText  string
	StashText string
	FilesText string
}

type statusTableContext struct {
	Rows      []statusRow
	MaxRepo   int
	MaxBranch int
	MaxStatus int
	MaxSync   int
	MaxStash  int
	MaxFiles  int
}

func newStatusTableContext() *statusTableContext {
	return &statusTableContext{
		Rows:      make([]statusRow, 0),
		MaxRepo:   len(constants.StatusTableColumns[0]),
		MaxBranch: len(constants.StatusTableColumns[1]),
		MaxStatus: len(constants.StatusTableColumns[2]),
		MaxSync:   len(constants.StatusTableColumns[3]),
		MaxStash:  len(constants.StatusTableColumns[4]),
		MaxFiles:  len(constants.StatusTableColumns[5]),
	}
}

func (c *statusTableContext) addRow(r statusRow) {
	c.Rows = append(c.Rows, r)
	if l := len(r.RepoName); l > c.MaxRepo {
		c.MaxRepo = l
	}
	if !r.Missing {
		if l := len(r.Branch); l > c.MaxBranch {
			c.MaxBranch = l
		}
		if l := len(stripANSI(r.StateIcon)); l > c.MaxStatus {
			c.MaxStatus = l
		}
		if l := len(stripANSI(r.SyncText)); l > c.MaxSync {
			c.MaxSync = l
		}
		if l := len(stripANSI(r.StashText)); l > c.MaxStash {
			c.MaxStash = l
		}
		if l := len(stripANSI(r.FilesText)); l > c.MaxFiles {
			c.MaxFiles = l
		}
	}
}

// printStatusTable prints each repo's status and returns a summary.
func printStatusTable(records []model.ScanRecord) statusSummary {
	s := statusSummary{Total: len(records)}
	
	tableCtx := newStatusTableContext()
	for _, rec := range records {
		row := computeOneStatus(rec, &s)
		tableCtx.addRow(row)
	}

	printStatusTableWithContext(tableCtx)
	return s
}

// printStatusTableTracked prints each repo's status with batch progress tracking.
func printStatusTableTracked(records []model.ScanRecord, prog *cloner.BatchProgress) statusSummary {
	s := statusSummary{Total: len(records)}
	
	tableCtx := newStatusTableContext()
	for _, rec := range records {
		prog.BeginItem(rec.RepoName)
		row := computeOneStatus(rec, &s)
		tableCtx.addRow(row)
		prog.Succeed(rec.RepoName)
	}

	printStatusTableWithContext(tableCtx)
	return s
}

func printStatusTableWithContext(c *statusTableContext) {
	const colGap = "   "

	fmt.Printf("  %s%-*s%s%-*s%s%-*s%s%-*s%s%-*s%s%-*s%s\n",
		constants.ColorWhite,
		c.MaxRepo, constants.StatusTableColumns[0], colGap,
		c.MaxBranch, constants.StatusTableColumns[1], colGap,
		c.MaxStatus, constants.StatusTableColumns[2], colGap,
		c.MaxSync, constants.StatusTableColumns[3], colGap,
		c.MaxStash, constants.StatusTableColumns[4], colGap,
		c.MaxFiles, constants.StatusTableColumns[5],
		constants.ColorReset)
	
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)

	for _, r := range c.Rows {
		if r.Missing {
			fmt.Printf("  %s%-*s%s%s⊘ not found%s\n",
				constants.ColorDim, c.MaxRepo, r.RepoName, colGap,
				constants.ColorYellow, constants.ColorReset)
			continue
		}
		
		branchStr := fmt.Sprintf("%s%s%s", constants.ColorCyan, r.Branch, constants.ColorReset)

		padBranch := c.MaxBranch + (len(branchStr) - len(stripANSI(branchStr)))
		padStatus := c.MaxStatus + (len(r.StateIcon) - len(stripANSI(r.StateIcon)))
		padSync := c.MaxSync + (len(r.SyncText) - len(stripANSI(r.SyncText)))
		padStash := c.MaxStash + (len(r.StashText) - len(stripANSI(r.StashText)))
		padFiles := c.MaxFiles + (len(r.FilesText) - len(stripANSI(r.FilesText)))

		fmt.Printf("  %-*s%s%-*s%s%-*s%s%-*s%s%-*s%s%-*s\n",
			c.MaxRepo, r.RepoName, colGap,
			padBranch, branchStr, colGap,
			padStatus, r.StateIcon, colGap,
			padSync, r.SyncText, colGap,
			padStash, r.StashText, colGap,
			padFiles, r.FilesText)
	}
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
