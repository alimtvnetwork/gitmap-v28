package cmdinstall_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
)

// TestAgyInstallCommandRegistered lives in an external test package because
// cmdagy imports cmdinstall; an internal test importing cmdagy would create
// an import cycle.
func TestAgyInstallCommandRegistered(t *testing.T) {
	isFound := false
	for _, sub := range cmdagy.AgyCmd.Commands() {
		if sub.Name() == "install" {
			isFound = true

			break
		}
	}

	if !isFound {
		t.Errorf("expected 'install' subcommand to be registered in AgyCmd")
	}
}
