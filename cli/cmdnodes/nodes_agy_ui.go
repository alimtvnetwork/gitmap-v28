// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
)

// RunNodesAgyUI starts the embedded Antigravity studio web dashboard.
func RunNodesAgyUI(args []string) error {
	return cmdagy.RunAgyUI(args)
}
