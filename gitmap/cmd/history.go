package cmd

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// runHistory handles the "history" subcommand.
func runHistory(args []string) error {
	checkHelp("history", args)
	detail, cmdFilter, limit, jsonOut := parseHistoryFlags(args)
	records := loadHistory(cmdFilter)
	records = applyHistoryLimit(records, limit)

	if jsonOut {
		printHistoryJSON(records)
		return nil
	}

	printHistoryTerminal(records, detail)
	return nil
}

// parseHistoryFlags parses --detail, --command, --limit, --json flags.
func parseHistoryFlags(args []string) (string, string, int, bool) {
	fs := flag.NewFlagSet(constants.CmdHistory, flag.ExitOnError)
	detail := fs.String("detail", constants.DetailStandard, constants.FlagDescDetail)
	command := fs.String("command", "", constants.FlagDescCommand)
	limit := fs.Int("limit", 0, constants.FlagDescLimit)
	jsonFlag := fs.Bool("json", false, constants.FlagDescLBJSON)
	fs.Parse(args)

	return *detail, *command, *limit, *jsonFlag
}

// loadHistory fetches history from the database.
func loadHistory(cmdFilter string) []model.CommandHistoryRecord {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrHistoryQuery+"\n", err)
		cliexit.HandleError(nil, 1)
	}
	defer db.Close()

	records, err := queryHistoryRecords(db, cmdFilter)
	if err != nil {
		handleHistoryError(err)
	}
	return records
}

func queryHistoryRecords(db *store.DB, cmdFilter string) ([]model.CommandHistoryRecord, error) {
	if cmdFilter != "" {
		return db.ListHistoryByCommand(cmdFilter)
	}
	return db.ListHistory()
}

// handleHistoryError processes errors returned from the history query.
func handleHistoryError(err error) {
	if isLegacyDataError(err) {
		fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
		cliexit.HandleError(nil, 1)
	}
	fmt.Fprintf(os.Stderr, constants.ErrHistoryQuery+"\n", err)
	cliexit.HandleError(nil, 1)
}

// applyHistoryLimit truncates results to the given limit.
func applyHistoryLimit(records []model.CommandHistoryRecord, limit int) []model.CommandHistoryRecord {
	if limit > 0 && limit < len(records) {
		return records[:limit]
	}

	return records
}

// printHistoryTerminal prints history in table format based on detail level,
// followed by a "Revert points" section enumerating undo commands for any
// row whose Command has a known inverse.
func printHistoryTerminal(records []model.CommandHistoryRecord, detail string) {
	if len(records) == 0 {
		fmt.Print(constants.MsgHistoryEmpty)
		return
	}

	printHistoryHeader(detail)
	for _, r := range records {
		printHistoryRow(r, detail)
	}
	printHistoryRevertSection(records)
}

// printHistoryHeader prints the colored column header. "LAST" (relative
// time) is always the right-most column so the eye lands on "when did
// this happen" without scanning past durations + flags.
func printHistoryHeader(detail string) {
	c := constants.ColorMagenta
	r := constants.ColorReset
	switch detail {
	case constants.DetailBasic:
		fmt.Printf("%s%-16s %-8s %s%s\n", c, "COMMAND", "STATUS", "LAST", r)
	case constants.DetailDetailed:
		fmt.Printf("%s%-16s %-18s %-22s %-8s %-10s %-6s %-30s %s%s\n",
			c, "COMMAND", "ARGS", "FLAGS", "STATUS", "DURATION", "REPOS", "SUMMARY", "LAST", r)
	default:
		fmt.Printf("%s%-16s %-22s %-8s %-10s %s%s\n",
			c, "COMMAND", "FLAGS", "STATUS", "DURATION", "LAST", r)
	}
}

type historyRowTokens struct {
	cmd    string
	flags  string
	status string
	dur    string
	last   string
}

func formatHistoryRowTokens(r model.CommandHistoryRecord) historyRowTokens {
	return historyRowTokens{
		cmd:    colorize(constants.ColorCyan, padRight(r.Command, 16)),
		flags:  colorize(constants.ColorDim, padRight(truncateHist(r.Flags, 22), 22)),
		status: colorizedStatus(r.ExitCode),
		dur:    colorize(constants.ColorYellow, padRight(strconv.FormatInt(r.DurationMs, 10)+"ms", 10)),
		last:   colorize(constants.ColorDim, relativeHistoryTime(r)),
	}
}

// printHistoryRow prints one row at the chosen detail level with ANSI
// colors: cyan command, dim flags, green OK / red FAIL, yellow
// duration, dim relative-time on the right.
func printHistoryRow(r model.CommandHistoryRecord, detail string) {
	tok := formatHistoryRowTokens(r)
	if detail == constants.DetailDetailed {
		printDetailedHistoryRow(tok, r)
		return
	}
	if detail == constants.DetailBasic {
		fmt.Printf("%s %s %s\n", tok.cmd, tok.status, tok.last)
		return
	}
	fmt.Printf("%s %s %s %s %s\n", tok.cmd, tok.flags, tok.status, tok.dur, tok.last)
}

func printDetailedHistoryRow(tok historyRowTokens, r model.CommandHistoryRecord) {
	args := colorize(constants.ColorWhite, padRight(truncateHist(r.Args, 18), 18))
	repos := padRight(strconv.Itoa(r.RepoCount), 6)
	summary := padRight(truncateHist(r.Summary, 30), 30)
	fmt.Printf("%s %s %s %s %s %s %s %s\n",
		tok.cmd, args, tok.flags, tok.status, tok.dur, repos, summary, tok.last)
}

