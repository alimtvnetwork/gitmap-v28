package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runBookmarkList displays all saved bookmarks.
func runBookmarkList(args []string) error {
	jsonOut := hasJSONFlag(args)
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, "constants.ErrBookmarkQuery+")
	}
	defer db.Close()

	records, err := db.ListBookmarks()
	if err != nil {
		return apperror.WrapSimple(err, "constants.ErrBookmarkQuery+")
	}

	if jsonOut {
		printBookmarkJSON(records)

		return nil
	}

	printBookmarkTerminal(records)
	return nil
}

// hasJSONFlag checks if --json is present in args.
func hasJSONFlag(args []string) bool {
	for _, a := range args {
		if a == "--json" {
			return true
		}
	}

	return false
}

// printBookmarkTerminal prints bookmarks as a table.
func printBookmarkTerminal(records []model.BookmarkRecord) {
	if len(records) == 0 {
		fmt.Print(constants.MsgBookmarkEmpty)

		return
	}

	fmt.Println(constants.MsgBookmarkColumns)
	for _, r := range records {
		fmt.Printf(constants.MsgBookmarkRowFmt, r.Name, r.Command, r.Args, r.Flags)
	}
}

// printBookmarkJSON outputs bookmarks as JSON.
func printBookmarkJSON(records []model.BookmarkRecord) {
	if err := encodeBookmarkListJSON(os.Stdout, records); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to encode bookmarks to JSON: %v\n", err)
	}
}

// runBookmarkDelete removes a saved bookmark by name.
func runBookmarkDelete(args []string) error {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrBookmarkDelUsage)
		return apperror.NewSimple("fatal error", "E9000")
	}

	name := args[0]
	return deleteBookmarkFromDB(name)
}

// deleteBookmarkFromDB removes the bookmark and prints confirmation.
func deleteBookmarkFromDB(name string) *apperror.AppError {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrBookmarkDelete)
	}
	defer db.Close()

	_, findErr := db.FindBookmarkByName(name)
	if findErr != nil {
		return apperror.WrapSimple(findErr, constants.ErrBookmarkNotFound)
	}

	err = db.DeleteBookmark(name)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrBookmarkDelete)
	}

	fmt.Printf(constants.MsgBookmarkDeleted, name)
	return nil
}
