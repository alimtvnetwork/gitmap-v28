package cmdpipeline

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func dispatchDbMutateSubcmd(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "clear", "cl":
		return result.MatchWrapper(runPipelineDBClear(rest))
	case "reset":
		return result.MatchWrapper(runPipelineDBReset(rest))
	case "optimize", "opt":
		return result.MatchWrapper(runPipelineDBOptimize(rest))
	default:
		return result.UnmatchedWrapper()
	}
}

func dispatchDbQuerySubcmd(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "status", "st", "s", "info":
		return result.MatchWrapper(runPipelineDBStatusWithTelemetry(rest))
	case "errorlogs", "error-logs", "errors", "err":
		return result.MatchWrapper(runPipelineDBErrorLogs(rest))
	case "help", "-h", "--help":
		printPipelineDBHelp()

		return result.MatchWrapper(nil)
	default:
		return result.UnmatchedWrapper()
	}
}

func dispatchPipelineDBSubcmd(sub string, rest []string) result.ErrorWrapper {
	resQuery := dispatchDbQuerySubcmd(sub, rest)
	if resQuery.IsMatched() {
		return resQuery
	}

	return dispatchDbMutateSubcmd(sub, rest)
}

// handlePipelineDB routes the gitmap pipeline db subcommands.
func handlePipelineDB(args []string) error {
	if len(args) == 0 {
		return runPipelineDBStatusWithTelemetry(nil)
	}

	sub := strings.ToLower(strings.TrimSpace(args[0]))
	res := dispatchPipelineDBSubcmd(sub, args[1:])
	if res.IsMatched() {
		return res.AsError()
	}

	printPipelineDBHelp()

	return apperror.NewValidationError("unknown pipeline db subcommand: " + sub)
}

func printPipelineDBUsage() {
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
}

func printPipelineDBFlagsAndExamples() {
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

func printPipelineDBHelp() {
	printPipelineDBUsage()
	printPipelineDBFlagsAndExamples()
}

func renderPipelineDbInfo(args []string, repo string, info pipelinedb.PipelineDatabaseInfo) error {
	if hasArgFlag(args, "--json") {
		return printJSON(info)
	}

	printPipelineDBStatusMetrics(repo, info)

	return nil
}

func runPipelineDBStatusWithTelemetry(args []string) error {
	repo := resolveCurrentRepoSlug()
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return apperror.WrapSimple(err, "open pipeline split db for status")
	}

	defer db.Close()

	return renderPipelineDbInfo(args, repo, db.GetDatabaseInfo())
}

func printPipelineDBStatusMetrics(repo string, info pipelinedb.PipelineDatabaseInfo) {
	fmt.Println(constants.ColorCyan + "● Pipeline Split Database Summary:" + constants.ColorReset)
	fmt.Printf("  • %-20s %s\n", "Repository:", repo)
	fmt.Printf("  • %-20s %s\n", "Database File:", info.Path)
	fmt.Printf("  • %-20s %s\n", "File Size:", info.HumanSize)
	fmt.Printf("  • %-20s %s\n", "Summary:", pipelinedb.FormatPipelineDbSummary(info))
	fmt.Printf("  • %-20s %d\n", "Total Runs:", info.TotalRuns)
	fmt.Printf("  • %-20s %s%d%s\n", "Failed Runs:", constants.ColorRed, info.FailedRuns, constants.ColorReset)
	fmt.Printf("  • %-20s %d\n", "Error Logs:", info.ErrorCount)
	printPipelineDBLastUpdated(info.LastUpdated)
}

func printPipelineDBLastUpdated(lastUpdated string) {
	if lastUpdated != "" {
		fmt.Printf("  • %-20s %s\n", "Last Synced:", lastUpdated)
	}
}

func resolveTargetRepoSlug(repoSlug string) string {
	if repoSlug != "" {
		return repoSlug
	}

	return resolveCurrentRepoSlug()
}

// GetPipelineDbSummary returns the formatted summary string for the repo's pipeline database.
func GetPipelineDbSummary(repoSlug string) string {
	targetSlug := resolveTargetRepoSlug(repoSlug)
	db, err := pipelinedb.OpenPipelineSplitDb(targetSlug)
	if err != nil {
		return ""
	}

	defer db.Close()

	info := db.GetDatabaseInfo()

	return pipelinedb.FormatPipelineDbSummary(info)
}

// FormatPipelineDbStatusSummary formats a status line for pipeline status output.
func FormatPipelineDbStatusSummary(repoSlug string) string {
	summary := GetPipelineDbSummary(repoSlug)
	if summary == "" {
		return ""
	}

	return fmt.Sprintf("Pipeline Database: %s", summary)
}

// FormatPipelineDbErrorsHeader formats the error logs header with database info.
func FormatPipelineDbErrorsHeader(repoSlug string) string {
	summary := GetPipelineDbSummary(repoSlug)
	if summary == "" {
		return ""
	}

	return fmt.Sprintf("Saved in: %s", summary)
}

// FormatPipelineDbHistoryFooter formats the footer summary for pipeline history.
func FormatPipelineDbHistoryFooter(repoSlug string) string {
	summary := GetPipelineDbSummary(repoSlug)
	if summary == "" {
		return ""
	}

	return fmt.Sprintf("SQLite Telemetry: %s", summary)
}

// PrintPipelineDbTelemetryHeader prints formatted database telemetry header line.
func PrintPipelineDbTelemetryHeader(repoSlug string) {
	header := FormatPipelineDbErrorsHeader(repoSlug)
	if header != "" {
		fmt.Printf("  %s%s%s\n", constants.ColorDim, header, constants.ColorReset)
	}
}

// PrintPipelineDbTelemetryFooter prints formatted database telemetry footer line.
func PrintPipelineDbTelemetryFooter(repoSlug string) {
	footer := FormatPipelineDbHistoryFooter(repoSlug)
	if footer != "" {
		fmt.Printf("\n  %s%s%s\n", constants.ColorDim, footer, constants.ColorReset)
	}
}
