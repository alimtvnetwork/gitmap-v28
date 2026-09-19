package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	versionSyncCmd = &cobra.Command{
		Use:     "version-sync [dir]",
		Aliases: []string{"sync-version", "ver-sync"},
		Short:   "Audit and synchronize repository version manifests",
		RunE:    runVersionSyncCmd,
	}

	versionSyncOpts VersionSyncOptions
)

func runVersionSyncCmd(cmd *cobra.Command, args []string) error {
	hasArgs := len(args) > 0
	if hasArgs {
		versionSyncOpts.Dir = args[0]
	}
	monad := RunVersionSync(versionSyncOpts)
	isFail := monad.IsFailure()
	if isFail {
		return monad.Err
	}
	res := monad.Value
	renderVersionSyncResult(res, versionSyncOpts.IsJson)
	return nil
}

func renderVersionSyncResult(res VersionSyncResult, isJson bool) {
	if isJson {
		printVersionSyncJson(res)
		return
	}
	printVersionSyncTerminal(res)
}

func printVersionSyncJson(res VersionSyncResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	isSuccess := err == nil
	if isSuccess {
		fmt.Println(string(bytes))
	}
}

func printVersionSyncTerminal(res VersionSyncResult) {
	printVersionSyncHeader(res)
	printVersionCheckItems(res.Items)
	printVersionSyncSummary(res)
}

func printVersionSyncHeader(res VersionSyncResult) {
	fmt.Printf("\n%s[Version Synchronization Guard]%s (Canonical: v%s)\n",
		constants.ColorBold, constants.ColorReset, res.CanonicalVersion)
	fmt.Printf("  Total Checked:  %d\n", res.TotalChecked)
	fmt.Printf("  Mismatches:     %d\n", res.MismatchCount)
	fmt.Printf("  Fixed:          %d\n", res.FixedCount)
	fmt.Printf("  Duration:       %s\n", res.Duration)
}

func printVersionCheckItems(items []VersionCheckItem) {
	fmt.Println()
	for _, item := range items {
		printSingleCheckItem(item)
	}
}

func printSingleCheckItem(item VersionCheckItem) {
	isMatch := item.IsMatch
	if isMatch {
		fmt.Printf("  %s✓%s %s: %s\n", constants.ColorGreen, constants.ColorReset, item.Name, item.Message)
		return
	}
	fmt.Printf("  %s✗%s %s: %s\n", constants.ColorRed, constants.ColorReset, item.Name, item.Message)
}

func printVersionSyncSummary(res VersionSyncResult) {
	isClean := res.IsClean
	if isClean {
		fmt.Printf("\n%s✅ PASS: All manifests are 100%% synchronized to v%s.%s\n\n",
			constants.ColorGreen, res.CanonicalVersion, constants.ColorReset)
		return
	}
	fmt.Printf("\n%s⚠️ FAIL: %d manifest(s) out of sync. Use --fix to harmonize.%s\n\n",
		constants.ColorYellow, res.MismatchCount, constants.ColorReset)
}

func getVersionSyncOptions() VersionSyncOptions {
	return versionSyncOpts
}

func init() {
	AutomationCmd.AddCommand(versionSyncCmd)
	initVersionSyncFlags()
}

func initVersionSyncFlags() {
	versionSyncCmd.Flags().BoolVarP(&versionSyncOpts.IsFixMode, "fix", "f", false, "Automatically harmonize mismatched manifests")
	versionSyncCmd.Flags().StringVarP(&versionSyncOpts.TargetVersion, "target-version", "t", "", "Override canonical version to harmonize")
	versionSyncCmd.Flags().BoolVar(&versionSyncOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
