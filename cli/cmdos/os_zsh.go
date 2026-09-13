package cmdos

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdzsh"
)

func runOSZsh(args []string) error {
	return cmdzsh.RunZsh(args)
}
