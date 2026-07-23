package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runZipGroupList handles "zip-group list".
func runZipGroupList() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	groups, err := db.ListZipGroupsWithCount()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	printZipGroupList(groups)
	printHints(zipGroupListHints())
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
func runZipGroupShow(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}

	name := args[0]
	executeZipGroupShow(name)
}

// executeZipGroupShow opens the DB and displays group items with dynamic expansion.
func executeZipGroupShow(name string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	items, err := db.ListZipGroupItems(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
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
		if item.IsFolder {
			fmt.Printf(constants.MsgZGShowFolder, item.RelativePath)
			fmt.Printf(constants.MsgZGShowPaths, item.RepoPath, item.RelativePath, item.FullPath)

			// Dynamically expand folder contents.
			files := expandFolder(item.FullPath)
			if len(files) > 0 {
				fmt.Printf(constants.MsgZGShowExpanded, len(files))
				for _, f := range files {
					fmt.Printf(constants.MsgZGShowExpFile, f)
				}
			}
		} else {
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
