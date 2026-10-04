package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// runReset handles the "reset" subcommand: deletes all database files from disk,
// recreates the schemas, and reseeds clean baseline data.
func runReset(args []string) error {
	checkHelp(constants.CmdReset, args)
	opts := cmddb.ParseResetOptionsExport(args)
	if !opts.IsDryRun && !opts.IsConfirm {
		msg := fmt.Sprintf("\n  %s⚠ Warning:%s This will permanently delete all database files and rebuild them from scratch. [y/N]: ", constants.ColorYellow, constants.ColorReset)
		if !cmddb.ConfirmOrSkip(msg, args) {
			printResetNoConfirm(opts.IsRescan)

			return nil
		}
		opts.IsConfirm = true
	}

	if err := cmddb.PerformComprehensiveResetExport(opts); err != nil {
		return err
	}

	if opts.IsRescan && !opts.IsDryRun {
		_ = runRescan()
	}

	return nil
}

func printResetNoConfirm(isRescan bool) {
	cmd := "gitmap reset --confirm"
	if isRescan {
		cmd += " --rescan"
	}

	fmt.Println()
	fmt.Printf("  %s⚠ Warning:%s This will permanently delete the database files and rebuild them from scratch.\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("  Run with %s--confirm%s (or %s-y%s) to proceed:\n\n", constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %s%s%s\n\n", constants.ColorGreen, cmd, constants.ColorReset)
}

// parseResetFlags parses the --confirm flag for the reset command.
func parseResetFlags(args []string) (bool, bool) {
	opts := cmddb.ParseResetOptionsExport(args)

	return opts.IsConfirm, opts.IsRescan
}

// executeReset removes the active DB file, reopens to rebuild schema, then
// reapplies any JSON-based seeds.
func executeReset() {
	opts := cmddb.ResetOptions{IsConfirm: true}
	_ = cmddb.PerformComprehensiveResetExport(opts)
}

// removeActiveDbFile deletes the SQLite file for the active profile.
// Missing file is treated as success (already reset).
func removeActiveDbFile() error {
	path := activeDbPath()
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

// activeDbPath returns the absolute path to the active profile's DB file.
func activeDbPath() string {
	dbFile := store.ActiveProfileDBFile(constants.DefaultOutputFolder)

	return filepath.Join(constants.DefaultOutputFolder, constants.DBDir, dbFile)
}

// Backwards-compatible aliases
//
//nolint:unused
//nolint:unused
var (
	removeActiveDBFile = removeActiveDbFile
	activeDBPath       = activeDbPath
)

// reseedFromJSON reapplies optional JSON-backed seeds. The schema-level
// seeds (ProjectTypes, TaskTypes) are reapplied automatically by Migrate()
// when openDb is called — this only handles file-based seed sources.
func reseedFromJSON(db *store.DB) {
	if _, err := os.Stat(constants.SEOSeedFile); err != nil {
		return
	}

	seedFromFile(db, constants.SEOSeedFile)
	fmt.Printf(constants.MsgResetReseeded, constants.SEOSeedFile)
}
