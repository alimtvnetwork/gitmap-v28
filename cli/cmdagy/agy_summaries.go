package cmdagy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AgySummaryRow represents a projected conversation summary row.
type AgySummaryRow struct {
	ConversationID   string `json:"conversationId"`
	Title            string `json:"title"`
	Status           string `json:"status"`
	NotFullyIdle     int    `json:"notFullyIdle"`
	LastModifiedTime string `json:"lastModifiedTime"`
	WorkspaceURI     string `json:"workspaceUri"`
}

type agySummariesOptions struct {
	limit    int
	query    string
	isJSON   bool
	isActive bool
	isHelp   bool
}

// ResolveAgyConversationDbPath resolves ~/.gemini/antigravity/conversation_summaries.db.
func ResolveAgyConversationDbPath() (string, *apperror.AppError) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", apperror.WrapSimple(homeErr, "resolve user home dir")
	}

	dbPath := filepath.Join(homeDir, ".gemini", "antigravity", "conversation_summaries.db")
	if _, err := os.Stat(dbPath); err != nil {
		return "", apperror.NewNotFound("agy_summaries", "E4040", fmt.Sprintf("conversation_summaries.db not found at %s", dbPath))
	}

	return dbPath, nil
}

func cleanWorkspaceURI(raw string) string {
	cleaned := strings.Trim(raw, "[]\"'")
	cleaned = strings.TrimPrefix(cleaned, "file:///")
	cleaned = strings.ReplaceAll(cleaned, "%3A", ":")
	cleaned = strings.ReplaceAll(cleaned, "%20", " ")

	return filepath.ToSlash(cleaned)
}

func parseAgySummariesOptions(args []string) agySummariesOptions {
	opts := agySummariesOptions{limit: 10}
	for idx := 0; idx < len(args); idx++ {
		opts = applySingleSummaryArg(args[idx], idx, args, opts)
	}

	return opts
}

func applySingleSummaryArg(arg string, idx int, args []string, opts agySummariesOptions) agySummariesOptions {
	if arg == "-j" || arg == "--json" {
		opts.isJSON = true
		return opts
	}
	if arg == "-a" || arg == "--active" {
		opts.isActive = true
		return opts
	}
	if arg == "-h" || arg == "--help" || arg == "help" {
		opts.isHelp = true
		return opts
	}

	return applyValueSummaryArg(arg, idx, args, opts)
}

func applyValueSummaryArg(arg string, idx int, args []string, opts agySummariesOptions) agySummariesOptions {
	if (arg == "-l" || arg == "--limit") && idx+1 < len(args) {
		if val, err := strconv.Atoi(args[idx+1]); err == nil && val > 0 {
			opts.limit = val
		}
		return opts
	}
	if (arg == "-q" || arg == "--query") && idx+1 < len(args) {
		opts.query = strings.TrimSpace(args[idx+1])
		return opts
	}
	if val, err := strconv.Atoi(arg); err == nil && val > 0 {
		opts.limit = val
	}

	return opts
}

func buildAgySummariesQuery(opts agySummariesOptions) (string, []any) {
	query := "SELECT conversation_id, title, status, not_fully_idle, last_modified_time, workspace_uris FROM conversation_summaries"
	clauses, params := []string{}, []any{}
	if opts.isActive {
		clauses = append(clauses, "(not_fully_idle = 1 OR status = 'CASCADE_RUN_STATUS_RUNNING')")
	}
	if opts.query != "" {
		clauses = append(clauses, "(title LIKE ? OR conversation_id LIKE ? OR workspace_uris LIKE ?)")
		pat := "%" + opts.query + "%"
		params = append(params, pat, pat, pat)
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	query += " ORDER BY last_modified_time DESC LIMIT ?"
	params = append(params, opts.limit)

	return query, params
}

func fetchAgySummaries(dbPath string, opts agySummariesOptions) ([]AgySummaryRow, *apperror.AppError) {
	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return nil, openRes.AppError()
	}
	conn := openRes.Value
	defer conn.Close()

	sqlQuery, params := buildAgySummariesQuery(opts)
	rows, err := conn.Query(sqlQuery, params...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query conversation_summaries")
	}
	defer rows.Close()

	return scanAgySummaryRows(rows)
}

