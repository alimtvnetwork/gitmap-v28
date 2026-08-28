package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runGroupDelete handles "group delete <name>".
func runGroupDelete(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple(constants.ErrGroupNameReq, "E9000")
	}
	name := args[0]
	executeGroupDelete(name)
	return nil
}

// executeGroupDelete opens the DB and deletes the group.
func executeGroupDelete(name string) {
	db, err := openDB()
	if err != nil {
		apperror.WrapSimple(err, constants.ErrListDBFailed)
		return
	}
	defer db.Close()

	err = db.DeleteGroup(name)
	if err != nil {
		apperror.WrapSimple(err, constants.ErrBareFmt)
		return
	}
	fmt.Printf(constants.MsgGroupDeleted, name)
}
