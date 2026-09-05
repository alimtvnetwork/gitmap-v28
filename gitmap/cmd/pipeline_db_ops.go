package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/pipelinedb"
)

func runPipelineDBStatus(args []string) error {
	repo := resolveCurrentRepoSlug()
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return err
	}
	defer db.Close()

	stats, err := db.GetStats()
	if err != nil {
		return err
	}

	if hasArgFlag(args, "--json") {
		return printJSON(stats)
	}

	fmt.Println(constants.ColorCyan + "● Pipeline Split Database Summary:" + constants.ColorReset)
	fmt.Printf("  • %-20s %s\n", "Repository:", repo)
	fmt.Printf("  • %-20s %s\n", "Database File:", stats.Path)
	fmt.Printf("  • %-20s %s\n", "File Size:", formatBytes(int64(stats.Size)))
	fmt.Printf("  • %-20s %d\n", "Total Runs:", stats.TotalRuns)
	fmt.Printf("  • %-20s %s%d%s\n", "Success Runs:", constants.ColorGreen, stats.SuccessRuns, constants.ColorReset)
	fmt.Printf("  • %-20s %s%d%s\n", "Failed Runs:", constants.ColorRed, stats.FailedRuns, constants.ColorReset)
	fmt.Printf("  • %-20s %d\n", "Error Logs:", stats.ErrorLogCount)
	fmt.Printf("  • %-20s %d\n", "Tracked Segments:", stats.SegmentCount)
	if stats.LastUpdated != "" {
		fmt.Printf("  • %-20s %s\n", "Last Synced:", stats.LastUpdated)
	}
	return nil
}

func runPipelineDBClear(args []string) error {
	repo := resolveCurrentRepoSlug()
	msg := fmt.Sprintf("Clear all pipeline runs and error logs for %s? [y/N]: ", repo)
	if !confirmOrSkip(msg, args) {
		fmt.Println("Clear operation canceled.")
		return nil
	}
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Clear(); err != nil {
		return err
	}
	fmt.Printf("%s✓ Pipeline split database cleared for %s.%s\n", constants.ColorGreen, repo, constants.ColorReset)
	return nil
}

func runPipelineDBReset(args []string) error {
	repo := resolveCurrentRepoSlug()
	msg := fmt.Sprintf("Reset and re-create pipeline schema for %s? [y/N]: ", repo)
	if !confirmOrSkip(msg, args) {
		fmt.Println("Reset operation canceled.")
		return nil
	}
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Reset(); err != nil {
		return err
	}
	fmt.Printf("%s✓ Pipeline split database reset for %s.%s\n", constants.ColorGreen, repo, constants.ColorReset)
	return nil
}

func runPipelineDBOptimize(args []string) error {
	repo := resolveCurrentRepoSlug()
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return err
	}
	defer db.Close()

	reclaimed, err := db.Optimize()
	if err != nil {
		return err
	}
	fmt.Printf("%s✓ Pipeline split DB optimized.%s Reclaimed: %s (%s)\n",
		constants.ColorGreen, constants.ColorReset, formatBytes(reclaimed), db.Path)
	return nil
}

func runPipelineDBErrorLogs(args []string) error {
	repo := resolveCurrentRepoSlug()
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return err
	}
	defer db.Close()

	logs, err := db.QueryRecentErrorLogs(20)
	if err != nil {
		return err
	}
	if hasArgFlag(args, "--json") {
		return printJSON(logs)
	}

	if len(logs) == 0 {
		fmt.Printf("No error logs recorded in pipeline database for %s.\n", repo)
		return nil
	}

	fmt.Printf("\n  %s● Recorded CI/CD Error Logs for %s (%d records):%s\n",
		constants.ColorRed, repo, len(logs), constants.ColorReset)
	fmt.Printf("    %s\n", strings.Repeat("─", 78))
	for _, l := range logs {
		fmt.Printf("    [%s] Run #%d (%s - %s)\n", l.CreatedAt, l.RunId, l.WorkflowName, l.StepName)
		for _, line := range strings.Split(l.ErrorText, "\n") {
			fmt.Printf("      %s\n", line)
		}
		fmt.Println()
	}
	fmt.Printf("    %s\n", strings.Repeat("─", 78))
	return nil
}
