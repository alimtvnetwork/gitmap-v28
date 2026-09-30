// Package cmderrors — failed_commands_cmd.go provides the CLI command to inspect and manage failed/undetected commands.
package cmderrors

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunFailedCommandsCLI is the entry point for 'gitmap failed-commands', 'gitmap fc', and 'gitmap e failed'.
func RunFailedCommandsCLI(args []string) error {
	if isErrorsHelpRequested(args) {
		renderFailedCommandsHelp()

		return nil
	}

	if isErrorsClearRequested(args) {
		return handleClearFailedCommands()
	}

	if isFailedCountRequested(args) {
		return handleCountFailedCommands(args)
	}

	return handleListFailedCommands(args)
}

func isFailedSubcommandToken(token string) bool {
	low := strings.ToLower(strings.TrimSpace(token))

	return low == "failed" || low == "failed-commands" || low == "failed-command" ||
		low == "fc" || low == "unknown" || low == "unknown-commands" || low == "failed-to-detect"
}

func isFailedCountRequested(args []string) bool {
	for _, arg := range args {
		low := strings.ToLower(strings.TrimSpace(arg))
		if low == "count" || low == "stats" || low == "summary" || low == "--count" {
			return true
		}
	}

	return false
}

func handleClearFailedCommands() error {
	db, err := store.OpenErrorsSplitDB()
	if err != nil {
		return fmt.Errorf("open errors db: %w", err)
	}

	defer db.Close()

	if err := db.ClearFailedCommands(); err != nil {
		return fmt.Errorf("clear failed commands: %w", err)
	}

	fmt.Println("  ✓ Cleared all recorded failed/undetected commands from FailedCommand table in gitmap-errors.db.")

	return nil
}

func handleCountFailedCommands(args []string) error {
	isJSON := hasFlag(args, "--json")

	db, err := store.OpenErrorsSplitDB()
	if err != nil {
		return fmt.Errorf("open errors db: %w", err)
	}

	defer db.Close()

	distinctCount, totalHits, countErr := db.CountFailedCommands()
	if countErr != nil {
		return fmt.Errorf("count failed commands: %w", countErr)
	}

	if isJSON {
		return emitFailedCountJSON(distinctCount, totalHits, db.Path)
	}

	fmt.Printf("\n%s%sFailed / Undetected Commands Summary:%s\n\n", constants.ColorCyan, constants.ColorBold, constants.ColorReset)
	fmt.Printf("  • Distinct Failed Commands: %s%d%s\n", constants.ColorBold, distinctCount, constants.ColorReset)
	fmt.Printf("  • Total Failed Attempts:    %s%d%s\n", constants.ColorBold, totalHits, constants.ColorReset)
	fmt.Printf("  • Database Table:           FailedCommand (%s)\n\n", db.Path)

	return nil
}

