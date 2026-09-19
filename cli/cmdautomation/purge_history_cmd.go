package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	purgeHistoryCmd = &cobra.Command{
		Use:     "purge-history",
		Aliases: []string{"trace-history", "clean-history"},
		Short:   "Traces and identifies large historical git blobs without rewriting recent commits",
		RunE:    runPurgeHistoryCmd,
	}

	purgeHistoryOpts PurgeHistoryOptions
)

func runPurgeHistoryCmd(cmd *cobra.Command, args []string) error {
	monad := RunPurgeHistory(purgeHistoryOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	renderPurgeHistoryResult(monad.Value, purgeHistoryOpts.IsJson)
	return nil
}

func renderPurgeHistoryResult(res PurgeHistoryResult, isJson bool) {
	if isJson {
		printPurgeHistoryJson(res)
		return
	}
	printPurgeHistoryTerminal(res)
}

func printPurgeHistoryJson(res PurgeHistoryResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printPurgeHistoryTerminal(res PurgeHistoryResult) {
	printPurgeHistoryHeader(res)
	if res.TotalFound == 0 {
		fmt.Printf("%s✅ No large historical git blobs found matching criteria.%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	printPurgeHistoryItems(res.Items)
}

func printPurgeHistoryHeader(res PurgeHistoryResult) {
	fmt.Printf("\n%s[Git History Blob Tracer & Purger]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Blobs Identified: %s%d%s\n", constants.ColorGreen, res.TotalFound, constants.ColorReset)
	fmt.Printf("  Total Size:       %s%.2f MB%s (%d bytes)\n",
		constants.ColorGreen, res.TotalMB, constants.ColorReset, res.TotalBytes)
	fmt.Printf("  Blobs Purged:     %d\n", res.TotalPurged)
	fmt.Printf("  Duration:         %s\n\n", res.Duration)
}

func printPurgeHistoryItems(items []PurgeHistoryItem) {
	limit := resolvePurgeDisplayLimit(len(items))
	for i := 0; i < limit; i++ {
		printSinglePurgeItem(i+1, items[i])
	}
	if len(items) > limit {
		fmt.Printf("  ...and %d more large blob(s).\n\n", len(items)-limit)
	}
	fmt.Println()
}

func resolvePurgeDisplayLimit(total int) int {
	if total > 25 {
		return 25
	}
	return total
}

func printSinglePurgeItem(idx int, it PurgeHistoryItem) {
	shortHash := truncateCommitHash(it.CommitHash)
	status := "IDENTIFIED"
	color := constants.ColorCyan
	if it.IsPurged {
		status = "PURGED"
		color = constants.ColorGreen
	}
	fmt.Printf("  %3d. %s[%-10s]%s %s | %6.2f MB | %s\n",
		idx, color, status, constants.ColorReset, shortHash, it.SizeMB, it.Path)
}

func truncateCommitHash(hash string) string {
	if len(hash) > 8 {
		return hash[:8]
	}
	return hash
}

func init() {
	AutomationCmd.AddCommand(purgeHistoryCmd)
	initPurgeHistoryFlags()
}

func initPurgeHistoryFlags() {
	purgeHistoryCmd.Flags().StringVar(&purgeHistoryOpts.Dir, "dir", ".", "Target repository directory")
	purgeHistoryCmd.Flags().Float64Var(&purgeHistoryOpts.MinSizeMB, "min-size-mb", 1.0, "Minimum blob size in MB to identify")
	purgeHistoryCmd.Flags().StringVar(&purgeHistoryOpts.PathPattern, "path", "", "Target file or path pattern to trace")
	purgeHistoryCmd.Flags().BoolVarP(&purgeHistoryOpts.IsDryRun, "dry-run", "n", true, "Preview large blobs without rewriting history")
	purgeHistoryCmd.Flags().BoolVarP(&purgeHistoryOpts.IsAutoConfirm, "confirm", "y", false, "Bypass confirmation prompt")
	purgeHistoryCmd.Flags().BoolVar(&purgeHistoryOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
