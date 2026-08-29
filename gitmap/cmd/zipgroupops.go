package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runZipGroupRemove handles "zip-group remove <group> <path>".
func runZipGroupRemove(args []string) error {
	if len(args) < 2 {
		err := apperror.NewWithDetails(
			"cmd.zipgroup.remove",
			"E1033",
			"insufficient arguments for zip-group remove; usage: gitmap zip-group remove <group> <path>",
			"cmd.zipgroup",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
		return nil
	}

	groupName := args[0]
	rawPath := args[1]

	db, err := openDB()
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.zipgroup.remove.openDB",
			"E1034",
			"failed to open database for zip-group remove",
			"cmd.zipgroup",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}
	defer db.Close()

	// Resolve to full path for matching.
	_, _, fullPath, _, _ := resolveZipPath(rawPath)

	err = db.RemoveZipGroupItem(groupName, fullPath)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.zipgroup.remove.item",
			"E1035",
			"failed to remove item from zip group",
			"cmd.zipgroup",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"group": groupName, "path": rawPath},
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}

	fmt.Printf(constants.MsgZGItemRemoved, rawPath, groupName)
	syncZipGroupJSON(db)
	return nil
}

// runZipGroupDelete handles "zip-group delete <name>".
func runZipGroupDelete(args []string) error {
	if len(args) == 0 {
		err := apperror.NewWithDetails(
			"cmd.zipgroup.delete",
			"E1036",
			"missing required group name for zip-group delete",
			"cmd.zipgroup",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
		return nil
	}

	name := args[0]

	db, err := openDB()
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.zipgroup.delete.openDB",
			"E1037",
			"failed to open database for zip-group delete",
			"cmd.zipgroup",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}
	defer db.Close()

	err = db.DeleteZipGroup(name)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.zipgroup.delete.group",
			"E1038",
			"failed to delete zip group from database",
			"cmd.zipgroup",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"name": name},
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}

	fmt.Printf(constants.MsgZGDeleted, name)
	syncZipGroupJSON(db)
	return nil
}

// runZipGroupRename handles "zip-group rename <group> --archive <name>".
func runZipGroupRename(args []string) error {
	name, archiveName := parseZipGroupRenameFlags(args)
	if len(name) == 0 {
		err := apperror.NewWithDetails(
			"cmd.zipgroup.rename",
			"E1039",
			"missing required group name for zip-group rename",
			"cmd.zipgroup",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
		return nil
	}
	if len(archiveName) == 0 {
		err := apperror.NewWithDetails(
			"cmd.zipgroup.rename.archive",
			"E1040",
			"missing required --archive flag for zip-group rename",
			"cmd.zipgroup",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
		return nil
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
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.zipgroup.rename.openDB",
			"E1041",
			"failed to open database for zip-group rename",
			"cmd.zipgroup",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return
	}
	defer db.Close()

	err = db.UpdateZipGroupArchive(name, archiveName)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.zipgroup.rename.update",
			"E1042",
			"failed to update zip group archive name",
			"cmd.zipgroup",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"name": name, "archive": archiveName},
		)
		cliexit.HandleError(appErr, 1)
		return
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
