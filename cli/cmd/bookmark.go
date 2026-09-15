package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// runBookmark handles the "bookmark" subcommand routing.
func runBookmark(args []string) error {
	checkHelp("bookmark", args)
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrBookmarkUsage)

		return apperror.NewSimple("fatal error", "E9000")
	}

	sub := args[0]
	rest := args[1:]

	return routeBookmarkSub(sub, rest).AsError()
}

// routeBookmarkSub routes to the appropriate bookmark subcommand.
func routeBookmarkSub(sub string, args []string) result.ErrorWrapper {
	if sub == constants.CmdBookmarkSave {
		return result.MatchWrapperAppErr(runBookmarkSave(args))
	}

	if sub == constants.CmdBookmarkList {
		return result.MatchWrapper(runBookmarkList(args))
	}

	if sub == constants.CmdBookmarkRun {
		return result.MatchWrapperAppErr(runBookmarkRun(args))
	}

	if sub == constants.CmdBookmarkDelete {
		return result.MatchWrapper(runBookmarkDelete(args))
	}

	fmt.Fprint(os.Stderr, constants.ErrBookmarkUsage)

	return result.FailureWrapper(apperror.NewSimple("fatal error", "E9000"))
}
