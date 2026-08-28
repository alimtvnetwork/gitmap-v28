package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runActiveGroupPull pulls all repos in the active group.
func runActiveGroupPull() error {
	name := requireActiveGroup()
	records := loadRecordsByGroup(name)

	for _, r := range records {
		pullOneRepo(r)
	}
	return nil
}

// runActiveGroupStatus shows status for active group repos.
func runActiveGroupStatus() error {
	name := requireActiveGroup()
	records := loadRecordsByGroup(name)

	printStatusBanner(len(records))
	summary := printStatusTable(records)
	printStatusSummary(summary)
	return nil
}

// runActiveGroupExec runs a git command across active group repos.
func runActiveGroupExec(args []string) error {
	if len(args) == 0 {
		return apperror.New(constants.ErrExecUsage, "E9000", nil)
	}

	name := requireActiveGroup()
	records := loadRecordsByGroup(name)

	printExecBanner(args, len(records))
	succeeded, failed, missing := execAllRepos(records, args)
	printExecSummary(succeeded, failed, missing, len(records))
	return nil
}

// requireActiveGroup returns the active group name or exits.
func requireActiveGroup() string {
	db, err := openDB()
	if err != nil {
		return apperror.Wrap(err, constants.ErrListDBFailed, nil)
	}
	defer db.Close()

	name := db.GetSetting(constants.SettingActiveGroup)
	if len(name) == 0 {
		return apperror.New(constants.MsgGroupNoActive, "E9000", nil)
	}

	return name
}
