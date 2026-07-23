package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runBookmarkRun loads a bookmark by name and dispatches the saved command.
func runBookmarkRun(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrBookmarkRunUsage)
		os.Exit(1)
	}

	name := args[0]
	loadAndDispatchBookmark(name)
}

// loadAndDispatchBookmark fetches the bookmark and runs it.
func loadAndDispatchBookmark(name string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkQuery+"\n", err)
		os.Exit(1)
	}
	defer db.Close()

	bk, err := db.FindBookmarkByName(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkNotFound, name)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgBookmarkRunning, bk.Name, bk.Command, bk.Args, bk.Flags)
	replayBookmark(bk.Command, bk.Args, bk.Flags)
}

// replayBookmark reconstructs os.Args and dispatches the command.
func replayBookmark(command, args, flags string) {
	var combined []string
	combined = append(combined, splitNonEmpty(args)...)
	combined = append(combined, splitNonEmpty(flags)...)

	os.Args = buildReplayArgs(command, combined)
	dispatch(command)
}

// splitNonEmpty splits a space-separated string, ignoring empty input.
func splitNonEmpty(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	return strings.Fields(s)
}

// buildReplayArgs constructs the full os.Args for replay.
func buildReplayArgs(command string, extra []string) []string {
	result := []string{"gitmap", command}
	result = append(result, extra...)

	return result
}
