package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// runBookmarkRun loads a bookmark by name and dispatches the saved command.
func runBookmarkRun(args []string) *apperror.AppError {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrBookmarkRunUsage)
		cliexit.HandleError(nil, 1)
	}

	name := args[0]
	return loadAndDispatchBookmark(name)
}

// loadAndDispatchBookmark fetches the bookmark and runs it.
func loadAndDispatchBookmark(name string) *apperror.AppError {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkQuery+"\n", err)
		cliexit.HandleError(nil, 1)
	}
	defer db.Close()

	bk, err := db.FindBookmarkByName(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkNotFound, name)
		cliexit.HandleError(nil, 1)
	}

	fmt.Printf(constants.MsgBookmarkRunning, bk.Name, bk.Command, bk.Args, bk.Flags)
	replayBookmark(bk.Command, bk.Args, bk.Flags)
	return nil
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
