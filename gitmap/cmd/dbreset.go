package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runDBReset handles the "db-reset" subcommand.
func runDBReset(args []string) {
	checkHelp("db-reset", args)
	confirm := parseDBResetFlags(args)
	if confirm {
		executeDBReset()

		return
	}

	fmt.Fprintln(os.Stderr, constants.ErrDBResetNoConfirm)
	os.Exit(1)
}

// parseDBResetFlags parses the --confirm flag.
func parseDBResetFlags(args []string) bool {
	fs := flag.NewFlagSet(constants.CmdDBReset, flag.ExitOnError)
	confirmFlag := fs.Bool("confirm", false, constants.FlagDescConfirm)
	fs.Parse(args)

	return *confirmFlag
}

// executeDBReset opens the database, resets it, and prints confirmation.
func executeDBReset() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDBResetFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	err = db.Reset()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDBResetFailed, err)
		os.Exit(1)
	}

	fmt.Print(constants.MsgDBResetDone)
}
