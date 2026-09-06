package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runDbReset handles the "db-reset" subcommand.
func runDbReset(args []string) error {
	checkHelp(constants.CmdDbReset, args)
	isConfirm := parseDbResetFlags(args)
	if !isConfirm {
		printDbResetNoConfirm()
		cliexit.Exit(1)

		return nil
	}

	executeDbReset()

	return nil
}

func printDbResetNoConfirm() {
	fmt.Println()
	fmt.Printf("  %s⚠ Warning:%s This will delete all tracked repos and groups.\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("  Run with %s--confirm%s to proceed:\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sgitmap db-reset --confirm%s\n\n", constants.ColorGreen, constants.ColorReset)
}

// parseDbResetFlags parses the --confirm flag.
func parseDbResetFlags(args []string) bool {
	return parseConfirmFlag(constants.CmdDbReset, args)
}

// executeDbReset opens the database, resets it, and prints confirmation.
func executeDbReset() {
	db, err := openDb()
	if err != nil {
		apperror.WrapSimple(err, constants.ErrDBResetFailed)

		return
	}
	defer db.Close()

	err = db.Reset()
	if err != nil {
		apperror.WrapSimple(err, constants.ErrDBResetFailed)

		return
	}

	fmt.Print(constants.MsgDBResetDone)
}

// Backwards-compatible aliases for PascalCase callers
var (
	runDBReset            = runDbReset
	parseDBResetFlags     = parseDbResetFlags
	printDBResetNoConfirm = printDbResetNoConfirm
	executeDBReset        = executeDbReset
)
