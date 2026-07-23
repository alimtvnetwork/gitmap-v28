package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runZipGroupCreate handles "zip-group create <name> [path...] [--archive <name>]".
// If paths are provided, they are resolved and added as items immediately.
func runZipGroupCreate(args []string) {
	name, archiveName, paths := parseZipGroupCreateFlags(args)
	if len(name) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}
	executeZipGroupCreate(name, archiveName, paths)
}

// parseZipGroupCreateFlags parses flags for zip-group create.
func parseZipGroupCreateFlags(args []string) (name, archive string, paths []string) {
	fs := flag.NewFlagSet(constants.SubCmdZGCreate, flag.ExitOnError)
	archiveFlag := fs.String("archive", "", constants.FlagDescZGArchive)
	fs.Parse(args)

	remaining := fs.Args()
	if len(remaining) > 0 {
		name = remaining[0]
	}
	if len(remaining) > 1 {
		paths = remaining[1:]
	}

	return name, *archiveFlag, paths
}

// executeZipGroupCreate opens the DB, creates the group, and optionally adds paths.
func executeZipGroupCreate(name, archiveName string, paths []string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.CreateZipGroup(name, archiveName)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	if len(paths) == 0 {
		fmt.Printf(constants.MsgZGCreated, name)
		printHints(zipGroupCreateHints())
		syncZipGroupJSON(db)

		return
	}

	// Add provided paths as items.
	for _, p := range paths {
		addResolvedZipGroupItem(db, name, p)
	}

	count, countErr := db.CountZipGroupItems(name)
	if countErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not count zip group items: %v\n", countErr)
	}
	fmt.Printf(constants.MsgZGCreatedPath, name, fmt.Sprintf("%d", count), "item(s)")
	printHints(zipGroupCreateHints())
	syncZipGroupJSON(db)
}

// runZipGroupAdd handles "zip-group add <group> <path...>".
func runZipGroupAdd(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, constants.ErrZGEmpty)
		os.Exit(1)
	}

	groupName := args[0]
	paths := args[1:]

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	for _, p := range paths {
		addResolvedZipGroupItem(db, groupName, p)
	}

	syncZipGroupJSON(db)
}

// addResolvedZipGroupItem resolves a path relative to CWD and adds it to the group.
func addResolvedZipGroupItem(db *store.DB, groupName, rawPath string) {
	repoPath, relativePath, fullPath, isFolder, err := resolveZipPath(rawPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrZGPathResolve+"\n", rawPath, err)

		return
	}

	err = db.AddZipGroupItem(groupName, repoPath, relativePath, fullPath, isFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

		return
	}

	pathType := constants.MsgZGTypeFile
	if isFolder {
		pathType = constants.MsgZGTypeFolder
	}

	fmt.Printf(constants.MsgZGItemAdded, rawPath, groupName, pathType)
}

// resolveZipPath resolves a raw path into repo path, relative path, full path, and type.
func resolveZipPath(rawPath string) (repoPath, relativePath, fullPath string, isFolder bool, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", "", false, err
	}

	repoPath = cwd
	relativePath = rawPath

	if filepath.IsAbs(rawPath) {
		fullPath = filepath.Clean(rawPath)
	} else {
		fullPath = filepath.Clean(filepath.Join(cwd, rawPath))
	}

	info, statErr := os.Stat(fullPath)
	if statErr != nil {
		// Path doesn't exist yet — store as file by default.
		return repoPath, relativePath, fullPath, false, nil
	}

	isFolder = info.IsDir()

	return repoPath, relativePath, fullPath, isFolder, nil
}
