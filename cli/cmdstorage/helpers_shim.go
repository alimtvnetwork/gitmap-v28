package cmdstorage

import (
	"strings"
)

func hasArgFlag(args []string, flagName string) bool {
	for _, a := range args {
		if a == flagName || strings.HasPrefix(a, flagName+"=") {
			return true
		}
	}
	return false
}

// RunBackupCloudRestoreFn is wired by cmd/di_hooks.go.
var RunBackupCloudRestoreFn func(args []string) error

func runBackupCloudRestore(args []string) error {
	if RunBackupCloudRestoreFn != nil {
		return RunBackupCloudRestoreFn(args)
	}
	return nil
}
