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
	if isStorageSpaceSubcommand(sub) {
		return routeStorageSpaceSubcommand(args[1:])
	}
	if isStorageListSubcommand(sub) {
		return result.MatchWrapper(runStorageListDatabases())
	}
	if isStorageCleanSubcommand(sub) {
		return result.MatchWrapper(runStorageClean(args[1:]))
	}

	return routeStorageSecondarySubcommand(sub, args)
}

func routeStorageSecondarySubcommand(sub string, args []string) result.ErrorWrapper {
	if isStorageResetErrorsSubcommand(sub) {
		return result.MatchWrapper(runStorageResetErrors(args[1:]))
	}
	if isStorageRestoreSubcommand(sub) {
		return result.MatchWrapper(runStorageRestoreDB(args[1:]))
	}

	return result.MatchWrapper(runStorageDriveReport(stripFirstArgIfStatus(sub, args)))
}

func routeStorageSpaceSubcommand(rest []string) result.ErrorWrapper {
	if len(rest) == 0 {
		return result.MatchWrapper(runStorageDriveReport(nil))
	}

	return routeStorageSubcommand(rest)
}

func isStorageSpaceSubcommand(sub string) bool {
	return sub == "space"
}

func isStorageRestoreSubcommand(sub string) bool {
	return sub == "restore-db" || sub == "restoredb" || sub == "restore"
}

func isStorageResetErrorsSubcommand(sub string) bool {
	return sub == "reset-errors" || sub == "error-reset" || sub == "reset" || sub == "clear-errors"
}

func isStorageListSubcommand(sub string) bool {
	return sub == "ls" || sub == "list" || sub == "db" || sub == "dbs"
}

func isStorageCleanSubcommand(sub string) bool {
	return sub == "clean" || sub == "clear" || sub == "prune"
}

func stripFirstArgIfStatus(sub string, args []string) []string {
	if isStatusSubcommand(sub) {
		return args[1:]
	}

	return args
}

func isStatusSubcommand(sub string) bool {
	return sub == "status" || sub == "st" || sub == "info"
}
