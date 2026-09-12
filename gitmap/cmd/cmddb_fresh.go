package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runDbFresh(args []string) error {
	msg := constants.ColorYellow + "Are you sure you want to start fresh? [y/N]: " + constants.ColorReset
	if !hasConfirmFlag(args) {
		printFreshWarning()
	}

	if !confirmOrSkip(msg, args) {
		fmt.Println(constants.ColorDim + "Start fresh operation canceled." + constants.ColorReset)

		return nil
	}

	return executeStartFresh()
}

// Backwards-compatible alias for runDbFresh

func runStartFresh(args []string) error {
	return runDbFresh(args)
}

func printFreshWarning() {
	fmt.Println()
	fmt.Println("  " + constants.ColorRed + "⚠ WARNING: Irreversible Database Transaction!" + constants.ColorReset)
	fmt.Println("  This will permanently delete all tracked repositories, scan histories,")
	fmt.Println("  search caches, profiles, and split databases across your entire system.")
	fmt.Println()
}

func executeStartFresh() error {
	removedCount := wipeAllDbFiles()
	recreateRepoSearchDir()

	freshDb, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "E9003")
	}

	defer freshDb.Close()

	if migrateErr := freshDb.Migrate(); migrateErr != nil {
		return apperror.WrapSimple(migrateErr, "E9004")
	}

	printFreshSuccess(removedCount, store.DefaultDBPath())

	return nil
}

func wipeAllDbFiles() int {
	binDir := store.BinaryDataDir()
	count := removeMatchingFiles(binDir)
	for _, splitDir := range findSplitDbDirs() {
		count += removeMatchingFiles(splitDir)
	}

	return count
}

// Backwards-compatible alias for wipeAllDbFiles

func removeMatchingFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	removed := 0
	for _, e := range entries {
		removed += tryRemoveDbFile(dir, e.Name())
	}

	return removed
}

func tryRemoveDbFile(dir, name string) int {
	if !isDbRelatedFile(name) {
		return 0
	}

	target := filepath.Join(dir, name)
	if err := os.Remove(target); err == nil {
		return 1
	}

	return 0
}

// Backwards-compatible alias for tryRemoveDbFile

func isDbRelatedFile(name string) bool {
	return strings.HasSuffix(name, ".db") ||
		strings.HasSuffix(name, ".db-wal") ||
		strings.HasSuffix(name, ".db-shm") ||
		strings.HasSuffix(name, ".db-journal")
}

// Backwards-compatible alias for isDbRelatedFile

func recreateRepoSearchDir() {
	binDir := store.BinaryDataDir()
	searchDir := filepath.Join(binDir, "repo_search")
	_ = os.MkdirAll(searchDir, 0755)
}

func printFreshSuccess(count int, dbPath string) {
	fmt.Println()
	fmt.Printf("  %s✓ All SQLite database files purged (%d file(s) removed)%s\n", constants.ColorGreen, count, constants.ColorReset)
	fmt.Printf("  %s✓ Clean schema migrations re-executed%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  %s✓ Master database ready: %s%s\n", constants.ColorGreen, dbPath, constants.ColorReset)
	fmt.Println()
	fmt.Println("  " + constants.ColorCyan + "Run 'gitmap add' to scan and index repositories from a fresh state." + constants.ColorReset)
	fmt.Println()
}
