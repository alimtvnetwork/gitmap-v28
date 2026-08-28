package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runHistoryReset handles the "history-reset" subcommand.
func runHistoryReset(args []string) error {
	checkHelp("history-reset", args)
	confirm := parseHistoryResetFlags(args)
	if confirm {
		executeHistoryReset()

		return nil
	}

	fmt.Fprint(os.Stderr, constants.ErrHistoryResetNoConfirm)
	return apperror.New("fatal error", "E9000", nil)
	return nil
}

// parseHistoryResetFlags parses the --confirm flag.
func parseHistoryResetFlags(args []string) bool {
	fs := flag.NewFlagSet(constants.CmdHistoryReset, flag.ExitOnError)
	confirmFlag := fs.Bool("confirm", false, constants.FlagDescConfirm)
	fs.Parse(args)

	return *confirmFlag
}

// executeHistoryReset opens the database and clears all history.
func executeHistoryReset() {
	db, err := openDB()
	if err != nil {
		return apperror.Wrap(err, constants.ErrHistoryResetFailed, nil)
	}
	defer db.Close()

	err = db.ClearHistory()
	if err != nil {
		return apperror.Wrap(err, constants.ErrHistoryResetFailed, nil)
	}

	fmt.Print(constants.MsgHistoryResetDone)
}
