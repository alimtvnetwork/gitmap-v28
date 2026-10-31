// Package cmd - roottooling_system.go: history system dispatch entries.
// Split from roottooling_util.go on 31-Oct-26 so roottooling_util.go stays under the 300-line file-size gate.
package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
)

func toolingSystemEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{"history clean", "history-clean"}, func() error {
			return cmdhistory.RunHistoryPurgeCLI(argsTail())
		}},
		{[]string{"history-undo", "hu", "history restore", "history-restore", "undo-history"}, func() error {
			return cmdhistory.RunHistoryUndoCLI(argsTail())
		}},
		{[]string{"history"}, func() error {
			return handleHistoryParentDispatcher(argsTail())
		}},
	}
}

func handleHistoryParentDispatcher(args []string) error {
	if len(args) == 0 {
		return cmdhistory.RunHistory(args)
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "purge", "clean", "p":
		return cmdhistory.RunHistoryPurgeCLI(args[1:])
	case "undo", "restore", "u":
		return cmdhistory.RunHistoryUndoCLI(args[1:])
	default:
		return cmdhistory.RunHistory(args)
	}
}
