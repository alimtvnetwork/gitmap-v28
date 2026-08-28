package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runBookmark handles the "bookmark" subcommand routing.
func runBookmark(args []string) *apperror.AppError {
	checkHelp("bookmark", args)
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrBookmarkUsage)
		return apperror.New("fatal error", "E9000", nil)
	}

	sub := args[0]
	rest := args[1:]

	return routeBookmarkSub(sub, rest)
}

// routeBookmarkSub routes to the appropriate bookmark subcommand.
func routeBookmarkSub(sub string, args []string) *apperror.AppError {
	if sub == constants.CmdBookmarkSave {
		runBookmarkSave(args)
		return nil
	}
	if sub == constants.CmdBookmarkList {
		runBookmarkList(args)
		return nil
	}
	if sub == constants.CmdBookmarkRun {
		runBookmarkRun(args)
		return nil
	}
	if sub == constants.CmdBookmarkDelete {
		runBookmarkDelete(args)
		return nil
	}

	fmt.Fprint(os.Stderr, constants.ErrBookmarkUsage)
	return apperror.New("fatal error", "E9000", nil)
}
