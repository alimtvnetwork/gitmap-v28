package cmdstatus

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// isLegacyDataError checks if an error indicates legacy UUID-format data.
func isLegacyDataError(err error) bool {
	return strings.Contains(err.Error(), "Scan error") ||
		strings.Contains(err.Error(), "converting driver.Value type string")
}

func openDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return store.OpenGlobalDefault()
	}
	_ = db.Migrate()
	if countRegisteredSSHHosts(db) == 0 {
		return resolveGlobalSSHDBFallback(db), nil
	}
	return db, nil
}

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

func countRegisteredSSHHosts(db *store.DB) int {
	var count int
	row := db.SQL().QueryRow("SELECT count(*) FROM ssh_hosts")
	if scanErr := row.Scan(&count); scanErr != nil {
		return 0
	}
	return count
}
func resolveGlobalSSHDBFallback(localDB *store.DB) *store.DB {
	globalDB, err := store.OpenGlobalDefault()
	if err != nil {
		return localDB
	}
	_ = globalDB.Migrate()
	if countRegisteredSSHHosts(globalDB) > 0 {
		localDB.Close()
		return globalDB
	}
	globalDB.Close()
	return localDB
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

// stripANSI removes ESC[...m color codes so JSON consumers get clean text.
func stripANSI(s string) string {
	var out strings.Builder
	out.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}

			i = j

			continue
		}

		out.WriteByte(s[i])
	}

	return out.String()
}

func visualWidth(s string) int {
	return runewidth.StringWidth(stripANSI(glyphs.FilterString(s)))
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

// printMissingRepoRemediation renders guidance when repos are missing on disk.
func printMissingRepoRemediation(c *statusTableContext) {
	missingSlugs := collectMissingSlugs(c)
	hasMissing := len(missingSlugs) > 0
	if !hasMissing {
		return
	}

	printMissingBanner(len(missingSlugs))
	printMissingItems(missingSlugs)
	printMissingRemediationSteps(missingSlugs)
}

func collectMissingSlugs(c *statusTableContext) []string {
	var missingSlugs []string
	if c == nil {
		return nil
	}

	for _, r := range c.Rows {
		if r.Missing {
			missingSlugs = append(missingSlugs, r.RepoName)
		}
	}

	return missingSlugs
}

func printMissingBanner(count int) {
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	fmt.Printf("  %s▲ %d missing repository(ies) detected:%s\n",
		constants.ColorYellow, count, constants.ColorReset)
}

func printMissingItems(slugs []string) {
	for _, slug := range slugs {
		fmt.Printf("     • %s%s%s\n", constants.ColorRed, slug, constants.ColorReset)
	}

	fmt.Println()
}

func printMissingRemediationSteps(slugs []string) {
	fmt.Println("  To resolve missing repositories:")
	fmt.Println("  1. Relocate to a new folder:")
	fmt.Println("     $ gitmap scan-folder update <slug> <new-path>")
	fmt.Println("  2. Untrack from database:")
	fmt.Printf("     $ gitmap rm %s\n", slugList(slugs))
	fmt.Println("  3. Clear or reset tracking database:")
	fmt.Println("     $ gitmap db-reset  (or: $ gitmap storage reset-errors)")
	fmt.Println()
}

func slugList(slugs []string) string {
	if len(slugs) > 3 {
		return fmt.Sprintf("%s %s ... (+%d more)", slugs[0], slugs[1], len(slugs)-2)
	}

	return strings.Join(slugs, " ")
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
