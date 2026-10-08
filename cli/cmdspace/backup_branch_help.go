package cmdspace

import "fmt"

// backupBranchUsage is printed for `gitmap space backup-branch --help`.
// It follows the existing space usage-string pattern until spec 243.3's
// displayer lands, then migrates like any other command.
const backupBranchUsage = `Usage: gitmap space backup-branch "<task string>" [flags]

Creates a backup branch (backup/<slug>) from current HEAD. The slug is
derived from the task string: lowercase, spaces/underscores become
hyphens, only [a-z0-9-] kept, repeats collapsed, trimmed, max 60 chars.

Behavior:
  - Refuses to run when the working tree has uncommitted changes —
    backups snapshot HEAD only. There is no override.
  - Refuses when backup/<slug> already exists unless --force is passed.
  - Pushes backup/<slug> to origin by default; --no-push keeps it local.
  - Prints backup/<slug> @ <short-sha> on success.

Flags:
  --no-push      Skip pushing backup/<slug> to origin
  --force, -f    Recreate backup/<slug> at current HEAD if it exists

Examples:
  gitmap space backup-branch "CLI help displayer overhaul"
  gitmap space backup-branch "pre-release sweep" --no-push
  gitmap space backup-branch "CLI help displayer overhaul" --force
`

// printBackupBranchHelp prints the backup-branch usage text.
func printBackupBranchHelp() {
	fmt.Print(backupBranchUsage)
}
