package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runBookmarkRun loads a bookmark by name and dispatches the saved command.
func runBookmarkRun(args []string) *apperror.AppError {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrBookmarkRunUsage)

		return apperror.NewValidationError(constants.ErrBookmarkRunUsage)
	}

	return loadAndDispatchBookmark(args[0])
}

// loadAndDispatchBookmark fetches the bookmark and runs it.
func loadAndDispatchBookmark(name string) *apperror.AppError {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkQuery+"\n", err)

		return apperror.WrapSimple(err, constants.ErrBookmarkQuery)
	}

	defer db.Close()

	return findAndReplayBookmark(db, name)
}

func findAndReplayBookmark(db storeBookmarkReader, name string) *apperror.AppError {
	bk, err := db.FindBookmarkByName(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkNotFound, name)

		return apperror.NewSimple(fmt.Sprintf(constants.ErrBookmarkNotFound, name), "E9000")
	}

	fmt.Printf(constants.MsgBookmarkRunning, bk.Name, bk.Command, bk.Args, bk.Flags)
	replayBookmark(bk.Command, bk.Args, bk.Flags)

	return nil
}

type storeBookmarkReader interface {
	FindBookmarkByName(string) (model.BookmarkRecord, error)
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
