package cmd

// CLI input validation helpers for `gitmap regoldens`. Split out so
// regoldens.go stays under the 200-line file cap.

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// validateDiffMode rejects unknown --diff values up front so the
// orchestrator can trust cfg.diffMode is "", "short", or "full".
func validateDiffMode(mode string) {
	if mode == "" || mode == constants.RegoldensDiffModeShort ||
		mode == constants.RegoldensDiffModeFull {
		return
	}
	fmt.Fprintf(os.Stderr, constants.ErrRegoldensDiffMode+"\n", mode)
	os.Exit(2)
}
