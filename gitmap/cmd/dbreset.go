package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runDbReset handles the "db-reset" subcommand.
func runDbReset(args []string) error {
	checkHelp(constants.CmdDbReset, args)
	isConfirm := parseDbResetFlags(args)
	if !isConfirm {
		return abortDbReset()
	}

	appErr := executeDbReset()
	if appErr != nil {
		return appErr
	}

	return nil
}

func abortDbReset() *apperror.AppError {
	printDbResetNoConfirm()

	return apperror.NewWithDetails(
		"dbreset",
		"E4001",
		"db reset canceled: --confirm flag required",
		"cli",
		apperror.ErrorTypeAbort,
		apperror.SeverityWarn,
		nil,
	)
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
func executeDbReset() *apperror.AppError {
	db, err := openDb()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrDBResetFailed)
	}
	defer db.Close()

	err = db.Reset()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrDBResetFailed)
	}

	fmt.Print(constants.MsgDBResetDone)

	return nil
}

// Backwards-compatible aliases for PascalCase callers
//
//nolint:unused
//nolint:unused
var (
	runDBReset            = runDbReset
	parseDBResetFlags     = parseDbResetFlags
	printDBResetNoConfirm = printDbResetNoConfirm
	executeDBReset        = executeDbReset
)
