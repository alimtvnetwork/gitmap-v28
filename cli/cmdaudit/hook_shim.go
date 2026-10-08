package cmdaudit

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// CreatePendingTaskFn is wired by cmd/di_hooks.go to the canonical implementation.
var CreatePendingTaskFn func(typeName, targetPath, workDir, sourceCmd, cmdArgs string) (int64, *store.DB)

func CreatePendingTask(typeName, targetPath, workDir, sourceCmd, cmdArgs string) (int64, *store.DB) {
	if CreatePendingTaskFn != nil {
		return CreatePendingTaskFn(typeName, targetPath, workDir, sourceCmd, cmdArgs)
	}
	return 0, nil
}
