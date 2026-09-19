package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	specMigrateCmd = &cobra.Command{
		Use:     "spec-migrate [dir]",
		Aliases: []string{"resequence-spec", "migrate-spec"},
		Short:   "Re-sequence specification file prefixes and update cross-links",
		RunE:    runSpecMigrateCmd,
	}

	specMigrateOpts SpecMigrateOptions
)

func runSpecMigrateCmd(cmd *cobra.Command, args []string) error {
	hasArgs := len(args) > 0
	if hasArgs {
		specMigrateOpts.Dir = args[0]
	}
	monad := RunSpecMigrate(specMigrateOpts)
	isFail := monad.IsFailure()
	if isFail {
		return monad.Err
	}
	res := monad.Value
	renderSpecMigrateResult(res, specMigrateOpts.IsJson, specMigrateOpts.IsDryRun)
	return nil
}

func renderSpecMigrateResult(res SpecMigrateResult, isJson, isDryRun bool) {
	if isJson {
		printSpecMigrateJson(res)
		return
	}
	printSpecMigrateTerminal(res, isDryRun)
}

func printSpecMigrateJson(res SpecMigrateResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	hasNoErr := err == nil
	if hasNoErr {
		fmt.Println(string(bytes))
	}
}

func printSpecMigrateTerminal(res SpecMigrateResult, isDryRun bool) {
	printSpecMigrateHeader(res, isDryRun)
	printSpecMigrateChanges(res.Migrated)
	printSpecMigrateStatus(res, isDryRun)
}

func printSpecMigrateHeader(res SpecMigrateResult, isDryRun bool) {
	title := "SPECIFICATION RE-SEQUENCING & CROSS-LINK MIGRATOR"
	if isDryRun {
		title = "SPECIFICATION RE-SEQUENCING (DRY RUN)"
	}
	fmt.Printf("\n%s[%s]%s\n", constants.ColorBold, title, constants.ColorReset)
	fmt.Printf("  Total Specs:    %d\n", res.TotalSpecs)
	fmt.Printf("  Files Migrated: %d\n", len(res.Migrated))
	fmt.Printf("  Refs Updated:   %d\n", res.FilesUpdated)
	fmt.Printf("  Duration:       %s\n", res.Duration)
}

func printSpecMigrateChanges(changes []SpecMigrateChange) {
	hasChanges := len(changes) > 0
	if !hasChanges {
		return
	}
	fmt.Printf("\n%sPlanned Spec Renames & Cross-References:%s\n", constants.ColorCyan, constants.ColorReset)
	for _, c := range changes {
		fmt.Printf("  • %s -> %s (%d references updated)\n", c.SourceFile, c.TargetFile, c.References)
	}
}

func printSpecMigrateStatus(res SpecMigrateResult, isDryRun bool) {
	if isDryRun {
		fmt.Printf("\n%sℹ️ DRY RUN COMPLETE: Spec migrations previewed without file modifications.%s\n\n",
			constants.ColorCyan, constants.ColorReset)
		return
	}
	fmt.Printf("\n%s✅ SUCCESS: Specification file re-sequenced and cross-references updated.%s\n\n",
		constants.ColorGreen, constants.ColorReset)
}

func init() {
	AutomationCmd.AddCommand(specMigrateCmd)
	initSpecMigrateFlags()
}

func initSpecMigrateFlags() {
	specMigrateCmd.Flags().IntVar(&specMigrateOpts.FromNum, "from", 0, "Source spec number to migrate from")
	specMigrateCmd.Flags().IntVar(&specMigrateOpts.ToNum, "to", 0, "Target spec number to migrate to")
	specMigrateCmd.Flags().BoolVarP(&specMigrateOpts.IsDryRun, "dry-run", "d", false, "Preview migration without disk changes")
	specMigrateCmd.Flags().BoolVar(&specMigrateOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
