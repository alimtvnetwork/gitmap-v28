package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runGroupDelete handles "group delete <name>".
func runGroupDelete(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple(constants.ErrGroupNameReq, "E9000")
	}

	name := args[0]
	if appErr := executeGroupDelete(name); appErr != nil {
		return appErr
	}

	return nil
}

// executeGroupDelete opens the DB and deletes the group.
func executeGroupDelete(name string) *apperror.AppError {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}

	defer db.Close()

	err = db.DeleteGroup(name)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrBareFmt)
	}

	fmt.Printf(constants.MsgGroupDeleted, name)

	return nil
}
