package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runBookmark handles the "bookmark" subcommand routing.
func runBookmark(args []string) {
	checkHelp("bookmark", args)
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrBookmarkUsage)
		os.Exit(1)
	}

	sub := args[0]
	rest := args[1:]

	routeBookmarkSub(sub, rest)
}

// routeBookmarkSub routes to the appropriate bookmark subcommand.
func routeBookmarkSub(sub string, args []string) {
	if sub == constants.CmdBookmarkSave {
		runBookmarkSave(args)

		return
	}
	if sub == constants.CmdBookmarkList {
		runBookmarkList(args)

		return
	}
	if sub == constants.CmdBookmarkRun {
		runBookmarkRun(args)

		return
	}
	if sub == constants.CmdBookmarkDelete {
		runBookmarkDelete(args)

		return
	}

	fmt.Fprint(os.Stderr, constants.ErrBookmarkUsage)
	os.Exit(1)
}
