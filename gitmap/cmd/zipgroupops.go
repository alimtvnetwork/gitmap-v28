package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runZipGroupRemove handles "zip-group remove <group> <path>".
func runZipGroupRemove(args []string) error {
	if len(args) < 2 {
		return apperror.New(constants.ErrZGEmpty, "E9000", nil)
	}

	groupName := args[0]
	rawPath := args[1]

	db, err := openDB()
	if err != nil {
		return apperror.Wrap(err, constants.ErrListDBFailed, nil)
	}
	defer db.Close()

	// Resolve to full path for matching.
	_, _, fullPath, _, _ := resolveZipPath(rawPath)

	err = db.RemoveZipGroupItem(groupName, fullPath)
	if err != nil {
		return apperror.Wrap(err, constants.ErrBareFmt, nil)
	}

	fmt.Printf(constants.MsgZGItemRemoved, rawPath, groupName)
	syncZipGroupJSON(db)
	return nil
}

// runZipGroupDelete handles "zip-group delete <name>".
func runZipGroupDelete(args []string) error {
	if len(args) == 0 {
		return apperror.New(constants.ErrZGEmpty, "E9000", nil)
	}

	name := args[0]

	db, err := openDB()
	if err != nil {
		return apperror.Wrap(err, constants.ErrListDBFailed, nil)
	}
	defer db.Close()

	err = db.DeleteZipGroup(name)
	if err != nil {
		return apperror.Wrap(err, constants.ErrBareFmt, nil)
	}

	fmt.Printf(constants.MsgZGDeleted, name)
	syncZipGroupJSON(db)
	return nil
}

// runZipGroupRename handles "zip-group rename <group> --archive <name>".
func runZipGroupRename(args []string) error {
	name, archiveName := parseZipGroupRenameFlags(args)
	if len(name) == 0 {
		return apperror.New(constants.ErrZGEmpty, "E9000", nil)
	}
	if len(archiveName) == 0 {
		return apperror.New(constants.FlagDescZGArchive, "E9000", nil)
	}
	executeZipGroupRename(name, archiveName)
	return nil
}

// parseZipGroupRenameFlags parses flags for zip-group rename.
func parseZipGroupRenameFlags(args []string) (name, archive string) {
	fs := flag.NewFlagSet(constants.SubCmdZGRename, flag.ExitOnError)
	archiveFlag := fs.String("archive", "", constants.FlagDescZGArchive)
	fs.Parse(args)

	if fs.NArg() > 0 {
		name = fs.Arg(0)
	}

	return name, *archiveFlag
}

// executeZipGroupRename sets a custom archive name for a group.
func executeZipGroupRename(name, archiveName string) {
	db, err := openDB()
	if err != nil {
		return apperror.Wrap(err, constants.ErrListDBFailed, nil)
	}
	defer db.Close()

	err = db.UpdateZipGroupArchive(name, archiveName)
	if err != nil {
		return apperror.Wrap(err, constants.ErrBareFmt, nil)
	}

	fmt.Printf(constants.MsgZGArchiveSet, archiveName, name)
	syncZipGroupJSON(db)
}

// syncZipGroupJSON writes zip group data to .gitmap/zip-groups.json.
func syncZipGroupJSON(db *store.DB) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to get working directory: %v\n", err)

		return
	}

	err = db.WriteZipGroupsJSON(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrZGJSONWrite+"\n", "zip-groups.json", err)
	}
}
