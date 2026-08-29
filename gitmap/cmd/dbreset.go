package cmd

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runDBReset handles the "db-reset" subcommand.
func runDBReset(args []string) error {
	checkHelp("db-reset", args)
	confirm := parseDBResetFlags(args)
	if confirm {
		executeDBReset()

		return nil
	}

	return apperror.NewSimple(constants.ErrDBResetNoConfirm, "E9000")
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
