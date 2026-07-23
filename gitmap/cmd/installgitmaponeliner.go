package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runInstallGitmapOneliner prints the canonical Windows (PowerShell) and
// macOS/Linux (bash) one-liners that bootstrap the gitmap installer.
//
// The install URLs are intentionally fixed (they always point at the
// canonical `main` branch installer scripts in the alimtvnetwork repo)
// while the surrounding header, icons, and platform sections are emitted
// dynamically by Go so the output stays in sync with the current binary
// version reported by constants.Version.
//
// Spec: spec/01-app/109-install-gitmap-oneliner.md
func runInstallGitmapOneliner() {
	w := os.Stdout

	fmt.Fprintf(w, constants.MsgInstallHintHeader, constants.Version)
	fmt.Fprint(w, constants.MsgInstallHintWindows)
	fmt.Fprint(w, constants.MsgInstallHintUnix)
}
