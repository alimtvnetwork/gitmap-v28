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

func runStartFresh(args []string) error {
	isConfirmed := hasConfirmFlag(args)
	if !isConfirmed {
		printFreshWarning()
		msg := constants.ColorYellow + "Are you sure you want to start fresh? [y/N]: " + constants.ColorReset
		ok, err := promptConfirm(msg)
		if err != nil || !ok {
			fmt.Println(constants.ColorDim + "Start fresh operation canceled." + constants.ColorReset)
			return nil
		}
	}
	return executeStartFresh()
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
		name := e.Name()
		if isDBRelatedFile(name) {
			target := filepath.Join(dir, name)
			if err := os.Remove(target); err == nil {
				removed++
			}
		}
	}
	return removed
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

func printFreshSuccess(removed int, mainPath string) {
	fmt.Println()
	fmt.Printf("  %s✓ Removed %d previous database and cache file(s).%s\n", constants.ColorGreen, removed, constants.ColorReset)
	fmt.Printf("  %s✓ Fresh master database initialized: %s%s\n", constants.ColorGreen, mainPath, constants.ColorReset)
	fmt.Printf("  %s✓ Rebuilt clean database schemas and migrations.%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  %s✓ Fresh split database directory prepared.%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Println()
	fmt.Println("  " + constants.ColorCyan + "★ Gitmap is ready for a fresh start! Run 'gitmap scan' to begin tracking repositories." + constants.ColorReset)
	fmt.Println()
}
