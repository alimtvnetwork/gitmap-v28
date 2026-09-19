package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	smokeTestCmd = &cobra.Command{
		Use:     "smoke-test [dir]",
		Aliases: []string{"installer-smoke", "smoke"},
		Short:   "Cross-platform dry-run validator for installed tools and installer scripts",
		RunE:    runSmokeTestCmd,
	}

	smokeTestOpts SmokeTestOptions
)

func runSmokeTestCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		smokeTestOpts.Dir = args[0]
	}
	monad := RunSmokeTest(smokeTestOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	renderSmokeResult(monad.Value, smokeTestOpts.IsJson)
	if monad.Value.FailedItems > 0 {
		return fmt.Errorf("smoke test failed with %d error(s)", monad.Value.FailedItems)
	}
	return nil
}

func renderSmokeResult(res SmokeTestResult, isJson bool) {
	if isJson {
		printSmokeJson(res)
		return
	}
	printSmokeTerminal(res)
}

func printSmokeJson(res SmokeTestResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printSmokeTerminal(res SmokeTestResult) {
	fmt.Printf("\n%s[Installer & Tool Smoke Tester]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Total Items:  %d\n", res.TotalItems)
	fmt.Printf("  Passed Items: %s%d%s\n", constants.ColorGreen, res.PassedItems, constants.ColorReset)
	fmt.Printf("  Failed Items: %s%d%s\n", resolveSmokeFailColor(res.FailedItems), res.FailedItems, constants.ColorReset)
	fmt.Printf("  Duration:     %s\n\n", res.Duration)
	for _, it := range res.Items {
		printSmokeItem(it)
	}
	printSmokeSummary(res)
}

func resolveSmokeFailColor(failed int) string {
	if failed > 0 {
		return constants.ColorRed
	}
	return constants.ColorReset
}

func printSmokeItem(it SmokeTestItem) {
	symbol := "✅ PASS"
	color := constants.ColorGreen
	if it.IsFail {
		symbol = "❌ FAIL"
		color = constants.ColorRed
	}
	fmt.Printf("  %s%s%s: %-30s [%s] %s\n", color, symbol, constants.ColorReset, it.Name, it.Duration, it.Message)
	for _, issue := range it.Issues {
		fmt.Printf("    %s• %s%s\n", constants.ColorYellow, issue, constants.ColorReset)
	}
}

func printSmokeSummary(res SmokeTestResult) {
	if res.IsPass {
		fmt.Printf("\n%s✅ PASS: All smoke test targets verified successfully.%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	fmt.Printf("\n%s❌ FAIL: Smoke test detected invalid patterns or missing tools.%s\n\n",
		constants.ColorRed, constants.ColorReset)
}

func init() {
	AutomationCmd.AddCommand(smokeTestCmd)
	initSmokeTestFlags()
}

func initSmokeTestFlags() {
	smokeTestCmd.Flags().StringSliceVar(&smokeTestOpts.Tools, "tools", nil, "List of tool binaries to verify in PATH")
	smokeTestCmd.Flags().BoolVar(&smokeTestOpts.IsJson, "json", false, "Output results as machine-readable JSON")
	smokeTestCmd.Flags().IntVarP(&smokeTestOpts.Workers, "workers", "w", 0, "Number of worker threads (default: CPU cores)")
	smokeTestCmd.Flags().StringVarP(&smokeTestOpts.Filter, "filter", "k", "", "Filter targets by name substring")
}
