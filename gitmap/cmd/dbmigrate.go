package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runDBMigrate handles the "db-migrate" (alias "dbm") subcommand.
//
// It opens the active-profile database, runs Migrate() (which is idempotent
// and safe to invoke repeatedly), and prints a single-line summary. The
// --verbose flag prints every migration step that ran.
func runDBMigrate(args []string) {
	checkHelp(constants.CmdDBMigrate, args)
	verbose := parseDBMigrateFlags(args)

	fmt.Print(constants.MsgDBMigrateRunning)

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDBMigrateFailFmt, err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDBMigrateFailFmt, err)
		os.Exit(1)
	}

	printDBMigrateSummary(verbose)
}

// parseDBMigrateFlags extracts the --verbose flag.
func parseDBMigrateFlags(args []string) bool {
	fs := flag.NewFlagSet(constants.CmdDBMigrate, flag.ExitOnError)
	v := fs.Bool(constants.FlagDBMigrateVerbose, false, constants.FlagDescDBMigrateV)

	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		os.Exit(2)
	}

	return *v
}

// printDBMigrateSummary writes the post-run summary line.
//
// Migrate() streams every per-step warning to os.Stderr already (with the
// table + column + action context). If any warning was printed, the user
// has already seen it; here we just confirm the run reached the end.
func printDBMigrateSummary(verbose bool) {
	fmt.Print(constants.MsgDBMigrateNoWork)

	if verbose {
		fmt.Println("    (verbose: every CREATE/ALTER is idempotent — re-running has no effect)")
		fmt.Println("    (any per-step warnings above include the offending table + column)")
	}
}

// runPostUpdateMigrate is invoked from the update flow after the binary is
// replaced. It is best-effort: any failure is warned, never fatal, since the
// user may have an in-flight DB lock or read-only environment.
func runPostUpdateMigrate() {
	fmt.Print(constants.MsgDBMigratePostUpdate)

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnDBMigratePostFail, err)

		return
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnDBMigratePostFail, err)

		return
	}

	fmt.Println("  ✓ Schema migrations complete.")
}
