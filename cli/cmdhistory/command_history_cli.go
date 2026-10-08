// Package cmd — command_history_cli.go implements gitmap history subcommands and suggestions.
package cmdhistory

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunCommandHistoryCLI executes history subcommands (ls, suggest, clear).
func RunCommandHistoryCLI(args []string) error {
	if len(args) == 0 {
		return runHistoryListCLI(args)
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "ls", "list":
		return runHistoryListCLI(args[1:])
	case "suggest":
		return runHistorySuggestCLI(args[1:])
	case "clear", "clean", "reset":
		return runHistoryClearCLI()
	case "help", "-h", "--help":
		printHistoryCLIHelp()
		return nil
	default:
		return runHistoryListCLI(args)
	}
}

func runHistoryListCLI(args []string) error {
	fs := flag.NewFlagSet("history ls", flag.ContinueOnError)
	limit := fs.Int("limit", 20, "Maximum commands to show")
	isJSON := fs.Bool("json", false, "Output in JSON format")
	_ = fs.Parse(args)

	db, err := store.OpenCommandHistorySplitDB("")
	if err != nil {
		return err
	}
	defer db.Close()

	entries, err := db.ListRecentCommands(*limit)
	if err != nil {
		return err
	}

	if *isJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}

	renderHistoryTable(entries)
	return nil
}

func renderHistoryTable(entries []store.CommandHistoryEntry) {
	if len(entries) == 0 {
		fmt.Printf("  %sNo command history recorded yet.%s\n", constants.ColorYellow, constants.ColorReset)
		return
	}

	fmt.Printf("\n  %s%sCommand Execution History%s (%d entries)\n\n", constants.ColorBold, constants.ColorCyan, constants.ColorReset, len(entries))
	fmt.Printf("  %-6s  %-20s  %-10s  %-6s  %s\n", "ID", "EXECUTED AT", "DURATION", "STATUS", "COMMAND")
	fmt.Printf("  %-6s  %-20s  %-10s  %-6s  %s\n", "------", "--------------------", "----------", "------", "----------------------------------------")

	for _, e := range entries {
		timeStr := e.ExecutedAt
		if len(timeStr) > 19 {
			timeStr = timeStr[:19]
		}
		status := fmt.Sprintf("%s✓ 0%s", constants.ColorGreen, constants.ColorReset)
		if e.ExitCode != 0 {
			status = fmt.Sprintf("%s✗ %d%s", constants.ColorRed, e.ExitCode, constants.ColorReset)
		}
		durStr := fmt.Sprintf("%dms", e.DurationMs)
		fmt.Printf("  %-6d  %-20s  %-10s  %-6s  %s\n", e.ID, timeStr, durStr, status, e.Line)
	}
	fmt.Println()
}

func runHistorySuggestCLI(args []string) error {
	prefix := ""
	if len(args) > 0 {
		prefix = args[0]
	}

	db, err := store.OpenCommandHistorySplitDB("")
	if err != nil {
		return err
	}
	defer db.Close()

	suggestions, err := db.SuggestCommands(prefix, 10)
	if err != nil {
		return err
	}

	for _, s := range suggestions {
		fmt.Println(s)
	}
	return nil
}

func runHistoryClearCLI() error {
	db, err := store.OpenCommandHistorySplitDB("")
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.ClearCommands(); err != nil {
		return err
	}

	fmt.Printf("  %s✓%s Command history database cleared.\n", constants.ColorGreen, constants.ColorReset)
	return nil
}

func printHistoryCLIHelp() {
	fmt.Printf("\n  %s%sgitmap history — Command Execution History & Reuse Suggestions%s\n\n", constants.ColorBold, constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Usage:")
	fmt.Println("    gitmap history [ls] [--limit 20] [--json]    List recent command executions")
	fmt.Println("    gitmap history suggest <prefix>               Suggest past commands matching prefix")
	fmt.Println("    gitmap history clear                          Clear command execution history")
	fmt.Println("    gitmap history help                           Show this help guide")
	fmt.Println()
}
