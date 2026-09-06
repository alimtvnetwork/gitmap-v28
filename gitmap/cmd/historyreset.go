package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runHistoryReset handles the "history-reset" subcommand.
func runHistoryReset(args []string) error {
	checkHelp(constants.CmdHistoryReset, args)
	isConfirm := parseHistoryResetFlags(args)
	if !isConfirm {
		printHistoryResetNoConfirm()
		cliexit.Exit(1)

		return nil
	}

	executeHistoryReset()

	return nil
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
func executeHistoryReset() {
	db, err := openDb()
	if err != nil {
		apperror.WrapSimple(err, constants.ErrHistoryResetFailed)

		return
	}
	defer db.Close()

	err = db.ClearHistory()
	if err != nil {
		apperror.WrapSimple(err, constants.ErrHistoryResetFailed)

		return
	}

	fmt.Print(constants.MsgHistoryResetDone)
}
