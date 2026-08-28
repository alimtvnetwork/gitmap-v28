package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runMultiGroupPull pulls all repos across active multi-groups.
func runMultiGroupPull() error {
	db, records := loadMultiGroupRepos()
	defer db.Close()

	for _, r := range records {
		pullOneRepo(r)
	}
	return nil
}

// runMultiGroupStatus shows status for active multi-group repos.
func runMultiGroupStatus() error {
	db, records := loadMultiGroupRepos()
	defer db.Close()

	printStatusBanner(len(records))
	summary := printStatusTable(records)
	printStatusSummary(summary)
	return nil
}

// runMultiGroupExec runs a git command across active multi-group repos.
func runMultiGroupExec(args []string) error {
	if len(args) == 0 {
		return apperror.New(constants.ErrExecUsage, "E9000", nil)
	}

	db, records := loadMultiGroupRepos()
	defer db.Close()

	printExecBanner(args, len(records))
	succeeded, failed, missing := execAllRepos(records, args)
	printExecSummary(succeeded, failed, missing, len(records))
	return nil
}

// loadMultiGroupRepos loads all repos from the active multi-group.
func loadMultiGroupRepos() (*store.DB, []model.ScanRecord) {
	db, err := openDB()
	if err != nil {
		return apperror.Wrap(err, constants.ErrListDBFailed, nil)
	}

	names := loadMultiGroupNames(db)
	records := collectGroupRepos(db, names)

	return db, records
}

// collectGroupRepos gathers repos from multiple groups.
func collectGroupRepos(db *store.DB, names []string) []model.ScanRecord {
	var all []model.ScanRecord
	seen := make(map[int64]bool)

	for _, name := range names {
		repos := loadGroupReposSafe(db, name)
		for _, r := range repos {
			if seen[r.ID] {
				continue
			}
			seen[r.ID] = true
			all = append(all, r)
		}
	}

	return all
}

// loadGroupReposSafe loads repos for a group, printing errors.
func loadGroupReposSafe(db *store.DB, name string) []model.ScanRecord {
	repos, err := db.ShowGroup(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrMGGroupMissing, name)

		return nil
	}

	return repos
}
