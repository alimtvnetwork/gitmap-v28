package cmdinstall

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestProbeToolZsh_Registration(t *testing.T) {
	bins, args := resolveToolCandidates(constants.ToolZsh)
	isBinMatch := len(bins) == 1 && bins[0] == "zsh"
	isArgMatch := len(args) == 1 && args[0] == "--version"
	isValidConfig := isBinMatch && isArgMatch
	if isValidConfig {
		return
	}

	t.Fatalf("unexpected probe config for ToolZsh: bins=%v, args=%v", bins, args)
}
