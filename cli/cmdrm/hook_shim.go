package cmdrm

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// ResolveMultiReposFn is wired by cmd/di_hooks.go to the canonical implementation.
var ResolveMultiReposFn func(db *store.DB, targets []string) ([]model.ScanRecord, []string)
