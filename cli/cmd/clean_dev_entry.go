package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
)

func runCleanDevTopLevel(args []string) error {
	return cmdos.RunOSDevClean(args)
}
