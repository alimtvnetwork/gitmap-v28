package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runRepoDB routes gitmap repo db subcommands.
func runRepoDB(args []string) error {
	if len(args) == 0 {
		return handleRepoDBStatus(nil)
	}

	sub := strings.ToLower(strings.TrimSpace(args[0]))
	switch sub {
	case "status", "st", "s", "info":
		return handleRepoDBStatus(args[1:])
	case "log", "logs", "l":
		return handleRepoDBLog(args[1:])
	case "errorlogs", "error-logs", "errors", "err":
		return handleRepoDBErrorLogs(args[1:])
	case "clear", "cl":
		return handleRepoDBClear(args[1:])
	case "reset":
		return handleRepoDBReset(args[1:])
	case "optimize", "opt":
		return handleRepoDBOptimize(args[1:])
	case "help", "-h", "--help":
		printRepoDBHelp()
		return nil
	default:
		printRepoDBHelp()
		return fmt.Errorf("unknown repo db subcommand: %s", sub)
	}
}

func printRepoDBHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap repo db [command] [flags]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Commands:" + constants.ColorReset)
	fmt.Printf("  %-22s %s\n", "status", "Show repository split database location, path, size, and table row counts (default)")
	fmt.Printf("  %-22s %s\n", "log", "Show recent indexing and scan activity logs from RepoScanLog")
	fmt.Printf("  %-22s %s\n", "error-logs", "Show failure entries from RepoScanLog (alias: errorlogs)")
	fmt.Printf("  %-22s %s\n", "clear", "Clear search cache and file sequence index (supports -y)")
	fmt.Printf("  %-22s %s\n", "reset", "Drop all tables and recreate clean repository schema")
	fmt.Printf("  %-22s %s\n", "optimize", "Execute VACUUM and optimize database file, reporting reclaimed space")
	fmt.Printf("  %-22s %s\n", "help", "Show repository database help")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Flags:" + constants.ColorReset)
	fmt.Printf("  %-22s %s\n", "-y, --yes", "Skip confirmation prompts for clear and reset")
	fmt.Printf("  %-22s %s\n", "--json", "Output data in structured JSON format")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("  gitmap repo db")
	fmt.Println("  gitmap repo db status")
	fmt.Println("  gitmap repo db log")
	fmt.Println("  gitmap repo db error-logs")
	fmt.Println("  gitmap repo db optimize")
	fmt.Println("  gitmap repo db clear -y")
}
