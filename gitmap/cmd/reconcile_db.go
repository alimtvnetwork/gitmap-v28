package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runPruneStaleDB(dir string, currentRecords []model.ScanRecord) error {
	fmt.Printf("  " + constants.ColorDim + "→ pruning missing/stale db entries..." + constants.ColorReset)
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Printf(" [failed: %v]\n", err)

		return nil
	}
	defer db.Close()
	removed := pruneStaleRecords(db, dir, currentRecords)
	fmt.Printf(" [pruned: "+constants.ColorGreen+"ok"+constants.ColorReset+" - removed %d stale entries]\n", removed)

	return nil
}

func pruneStaleRecords(db *store.DB, dir string, currentRecords []model.ScanRecord) int {
	validPaths := make(map[string]bool, len(currentRecords))
	for _, r := range currentRecords {
		validPaths[r.AbsolutePath] = true
	}
	allRepos, err := db.ListRepos()
	if err != nil {
		fmt.Printf(" [failed: load repos]\n")

		return 0
	}

	return deleteStaleEntries(db, dir, allRepos, validPaths)
}

func deleteStaleEntries(db *store.DB, dir string, allRepos []model.ScanRecord, validPaths map[string]bool) int {
	removed := 0
	for _, repo := range allRepos {
		if !isSubPath(dir, repo.AbsolutePath) || validPaths[repo.AbsolutePath] {
			continue
		}
		if _, err := db.DeleteByPath(repo.AbsolutePath); err == nil {
			removed++
		}
	}

	return removed
}

func runReconcile(dir string, currentRecords []model.ScanRecord) error {
	return runPruneStaleDB(dir, currentRecords)
}

func isSubPath(parent, child string) bool {
	return len(child) > len(parent) && child[:len(parent)] == parent
}
