package cmd

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runGroupCreate handles "group create <name>".
func runGroupCreate(args []string) error {
	name, desc, color := parseGroupCreateFlags(args)
	if len(name) == 0 {
		return apperror.New(constants.ErrGroupNameReq, "E9000", nil)
	}
	executeGroupCreate(name, desc, color)
	return nil
}

// executeGroupCreate opens the DB and creates the group.
func executeGroupCreate(name, desc, color string) {
	db, err := openDB()
	if err != nil {
		return apperror.Wrap(err, constants.ErrListDBFailed, nil)
	}
	defer db.Close()

	_, err = db.CreateGroup(name, desc, color)
	if err != nil {
		return apperror.Wrap(err, constants.ErrBareFmt, nil)
	}
	fmt.Printf(constants.MsgGroupCreated, name)
}

// parseGroupCreateFlags parses flags for group create.
func parseGroupCreateFlags(args []string) (name, desc, color string) {
	fs := flag.NewFlagSet(constants.CmdGroupCreate, flag.ExitOnError)
	descFlag := fs.String("description", "", constants.FlagDescGroupDesc)
	colorFlag := fs.String("color", "", constants.FlagDescGroupColor)
	fs.Parse(args)

	if fs.NArg() > 0 {
		name = fs.Arg(0)
	}

	return name, *descFlag, *colorFlag
}
