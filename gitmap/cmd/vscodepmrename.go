// Package cmd — vscodepmrename.go: thin wrapper around
// vscodepm.RenameByPath used by `gitmap as`. Soft-fails so an alias
// rename never aborts on a missing VS Code install.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodepm"
)

// renameVSCodePMByPath updates the projects.json entry whose rootPath
// matches absPath to use newName. Errors (missing user-data root,
// extension not installed, I/O) are routed through the existing
// soft-error reporter and never propagate.
func renameVSCodePMByPath(absPath, newName string) {
	if isVSCodeSyncDisabled() {
		return
	}

	if _, err := vscodepm.RenameByPath(absPath, newName); err != nil {
		reportVSCodePMSoftError(err)
	}
}
