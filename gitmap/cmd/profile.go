package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runProfile handles the "profile" subcommand routing.
func runProfile(args []string) error {
	checkHelp("profile", args)
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrProfileUsage)
		return apperror.New("fatal error", "E9000", nil)
	}

	sub := args[0]
	rest := args[1:]

	routeProfileSub(sub, rest)
	return nil
}

// routeProfileSub routes to the appropriate profile subcommand.
func routeProfileSub(sub string, args []string) {
	if sub == constants.CmdProfileCreate {
		runProfileCreate(args)

		return
	}
	if sub == constants.CmdProfileList {
		runProfileList()

		return
	}
	if sub == constants.CmdProfileSwitch {
		runProfileSwitch(args)

		return
	}
	if sub == constants.CmdProfileDelete {
		runProfileDelete(args)

		return
	}
	if sub == constants.CmdProfileShow {
		runProfileShow()

		return
	}

	fmt.Fprint(os.Stderr, constants.ErrProfileUsage)
	return apperror.New("fatal error", "E9000", nil)
}
