package cmdhistory

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

func isHistoryPurgeRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}

	return args[0] == "purge"
}

func isHistorySubcommand(sub string) bool {
	return sub == "ls" || sub == "list" || sub == "suggest" || sub == "clear" || sub == "clean"
}

// runHistory handles the "history" subcommand.
func RunHistory(args []string) error {
	if isHistoryPurgeRequest(args) {
		return RunHistoryPurgeCLI(args[1:])
	}
	checkHelp("history", args)
	if len(args) > 0 && (strings.EqualFold(args[0], "ssh") || strings.EqualFold(args[0], "--ssh")) {
		return cmdssh.RunFleetPASCommand("history", "gitmap history --limit 10", func() error {
			return RunHistory(args[1:])
		})
	}
	if len(args) > 0 && isHistorySubcommand(strings.ToLower(args[0])) {
		return RunCommandHistoryCLI(args)
	}
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
		cliexit.HandleError(err, 1)
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
	paddedCmd := termout.PadRightMax(r.Command, 16)
	truncFlags := termout.PadRightMax(termout.TruncateRightWithEllipsis(r.Flags, 22), 22)
	durText := termout.PadRightMax(strconv.FormatInt(r.DurationMs, 10)+"ms", 10)

	return historyRowTokens{
		cmd:    termout.SafeTokyoStyle(termout.TokyoPrimary).Render(paddedCmd),
		flags:  termout.SafeTokyoStyle(termout.TokyoMuted).Render(truncFlags),
		status: colorizedStatus(r.ExitCode),
		dur:    termout.SafeTokyoStyle(termout.TokyoWarn).Render(durText),
		last:   termout.SafeTokyoStyle(termout.TokyoMuted).Render(relativeHistoryTime(r)),
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
	truncArgs := termout.PadRightMax(termout.TruncateRightWithEllipsis(r.Args, 18), 18)
	args := termout.SafeTokyoStyle(termout.TokyoFG).Render(truncArgs)
	repos := termout.PadRightMax(strconv.Itoa(r.RepoCount), 6)
	summary := termout.PadRightMax(termout.TruncateRightWithEllipsis(r.Summary, 30), 30)
	fmt.Printf("%s %s %s %s %s %s %s %s\n",
		tok.cmd, args, tok.flags, tok.status, tok.dur, repos, summary, tok.last)
}

// colorizedStatus renders a colored OK / FAIL token.
func colorizedStatus(code int) string {
	if code == 0 {
		return termout.TokyoOKPill("✓ OK")
	}

	return termout.TokyoErrorPill("✗ FAIL")
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
	fmt.Println(termout.SafeTokyoStyle(termout.TokyoAccent).Render("Revert points"))
	fmt.Println(termout.SafeTokyoStyle(termout.TokyoMuted).Render("  Run the suggested command to undo the referenced state."))
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
