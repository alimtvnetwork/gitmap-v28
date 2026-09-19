package cmdautomation

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	dbMigrateCmd = &cobra.Command{
		Use:     "db-migrate [db-path] [migrations-dir]",
		Aliases: []string{"migrate", "db-up"},
		Short:   "Execute ordered SQL migrations with rollback protection and tracking",
		RunE:    runDbMigrateCmd,
	}

	dbMigrateOpts DbMigrateOptions
)

func runDbMigrateCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		dbMigrateOpts.DbPath = args[0]
	}
	if len(args) > 1 {
		dbMigrateOpts.MigrationsDir = args[1]
	}
	monad := RunDbMigrate(dbMigrateOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	renderDbMigrateResult(monad.Value)
	return nil
}

func renderDbMigrateResult(res DbMigrateResult) {
	fmt.Printf("\n%s[SQLite Migration Engine]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Total Applied:   %d\n", res.TotalApplied)
	fmt.Printf("  Latest Version:  %s\n", res.CurrentVersion)
	fmt.Printf("  Duration:        %s\n", res.Duration)

	if len(res.AppliedMigrations) > 0 {
		printAppliedMigrationsList(res.AppliedMigrations)
	}

	if res.IsDryRun {
		fmt.Printf("%s[DRY RUN] Preview complete - no database changes committed.%s\n\n",
			constants.ColorYellow, constants.ColorReset)
		return
	}
	fmt.Printf("%s✅ Migrations successfully synchronized.%s\n\n",
		constants.ColorGreen, constants.ColorReset)
}

func printAppliedMigrationsList(migrations []string) {
	fmt.Printf("\n%sApplied Migrations:%s\n", constants.ColorCyan, constants.ColorReset)
	for _, m := range migrations {
		fmt.Printf("  • %s\n", m)
	}
	fmt.Println()
}

func initDbMigrateFlags() {
	dbMigrateCmd.Flags().BoolVar(&dbMigrateOpts.IsDryRun, "dry-run", false, "Preview migrations without writing")
	dbMigrateCmd.Flags().BoolVar(&dbMigrateOpts.IsStatus, "status", false, "Display migration history and current version")
	dbMigrateCmd.Flags().IntVar(&dbMigrateOpts.RollbackStep, "rollback", 0, "Rollback last N migrations")
	dbMigrateCmd.Flags().StringVar(&dbMigrateOpts.SqlScript, "sql", "", "Execute inline SQL migration statement")
}

func init() {
	AutomationCmd.AddCommand(dbMigrateCmd)
	initDbMigrateFlags()
}
