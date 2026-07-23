// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
)

// runReleaseBranch handles the 'release-branch' command.
func runReleaseBranch(args []string) {
	checkHelp("release-branch", args)
	branch, assets, notes, draft, dryRun, verbose, noCommit, yes := parseReleaseBranchFlags(args)
	_ = verbose

	if len(branch) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrReleaseBranchUsage)
		os.Exit(1)
	}

	err := release.ExecuteFromBranch(branch, assets, notes, draft, dryRun, noCommit, yes)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}
}

// parseReleaseBranchFlags parses flags for the release-branch command.
func parseReleaseBranchFlags(args []string) (branch, assets, notes string, draft, dryRun, verbose, noCommit, yes bool) {
	fs := flag.NewFlagSet(constants.CmdReleaseBranch, flag.ExitOnError)
	assetsFlag := fs.String("assets", "", constants.FlagDescAssets)
	notesFlag := fs.String("notes", "", constants.FlagDescNotes)
	draftFlag := fs.Bool("draft", false, constants.FlagDescDraft)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	noCommitFlag := fs.Bool("no-commit", false, constants.FlagDescNoCommit)
	yesFlag := fs.Bool("yes", false, constants.FlagDescYes)

	// Register -N as shorthand for --notes, -y as shorthand for --yes.
	fs.StringVar(notesFlag, "N", "", constants.FlagDescNotes)
	fs.BoolVar(yesFlag, "y", false, constants.FlagDescYes)

	fs.Parse(args)

	branch = ""
	if fs.NArg() > 0 {
		branch = fs.Arg(0)
	}

	return branch, *assetsFlag, *notesFlag, *draftFlag, *dryRunFlag, *verboseFlag, *noCommitFlag, *yesFlag
}
