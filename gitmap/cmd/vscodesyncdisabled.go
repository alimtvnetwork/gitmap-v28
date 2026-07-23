// Package cmd — vscodesyncdisabled.go: global kill switch for the
// VS Code Project Manager projects.json sync.
//
// The per-command `--no-vscode-sync` flag opts out of a single
// invocation. This file adds a process-wide lever that disables the
// sync for the current gitmap run AND every subprocess it spawns
// (clone-fix-repo, reclone batches, etc.) by flipping
// GITMAP_VSCODE_SYNC_DISABLED=1.
//
// Activation paths (any one is enough):
//
//  1. Pass `--vscode-sync-disabled` (or `-vscode-sync-disabled`)
//     anywhere on the command line. It is stripped from os.Args
//     before subcommand dispatch so individual flag.FlagSets never
//     see an unknown flag.
//  2. Export GITMAP_VSCODE_SYNC_DISABLED=1 in the shell. Useful for
//     CI / headless boxes that should never touch projects.json
//     regardless of which gitmap command runs.
//
// Honored centrally by syncClonedReposToVSCodePM in clonepmsync.go,
// so every present and future clone variant inherits the behavior
// without per-call wiring.
package cmd

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// stripVSCodeSyncDisabledFlag removes every occurrence of the global
// `--vscode-sync-disabled` / `-vscode-sync-disabled` token from args
// and, if any was present, flips GITMAP_VSCODE_SYNC_DISABLED=1 for
// the current process. Returns the cleaned argv slice so callers can
// reassign os.Args before subcommand flag parsing.
func stripVSCodeSyncDisabledFlag(args []string) []string {
	short := "-" + constants.FlagVSCodeSyncDisabled
	long := "--" + constants.FlagVSCodeSyncDisabled

	out := make([]string, 0, len(args))
	wasFound := false
	for _, a := range args {
		if a == short || a == long {
			wasFound = true

			continue
		}
		out = append(out, a)
	}

	if wasFound {
		os.Setenv(constants.EnvVSCodeSyncDisabled, constants.EnvVSCodeSyncDisabledOn)
	}

	return out
}

// isVSCodeSyncDisabled reports whether the global kill switch is
// active for this process. Checked once per sync attempt.
func isVSCodeSyncDisabled() bool {
	return os.Getenv(constants.EnvVSCodeSyncDisabled) == constants.EnvVSCodeSyncDisabledOn
}
