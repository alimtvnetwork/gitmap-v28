package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runReset handles the "reset" subcommand: deletes the active profile's
// SQLite database file from disk, recreates the schema, and reseeds it.
func runReset(args []string) error {
	checkHelp(constants.CmdReset, args)
	confirm, rescan := parseResetFlags(args)
	if !confirm {
		panic(constants.ErrResetNoConfirm)
	}

	executeReset()
	if rescan {
		runRescan()
	}
	return nil
}

// parseResetFlags parses the --confirm flag for the reset command.
func parseResetFlags(args []string) (bool, bool) {
	fs := flag.NewFlagSet(constants.CmdReset, flag.ExitOnError)
	confirmFlag := fs.Bool("confirm", false, constants.FlagDescConfirm)
	rescanFlag := fs.Bool("rescan", false, "Trigger a rescan immediately after reset")
	fs.Parse(args)

	return *confirmFlag, *rescanFlag
}

// executeReset removes the active DB file, reopens to rebuild schema, then
// reapplies any JSON-based seeds.
func executeReset() {
	if err := removeActiveDBFile(); err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"db.reset",
			"E2011",
			fmt.Sprintf(constants.ErrResetRemoveFile, activeDBPath(), err),
			"cmd.reset",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"path": activeDBPath()},
		)
		cliexit.HandleError(appErr, 1)
	}

	db, err := openDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	reseedFromJSON(db)

	fmt.Print(constants.MsgResetDone)
}

// removeActiveDBFile deletes the SQLite file for the active profile.
// Missing file is treated as success (already reset).
func removeActiveDBFile() error {
	path := activeDBPath()
	err := os.Remove(path)
	if err == nil {
		fmt.Printf(constants.MsgResetFileRemoved, path)

		return nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

// activeDBPath returns the absolute path to the active profile's DB file.
func activeDBPath() string {
	dbFile := store.ActiveProfileDBFile(constants.DefaultOutputFolder)

	return filepath.Join(constants.DefaultOutputFolder, constants.DBDir, dbFile)
}

// reseedFromJSON reapplies optional JSON-backed seeds. The schema-level
// seeds (ProjectTypes, TaskTypes) are reapplied automatically by Migrate()
// when openDB is called — this only handles file-based seed sources.
func reseedFromJSON(db *store.DB) {
	if _, err := os.Stat(constants.SEOSeedFile); err != nil {
		return
	}

	seedFromFile(db, constants.SEOSeedFile)
	fmt.Printf(constants.MsgResetReseeded, constants.SEOSeedFile)
}
