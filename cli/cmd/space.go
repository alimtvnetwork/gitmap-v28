package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// spaceUsage is printed for `gitmap space --help` or bare `gitmap space`.
const spaceUsage = `Usage: gitmap space <subcommand> [flags]

Space operations namespace. Subcommands:

  common      Apply the curated common baselines (.gitignore,
              .gitattributes, .prettierignore, .prettierrc) and run
              'git lfs install --local' — same logic as 'gitmap commons'.

Flags (with common):
  --dry-run, -n    Print planned additions without touching disk
  --force,   -f    Overwrite conflicting JSON values in .prettierrc

Examples:
  gitmap space common
  gitmap space common --dry-run
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

// dispatchSpace routes `gitmap space <subcommand>`; the `common`
// subcommand reuses the same baseline logic as `gitmap commons`.
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
