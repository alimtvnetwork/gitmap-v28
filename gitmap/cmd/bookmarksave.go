package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// runBookmarkSave saves a new bookmark from name + command + args/flags.
func runBookmarkSave(args []string) {
	if len(args) < 2 {
		fmt.Fprint(os.Stderr, constants.ErrBookmarkSaveUsage)
		os.Exit(1)
	}

	name := args[0]
	command := args[1]
	flags, positional := splitBookmarkArgs(args[2:])

	saveBookmarkToDB(name, command, positional, flags)
}

// splitBookmarkArgs separates flags from positional arguments.
func splitBookmarkArgs(args []string) (string, string) {
	var flags, positional []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			positional = append(positional, arg)
		}
	}

	return strings.Join(flags, " "), strings.Join(positional, " ")
}

// saveBookmarkToDB persists the bookmark record.
func saveBookmarkToDB(name, command, args, flags string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkSave, err)
		os.Exit(1)
	}
	defer db.Close()

	checkBookmarkNotExists(db, name)

	record := model.BookmarkRecord{
		Name:    name,
		Command: command,
		Args:    args,
		Flags:   flags,
	}

	err = db.InsertBookmark(record)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkSave, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgBookmarkSaved, name, command, args, flags)
}

// checkBookmarkNotExists exits with error if a bookmark name is taken.
func checkBookmarkNotExists(db interface {
	FindBookmarkByName(string) (model.BookmarkRecord, error)
}, name string) {
	_, err := db.FindBookmarkByName(name)
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrBookmarkExists, name)
	os.Exit(1)
}
