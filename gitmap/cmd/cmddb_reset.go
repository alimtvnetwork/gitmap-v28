package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runDBResetAction(args []string) error {
	isConfirmed := hasConfirmFlag(args)
	if !isConfirmed {
		msg := constants.ColorYellow + "Are you sure you want to reset the database? All tracked repository records and split databases will be cleared. [y/N]: " + constants.ColorReset
		ok, err := promptConfirm(msg)
		if err != nil || !ok {
			fmt.Println(constants.ColorDim + "Database reset canceled." + constants.ColorReset)
			return nil
		}
	}
	return performDBReset()
}

func performDBReset() error {
	mainDB, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "E9001")
	}
	defer mainDB.Close()

	if resetErr := mainDB.Reset(); resetErr != nil {
		return apperror.WrapSimple(resetErr, "E9002")
	}
	removedSplit := clearSplitDBFiles()

	fmt.Printf("%s✓ Main database reset: %s%s\n", constants.ColorGreen, store.DefaultDBPath(), constants.ColorReset)
	fmt.Printf("%s✓ Cleared %d split repository database(s)%s\n", constants.ColorGreen, removedSplit, constants.ColorReset)
	fmt.Printf("%s✓ Database reset successfully completed.%s\n", constants.ColorGreen, constants.ColorReset)
	return nil
}

func clearSplitDBFiles() int {
	dirs := findSplitDBDirs()
	count := 0
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				target := dir + string(os.PathSeparator) + e.Name()
				if rmErr := os.Remove(target); rmErr == nil {
					count++
				}
			}
		}
	}
	return count
}
