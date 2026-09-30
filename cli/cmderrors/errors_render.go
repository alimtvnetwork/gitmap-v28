// Package cmderrors provides rendering utilities for internal error inspection.
package cmderrors

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func renderErrorsHelp() {
	fmt.Printf("\n%s%sGitMap Internal Errors Inspector (gitmap errors / gitmap e)%s\n\n", constants.ColorCyan, constants.ColorBold, constants.ColorReset)
	fmt.Println("  Usage:")
	fmt.Println("    gitmap errors [flags]         List recorded internal errors")
	fmt.Println("    gitmap e [flags]              Short alias to list recorded internal errors")
	fmt.Println("    gitmap e <id>                 Show detailed information for a specific error")
	fmt.Println("    gitmap e clear                Clear all recorded internal errors")
	fmt.Println("    gitmap e failed               List failed/undetected commands (alias: gitmap fc)")
	fmt.Println("    gitmap failed-commands count  Show count of failed/undetected commands")
	fmt.Println()
	fmt.Println("  Flags:")
	fmt.Println("    --json                        Output error records as JSON")
	fmt.Println("    --limit, -l <n>               Maximum number of errors to list (default: 25)")
	fmt.Println("    --clear                       Purge all errors from gitmap-errors.db")
	fmt.Println("    --help, -h                    Show this help message")
	fmt.Println()
}

func renderErrorsTable(records []store.InternalErrorRecord, dbPath string) {
	if len(records) == 0 {
		fmt.Printf("\n  %s✓%s %sNo internal errors recorded in gitmap-errors.db. System is clean.%s\n\n",
			constants.ColorGreen, constants.ColorReset, constants.ColorBold, constants.ColorReset)
		fmt.Printf("  • Database: %s\n\n", dbPath)

		return
	}

	fmt.Printf("\n%s%sInternal Errors Log (gitmap-errors.db):%s\n\n", constants.ColorYellow, constants.ColorBold, constants.ColorReset)
	fmt.Printf("  %-6s  %-20s  %-20s  %-16s  %s\n", "ID", "TIMESTAMP (UTC)", "TYPE", "CODE", "MESSAGE")
	fmt.Printf("  %-6s  %-20s  %-20s  %-16s  %s\n", "------", "--------------------", "--------------------", "----------------", "----------------------------------------")

	for _, rec := range records {
		ts := rec.CreatedAt
		if len(ts) > 19 {
			ts = ts[:19]
		}
		msg := strings.ReplaceAll(rec.Message, "\n", " ")
		if len(msg) > 60 {
			msg = msg[:57] + "..."
		}
		fmt.Printf("  %-6d  %-20s  %-20s  %-16s  %s\n", rec.ID, ts, truncateStr(rec.ErrorType, 20), truncateStr(rec.ErrorCode, 16), msg)
	}

	fmt.Printf("\n  Showing %d error(s).\n", len(records))
	fmt.Printf("  • Database: %s\n", dbPath)
	fmt.Printf("  • Hint:     Run 'gitmap e <id>' to inspect full error details, or 'gitmap e clear' to clear.\n\n")
}

func renderSingleError(rec *store.InternalErrorRecord) {
	if rec == nil {
		return
	}

	fmt.Printf("\n%s%sInternal Error #%d Details:%s\n\n", constants.ColorRed, constants.ColorBold, rec.ID, constants.ColorReset)
	fmt.Printf("  • ID:          %d\n", rec.ID)
	fmt.Printf("  • Timestamp:   %s\n", rec.CreatedAt)
	fmt.Printf("  • Error Type:  %s\n", rec.ErrorType)
	fmt.Printf("  • Error Code:  %s\n", rec.ErrorCode)
	if rec.Command != "" {
		fmt.Printf("  • Command:     %s\n", rec.Command)
	}
	fmt.Printf("  • Message:     %s\n", rec.Message)
	if rec.Details != "" {
		fmt.Printf("  • Details:     %s\n", rec.Details)
	}
	if rec.SourceFile != "" {
		fmt.Printf("  • Source File: %s\n", rec.SourceFile)
	}
	if rec.StackTrace != "" {
		fmt.Printf("  • Stack Trace:\n%s\n", indentStr(rec.StackTrace, "      "))
	}
	fmt.Println()
}

func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen-3] + "..."
}

func indentStr(s, indent string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}

	return strings.Join(lines, "\n")
}
