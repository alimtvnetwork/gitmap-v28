package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunStorageCmd handles gitmap storage inspection and database listing.
func RunStorageCmd(args []string) error {
	checkHelp(constants.CmdStorage, args)
	if len(args) > 0 {
		sub := strings.ToLower(args[0])
		if isStorageListSubcommand(sub) {
			return runStorageListDatabases()
		}
	}

	return runStorageDriveReport(args)
}

func isStorageListSubcommand(sub string) bool {
	return sub == "ls" || sub == "list" || sub == "db" || sub == "dbs"
}
