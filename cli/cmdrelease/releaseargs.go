package cmdrelease

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
)

// reorderFlagsBeforeArgs moves flag-like arguments (starting with "-")
// before positional arguments. Go's flag package stops parsing at the
// first non-flag argument, so "gitmap release v2.55 -y" would silently
// ignore -y. This reorders to "-y v2.55" so all flags are parsed.
//
// Flags that take a value (e.g. --bump patch, -N "note") are kept
// together with their value argument.
func reorderFlagsBeforeArgs(args []string) []string {
	return flagutil.ReorderFlagsBeforeArgs(args)
}
