package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	helpAuditCmd = &cobra.Command{
		Use:     "help-audit [dir]",
		Aliases: []string{"help-check", "doc-audit"},
		Short:   "Audit CLI commands for help descriptions and documentation parity",
		RunE:    runHelpAuditCmd,
	}

	helpAuditOpts HelpAuditOptions
)

func runHelpAuditCmd(cmd *cobra.Command, args []string) error {
	hasArgs := len(args) > 0
	if hasArgs {
		helpAuditOpts.Dir = args[0]
	}
	monad := RunHelpAudit(helpAuditOpts)
	isFail := monad.IsFailure()
	if isFail {
		return monad.Err
	}
	res := monad.Value
	renderHelpAuditResult(res, helpAuditOpts.IsJson)
	return evaluateHelpAuditExit(res, helpAuditOpts.IsStrict)
}

func evaluateHelpAuditExit(res HelpAuditResult, isStrict bool) error {
	hasViolations := len(res.Violations) > 0
	isStrictViolation := isStrict && hasViolations
	if isStrictViolation {
		return apperror.NewValidationError("CLI help audit violations found under --strict")
	}
	return nil
}

func renderHelpAuditResult(res HelpAuditResult, isJson bool) {
	if isJson {
		printHelpAuditJson(res)
		return
	}
	printHelpAuditTerminal(res)
}

func printHelpAuditJson(res HelpAuditResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	isSuccess := err == nil
	if isSuccess {
		fmt.Println(string(bytes))
	}
}

func printHelpAuditTerminal(res HelpAuditResult) {
	printHelpAuditHeader(res)
	printHelpAuditViolations(res.Violations)
	printHelpAuditStatus(res)
}

func printHelpAuditHeader(res HelpAuditResult) {
	fmt.Printf("\n%s[CLI Command Help & Parity Auditor]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Scanned Files:  %d\n", res.ScannedFiles)
	fmt.Printf("  Total Commands: %d\n", res.TotalCommands)
	fmt.Printf("  Violations:     %d\n", len(res.Violations))
	fmt.Printf("  Duration:       %s\n", res.Duration)
}

func printHelpAuditViolations(vios []HelpViolation) {
	hasVios := len(vios) > 0
	if !hasVios {
		return
	}
	fmt.Printf("\n%sDetected Parity Issues:%s\n", constants.ColorYellow, constants.ColorReset)
	limit := 10
	if len(vios) < limit {
		limit = len(vios)
	}
	for i := 0; i < limit; i++ {
		v := vios[i]
		fmt.Printf("  - %s (%s): %s\n", v.Command, v.File, v.Issue)
	}
	hasMore := len(vios) > limit
	if hasMore {
		fmt.Printf("  ...and %d more issue(s)\n", len(vios)-limit)
	}
}

func printHelpAuditStatus(res HelpAuditResult) {
	if res.IsClean {
		fmt.Printf("\n%s✅ SUCCESS: All CLI commands have complete help text and descriptions.%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	fmt.Printf("\n%s⚠️ WARNING: Found %d CLI command parity violation(s).%s\n\n",
		constants.ColorYellow, len(res.Violations), constants.ColorReset)
}

func init() {
	AutomationCmd.AddCommand(helpAuditCmd)
	initHelpAuditFlags()
}

func initHelpAuditFlags() {
	helpAuditCmd.Flags().BoolVarP(&helpAuditOpts.IsStrict, "strict", "s", false, "Fail with exit code 1 if violations exist")
	helpAuditCmd.Flags().BoolVar(&helpAuditOpts.IsJson, "json", false, "Output results as machine-readable JSON")
	helpAuditCmd.Flags().StringSliceVarP(&helpAuditOpts.Extensions, "ext", "e", nil, "Filter by file extensions")
}
