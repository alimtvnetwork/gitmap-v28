package cmd

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
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

func visualWidth(s string) int {
	return runewidth.StringWidth(stripANSI(s))
}

func (c *statusTableContext) addRow(r statusRow) {
	c.Rows = append(c.Rows, r)
	if l := visualWidth(r.RepoName); l > c.MaxRepo {
		c.MaxRepo = l
	}

	if r.Missing {
		return
	}

	c.updateColumnWidths(r)
}

func (c *statusTableContext) updateColumnWidths(r statusRow) {
	if l := visualWidth(r.Branch); l > c.MaxBranch {
		c.MaxBranch = l
	}
	if l := visualWidth(r.StateIcon); l > c.MaxStatus {
		c.MaxStatus = l
	}
	if l := visualWidth(r.SyncText); l > c.MaxSync {
		c.MaxSync = l
	}
	if l := visualWidth(r.StashText); l > c.MaxStash {
		c.MaxStash = l
	}
	if l := visualWidth(r.FilesText); l > c.MaxFiles {
		c.MaxFiles = l
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

func printStatusTableHeader(c *statusTableContext) {
	const colGap = "   "
	dividerLen := calculateStatusDividerLen(c)

	fmt.Printf("  %s%s%s%s%s%s%s%s%s%s%s%s%s\n",
		constants.ColorWhite,
		cmdpull.PadVisual(constants.StatusTableColumns[0], c.MaxRepo), colGap,
		cmdpull.PadVisual(constants.StatusTableColumns[1], c.MaxBranch), colGap,
		cmdpull.PadVisual(constants.StatusTableColumns[2], c.MaxStatus), colGap,
		cmdpull.PadVisual(constants.StatusTableColumns[3], c.MaxSync), colGap,
		cmdpull.PadVisual(constants.StatusTableColumns[4], c.MaxStash), colGap,
		cmdpull.PadVisual(constants.StatusTableColumns[5], c.MaxFiles),
		constants.ColorReset)

	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("-", dividerLen), constants.ColorReset)
}

func calculateStatusDividerLen(c *statusTableContext) int {
	return c.MaxRepo + c.MaxBranch + c.MaxStatus + c.MaxSync + c.MaxStash + c.MaxFiles + (5 * 3)
}

func printStatusTableRow(c *statusTableContext, r statusRow, pastelColor string) {
	const colGap = "   "
	if r.Missing {
		fmt.Printf("  %s%s%s%s%s✖ not found%s\n",
			pastelColor, cmdpull.PadVisual(r.RepoName, c.MaxRepo), constants.ColorReset, colGap,
			constants.ColorRed, constants.ColorReset)

		return
	}

	branchStr := fmt.Sprintf("%s%s%s", constants.ColorCyan, r.Branch, constants.ColorReset)
	repoStr := fmt.Sprintf("%s%s%s", pastelColor, r.RepoName, constants.ColorReset)

	fmt.Printf("  %s%s%s%s%s%s%s%s%s%s%s\n",
		cmdpull.PadVisual(repoStr, c.MaxRepo), colGap,
		cmdpull.PadVisual(branchStr, c.MaxBranch), colGap,
		cmdpull.PadVisual(r.StateIcon, c.MaxStatus), colGap,
		cmdpull.PadVisual(r.SyncText, c.MaxSync), colGap,
		cmdpull.PadVisual(r.StashText, c.MaxStash), colGap,
		cmdpull.PadVisual(r.FilesText, c.MaxFiles))
}

func printStatusTableWithContext(c *statusTableContext) {
	printStatusTableHeader(c)
	for i, r := range c.Rows {
		pastelColor := constants.ColorCycle[i%len(constants.ColorCycle)]
		printStatusTableRow(c, r, pastelColor)
	}

	printMissingRepoRemediation(c)
}

// printStatusSummary shows the final totals.
func printStatusSummary(s statusSummary) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	parts := buildSummaryParts(s)
	line := strings.Join(parts, constants.SummaryJoinSep)
	fmt.Printf("  %s\n", line)
	printDirtySummaryTip(s.Dirty)
	fmt.Println()
}

func printDirtySummaryTip(dirtyCount int) {
	if dirtyCount <= 0 {
		return
	}

	fmt.Printf("  %sTip: %d repository(ies) dirty. Run 'gitmap fix' or 'gitmap pull --fix' to remediate.%s\n",
		constants.ColorDim, dirtyCount, constants.ColorReset)
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
