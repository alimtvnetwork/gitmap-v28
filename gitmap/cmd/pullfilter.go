package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// filterByAvailableUpdates keeps only records whose latest VersionProbe row
// has IsAvailable=1. Used by `gitmap pull --only-available`. Returns the
// records unchanged if the filter cannot be applied (DB open fails, etc.) —
// failing closed would surprise users who expect pull to do something.
func filterByAvailableUpdates(records []model.ScanRecord) []model.ScanRecord {
	if len(records) == 0 {
		return records
	}

	available, ok := loadAvailableRepoIDs()
	if !ok {
		fmt.Fprintln(os.Stderr, constants.WarnPullFilterFallback)
		return records
	}

	return intersectByID(records, available)
}

// loadAvailableRepoIDs returns a set of RepoIds with available updates.
// ok=false when the find-next query fails for any reason.
func loadAvailableRepoIDs() (map[int64]bool, bool) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, false
	}
	defer db.Close()

	rows, err := db.FindNext(0)
	if err != nil {
		return nil, false
	}

	set := make(map[int64]bool, len(rows))
	for _, r := range rows {
		set[r.Repo.ID] = true
	}

	return set, true
}

// intersectByID keeps only records whose ID is present in the available set.
func intersectByID(records []model.ScanRecord, available map[int64]bool) []model.ScanRecord {
	kept := make([]model.ScanRecord, 0, len(records))
	for _, rec := range records {
		if available[rec.ID] {
			kept = append(kept, rec)
		}
	}

	return kept
}
