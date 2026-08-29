package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runZipGroupList handles "zip-group list".
func runZipGroupList() error {
	db, err := openDB()
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.zipgroup.list.openDB",
			"E1043",
			"failed to open database for zip-group list",
			"cmd.zipgroup",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}
	defer db.Close()

	groups, err := db.ListZipGroupsWithCount()
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.zipgroup.list.groups",
			"E1044",
			"failed to list zip groups from database",
			"cmd.zipgroup",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}

	printZipGroupList(groups)
	printHints(zipGroupListHints())
	return nil
}

// printZipGroupList renders the zip group table to stdout.
func printZipGroupList(groups []store.ZipGroupWithCount) {
	if len(groups) == 0 {
		fmt.Println("  No zip groups defined.")

		return
	}

	fmt.Printf(constants.MsgZGListHeader, len(groups))

	for _, g := range groups {
		archive := g.ArchiveName
		if len(archive) == 0 {
			archive = g.Name + ".zip"
		}

		fmt.Printf(constants.MsgZGListRow, g.Name, g.ItemCount, archive)
	}
}

// runZipGroupShow handles "zip-group show <name>".
func runZipGroupShow(args []string) error {
	if len(args) == 0 {
		err := apperror.NewWithDetails(
			"cmd.zipgroup.show",
			"E1045",
			"missing required group name for zip-group show",
			"cmd.zipgroup",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
		return nil
	}

	name := args[0]
	executeZipGroupShow(name)
	return nil
}

// executeZipGroupShow opens the DB and displays group items with dynamic expansion.
func executeZipGroupShow(name string) {
	db, err := openDB()
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.zipgroup.show.openDB",
			"E1046",
			"failed to open database for zip-group show",
			"cmd.zipgroup",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return
	}
	defer db.Close()

	items, err := db.ListZipGroupItems(name)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.zipgroup.show.items",
			"E1047",
			"failed to list zip group items from database",
			"cmd.zipgroup",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"group": name},
		)
		cliexit.HandleError(appErr, 1)
		return
	}

	group, grpErr := db.FindZipGroupByName(name)
	if grpErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not load zip group %s: %v\n", name, grpErr)
	}
	printZipGroupShow(group, items)
}

// printZipGroupShow renders items in a zip group with dynamic folder expansion.
func printZipGroupShow(group model.ZipGroup, items []model.ZipGroupItem) {
	fmt.Printf(constants.MsgZGShowHeader, group.Name, len(items))

	for _, item := range items {
		isFolder := item.IsFolder == true
		if isFolder == true {
			fmt.Printf(constants.MsgZGShowFolder, item.RelativePath)
			fmt.Printf(constants.MsgZGShowPaths, item.RepoPath, item.RelativePath, item.FullPath)
		}

		var files []string
		if isFolder == true {
			files = expandFolder(item.FullPath)
		}
		hasFiles := len(files) > 0
		if isFolder == true && hasFiles == true {
			fmt.Printf(constants.MsgZGShowExpanded, len(files))
			for _, f := range files {
				fmt.Printf(constants.MsgZGShowExpFile, f)
			}
		}
		if isFolder == false {
			fmt.Printf(constants.MsgZGShowFile, item.RelativePath)
			fmt.Printf(constants.MsgZGShowPaths, item.RepoPath, item.RelativePath, item.FullPath)
		}
	}

	if len(group.ArchiveName) > 0 {
		fmt.Printf(constants.MsgZGShowArchive, group.ArchiveName)
	}

	printHints(zipGroupShowHints())
}

// expandFolder returns relative file paths inside a folder for display.
func expandFolder(folderPath string) []string {
	var files []string

	walkErr := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		rel, relErr := filepath.Rel(folderPath, path)
		if relErr != nil {
			rel = path
		}

		files = append(files, rel)

		return nil
	})

	if walkErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not walk folder %s: %v\n", folderPath, walkErr)
	}

	return files
}
