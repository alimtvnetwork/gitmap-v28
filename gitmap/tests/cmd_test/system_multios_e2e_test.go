package cmd_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin"
)

func TestSystemMultiOSE2E(t *testing.T) {
	helpText := commitin.PrintCommitInHelp()
	if len(helpText) == 0 {
		t.Fatal("expected non-empty commit-in help text")
	}
}
