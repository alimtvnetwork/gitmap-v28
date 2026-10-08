package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdspace"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// spaceUsage is printed for `gitmap space --help` or bare `gitmap space`.
const spaceUsage = `Usage: gitmap space <subcommand> [flags]

Space operations namespace. Subcommands:

  common      Apply the curated common baselines (.gitignore,
              .gitattributes, .prettierignore, .prettierrc) and run
              'git lfs install --local' — same logic as 'gitmap commons'.
  backup-branch
              Create a backup branch (backup/<slug>) from the current
              HEAD for the given task string, then push it to origin.

Flags (with common):
  --dry-run, -n    Print planned additions without touching disk
  --force,   -f    Overwrite conflicting JSON values in .prettierrc

Flags (with backup-branch):
  --no-push       Skip pushing the new branch to origin
  --force         Recreate the branch at current HEAD if it already exists

Examples:
  gitmap space common
  gitmap space common --dry-run
  gitmap space backup-branch "CLI help displayer overhaul"
  gitmap space backup-branch "hotfix" --no-push
  gitmap space --help
`

// spaceCommonUsage is printed for `gitmap space common --help`.
const spaceCommonUsage = `Usage: gitmap space common [flags]

Applies the curated common baselines (.gitignore, .gitattributes,
.prettierignore, .prettierrc) and runs 'git lfs install --local',
all in one pass. Same logic as 'gitmap commons'.

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
  gitmap space common
  gitmap space common --dry-run
  gitmap space common --force
`

// spaceBackupBranchUsage is printed for `gitmap space backup-branch --help`.
const spaceBackupBranchUsage = `Usage: gitmap space backup-branch "<task string>" [flags]

Creates a backup branch named backup/<slug> from the current HEAD and
pushes it to origin. The slug is derived from the task string:
lowercase, spaces/underscores become hyphens, only [a-z0-9-] kept,
repeats collapsed, max 60 chars.

Behavior:
  - The working tree must be clean; a dirty tree is always refused.
  - Refuses to overwrite an existing branch unless --force is passed.
  - Prints backup/<slug> @ <short-sha> on success.

Flags:
  --no-push       Skip pushing the new branch to origin
  --force         Recreate the branch at current HEAD if it already exists

Examples:
  gitmap space backup-branch "CLI help displayer overhaul"
  gitmap space backup-branch "hotfix rollout" --no-push
`

// dispatchSpace routes `gitmap space <subcommand>`; the `common`
// subcommand reuses the same baseline logic as `gitmap commons`, and
// `backup-branch` delegates to the cmdspace implementation.
func dispatchSpace(command string) (bool, error) {
	if command != constants.CmdSpace {
		return false, nil
	}

	rest := os.Args[2:]
	if len(rest) == 0 {
		fmt.Print(spaceUsage)

		return true, nil
	}

	sub := rest[0]
	if sub == "--help" || sub == "-h" {
		fmt.Print(spaceUsage)

		return true, nil
	}

	if sub == constants.CmdSpaceBackupBranch {
		return true, runSpaceBackupBranch(rest[1:])
	}

	if sub != "common" {
		fmt.Print(spaceUsage)

		return true, fmt.Errorf("unknown space subcommand: %s", sub)
	}

	subRest := rest[1:]
	for _, a := range subRest {
		if a == "--help" || a == "-h" {
			fmt.Print(spaceCommonUsage)

			return true, nil
		}
	}

	dry, force := parseCommonFlags(subRest)

	runCommonLines(".gitignore", defaultGitignoreBaseline, dry)
	runCommonLines(".gitattributes", defaultGitattributesBaseline, dry)
	runCommonLFSInstall(dry)
	runCommonLines(".prettierignore", defaultPrettierignoreBaseline, dry)
	runCommonPrettierRC(dry, force)

	return true, nil
}

// runSpaceBackupBranch prints backup-branch help on --help/-h, else
// delegates to the cmdspace implementation.
func runSpaceBackupBranch(args []string) error {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(spaceBackupBranchUsage)

			return nil
		}
	}

	return cmdspace.RunBackupBranch(args)
}
