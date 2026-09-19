package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	testInventoryCmd = &cobra.Command{
		Use:     "test-inventory [dir]",
		Aliases: []string{"tests-inv", "inventory"},
		Short:   "Discover unit and integration tests, output inventory manifest with durations",
		RunE:    runTestInventoryCmd,
	}

	testInventoryOpts TestInventoryOptions
)

func runTestInventoryCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		testInventoryOpts.Dir = args[0]
	}
	monad := RunTestInventory(testInventoryOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	renderTestInventoryResult(monad.Value, testInventoryOpts.IsJson)
	return nil
}

func renderTestInventoryResult(res TestInventoryResult, isJson bool) {
	if isJson {
		printInventoryJson(res)
		return
	}
	printInventoryTerminal(res)
}

func printInventoryJson(res TestInventoryResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printInventoryTerminal(res TestInventoryResult) {
	fmt.Printf("\n%s[Test Inventory Manifest Generator]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Total Tests:    %d\n", res.TotalTests)
	fmt.Printf("  Packages:       %d\n", res.Summary.Packages)
	fmt.Printf("  Fast Tests:     %s%d%s (est. %.2fs)\n", constants.ColorGreen, res.Summary.FastTests, constants.ColorReset, res.Summary.EstimatedFastSec)
	fmt.Printf("  Slow Tests:     %s%d%s (est. %.2fs, >%.1fs threshold)\n", constants.ColorYellow, res.Summary.SlowTests, constants.ColorReset, res.Summary.EstimatedSlowSec, res.Summary.SlowThresholdSec)
	fmt.Printf("  Manifest Saved: %s%s%s\n", constants.ColorCyan, res.OutPath, constants.ColorReset)
	fmt.Printf("  Duration:       %s\n\n", res.Duration)
	printInventorySample(res)
}

func printInventorySample(res TestInventoryResult) {
	limit := 10
	count := 0
	for _, t := range res.Tests {
		if count >= limit {
			break
		}
		color := constants.ColorGreen
		if t.IsSlow {
			color = constants.ColorYellow
		}
		fmt.Printf("  %s[%-4s]%s %s (%.3fs)\n", color, t.Tier, constants.ColorReset, t.Id, t.DurationSec)
		count++
	}
	if len(res.Tests) > limit {
		fmt.Printf("  ...and %d more test(s)\n\n", len(res.Tests)-limit)
	}
}

func init() {
	AutomationCmd.AddCommand(testInventoryCmd)
	initTestInventoryFlags()
}

func initTestInventoryFlags() {
	testInventoryCmd.Flags().StringVar(&testInventoryOpts.OutPath, "out", "", "Output manifest path (default: .ai-memory/test-inventory.json)")
	testInventoryCmd.Flags().BoolVar(&testInventoryOpts.IsRefresh, "refresh", false, "Force refresh existing inventory manifest")
	testInventoryCmd.Flags().BoolVar(&testInventoryOpts.IsJson, "json", false, "Output results as machine-readable JSON")
	testInventoryCmd.Flags().Float64Var(&testInventoryOpts.SlowThreshold, "slow-threshold", 4.0, "Threshold in seconds to categorize slow tests")
	testInventoryCmd.Flags().BoolVar(&testInventoryOpts.IsForceRunAll, "force-run-all", false, "Mark all discovered tests as needing run")
	testInventoryCmd.Flags().StringSliceVar(&testInventoryOpts.RecordPaths, "record", nil, "Record modified file paths safely")
	testInventoryCmd.Flags().BoolVar(&testInventoryOpts.IsQueryRecent, "query-recent", false, "Query tests associated with recent changes")
	testInventoryCmd.Flags().BoolVar(&testInventoryOpts.IsClear, "clear", false, "Clear recent changes log")
}
