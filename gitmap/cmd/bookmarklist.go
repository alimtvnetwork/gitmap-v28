package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// runBookmarkList displays all saved bookmarks.
func runBookmarkList(args []string) {
	jsonOut := hasJSONFlag(args)
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkQuery+"\n", err)
		os.Exit(1)
	}
	defer db.Close()

	records, err := db.ListBookmarks()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkQuery+"\n", err)
		os.Exit(1)
	}

	if jsonOut {
		printBookmarkJSON(records)

		return
	}

	printBookmarkTerminal(records)
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
func runBookmarkDelete(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrBookmarkDelUsage)
		os.Exit(1)
	}

	name := args[0]
	deleteBookmarkFromDB(name)
}

// deleteBookmarkFromDB removes the bookmark and prints confirmation.
func deleteBookmarkFromDB(name string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkDelete, err)
		os.Exit(1)
	}
	defer db.Close()

	_, findErr := db.FindBookmarkByName(name)
	if findErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkNotFound, name)
		os.Exit(1)
	}

	err = db.DeleteBookmark(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkDelete, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgBookmarkDeleted, name)
}
