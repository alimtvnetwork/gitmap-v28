package cmdselfinstall

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
)

func reorderFlagsBeforeArgs(args []string) []string {
	return flagutil.ReorderFlagsBeforeArgs(args)
}
