package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runDBHelp() error {
	printDBHelpHeader()
	printDBHelpSubcommands()
	printDBHelpFlags()
	printDBHelpExamples()
	return nil
}

func printDBHelpHeader() {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + "Gitmap Database Management (gitmap db)" + constants.ColorReset)
	fmt.Println()
	fmt.Println("  Manage, inspect, and reset Gitmap's primary SQLite database and per-repo split databases.")
	fmt.Println()
	fmt.Println("  " + constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("    gitmap db [subcommand] [flags]")
	fmt.Println("    gitmap start-fresh [flags]")
	fmt.Println()
}

func printDBHelpSubcommands() {
	fmt.Println("  " + constants.ColorCyan + "Subcommands:" + constants.ColorReset)
	fmt.Printf("    %-24s %s\n", "status", "Show consolidated SQLite location, path, size, and split database counts")
	fmt.Printf("    %-24s %s\n", "optimize", "Execute VACUUM and optimize across Master, Repo, and Pipeline split databases")
	fmt.Printf("    %-24s %s\n", "ls, list", "List all databases (main and split DBs) with paths, sizes, and architectural purpose")
	fmt.Printf("    %-24s %s\n", "repo-db list", "List per-repository split databases with table metrics and tracking status")
	fmt.Printf("    %-24s %s\n", "sizes list", "Show compact disk usage and size breakdown for all database files")
	fmt.Printf("    %-24s %s\n", "clear", "Clear search caches across split databases (--yes to skip prompt)")
	fmt.Printf("    %-24s %s\n", "reset", "Reset database records and remove split DBs (--yes to skip prompt)")
	fmt.Printf("    %-24s %s\n", "help", "Show this help message")
	fmt.Println()
}

func printDBHelpFlags() {
	fmt.Println("  " + constants.ColorCyan + "Flags:" + constants.ColorReset)
	fmt.Printf("    %-24s %s\n", "-y, --yes, -f, --force", "Bypass confirmation prompt for reset or start-fresh")
	fmt.Printf("    %-24s %s\n", "-h, --help", "Show help for the command")
	fmt.Println()
}

func printDBHelpExamples() {
	fmt.Println("  " + constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("    gitmap db ls                     # View architectural database overview")
	fmt.Println("    gitmap db repo-db list           # View split DB table metrics")
	fmt.Println("    gitmap db sizes list             # View database disk usage table")
	fmt.Println("    gitmap db reset --yes            # Reset DB non-interactively")
	fmt.Println("    gitmap start-fresh               # Clear all data & rebuild clean schemas")
	fmt.Println()
}
