package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runReconcile(dir string, currentRecords []model.ScanRecord) error {
	fmt.Printf("  " + constants.ColorDim + "→ reconciling missing/stale db entries..." + constants.ColorReset)

	db, err := store.OpenDefault()
	if err != nil {
		fmt.Printf(" [failed: %v]\n", err)
		return nil
	}
	defer db.Close()

	validPaths := make(map[string]bool)
	for _, r := range currentRecords {
		validPaths[r.AbsolutePath] = true
	}

	allRepos, err := db.ListRepos()
	if err != nil {
		fmt.Printf(" [failed: load repos]\n")
		return nil
	}

	removed := 0
	for _, repo := range allRepos {
		if !isSubPath(dir, repo.AbsolutePath) || validPaths[repo.AbsolutePath] {
			continue
		}
		if _, err := db.DeleteByPath(repo.AbsolutePath); err == nil {
			removed++
		}
	}

	fmt.Printf(" [reconciled: "+constants.ColorGreen+"ok"+constants.ColorReset+" - removed %d stale entries]\n", removed)
	return nil
}

func isSubPath(parent, child string) bool {
	return len(child) > len(parent) && child[:len(parent)] == parent
}
