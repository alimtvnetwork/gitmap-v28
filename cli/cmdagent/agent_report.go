package cmdagent

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

const agentReportHelp = `agent report — shift report from the agent task database.

Description:
  Scans agent task run databases for the window and reports per-task
  claimed/completed/blocked counts plus open blockers. Prints a clean
  "no data" message when no task database exists.

Usage:
  gitmap agent report [--since <duration>]

Flags:
  --since <duration>  window like 30m, 12h, 2d (default: 24h)
  --help, -h          show this help

Example:
  gitmap agent report --since 12h
`

// reportCmd wires `gitmap agent report` into the agent cobra tree.
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Shift report of agent task activity over a time window",
	RunE: func(cmd *cobra.Command, args []string) error {
		since, _ := cmd.Flags().GetString("since")
		return RunAgentReport([]string{"--since", since})
	},
}

func init() {
	reportCmd.Flags().String("since", "24h", "window like 30m, 12h, 2d")
	AgentCmd.AddCommand(reportCmd)
}

// RunAgentReport implements `gitmap agent report --since <duration>`.
func RunAgentReport(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(agentReportHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("agent report", flag.ContinueOnError)
	sinceFlag := fs.String("since", "24h", "window like 30m, 12h, 2d")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(agentReportHelp)
		return apperror.WrapSimple(parseErr, "agent.report.flags")
	}
	window, parseErr := parseReportDuration(*sinceFlag)
	if parseErr != nil {
		fmt.Print(agentReportHelp)
		return appfault.NewValidationError(fmt.Sprintf("invalid --since value %q: use like 30m, 12h, 2d", *sinceFlag))
	}

	tempDir := store.ResolveAgentTempDir("")
	if _, statErr := os.Stat(tempDir); statErr != nil {
		fmt.Printf("agent report: no agent task data (no task dir at %s).\n", tempDir)
		return nil
	}
	taskDBs := findTaskDBs(tempDir)
	if len(taskDBs) == 0 {
		fmt.Println("agent report: no agent task databases found; nothing to report.")
		return nil
	}

	cutoff := time.Now().Add(-window)
	fmt.Printf("agent report (since %s):\n", *sinceFlag)
	reported := 0
	for _, taskDB := range taskDBs {
		summary, probeErr := probeTaskDB(taskDB.path, cutoff)
		if probeErr != nil || !summary.hasTelemetry {
			continue
		}
		reported++
		fmt.Printf("  %s: %d in window · %d done · %d active · %d blocked\n",
			taskDB.name, summary.inWindow, summary.done, summary.active, summary.blocked)
		for _, b := range summary.blockers {
			fmt.Printf("    blocked: %s\n", b)
		}
	}
	if reported == 0 {
		fmt.Println("  no telemetry rows in window.")
	}

	return nil
}

// reportTaskDB is one tier-2 task database on disk.
type reportTaskDB struct {
	name string
	path string
}

// reportSummary is the probed rollup for one task database.
type reportSummary struct {
	hasTelemetry bool
	inWindow     int
	done         int
	active       int
	blocked      int
	blockers     []string
}

func findTaskDBs(tempDir string) []reportTaskDB {
	entries, readErr := os.ReadDir(tempDir)
	if readErr != nil {
		return nil
	}
	var out []reportTaskDB
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dbPath := filepath.Join(tempDir, e.Name(), "agent-task.db")
		if _, statErr := os.Stat(dbPath); statErr != nil {
			continue
		}
		out = append(out, reportTaskDB{name: e.Name(), path: dbPath})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })

	return out
}

// probeTaskDB opens a task DB and rolls up a telemetry table if one parses.
func probeTaskDB(dbPath string, cutoff time.Time) (reportSummary, *appfault.AppError) {
	var summary reportSummary
	conn, openErr := sql.Open("sqlite", "file:"+dbPath+"?mode=ro")
	if openErr != nil {
		return summary, apperror.WrapSimple(openErr, "agent.report.open")
	}
	defer conn.Close()

	for _, table := range []string{"subtasks", "tasks"} {
		if !hasReportTable(conn, table) {
			continue
		}
		rows, queryErr := conn.Query("SELECT * FROM " + table + " LIMIT 5000")
		if queryErr != nil {
			continue
		}
		summary = rollupReportRows(rows, cutoff)
		rows.Close()
		if summary.hasTelemetry {
			return summary, nil
		}
	}

	return summary, nil
}

func hasReportTable(conn *sql.DB, table string) bool {
	var name string
	queryErr := conn.QueryRow(
		"SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)

	return queryErr == nil && name == table
}

// rollupReportRows folds generic rows into done/active/blocked counts.
func rollupReportRows(rows *sql.Rows, cutoff time.Time) reportSummary {
	var summary reportSummary
	cols, colsErr := rows.Columns()
	if colsErr != nil {
		return summary
	}
	statusIdx := indexOfFold(cols, []string{"status", "state"})
	timeIdx := indexOfFold(cols, []string{"updated_at", "completed_at", "claimed_at", "created_at"})
	titleIdx := indexOfFold(cols, []string{"title", "name", "slug", "description"})

	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if scanErr := rows.Scan(ptrs...); scanErr != nil {
			continue
		}
		summary.hasTelemetry = true
		rowTime := time.Time{}
		hasTime := false
		if timeIdx >= 0 {
			rowTime, hasTime = parseReportTime(stringifyReportVal(vals[timeIdx]))
		}
		if hasTime && rowTime.Before(cutoff) {
			continue
		}
		summary.inWindow++
		status := ""
		if statusIdx >= 0 {
			status = strings.ToLower(stringifyReportVal(vals[statusIdx]))
		}
		title := ""
		if titleIdx >= 0 {
			title = stringifyReportVal(vals[titleIdx])
		}
		switch {
		case strings.Contains(status, "block") || strings.Contains(status, "fail"):
			summary.blocked++
			if title != "" && len(summary.blockers) < 5 {
				summary.blockers = append(summary.blockers, title)
			}
		case strings.Contains(status, "done") || strings.Contains(status, "complet") || strings.Contains(status, "success"):
			summary.done++
		default:
			summary.active++
		}
	}

	return summary
}

func indexOfFold(cols []string, candidates []string) int {
	for _, cand := range candidates {
		for i, c := range cols {
			if strings.EqualFold(c, cand) {
				return i
			}
		}
	}

	return -1
}

func stringifyReportVal(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case []byte:
		return string(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	default:
		return fmt.Sprintf("%v", t)
	}
}

func parseReportTime(raw string) (time.Time, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return time.Time{}, false
	}
	if unix, numErr := strconv.ParseInt(trimmed, 10, 64); numErr == nil {
		if unix > 1000000000000 {
			return time.UnixMilli(unix), true
		}
		if unix > 1000000000 {
			return time.Unix(unix, 0), true
		}
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006-01-02"} {
		if parsed, parseErr := time.Parse(layout, trimmed); parseErr == nil {
			return parsed, true
		}
	}

	return time.Time{}, false
}

func parseReportDuration(raw string) (time.Duration, error) {
	trimmed := strings.TrimSpace(raw)
	if strings.HasSuffix(trimmed, "d") {
		days, numErr := strconv.Atoi(strings.TrimSuffix(trimmed, "d"))
		if numErr != nil || days < 0 {
			return 0, apperror.NewValidationError("bad duration")
		}

		return time.Duration(days) * 24 * time.Hour, nil
	}

	return time.ParseDuration(trimmed)
}
