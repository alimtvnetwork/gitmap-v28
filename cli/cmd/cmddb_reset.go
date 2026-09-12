package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runDbResetAction(args []string) error {
	msg := constants.ColorYellow + "Are you sure you want to reset the database? All tracked repository records and split databases will be cleared. [y/N]: " + constants.ColorReset
	if !confirmOrSkip(msg, args) {
		fmt.Println(constants.ColorDim + "Database reset canceled." + constants.ColorReset)

		return nil
	}

	return performDbReset()
}

func performDbReset() error {
	mainDb, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "E9001")
	}

	defer mainDb.Close()

	if resetErr := mainDb.Reset(); resetErr != nil {
		return apperror.WrapSimple(resetErr, "E9002")
	}

	removedSplit := clearSplitDbFiles()

	fmt.Printf("%s✓ Main database reset: %s%s\n", constants.ColorGreen, store.DefaultDBPath(), constants.ColorReset)
	fmt.Printf("%s✓ Cleared %d split repository database(s)%s\n", constants.ColorGreen, removedSplit, constants.ColorReset)
	fmt.Printf("%s✓ Database reset successfully completed.%s\n", constants.ColorGreen, constants.ColorReset)

	return nil
}

func clearSplitDbFiles() int {
	dirs := findSplitDbDirs()
	count := 0
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, e := range entries {
			count += removeSingleSplitEntry(dir, e)
		}
	}

	return count
}

func removeSingleSplitEntry(dir string, e os.DirEntry) int {
	if e.IsDir() {
		return 0
	}

	target := dir + string(os.PathSeparator) + e.Name()
	if rmErr := os.Remove(target); rmErr == nil {
		return 1
	}

	return 0
}

// Backwards-compatible aliases
//
//nolint:unused
//nolint:unused
var (
	runDBResetAction  = runDbResetAction
	performDBReset    = performDbReset
	clearSplitDBFiles = clearSplitDbFiles
)