func scanAgySummaryRows(rows *sql.Rows) ([]AgySummaryRow, *apperror.AppError) {
	var results []AgySummaryRow
	for rows.Next() {
		var row AgySummaryRow
		var rawWS string
		scanErr := rows.Scan(&row.ConversationID, &row.Title, &row.Status, &row.NotFullyIdle, &row.LastModifiedTime, &rawWS)
		if scanErr != nil {
			return nil, apperror.WrapSimple(scanErr, "scan conversation_summaries row")
		}
		row.WorkspaceURI = cleanWorkspaceURI(rawWS)
		results = append(results, row)
	}

	return results, nil
}

func truncateSummaryText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}

	return text[:maxLen-3] + "..."
}

func formatSummaryTime(raw string) string {
	if len(raw) >= 19 {
		return raw[:19]
	}

	return raw
}

func printSummaryTableHeader(dbPath string) {
	fmt.Printf("\n%s● Antigravity Conversation Summaries%s (%s)\n", constants.ColorCyan, constants.ColorReset, dbPath)
	fmt.Printf("  ╭────┬──────────────────────────────────────┬──────────────────────────────┬─────────────────────────┬────────┬─────────────────────┬──────────────────────────────╮\n")
	fmt.Printf("  │  # │ CONVERSATION ID                      │ TITLE                        │ STATUS                  │ IDLE?  │ LAST MODIFIED       │ WORKSPACE                    │\n")
	fmt.Printf("  ├────┼──────────────────────────────────────┼──────────────────────────────┼─────────────────────────┼────────┼─────────────────────┼──────────────────────────────┤\n")
}

func printSummaryTableRow(idx int, r AgySummaryRow) {
	idleText := "Yes"
	if r.NotFullyIdle == 1 || r.Status == "CASCADE_RUN_STATUS_RUNNING" {
		idleText = "Active"
	}
	titleTrunc := truncateSummaryText(r.Title, 28)
	wsTrunc := truncateSummaryText(r.WorkspaceURI, 28)
	timeClean := formatSummaryTime(r.LastModifiedTime)
	fmt.Printf("  │ %2d │ %-36s │ %-28s │ %-23s │ %-6s │ %-19s │ %-28s │\n", idx+1, r.ConversationID, titleTrunc, r.Status, idleText, timeClean, wsTrunc)
}

func printSummaryTableFooter(count int, elapsed time.Duration) {
	fmt.Printf("  ╰────┴──────────────────────────────────────┴──────────────────────────────┴─────────────────────────┴────────┴─────────────────────┴──────────────────────────────╯\n")
	fmt.Printf("  %d conversation(s) returned (in %v).\n\n", count, elapsed)
}

func renderAgySummariesTable(rows []AgySummaryRow, elapsed time.Duration, dbPath string) {
	printSummaryTableHeader(dbPath)
	for idx, r := range rows {
		printSummaryTableRow(idx, r)
	}
	printSummaryTableFooter(len(rows), elapsed)
}

func renderAgySummariesJSON(rows []AgySummaryRow) error {
	bytes, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal summaries json")
	}
	fmt.Println(string(bytes))

	return nil
}

func printAgySummariesHelp() {
	fmt.Println()
	fmt.Printf("  %sgitmap agy summaries%s — Inspect Antigravity conversation states and history\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Usage:")
	fmt.Println("    gitmap agy summaries [N] [flags]")
	fmt.Println("    gitmap agy convs [N] [flags]")
	fmt.Println("    gitmap agy db [N] [flags]")
	fmt.Println()
	fmt.Println("  Flags:")
	fmt.Println("    -l, --limit <N>      Maximum rows to display (default: 10)")
	fmt.Println("    -q, --query <text>   Filter by title, conversation ID, or workspace")
	fmt.Println("    -a, --active         Filter for running / active conversations only")
	fmt.Println("    -j, --json           Output in structured JSON format")
	fmt.Println()
}

// RunAgySummariesCLI routes and executes gitmap agy summaries.
func RunAgySummariesCLI(args []string) error {
	opts := parseAgySummariesOptions(args)
	if opts.isHelp {
		printAgySummariesHelp()
		return nil
	}

	start := time.Now()
	dbPath, pathErr := ResolveAgyConversationDbPath()
	if pathErr != nil {
		return pathErr.Unwrap()
	}

	rows, fetchErr := fetchAgySummaries(dbPath, opts)
	if fetchErr != nil {
		return fetchErr.Unwrap()
	}

	if opts.isJSON {
		return renderAgySummariesJSON(rows)
	}

	renderAgySummariesTable(rows, time.Since(start), dbPath)

	return nil
}
