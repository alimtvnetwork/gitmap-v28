package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// handlePipelineDB routes the gitmap pipeline db subcommands.
func handlePipelineDB(args []string) error {
	if len(args) == 0 {
		return runPipelineDBStatus(nil)
	}

	sub := strings.ToLower(strings.TrimSpace(args[0]))
	switch sub {
	case "status", "st", "s", "info":
		return runPipelineDBStatus(args[1:])
	case "clear", "cl":
		return runPipelineDBClear(args[1:])
	case "reset":
		return runPipelineDBReset(args[1:])
	case "optimize", "opt":
		return runPipelineDBOptimize(args[1:])
	case "errorlogs", "error-logs", "errors", "err":
		return runPipelineDBErrorLogs(args[1:])
	case "help", "-h", "--help":
		printPipelineDBHelp()
		return nil
	default:
		printPipelineDBHelp()
		return fmt.Errorf("unknown pipeline db subcommand: %s", sub)
	}
}

func printPipelineDBHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap pipeline db [command] [flags]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Commands:" + constants.ColorReset)
	fmt.Printf("  %-22s %s\n", "status", "Show pipeline split database location, path, size, and run metrics (default)")
	fmt.Printf("  %-22s %s\n", "clear", "Clear recorded pipeline runs and error logs (supports -y)")
	fmt.Printf("  %-22s %s\n", "reset", "Drop all tables and recreate fresh pipeline schema")
	fmt.Printf("  %-22s %s\n", "optimize", "Execute VACUUM and optimize database file, reporting reclaimed space")
	fmt.Printf("  %-22s %s\n", "error-logs", "Query and display error logs stored in this repo's pipeline DB (alias: errorlogs)")
	fmt.Printf("  %-22s %s\n", "help", "Show pipeline database help")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Flags:" + constants.ColorReset)
	fmt.Printf("  %-22s %s\n", "-y, --yes", "Skip confirmation prompts for clear and reset")
	fmt.Printf("  %-22s %s\n", "--json", "Output data in structured JSON format")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("  gitmap pipeline db")
	fmt.Println("  gitmap pipeline db status")
	fmt.Println("  gitmap pipeline db error-logs")
	fmt.Println("  gitmap pipeline db optimize")
	fmt.Println("  gitmap pipeline db clear -y")
}
