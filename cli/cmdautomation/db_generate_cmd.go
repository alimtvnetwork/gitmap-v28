package cmdautomation

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	dbGenerateCmd = &cobra.Command{
		Use:     "db-generate [db-path]",
		Aliases: []string{"db-gen", "gen-db"},
		Short:   "Generate Go structs, TypeScript interfaces, and column enums from SQLite tables",
		RunE:    runDbGenerateCmd,
	}

	dbGenOpts DbGenerateOptions
)

func runDbGenerateCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		dbGenOpts.DbPath = args[0]
	}
	monad := RunDbGenerate(dbGenOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	renderDbGenerateResult(monad.Value)
	return nil
}

func renderDbGenerateResult(res DbGenerateResult) {
	fmt.Printf("\n%s[SQLite Struct & Enum Code Generator]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Tables Scanned:  %d\n", res.TableCount)
	fmt.Printf("  Structs Emitted: %d\n", res.StructCount)
	fmt.Printf("  Enums Emitted:   %d\n", res.EnumCount)
	fmt.Printf("  Files Generated: %d\n", len(res.GeneratedFiles))
	fmt.Printf("  Duration:        %s\n", res.Duration)

	if res.IsDryRun {
		fmt.Printf("%s[DRY RUN] Preview mode active - no files written to disk.%s\n\n",
			constants.ColorYellow, constants.ColorReset)
		return
	}
	fmt.Printf("%s✅ Successfully generated code artifacts.%s\n\n",
		constants.ColorGreen, constants.ColorReset)
}

func initDbGenerateFlags() {
	dbGenerateCmd.Flags().StringVar(&dbGenOpts.Lang, "lang", "all", "Target language: go, ts, or all")
	dbGenerateCmd.Flags().StringVar(&dbGenOpts.OutDir, "out", "", "Output directory for generated files")
	dbGenerateCmd.Flags().BoolVar(&dbGenOpts.IsDryRun, "dry-run", false, "Preview generated code without writing")
	dbGenerateCmd.Flags().StringVar(&dbGenOpts.StructDir, "struct-dir", "", "Optional directory of model structs")
}

func init() {
	AutomationCmd.AddCommand(dbGenerateCmd)
	initDbGenerateFlags()
}
