package cmddb

import (
	"context"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/repodb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runDBClearAction(args []string) error {
	if !confirmOrSkip("Clear all search caches across master and split databases? [y/N]: ", args) {
		fmt.Println("Clear operation canceled.")

		return nil
	}

	clearedCount := clearAllSplitDBs(collectSplitDBs())
	fmt.Printf("%s✓ Cleared search caches across %d split database(s).%s\n",
		constants.ColorGreen, clearedCount, constants.ColorReset)

	return nil
}

func clearAllSplitDBs(splitDBs []DBFileInfo) int {
	clearedCount := 0
	for _, s := range splitDBs {
		if clearSingleSplitDB(s.Path) {
			clearedCount++
		}
	}

	return clearedCount
}

func clearSingleSplitDB(path string) bool {
	db, err := store.OpenSQLiteDB(path)
	if err != nil {
		return false
	}

	defer db.Close()

	if err := repodb.ClearRepoDB(context.Background(), db); err != nil {
		return false
	}

	return true
}
