package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	cleanArtifactsCmd = &cobra.Command{
		Use:     "clean-artifacts",
		Aliases: []string{"clean-build", "rm-artifacts"},
		Short:   "Safely deletes build binaries, test dumps, pycache, and temporary files",
		RunE:    runCleanArtifactsCmd,
	}

	cleanArtifactsOpts CleanArtifactsOptions
)

func runCleanArtifactsCmd(cmd *cobra.Command, args []string) error {
	monad := RunCleanArtifacts(cleanArtifactsOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	renderCleanArtifactsResult(monad.Value, cleanArtifactsOpts.IsJson, cleanArtifactsOpts.IsVerbose)
	return nil
}

func renderCleanArtifactsResult(res CleanArtifactsResult, isJson bool, isVerbose bool) {
	if isJson {
		printCleanArtifactsJson(res)
		return
	}
	printCleanArtifactsTerminal(res, isVerbose)
}

func printCleanArtifactsJson(res CleanArtifactsResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printCleanArtifactsTerminal(res CleanArtifactsResult, isVerbose bool) {
	printCleanArtifactSummary(res)
	if res.TotalFound == 0 {
		fmt.Printf("%s✅ No unwanted artifacts found. Repository is clean!%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	printCleanArtifactItems(res.Items, isVerbose)
}

func printCleanArtifactSummary(res CleanArtifactsResult) {
	fmt.Printf("\n%s[Git Hygiene: Build & Test Artifact Remover]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Total Found:   %d item(s)\n", res.TotalFound)
	fmt.Printf("  Total Deleted: %s%d item(s)%s\n", constants.ColorGreen, res.TotalDeleted, constants.ColorReset)
	fmt.Printf("  Storage Freed: %s%.2f MB%s (%d bytes)\n",
		constants.ColorGreen, res.FreedMB, constants.ColorReset, res.FreedBytes)
	fmt.Printf("  Duration:      %s\n\n", res.Duration)
}

func printCleanArtifactItems(items []CleanArtifactItem, isVerbose bool) {
	limit := resolveCleanDisplayLimit(len(items), isVerbose)
	for i := 0; i < limit; i++ {
		printSingleArtifactRow(items[i])
	}
	if len(items) > limit {
		fmt.Printf("  ...and %d more item(s). Use --verbose to show all.\n\n", len(items)-limit)
	}
}

func printSingleArtifactRow(it CleanArtifactItem) {
	status := "DELETED"
	color := constants.ColorGreen
	if it.IsDeleted == false {
		status = "DRY-RUN"
		color = constants.ColorYellow
	}
	tracked := formatTrackedBadge(it.IsGitTracked)
	fmt.Printf("  %s[%-7s]%s %-10s %s (%.2f KB)%s\n",
		color, status, constants.ColorReset, it.Category, it.Path, float64(it.SizeBytes)/1024, tracked)
}

func formatTrackedBadge(isTracked bool) string {
	if isTracked {
		return " [git-tracked]"
	}
	return ""
}

func resolveCleanDisplayLimit(total int, isVerbose bool) int {
	if isVerbose || total <= 20 {
		return total
	}
	return 20
}

func init() {
	AutomationCmd.AddCommand(cleanArtifactsCmd)
	initCleanArtifactsFlags()
}

func initCleanArtifactsFlags() {
	cleanArtifactsCmd.Flags().StringVar(&cleanArtifactsOpts.Dir, "dir", ".", "Root directory to scan for artifacts")
	cleanArtifactsCmd.Flags().BoolVarP(&cleanArtifactsOpts.IsDryRun, "dry-run", "n", false, "Preview matching items without deleting")
	cleanArtifactsCmd.Flags().BoolVarP(&cleanArtifactsOpts.IsAll, "all", "a", false, "Clean all preset categories (pycache, temp, binaries)")
	cleanArtifactsCmd.Flags().BoolVarP(&cleanArtifactsOpts.IsVerbose, "verbose", "v", false, "Show all discovered artifact paths")
	cleanArtifactsCmd.Flags().BoolVar(&cleanArtifactsOpts.IsJson, "json", false, "Output results as machine-readable JSON")
	cleanArtifactsCmd.Flags().BoolVar(&cleanArtifactsOpts.IsCleanPycache, "clean-pycache", false, "Remove python bytecode and test cache")
	cleanArtifactsCmd.Flags().BoolVar(&cleanArtifactsOpts.IsCleanTemp, "clean-temp", false, "Remove temporary files (.tmp, .log, .swp)")
	cleanArtifactsCmd.Flags().BoolVar(&cleanArtifactsOpts.IsCleanBinaries, "clean-binaries", false, "Remove compiled binaries (.exe, .syso, etc.)")
}
