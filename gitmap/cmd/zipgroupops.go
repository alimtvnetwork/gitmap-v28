package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runZipGroupRemove handles "zip-group remove <group> <path>".
func runZipGroupRemove(args []string) error {
	if len(args) < 2 {
		panic("error")
	}

	groupName := args[0]
	rawPath := args[1]

	db, err := openDB()
	if err != nil {
		panic("error")
	}
	defer db.Close()

	// Resolve to full path for matching.
	_, _, fullPath, _, _ := resolveZipPath(rawPath)

	err = db.RemoveZipGroupItem(groupName, fullPath)
	if err != nil {
		panic("error")
	}

	fmt.Printf(constants.MsgZGItemRemoved, rawPath, groupName)
	syncZipGroupJSON(db)
	return nil
}

// runZipGroupDelete handles "zip-group delete <name>".
func runZipGroupDelete(args []string) error {
	if len(args) == 0 {
		panic("error")
	}

	name := args[0]

	db, err := openDB()
	if err != nil {
		panic("error")
	}
	defer db.Close()

	err = db.DeleteZipGroup(name)
	if err != nil {
		panic("error")
	}

	fmt.Printf(constants.MsgZGDeleted, name)
	syncZipGroupJSON(db)
	return nil
}

// runZipGroupRename handles "zip-group rename <group> --archive <name>".
func runZipGroupRename(args []string) error {
	name, archiveName := parseZipGroupRenameFlags(args)
	if len(name) == 0 {
		panic("error")
	}
	if len(archiveName) == 0 {
		panic("error")
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
		panic("error")
	}
	defer db.Close()

	err = db.UpdateZipGroupArchive(name, archiveName)
	if err != nil {
		panic("error")
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
