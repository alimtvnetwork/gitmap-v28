package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
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

	return result.AsError(routeStorageSubcommand(args))
}

func routeStorageSubcommand(args []string) result.ErrorWrapper {
	sub := strings.ToLower(args[0])
	if isStorageListSubcommand(sub) {
		return result.MatchWrapper(runStorageListDatabases())
	}
	if isStorageCleanSubcommand(sub) {
		return result.MatchWrapper(runStorageClean(args[1:]))
	}

	return result.MatchWrapper(runStorageDriveReport(stripFirstArgIfStatus(sub, args)))
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
