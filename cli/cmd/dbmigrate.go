package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

// runDbMigrate handles the "db-migrate" (alias "dbm") subcommand.
//
// It opens the active-profile database, runs Migrate() (which is idempotent
// and safe to invoke repeatedly), and prints a single-line summary. The
// --verbose flag prints every migration step that ran.
func runDbMigrate(args []string) error {
	checkHelp(constants.CmdDBMigrate, args)
	isVerbose := parseDbMigrateFlags(args)

	fmt.Print(constants.MsgDBMigrateRunning)

	db, err := openDb()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrDBMigrateFailFmt)
	}

	defer db.Close()

	if err := db.Migrate(); err != nil {
		return apperror.WrapSimple(err, constants.ErrDBMigrateFailFmt)
	}

	printDbMigrateSummary(isVerbose)

	return nil
}

// Backwards-compatible alias for runDbMigrate
var runDBMigrate = runDbMigrate

// parseDbMigrateFlags extracts the --verbose flag.
func parseDbMigrateFlags(args []string) bool {
	fs := flag.NewFlagSet(constants.CmdDBMigrate, flag.ExitOnError)
	var isVerbose bool
	fs.BoolVar(&isVerbose, constants.FlagDBMigrateVerbose, false, constants.FlagDescDBMigrateV)

	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		cliexit.HandleError(nil, 2)
	}

	return isVerbose
}

// Backwards-compatible alias for parseDbMigrateFlags

// printDbMigrateSummary writes the post-run summary line.
//
// Migrate() streams every per-step warning to os.Stderr already (with the
// table + column + action context). If any warning was printed, the user
// has already seen it; here we just confirm the run reached the end.
func printDbMigrateSummary(isVerbose bool) {
	fmt.Print(constants.MsgDBMigrateNoWork)

	if isVerbose {
		fmt.Println("    (verbose: every CREATE/ALTER is idempotent — re-running has no effect)")
		fmt.Println("    (any per-step warnings above include the offending table + column)")
	}
}

// Backwards-compatible alias for printDbMigrateSummary

// runPostUpdateMigrate is invoked from the update flow after the binary is
// replaced. It is best-effort: any failure is warned, never fatal, since the
// user may have an in-flight DB lock or read-only environment.
func runPostUpdateMigrate() error {
	fmt.Print(constants.MsgDBMigratePostUpdate)

	db, err := openDb()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnDBMigratePostFail, err)

		return nil
	}

	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnDBMigratePostFail, err)

		return nil
	}

	fmt.Println("  ✓ Schema migrations complete.")

	return nil
}
