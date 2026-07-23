package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runActiveGroupPull pulls all repos in the active group.
func runActiveGroupPull() {
	name := requireActiveGroup()
	records := loadRecordsByGroup(name)

	for _, r := range records {
		pullOneRepo(r)
	}
}

// runActiveGroupStatus shows status for active group repos.
func runActiveGroupStatus() {
	name := requireActiveGroup()
	records := loadRecordsByGroup(name)

	printStatusBanner(len(records))
	summary := printStatusTable(records)
	printStatusSummary(summary)
}

// runActiveGroupExec runs a git command across active group repos.
func runActiveGroupExec(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrExecUsage)
		os.Exit(1)
	}

	name := requireActiveGroup()
	records := loadRecordsByGroup(name)

	printExecBanner(args, len(records))
	succeeded, failed, missing := execAllRepos(records, args)
	printExecSummary(succeeded, failed, missing, len(records))
}

// requireActiveGroup returns the active group name or exits.
func requireActiveGroup() string {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	name := db.GetSetting(constants.SettingActiveGroup)
	if len(name) == 0 {
		fmt.Fprintln(os.Stderr, constants.MsgGroupNoActive)
		os.Exit(1)
	}

	return name
}
