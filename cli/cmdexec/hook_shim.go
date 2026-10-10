package cmdexec

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// GetAliasSlugFn is wired by cmd/di_hooks.go to the canonical implementation.
var GetAliasSlugFn func() string

func GetAliasSlug() string {
	if GetAliasSlugFn != nil {
		return GetAliasSlugFn()
	}
	return ""
}

// HasAliasFn is wired by cmd/di_hooks.go to the canonical implementation.
var HasAliasFn func() bool

func HasAlias() bool {
	if HasAliasFn != nil {
		return HasAliasFn()
	}
	return false
}

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// CompletePendingTaskFn is wired by cmd/di_hooks.go to the canonical implementation.
var CompletePendingTaskFn func(db *store.DB, taskID int64)

func completePendingTask(db *store.DB, taskID int64) {
	if CompletePendingTaskFn != nil {
		CompletePendingTaskFn(db, taskID)
	}
}

// CreatePendingTaskFn is wired by cmd/di_hooks.go to the canonical implementation.
var CreatePendingTaskFn func(typeName, targetPath, workDir, sourceCmd, cmdArgs string) (int64, *store.DB)

func createPendingTask(typeName, targetPath, workDir, sourceCmd, cmdArgs string) (int64, *store.DB) {
	if CreatePendingTaskFn != nil {
		return CreatePendingTaskFn(typeName, targetPath, workDir, sourceCmd, cmdArgs)
	}
	return 0, nil
}

// FailPendingTaskFn is wired by cmd/di_hooks.go to the canonical implementation.
var FailPendingTaskFn func(db *store.DB, taskID int64, reason string)

func failPendingTask(db *store.DB, taskID int64, reason string) {
	if FailPendingTaskFn != nil {
		FailPendingTaskFn(db, taskID, reason)
	}
}

// GetAliasPathFn is wired by cmd/di_hooks.go to the canonical implementation.
var GetAliasPathFn func() string

func GetAliasPath() string {
	if GetAliasPathFn != nil {
		return GetAliasPathFn()
	}
	return ""
}

// LoadRecordsByGroupFn is wired by cmd/di_hooks.go to the canonical implementation.
var LoadRecordsByGroupFn func(group string) []model.ScanRecord

func loadRecordsByGroup(group string) []model.ScanRecord {
	if LoadRecordsByGroupFn != nil {
		return LoadRecordsByGroupFn(group)
	}
	return nil
}

// LoadAllRecordsDBFn is wired by cmd/di_hooks.go to the canonical implementation.
var LoadAllRecordsDBFn func() []model.ScanRecord
