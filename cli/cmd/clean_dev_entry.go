package cmd

import (
	"strings"

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

func runDevTopLevel(args []string) error {
	if len(args) > 0 {
		sub := strings.ToLower(args[0])
		if sub == "clean" || sub == "clear" || sub == "cleanup" {
			return cmdos.RunOSDevClean(args[1:])
		}
	}
	return cmdos.RunDevToolCLI(args)
}

func runCleanTopLevel(args []string) error {
	if len(args) > 0 {
		sub := strings.ToLower(args[0])
		if sub == "dev" || sub == "devs" || sub == "developer" {
			return cmdos.RunOSDevClean(args[1:])
		}
	}
	return cmdos.RunOSDevClean(args)
}
