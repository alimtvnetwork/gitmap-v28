package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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
		key := fsutil.CanonicalPathKey(r.AbsolutePath)
		validPaths[key] = true
	}

	allRepos, err := db.ListRepos()
	if err != nil {
		fmt.Printf(" [failed: load repos]\n")

		return 0
	}

	return deleteStaleEntries(db, dir, allRepos, validPaths)
}

func tryDeleteStaleRepo(db *store.DB, path string) int {
	if _, err := db.DeleteByPath(path); err == nil {
		return 1
	}

	return 0
}

func isStaleCandidate(dir, path string, validPaths map[string]bool) bool {
	isChild := isSubPath(dir, path)
	if !isChild {
		return false
	}

	key := fsutil.CanonicalPathKey(path)
	isValid := validPaths[key]

	return !isValid
}

func deleteStaleEntries(db *store.DB, dir string, allRepos []model.ScanRecord, validPaths map[string]bool) int {
	removed := 0
	for _, repo := range allRepos {
		isCandidate := isStaleCandidate(dir, repo.AbsolutePath, validPaths)
		if isCandidate {
			removed += tryDeleteStaleRepo(db, repo.AbsolutePath)
		}
	}

	return removed
}

//nolint:unused
func runReconcile(dir string, currentRecords []model.ScanRecord) error {
	return runPruneStaleDB(dir, currentRecords)
}

func isSubPath(parent, child string) bool {
	return fsutil.IsSubdirectory(parent, child)
}
