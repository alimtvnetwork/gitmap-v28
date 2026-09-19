package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	preflightCmd = &cobra.Command{
		Use:     "preflight [dir]",
		Aliases: []string{"check-all", "ci-local", "pre-commit"},
		Short:   "Parallel multi-core CI preflight checks (gofmt, relpaths, nested-if, naming)",
		RunE:    runPreflightCmd,
	}

	preflightOpts PreflightOptions
)

func runPreflightCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		preflightOpts.Dir = args[0]
	}
	monad := RunPreflight(preflightOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	renderPreflightResult(monad.Value, preflightOpts.IsJson)
	if monad.Value.FailedChecks > 0 {
		return fmt.Errorf("preflight failed with %d error(s)", monad.Value.FailedChecks)
	}
	return nil
}

func renderPreflightResult(res PreflightResult, isJson bool) {
	if isJson {
		printPreflightJson(res)
		return
	}
	printPreflightTerminal(res)
}

func printPreflightJson(res PreflightResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printPreflightTerminal(res PreflightResult) {
	fmt.Printf("\n%s[Local CI/CD Preflight Multi-Core Checkers]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Total Checks:  %d\n", res.TotalChecks)
	fmt.Printf("  Passed Checks: %s%d%s\n", constants.ColorGreen, res.PassedChecks, constants.ColorReset)
	fmt.Printf("  Failed Checks: %s%d%s\n", resolveFailColor(res.FailedChecks), res.FailedChecks, constants.ColorReset)
	fmt.Printf("  Duration:      %s\n\n", res.Duration)
	for _, c := range res.Checks {
		printPreflightCheck(c)
	}
	printPreflightSummary(res)
}

func resolveFailColor(failed int) string {
	if failed > 0 {
		return constants.ColorRed
	}
	return constants.ColorReset
}

func printPreflightCheck(c PreflightCheck) {
	symbol := "✅ PASS"
	color := constants.ColorGreen
	if c.IsFail {
		symbol = "❌ FAIL"
		color = constants.ColorRed
	}
	fmt.Printf("  %s%s%s: %-25s [%s] %s\n", color, symbol, constants.ColorReset, c.Name, c.Duration, c.Message)
	for _, d := range c.Details {
		fmt.Printf("    %s• %s%s\n", constants.ColorCyan, d, constants.ColorReset)
	}
}

func printPreflightSummary(res PreflightResult) {
	if res.IsPass {
		fmt.Printf("\n%s✅ PASS: All preflight quality gates cleared successfully.%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	fmt.Printf("\n%s❌ FAIL: Preflight failed. Remediate flagged issues before committing.%s\n\n",
		constants.ColorRed, constants.ColorReset)
}

func init() {
	AutomationCmd.AddCommand(preflightCmd)
	initPreflightFlags()
}

func initPreflightFlags() {
	preflightCmd.Flags().BoolVar(&preflightOpts.IsFailFast, "fail-fast", false, "Halt on first check failure")
	preflightCmd.Flags().BoolVar(&preflightOpts.IsJson, "json", false, "Output results as machine-readable JSON")
	preflightCmd.Flags().IntVarP(&preflightOpts.Workers, "workers", "w", 0, "Number of worker threads (default: CPU cores)")
	preflightCmd.Flags().StringVarP(&preflightOpts.Phase, "phase", "p", "all", "Execution phase (all, test, lint)")
	preflightCmd.Flags().StringVarP(&preflightOpts.Filter, "filter", "k", "", "Filter checks by name substring")
}
