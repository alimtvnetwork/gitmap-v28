package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/repodb"
)

func runDBClearAction(args []string) error {
	if !confirmOrSkip("Clear all search caches across master and split databases? [y/N]: ", args) {
		fmt.Println("Clear operation canceled.")
		return nil
	}

	splitDBs := collectSplitDBs()
	clearedCount := 0
	for _, s := range splitDBs {
		db, err := sql.Open("sqlite", s.Path)
		if err == nil {
			_ = repodb.ClearRepoDB(context.Background(), db)
			db.Close()
			clearedCount++
		}
	}

	fmt.Printf("%s✓ Cleared search caches across %d split database(s).%s\n",
		constants.ColorGreen, clearedCount, constants.ColorReset)
	return nil
}