// colorizedStatus renders an 8-wide colored OK / FAIL token.
func colorizedStatus(code int) string {
	if code == 0 {
		return colorize(constants.ColorGreen, padRight("✓ "+constants.MsgHistoryStatusOK, 8))
	}
	return colorize(constants.ColorRed, padRight("✗ "+constants.MsgHistoryStatusFail, 8))
}

// printHistoryJSON outputs history as stable JSON via the encoder
// in historyrender.go.
func printHistoryJSON(records []model.CommandHistoryRecord) {
	if err := encodeHistoryJSON(os.Stdout, records); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to encode history to JSON: %v\n", err)
	}
}

type revertRow struct {
	idx     int
	command string
	when    string
	hint    string
}

func collectRevertRows(records []model.CommandHistoryRecord) []revertRow {
	var hints []revertRow
	for i, r := range records {
		if h := revertHintFor(r); h != "" {
			hints = append(hints, revertRow{idx: i + 1, command: r.Command, when: relativeHistoryTime(r), hint: h})
		}
	}
	return hints
}

// printHistoryRevertSection enumerates revert commands for every row
// whose Command has a known inverse. Rows with no known revert are
// omitted so the section stays scannable. Empty when nothing is
// revertable (no header is printed in that case).
func printHistoryRevertSection(records []model.CommandHistoryRecord) {
	hints := collectRevertRows(records)
	if len(hints) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(colorize(constants.ColorMagenta, "Revert points"))
	fmt.Println(colorize(constants.ColorDim, "  Run the suggested command to undo the referenced state."))
	renderRevertRows(hints)
}

func renderRevertRows(hints []revertRow) {
	for _, h := range hints {
		fmt.Printf("  %s#%-3d%s %s%-16s%s %s%-18s%s  %s%s%s\n",
			constants.ColorDim, h.idx, constants.ColorReset,
			constants.ColorCyan, h.command, constants.ColorReset,
			constants.ColorDim, h.when, constants.ColorReset,
			constants.ColorYellow, h.hint, constants.ColorReset)
	}
}

var staticRevertHints = map[string]string{
	constants.CmdFixRepo:        "gitmap undo                  # restore latest fix-repo snapshot",
	"fix-repo-pub":              "gitmap undo                  # restore latest fix-repo snapshot",
	"fr":                        "gitmap undo                  # restore latest fix-repo snapshot",
	"frp":                       "gitmap undo                  # restore latest fix-repo snapshot",
	constants.CmdMakePublic:     "gitmap make-private",
	"mapub":                     "gitmap make-private",
	constants.CmdMakePrivate:    "gitmap make-public --yes",
	"mapri":                     "gitmap make-public --yes",
	constants.CmdMakeAllPublic:  "gitmap make-all-private",
	constants.CmdMakeAllPrivate: "gitmap make-all-public --yes",
}

// revertHintFor maps a history Command to a concrete inverse command
// the user can run to roll back. Returns "" when no inverse is known.
// Kept tiny + table-driven so adding a new revertable command is a
// one-line change.
func revertHintFor(r model.CommandHistoryRecord) string {
	if hint, ok := staticRevertHints[r.Command]; ok {
		return hint
	}
	if r.Command == "reclone-transport" {
		return recloneTransportHint(r)
	}
	return ""
}

func recloneTransportHint(r model.CommandHistoryRecord) string {
	if strings.Contains(r.Flags, "transport=ssh") {
		return "gitmap cfr " + r.Args + " --https"
	}
	if strings.Contains(r.Flags, "transport=https") {
		return "gitmap cfr " + r.Args + " --ssh"
	}
	return ""
}

// relativeHistoryTime renders the "how long ago did this finish" suffix
// shown in the right-most column. Falls back through FinishedAt →
// StartedAt → CreatedAt, returning "—" when every field is empty or
// unparseable so the column never collapses to whitespace.
func relativeHistoryTime(r model.CommandHistoryRecord) string {
	for _, s := range []string{r.FinishedAt, r.StartedAt, r.CreatedAt} {
		if t, ok := parseHistoryTime(s); ok {
			return humanizeDuration(time.Since(t)) + " ago"
		}
	}
	return "—"
}

// parseHistoryTime accepts both RFC3339 (gitmap-emitted) and SQLite
// CURRENT_TIMESTAMP ("YYYY-MM-DD HH:MM:SS") shapes.
func parseHistoryTime(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	layouts := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05"}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// humanizeDuration renders durations as the largest single unit
// ("3m", "2h", "5d") — matches the project's existing relative-time
// convention used in the usage footer's "X minutes ago" rows.
func humanizeDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

// padRight + truncate + colorize are tiny terminal-formatting helpers
// kept local to history rendering so the colored width math stays
// readable above. ANSI escapes are NOT counted by fmt's `%-Ns` width
// directive, so we pad BEFORE colorizing.
func padRight(s string, w int) string {
	if len(s) >= w {
		return s
	}
	return s + strings.Repeat(" ", w-len(s))
}

func truncateHist(s string, w int) string {
	if len(s) <= w {
		return s
	}
	if w <= 1 {
		return s[:w]
	}
	return s[:w-1] + "…"
}

func colorize(color, s string) string {
	if color == "" {
		return s
	}
	return color + s + constants.ColorReset
}
