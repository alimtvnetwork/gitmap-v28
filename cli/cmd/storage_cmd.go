package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunStorageCmd handles gitmap storage inspection and database listing.
func RunStorageCmd(args []string) error {
	checkHelp(constants.CmdStorage, args)
	if isStorageListRequested(args) {
		return runStorageListDatabases()
	}

	return runStorageDriveReport(args)
}

func isStorageListRequested(args []string) bool {
	if len(args) == 0 {
		return false
	}

	return isStorageListSubcommand(strings.ToLower(args[0]))
}

func isStorageListSubcommand(sub string) bool {
	return sub == "ls" || sub == "list" || sub == "db" || sub == "dbs"
}
