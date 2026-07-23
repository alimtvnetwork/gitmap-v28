package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runZipGroupRemove handles "zip-group remove <group> <path>".
func runZipGroupRemove(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}

	groupName := args[0]
	rawPath := args[1]

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	// Resolve to full path for matching.
	_, _, fullPath, _, _ := resolveZipPath(rawPath)

	err = db.RemoveZipGroupItem(groupName, fullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgZGItemRemoved, rawPath, groupName)
	syncZipGroupJSON(db)
}

// runZipGroupDelete handles "zip-group delete <name>".
func runZipGroupDelete(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}

	name := args[0]

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.DeleteZipGroup(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgZGDeleted, name)
	syncZipGroupJSON(db)
}

// runZipGroupRename handles "zip-group rename <group> --archive <name>".
func runZipGroupRename(args []string) {
	name, archiveName := parseZipGroupRenameFlags(args)
	if len(name) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}
	if len(archiveName) == 0 {
		fmt.Fprintln(os.Stderr, constants.FlagDescZGArchive)
		os.Exit(1)
	}
	executeZipGroupRename(name, archiveName)
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
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.UpdateZipGroupArchive(name, archiveName)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
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
