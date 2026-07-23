package cmd

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// runHistory handles the "history" subcommand.
func runHistory(args []string) {
	checkHelp("history", args)
	detail, cmdFilter, limit, jsonOut := parseHistoryFlags(args)
	records := loadHistory(cmdFilter)
	records = applyHistoryLimit(records, limit)

	if jsonOut {
		printHistoryJSON(records)

		return
	}

	printHistoryTerminal(records, detail)
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
		os.Exit(1)
	}
	defer db.Close()

	if cmdFilter != "" {
		records, err := db.ListHistoryByCommand(cmdFilter)
		if err != nil {
			if isLegacyDataError(err) {
				fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, constants.ErrHistoryQuery+"\n", err)
			os.Exit(1)
		}

		return records
	}

	records, err := db.ListHistory()
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrHistoryQuery+"\n", err)
		os.Exit(1)
	}

	return records
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

// printHistoryRow prints one row at the chosen detail level with ANSI
// colors: cyan command, dim flags, green OK / red FAIL, yellow
// duration, dim relative-time on the right.
func printHistoryRow(r model.CommandHistoryRecord, detail string) {
	cmd := colorize(constants.ColorCyan, padRight(r.Command, 16))
	flags := colorize(constants.ColorDim, padRight(truncateHist(r.Flags, 22), 22))
	status := colorizedStatus(r.ExitCode)
	dur := colorize(constants.ColorYellow, padRight(strconv.FormatInt(r.DurationMs, 10)+"ms", 10))
	last := colorize(constants.ColorDim, relativeHistoryTime(r))

	switch detail {
	case constants.DetailBasic:
		fmt.Printf("%s %s %s\n", cmd, status, last)
	case constants.DetailDetailed:
		args := colorize(constants.ColorWhite, padRight(truncateHist(r.Args, 18), 18))
		repos := padRight(strconv.Itoa(r.RepoCount), 6)
		summary := padRight(truncateHist(r.Summary, 30), 30)
		fmt.Printf("%s %s %s %s %s %s %s %s\n",
			cmd, args, flags, status, dur, repos, summary, last)
	default:
		fmt.Printf("%s %s %s %s %s\n", cmd, flags, status, dur, last)
	}
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

// printHistoryRevertSection enumerates revert commands for every row
// whose Command has a known inverse. Rows with no known revert are
// omitted so the section stays scannable. Empty when nothing is
// revertable (no header is printed in that case).
func printHistoryRevertSection(records []model.CommandHistoryRecord) {
	type revertRow struct {
		idx     int
		command string
		when    string
		hint    string
	}
	hints := make([]revertRow, 0, len(records))
	for i, r := range records {
		if h := revertHintFor(r); h != "" {
			hints = append(hints, revertRow{idx: i + 1, command: r.Command, when: relativeHistoryTime(r), hint: h})
		}
	}
	if len(hints) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(colorize(constants.ColorMagenta, "Revert points"))
	fmt.Println(colorize(constants.ColorDim,
		"  Run the suggested command to undo the referenced state."))
	for _, h := range hints {
		fmt.Printf("  %s#%-3d%s %s%-16s%s %s%-18s%s  %s%s%s\n",
			constants.ColorDim, h.idx, constants.ColorReset,
			constants.ColorCyan, h.command, constants.ColorReset,
			constants.ColorDim, h.when, constants.ColorReset,
			constants.ColorYellow, h.hint, constants.ColorReset)
	}
}

// revertHintFor maps a history Command to a concrete inverse command
// the user can run to roll back. Returns "" when no inverse is known.
// Kept tiny + table-driven so adding a new revertable command is a
// one-line change.
func revertHintFor(r model.CommandHistoryRecord) string {
	switch r.Command {
	case constants.CmdFixRepo, "fix-repo-pub", "fr", "frp":
		return "gitmap undo                  # restore latest fix-repo snapshot"
	case constants.CmdMakePublic, "mapub":
		return "gitmap make-private"
	case constants.CmdMakePrivate, "mapri":
		return "gitmap make-public --yes"
	case constants.CmdMakeAllPublic:
		return "gitmap make-all-private"
	case constants.CmdMakeAllPrivate:
		return "gitmap make-all-public --yes"
	case "reclone-transport":
		// Args holds the URL that was coerced; suggest re-running the
		// reclone with the OPPOSITE explicit transport so the user can
		// revisit the verdict.
		if strings.Contains(r.Flags, "transport=ssh") {
			return "gitmap cfr " + r.Args + " --https"
		}
		if strings.Contains(r.Flags, "transport=https") {
			return "gitmap cfr " + r.Args + " --ssh"
		}
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
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
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
