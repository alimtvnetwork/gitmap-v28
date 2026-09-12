package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runBookmarkSave saves a new bookmark from name + command + args/flags.
func runBookmarkSave(args []string) *apperror.AppError {
	if len(args) < 2 {
		fmt.Fprint(os.Stderr, constants.ErrBookmarkSaveUsage)

		return apperror.NewValidationError(constants.ErrBookmarkSaveUsage)
	}

	name := args[0]
	command := args[1]
	flags, positional := splitBookmarkArgs(args[2:])

	return saveBookmarkToDB(name, command, positional, flags)
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
func saveBookmarkToDB(name, command, args, flags string) *apperror.AppError {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkSave, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrBookmarkSave, err))
	}
	defer db.Close()

	if err := checkBookmarkNotExists(db, name); err != nil {
		return err
	}

	return insertBookmarkRecord(db, name, command, args, flags)
}

func insertBookmarkRecord(db storeBookmarkWriter, name, command, args, flags string) *apperror.AppError {
	record := model.BookmarkRecord{
		Name:    name,
		Command: command,
		Args:    args,
		Flags:   flags,
	}
	err := db.InsertBookmark(record)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBookmarkSave, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrBookmarkSave, err))
	}

	fmt.Printf(constants.MsgBookmarkSaved, name, command, args, flags)

	return nil
}

type storeBookmarkWriter interface {
	InsertBookmark(model.BookmarkRecord) error
}

// checkBookmarkNotExists exits with error if a bookmark name is taken.
func checkBookmarkNotExists(db interface {
	FindBookmarkByName(string) (model.BookmarkRecord, error)
}, name string) *apperror.AppError {
	_, err := db.FindBookmarkByName(name)
	if err != nil {
		return nil
	}

	return apperror.NewSimple("fatal error", "E9000")
}
