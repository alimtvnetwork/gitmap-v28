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

func runDBFresh(args []string) error {
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

func runStartFresh(args []string) error {
	return runDBFresh(args)
}

func printFreshWarning() {
	fmt.Println()
	fmt.Println("  " + constants.ColorRed + "⚠ WARNING: Irreversible Database Transaction!" + constants.ColorReset)
	fmt.Println("  This will permanently delete all tracked repositories, scan histories,")
	fmt.Println("  search caches, profiles, and split databases across your entire system.")
	fmt.Println()
}

func executeStartFresh() error {
	removedCount := wipeAllDBFiles()
	recreateRepoSearchDir()

	freshDB, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "E9003")
	}
	defer freshDB.Close()

	if migrateErr := freshDB.Migrate(); migrateErr != nil {
		return apperror.WrapSimple(migrateErr, "E9004")
	}
	printFreshSuccess(removedCount, store.DefaultDBPath())
	return nil
}

func wipeAllDBFiles() int {
	binDir := store.BinaryDataDir()
	count := removeMatchingFiles(binDir)
	for _, splitDir := range findSplitDBDirs() {
		count += removeMatchingFiles(splitDir)
	}
	return count
}

func removeMatchingFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	removed := 0
	for _, e := range entries {
		removed += tryRemoveDBFile(dir, e.Name())
	}
	return removed
}

func tryRemoveDBFile(dir, name string) int {
	if !isDBRelatedFile(name) {
		return 0
	}
	target := filepath.Join(dir, name)
	if err := os.Remove(target); err == nil {
		return 1
	}
	return 0
}

func isDBRelatedFile(name string) bool {
	return strings.HasSuffix(name, ".db") ||
		strings.HasSuffix(name, ".db-wal") ||
		strings.HasSuffix(name, ".db-shm") ||
		strings.HasSuffix(name, ".db-journal")
}

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
