package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	resultWrapperCmd = &cobra.Command{
		Use:     "result-wrapper [dir]",
		Aliases: []string{"res-wrapper", "result-audit"},
		Short:   "Audit Go functions returning multi-value map/slice error tuples",
		RunE:    runResultWrapperCmd,
	}

	paramsCmd = &cobra.Command{
		Use:     "params [dir]",
		Aliases: []string{"param-audit", "arity"},
		Short:   "Audit function signatures for argument reduction and parameter structs",
		RunE:    runParamsCmd,
	}

	enumsCmd = &cobra.Command{
		Use:     "enums [dir]",
		Aliases: []string{"enum-audit", "types-enum"},
		Short:   "Audit enums for *Type suffixes and ban raw numeric rune casts",
		RunE:    runEnumsCmd,
	}

	ruleOpts RuleAuditOptions
)

func runResultWrapperCmd(cmd *cobra.Command, args []string) error {
	ruleOpts.RuleName = "result-wrapper"
	return executeRuleAudit(args)
}

func runParamsCmd(cmd *cobra.Command, args []string) error {
	ruleOpts.RuleName = "params"
	return executeRuleAudit(args)
}

func runEnumsCmd(cmd *cobra.Command, args []string) error {
	ruleOpts.RuleName = "enums"
	return executeRuleAudit(args)
}

func executeRuleAudit(args []string) error {
	if len(args) > 0 {
		ruleOpts.Dir = args[0]
	}
	monad := RunRuleAudit(ruleOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	res := monad.Value
	renderRuleAuditResult(res, ruleOpts.IsJson)
	return nil
}

func renderRuleAuditResult(res RuleAuditResult, isJson bool) {
	if isJson {
		printRuleAuditJson(res)
		return
	}
	printRuleAuditTerminal(res)
}

func printRuleAuditJson(res RuleAuditResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printRuleAuditTerminal(res RuleAuditResult) {
	fmt.Printf("\n%s[Coding Guideline Rule Audit: %s]%s\n", constants.ColorBold, ruleOpts.RuleName, constants.ColorReset)
	fmt.Printf("  Scanned Files: %d\n", res.ScannedFiles)
	fmt.Printf("  Violations:    %d\n", len(res.Violations))
	fmt.Printf("  Duration:      %s\n", res.Duration)

	if len(res.Violations) > 0 {
		printRuleViolations(res.Violations)
		return
	}
	fmt.Printf("%s✅ PASS: Zero rule violations found for '%s'.%s\n\n",
		constants.ColorGreen, ruleOpts.RuleName, constants.ColorReset)
}

func printRuleViolations(vios []RuleViolation) {
	fmt.Printf("\n%s❌ Rule Violations:%s\n", constants.ColorRed, constants.ColorReset)
	limit := 10
	if len(vios) < limit {
		limit = len(vios)
	}
	for i := 0; i < limit; i++ {
		v := vios[i]
		fmt.Printf("  %s%s:%d%s [%s] %s\n", constants.ColorCyan, v.File, v.LineNumber,
			constants.ColorReset, v.Rule, v.Description)
	}
	if len(vios) > limit {
		fmt.Printf("  ...and %d more violation(s)\n", len(vios)-limit)
	}
	fmt.Println()
}

func initRuleFlags() {
	resultWrapperCmd.Flags().BoolVar(&ruleOpts.IsJson, "json", false, "Output results as machine-readable JSON")
	paramsCmd.Flags().BoolVar(&ruleOpts.IsJson, "json", false, "Output results as machine-readable JSON")
	enumsCmd.Flags().BoolVar(&ruleOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