func emitFailedCountJSON(distinctCount, totalHits int64, dbPath string) error {
	payload := map[string]any{
		"totalDistinctCommands": distinctCount,
		"totalFailedAttempts":   totalHits,
		"table":                 "FailedCommand",
		"databasePath":          dbPath,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(payload)
}

func handleListFailedCommands(args []string) error {
	isJSON := hasFlag(args, "--json")
	limit := parseLimitArg(args, 50)

	db, err := store.OpenErrorsSplitDB()
	if err != nil {
		return fmt.Errorf("open errors db: %w", err)
	}

	defer db.Close()

	summary, sumErr := db.GetFailedCommandSummary(limit)
	if sumErr != nil {
		return fmt.Errorf("list failed commands: %w", sumErr)
	}

	if isJSON {
		return emitFailedSummaryJSON(summary)
	}

	renderFailedCommandsTable(summary)

	return nil
}

func emitFailedSummaryJSON(summary store.FailedCommandSummary) error {
	if summary.Records == nil {
		summary.Records = []store.FailedCommandRecord{}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(summary)
}

func renderFailedCommandsHelp() {
	fmt.Printf("\n%s%sGitMap Failed / Undetected Commands Inspector (gitmap failed-commands / gitmap fc)%s\n\n",
		constants.ColorCyan, constants.ColorBold, constants.ColorReset)
	fmt.Println("  Usage:")
	fmt.Println("    gitmap failed-commands [flags]     List failed/undetected commands with hit counts & suggestions")
	fmt.Println("    gitmap fc [flags]                  Short alias for failed-commands")
	fmt.Println("    gitmap failed-commands count       Show total count of failed/undetected commands")
	fmt.Println("    gitmap failed-commands clear       Clear all recorded failed commands from FailedCommand table")
	fmt.Println("    gitmap e failed                    Inspect failed commands via errors subcommand")
	fmt.Println()
	fmt.Println("  Flags:")
	fmt.Println("    --json                             Output failed command summary & records as JSON")
	fmt.Println("    --limit, -l <n>                    Maximum number of failed commands to list (default: 50)")
	fmt.Println("    --clear                            Purge all rows from FailedCommand table")
	fmt.Println("    --help, -h                         Show this help message")
	fmt.Println()
}

func renderFailedCommandsTable(summary store.FailedCommandSummary) {
	if len(summary.Records) == 0 {
		fmt.Printf("\n  %s✓%s %sNo failed or undetected commands recorded in FailedCommand table (0 failed commands).%s\n\n",
			constants.ColorGreen, constants.ColorReset, constants.ColorBold, constants.ColorReset)
		fmt.Printf("  • Database: %s\n\n", summary.DatabasePath)

		return
	}

	fmt.Printf("\n%s%sFailed / Undetected Commands Log (FailedCommand table):%s\n",
		constants.ColorYellow, constants.ColorBold, constants.ColorReset)
	fmt.Printf("  • Distinct Failed Commands: %s%d%s   • Total Failed Attempts: %s%d%s\n\n",
		constants.ColorBold, summary.TotalDistinctCommands, constants.ColorReset,
		constants.ColorBold, summary.TotalFailedAttempts, constants.ColorReset)
	fmt.Printf("  %-5s  %-6s  %-10s  %-26s  %-34s  %s\n",
		"ID", "HITS", "DOMAIN", "FAILED COMMAND", "SUGGESTED ALTERNATIVES", "LAST SEEN (UTC)")
	fmt.Printf("  %-5s  %-6s  %-10s  %-26s  %-34s  %s\n",
		"-----", "------", "----------", "--------------------------", "----------------------------------", "-------------------")

	for _, rec := range summary.Records {
		renderFailedCommandRow(rec)
	}

	fmt.Printf("\n  Showing %d of %d distinct failed command(s) (%d total failed attempts).\n",
		len(summary.Records), summary.TotalDistinctCommands, summary.TotalFailedAttempts)
	fmt.Printf("  • Database: %s (table: FailedCommand)\n", summary.DatabasePath)
	fmt.Printf("  • 💡 Suggestions & Next Actions:\n")
	fmt.Printf("      - Create a custom shortcut for frequent commands: gitmap alias add <short> \"<target-cmd>\"\n")
	fmt.Printf("      - Check count only:                               gitmap failed-commands count\n")
	fmt.Printf("      - Clear recorded failed commands:                 gitmap failed-commands clear\n\n")
}

func renderFailedCommandRow(rec store.FailedCommandRecord) {
	ts := rec.LastSeenAt
	if len(ts) > 19 {
		ts = ts[:19]
	}
	sugg := rec.Suggestions
	if sugg == "" {
		sugg = "gitmap help"
	}
	fmt.Printf("  %-5d  %-6d  %-10s  %-26s  %-34s  %s\n",
		rec.ID,
		rec.HitCount,
		truncateStr(rec.Domain, 10),
		truncateStr(rec.Command, 26),
		truncateStr(sugg, 34),
		ts,
	)
}
