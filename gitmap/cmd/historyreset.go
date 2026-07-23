package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runHistoryReset handles the "history-reset" subcommand.
func runHistoryReset(args []string) {
	checkHelp("history-reset", args)
	confirm := parseHistoryResetFlags(args)
	if confirm {
		executeHistoryReset()

		return
	}

	fmt.Fprint(os.Stderr, constants.ErrHistoryResetNoConfirm)
	os.Exit(1)
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
		fmt.Fprintf(os.Stderr, constants.ErrHistoryResetFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.ClearHistory()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrHistoryResetFailed, err)
		os.Exit(1)
	}

	fmt.Print(constants.MsgHistoryResetDone)
}
