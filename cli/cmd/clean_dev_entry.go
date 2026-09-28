package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdwinutil"
)

func runCleanDevTopLevel(args []string) error {
	return cmdos.RunOSDevClean(args)
}

func runDevToolTopLevel(args []string) error {
	return cmdos.RunDevToolCLI(args)
}

func runWinUtilTopLevel(args []string) error {
	return cmdwinutil.RunWinUtil(args)
}
