package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	schemaAuditCmd = &cobra.Command{
		Use:     "schema-audit [db-path]",
		Aliases: []string{"audit-schema", "db-audit"},
		Short:   "Validate SQLite split-db schemas against naming conventions (PascalCase, affirmative booleans, PKs)",
		RunE:    runSchemaAuditCmd,
	}

	schemaAuditOpts SchemaAuditOptions
)

func runSchemaAuditCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		schemaAuditOpts.DbPath = args[0]
	}
	monad := RunSchemaAudit(schemaAuditOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	res := monad.Value
	renderSchemaAuditResult(res, schemaAuditOpts.IsJson)
	return nil
}

func renderSchemaAuditResult(res SchemaAuditResult, isJson bool) {
	if isJson {
		printSchemaAuditJson(res)
		return
	}
	printSchemaAuditTerminal(res)
}

func printSchemaAuditJson(res SchemaAuditResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printSchemaAuditTerminal(res SchemaAuditResult) {
	fmt.Printf("\n%s[SQLite Schema Convention Audit]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Files Scanned:   %d\n", res.ScannedFiles)
	fmt.Printf("  Tables Checked:  %d\n", res.TableCount)
	fmt.Printf("  Violations:      %d\n", len(res.Violations))
	fmt.Printf("  Duration:        %s\n", res.Duration)

	if len(res.Violations) > 0 {
		printSchemaViolationsList(res.Violations)
		return
	}
	fmt.Printf("%s✅ PASS: All schemas adhere to PascalCase, affirmative boolean, and PK rules.%s\n\n",
		constants.ColorGreen, constants.ColorReset)
}

func printSchemaViolationsList(vios []SchemaViolation) {
	fmt.Printf("\n%sSchema Violations Found:%s\n", constants.ColorRed, constants.ColorReset)
	for _, v := range vios {
		fmt.Printf("  • [%s] %s:%s ➔ %s\n", v.Severity, v.File, v.Table, v.Issue)
	}
	fmt.Println()
}

func initSchemaAuditFlags() {
	schemaAuditCmd.Flags().BoolVar(&schemaAuditOpts.IsJson, "json", false, "Output results as machine-readable JSON")
	schemaAuditCmd.Flags().BoolVar(&schemaAuditOpts.IsStrict, "strict", false, "Enforce strict column naming and PK rules")
	schemaAuditCmd.Flags().StringVar(&schemaAuditOpts.Dir, "dir", "", "Target directory containing SQL or Go schemas")
}

func init() {
	AutomationCmd.AddCommand(schemaAuditCmd)
	initSchemaAuditFlags()
}
