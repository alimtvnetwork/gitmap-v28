package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func init() {
	cmdos.StorageRunnerFn = RunStorageCmd
}

// RunStorageCmd handles gitmap storage inspection and database listing.
func RunStorageCmd(args []string) error {
	checkHelp(constants.CmdStorage, args)
	if len(args) == 0 {
		return runStorageDriveReport(args)
	}

	return routeStorageSubcommand(args)
}

func routeStorageSubcommand(args []string) error {
	sub := strings.ToLower(args[0])
	if isStorageListSubcommand(sub) {
		return runStorageListDatabases()
	}
	if isStorageCleanSubcommand(sub) {
		return runStorageClean(args[1:])
	}

	return runStorageDriveReport(stripFirstArgIfStatus(sub, args))
}

func isStorageListSubcommand(sub string) bool {
	return sub == "ls" || sub == "list" || sub == "db" || sub == "dbs"
}

func isStorageCleanSubcommand(sub string) bool {
	return sub == "clean" || sub == "clear" || sub == "prune"
}

func stripFirstArgIfStatus(sub string, args []string) []string {
	if sub == "status" || sub == "st" || sub == "info" {
		return args[1:]
	}

	return args
}
