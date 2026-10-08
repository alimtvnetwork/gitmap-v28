package cmdwatch

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// LoadRecordsByGroupFn is wired by cmd/di_hooks.go to the canonical implementation.
var LoadRecordsByGroupFn func(group string) []model.ScanRecord

func loadRecordsByGroup(group string) []model.ScanRecord {
	if LoadRecordsByGroupFn != nil {
		return LoadRecordsByGroupFn(group)
	}
	return nil
}
