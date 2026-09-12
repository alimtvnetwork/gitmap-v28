package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runHistoryReset handles the "history-reset" subcommand.
func runHistoryReset(args []string) error {
	checkHelp(constants.CmdHistoryReset, args)
	isConfirm := parseHistoryResetFlags(args)
	if !isConfirm {
		return abortHistoryReset()
	}

	appErr := executeHistoryReset()
	if appErr != nil {
		return appErr
	}

	return nil
}

func abortHistoryReset() *apperror.AppError {
	printHistoryResetNoConfirm()

	return apperror.NewWithDetails(
		"historyreset",
		"E4001",
		"history reset canceled: --confirm flag required",
		"cli",
		apperror.ErrorTypeAbort,
		apperror.SeverityWarn,
		nil,
	)
}

func printHistoryResetNoConfirm() {
	fmt.Println()
	fmt.Printf("  %s⚠ Warning:%s This will clear all command history.\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("  Run with %s--confirm%s to proceed:\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %sgitmap history-reset --confirm%s\n\n", constants.ColorGreen, constants.ColorReset)
}

// parseHistoryResetFlags parses the --confirm flag.
func parseHistoryResetFlags(args []string) bool {
	return parseConfirmFlag(constants.CmdHistoryReset, args)
}

// executeHistoryReset opens the database and clears all history.
func executeHistoryReset() *apperror.AppError {
	db, err := openDb()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrHistoryResetFailed)
	}

	defer db.Close()

	err = db.ClearHistory()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrHistoryResetFailed)
	}

	fmt.Print(constants.MsgHistoryResetDone)

	return nil
}
