package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// commonsUsage is printed for `gitmap commons --help` or on unknown flags.
const commonsUsage = `Usage: gitmap commons [flags]

Alias: co

Shortcut for 'gitmap sync all'. Adds or dedupe-merges the curated
baselines for .gitignore, .gitattributes, .prettierignore, .prettierrc,
and runs 'git lfs install --local' + the lfs/common .gitattributes
block, all in one pass.

Behavior:
  - Line-based targets append MISSING lines only; existing entries
    are preserved verbatim. Safe to re-run.
  - .prettierrc is JSON key-union: missing keys added, existing kept
    unless --force is passed.
  - Idempotent — a second run on an unchanged repo writes nothing.

Flags:
  --dry-run, -n    Print planned additions without touching disk
  --force,   -f    Overwrite conflicting JSON values in .prettierrc

Examples:
  gitmap commons
  gitmap commons --dry-run
  gitmap co -n
  gitmap commons --force
`

// dispatchCommons routes `gitmap commons` / `gitmap co` to the same
// logic as `gitmap sync all` without requiring a subcommand token.
func dispatchCommons(command string) bool {
	if command != constants.CmdCommons && command != constants.CmdCommonsAlias {
		return false
	}
	rest := os.Args[2:]
	for _, a := range rest {
		if a == "--help" || a == "-h" {
			fmt.Print(commonsUsage)

			return true
		}
	}
	dry, force := parseSyncFlags(rest)

	runSyncLines(".gitignore", defaultGitignoreBaseline, dry)
	runSyncLines(".gitattributes", defaultGitattributesBaseline, dry)
	runSyncLFSInstall(dry)
	runSyncLines(".prettierignore", defaultPrettierignoreBaseline, dry)
	runSyncPrettierRC(dry, force)

	return true
}
