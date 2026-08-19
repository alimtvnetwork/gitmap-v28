package cmd

import (
	"testing"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
)

func TestClusterCommand(t *testing.T) {
	args := []string{"--no-preflight", "ps", "echo hello"}

	// Disable pterm output in tests to prevent known data race in pterm.MultiPrinter.Stop()
	pterm.DisableOutput()
	defer pterm.EnableOutput()

	// runClusterCommand has an implicit dependency on cluster.Dispatch,
	// which currently stubs ExecPS, making this safe to run as an e2e test
	// without actually hitting the network or executing the command.
	runClusterCommand(cluster.ServersClients, args)
}
