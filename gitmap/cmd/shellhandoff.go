package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// WriteShellHandoff records `targetPath` so the shell wrapper function
// can `cd` to it after the binary exits.
//
// Mechanism: the wrapper function exports `GITMAP_HANDOFF_FILE=<tmp>`
// before invoking the binary. We write `targetPath` to that file. The
// wrapper then reads it and cds. If the env var is unset (binary called
// without the wrapper) this is a no-op — `cd` still prints the path to
// stdout for legacy capture.
//
// Spec: spec/04-generic-cli/21-post-install-shell-activation/01-contract.md
func WriteShellHandoff(targetPath string) {
	handoffFile := os.Getenv(constants.EnvGitmapHandoffFile)
	if len(handoffFile) == 0 {
		return
	}
	if len(targetPath) == 0 {
		return
	}

	if err := os.WriteFile(handoffFile, []byte(targetPath), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrShellHandoffWriteFmt, handoffFile, err)
	}
}
