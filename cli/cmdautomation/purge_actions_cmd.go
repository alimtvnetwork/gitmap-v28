package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	purgeActionsCmd = &cobra.Command{
		Use:     "purge-actions",
		Aliases: []string{"purge-artifacts", "clean-actions"},
		Short:   "Query GitHub Actions API and purge obsolete artifacts to maintain 0.0 GB footprint",
		RunE:    runPurgeActionsCmd,
	}

	purgeActionsOpts PurgeActionsOptions
)

func runPurgeActionsCmd(cmd *cobra.Command, args []string) error {
	monad := RunPurgeActions(purgeActionsOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	renderPurgeResult(monad.Value, purgeActionsOpts.IsJson)
	return nil
}

func renderPurgeResult(res PurgeActionsResult, isJson bool) {
	if isJson {
		printPurgeJson(res)
		return
	}
	printPurgeTerminal(res)
}

func printPurgeJson(res PurgeActionsResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printPurgeTerminal(res PurgeActionsResult) {
	fmt.Printf("\n%s[GitHub Actions Artifact Purge (0.0 GB Mandate)]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Target Repository: %s%s%s\n", constants.ColorCyan, res.Repo, constants.ColorReset)
	fmt.Printf("  Total Artifacts:   %d\n", res.TotalArtifacts)
	fmt.Printf("  Purged Artifacts:  %s%d%s\n", constants.ColorGreen, res.PurgedArtifacts, constants.ColorReset)
	fmt.Printf("  Storage Freed:     %s%.2f MB%s (%d bytes)\n", constants.ColorGreen, res.FreedMB, constants.ColorReset, res.FreedBytes)
	fmt.Printf("  Duration:          %s\n\n", res.Duration)
	if res.TotalArtifacts == 0 {
		fmt.Printf("%s✅ Repository is already at 0.0 GB Actions storage!%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	printPurgedArtifactsList(res.Artifacts)
}

func printPurgedArtifactsList(items []PurgeArtifactItem) {
	limit := 10
	if len(items) < limit {
		limit = len(items)
	}
	for i := 0; i < limit; i++ {
		a := items[i]
		status := "PURGED"
		color := constants.ColorGreen
		if !a.IsPurged {
			status = "SKIPPED"
			color = constants.ColorYellow
		}
		fmt.Printf("  %s[%s]%s ID: %d | %-30s | %.2f MB\n",
			color, status, constants.ColorReset, a.Id, a.Name, float64(a.SizeBytes)/(1024*1024))
	}
	if len(items) > limit {
		fmt.Printf("  ...and %d more artifact(s)\n\n", len(items)-limit)
	}
}

func init() {
	AutomationCmd.AddCommand(purgeActionsCmd)
	initPurgeActionsFlags()
}

func initPurgeActionsFlags() {
	purgeActionsCmd.Flags().StringVar(&purgeActionsOpts.Repo, "repo", "", "Target repository slug (owner/repo)")
	purgeActionsCmd.Flags().IntVar(&purgeActionsOpts.OlderThanDays, "older-than", 0, "Purge artifacts older than N days")
	purgeActionsCmd.Flags().BoolVar(&purgeActionsOpts.IsDryRun, "dry-run", false, "Preview artifacts without deleting")
	purgeActionsCmd.Flags().BoolVar(&purgeActionsOpts.IsJson, "json", false, "Output results as machine-readable JSON")
	purgeActionsCmd.Flags().IntVarP(&purgeActionsOpts.Workers, "workers", "w", 12, "Number of concurrent deletion threads")
}
