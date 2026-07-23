package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runGroupDelete handles "group delete <name>".
func runGroupDelete(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrGroupNameReq)
		os.Exit(1)
	}
	name := args[0]
	executeGroupDelete(name)
}

// executeGroupDelete opens the DB and deletes the group.
func executeGroupDelete(name string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.DeleteGroup(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgGroupDeleted, name)
}
